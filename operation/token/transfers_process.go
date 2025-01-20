package token

import (
	"context"
	"fmt"
	"sync"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/state"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	tstate "github.com/ProtoconNet/mitum-token/state"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var transfersItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(TransfersItemProcessor)
	},
}

var transfersProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(TransfersProcessor)
	},
}

func (Transfers) Process(
	_ context.Context, _ base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	// NOTE Process is nil func
	return nil, nil, nil
}

type TransfersItemProcessor struct {
	//h      util.Hash
	sender base.Address
	item   *TransfersItem
}

func (opp *TransfersItemProcessor) PreProcess(
	_ context.Context, _ base.Operation, getStateFunc base.GetStateFunc,
) error {
	e := util.StringError("preprocess TransfersItemProcessor")

	if err := opp.item.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	_, err := state.ExistsCurrencyPolicy(opp.item.Currency(), getStateFunc)
	if err != nil {
		return e.Wrap(common.ErrCurrencyNF.Wrap(errors.Errorf("currency id %v", opp.item.Currency())))
	}

	if _, _, _, cErr := state.ExistsCAccount(opp.item.Receiver(), "receiver", true, false, getStateFunc); cErr != nil {
		return e.Wrap(common.ErrCAccountNA.Wrap(errors.Errorf("%v: receiver %v is contract account", cErr, opp.item.Receiver())))
	}

	return nil
}

func (opp *TransfersItemProcessor) Process(
	_ context.Context, _ base.Operation, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, error) {
	e := util.StringError("preprocess TransfersItemProcessor")

	g := tstate.NewStateKeyGenerator(opp.item.Contract().String())
	var sts []base.StateMergeValue
	receiver := opp.item.Receiver()
	amount := opp.item.Amount()
	smv, err := state.CreateNotExistAccount(receiver, getStateFunc)
	if err != nil {
		return nil, e.Wrap(err)
	} else if smv != nil {
		sts = append(sts, smv)
	}

	switch st, found, err := getStateFunc(g.TokenBalance(receiver.String())); {
	case err != nil:
		return nil, e.Wrap(err)
	case found:
		_, err := tstate.StateTokenBalanceValue(st)
		if err != nil {
			return nil, e.Wrap(err)
		}
	}

	sts = append(sts, common.NewBaseStateMergeValue(
		g.TokenBalance(receiver.String()),
		tstate.NewAddTokenBalanceStateValue(amount),
		func(height base.Height, st base.State) base.StateValueMerger {
			return tstate.NewTokenBalanceStateValueMerger(height, g.TokenBalance(receiver.String()), st)
		},
	))

	return sts, nil
}

func (opp *TransfersItemProcessor) Close() {
	//opp.h = nil
	opp.item = nil

	transfersItemProcessorPool.Put(opp)
}

type TransfersProcessor struct {
	*base.BaseOperationProcessor
	required map[string]common.Big
}

func NewTransfersProcessor() types.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringError("create new TransfersProcessor")

		nopp := transfersProcessorPool.Get()
		opp, ok := nopp.(*TransfersProcessor)
		if !ok {
			return nil, e.Wrap(errors.Errorf("expected TransfersProcessor, not %T", nopp))
		}

		b, err := base.NewBaseOperationProcessor(
			height, getStateFunc, newPreProcessConstraintFunc, newProcessConstraintFunc)
		if err != nil {
			return nil, e.Wrap(err)
		}

		opp.BaseOperationProcessor = b
		opp.required = nil

		return opp, nil
	}
}

func (opp *TransfersProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	fact, ok := op.Fact().(TransfersFact)
	if !ok {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMTypeMismatch).Errorf("expected %T, not %T", TransfersFact{}, op.Fact()),
		), nil
	}

	var required = make(map[string]common.Big)
	for i := range fact.Items() {
		v, found := required[fact.Items()[i].contract.String()]
		if !found {
			required[fact.Items()[i].contract.String()] = fact.Items()[i].amount
		} else {
			required[fact.Items()[i].contract.String()] = v.Add(fact.Items()[i].amount)
		}
	}
	for ca, am := range required {
		g := tstate.NewStateKeyGenerator(ca)

		if err := state.CheckExistsState(g.Design(), getStateFunc); err != nil {
			return nil, base.NewBaseOperationProcessReasonError(
				common.ErrMPreProcess.
					Wrap(common.ErrMServiceNF).Errorf("token design for contract account %v",
					ca,
				)), nil
		}

		st, err := state.ExistsState(g.TokenBalance(fact.Sender().String()), "token balance", getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError(
				common.ErrMPreProcess.Wrap(common.ErrMStateNF).
					Errorf("token balance of sender %v in contract account %v", fact.Sender(), ca)), nil
		}

		tb, err := tstate.StateTokenBalanceValue(st)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError(
				common.ErrMPreProcess.Wrap(common.ErrMStateValInvalid).
					Errorf("token balance of sender %v in contract account %v", fact.Sender(), ca)), nil
		}

		if tb.Compare(am) < 0 {
			return nil, base.NewBaseOperationProcessReasonError(
				common.ErrMPreProcess.Wrap(common.ErrMValueInvalid).
					Errorf("token balance of sender %v is less than amount to transfer in contract account %v, %v < %v",
						fact.Sender(), ca, tb, am)), nil
		}
	}
	opp.required = required

	var wg sync.WaitGroup
	errChan := make(chan *base.BaseOperationProcessReasonError, len(fact.items))
	for i := range fact.items {
		wg.Add(1)
		go func(item TransfersItem) {
			defer wg.Done()
			tip := transfersItemProcessorPool.Get()
			t, ok := tip.(*TransfersItemProcessor)
			if !ok {
				err := base.NewBaseOperationProcessReasonError(
					common.ErrMPreProcess.Wrap(
						common.ErrMTypeMismatch).Errorf("expected %T, not %T", &TransfersItemProcessor{}, tip))
				errChan <- &err
				return
			}

			t.sender = fact.Sender()
			t.item = &item

			if err := t.PreProcess(ctx, op, getStateFunc); err != nil {
				err := base.NewBaseOperationProcessReasonError(common.ErrMPreProcess.Errorf("%v", err))
				errChan <- &err
				return
			}
			t.Close()
		}(fact.items[i])
	}
	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			return nil, *err, nil
		}
	}

	return ctx, nil, nil
}

func (opp *TransfersProcessor) Process( // nolint:dupl
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	fact, ok := op.Fact().(TransfersFact)
	if !ok {
		return nil, base.NewBaseOperationProcessReasonError("expected %T, not %T", TransfersFact{}, op.Fact()), nil
	}

	var stateMergeValues []base.StateMergeValue // nolint:prealloc
	var wg sync.WaitGroup
	var mu sync.Mutex
	errChan := make(chan *base.BaseOperationProcessReasonError, len(fact.items))
	for i := range fact.items {
		wg.Add(1)
		go func(item TransfersItem) {
			defer wg.Done()
			cip := transfersItemProcessorPool.Get()
			c, ok := cip.(*TransfersItemProcessor)
			if !ok {
				err := base.NewBaseOperationProcessReasonError("expected %T, not %T", &TransfersItemProcessor{}, cip)
				errChan <- &err
				return
			}

			c.sender = fact.Sender()
			c.item = &item

			s, err := c.Process(ctx, op, getStateFunc)
			if err != nil {
				err := base.NewBaseOperationProcessReasonError("process transfers item: %w", err)
				errChan <- &err
				return
			}
			mu.Lock()
			stateMergeValues = append(stateMergeValues, s...)
			mu.Unlock()
			c.Close()
		}(fact.items[i])
	}
	go func() {
		wg.Wait()
		close(errChan)
	}()
	for err := range errChan {
		if err != nil {
			return nil, *err, nil
		}
	}

	totalAmounts, err := PrepareSenderTotalAmounts(fact.Sender(), opp.required, getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("process Transfers; %w", err), nil
	}

	for key, total := range totalAmounts {
		stateMergeValues = append(
			stateMergeValues,
			common.NewBaseStateMergeValue(
				key,
				tstate.NewDeductTokenBalanceStateValue(total),
				func(height base.Height, st base.State) base.StateValueMerger {
					return tstate.NewTokenBalanceStateValueMerger(height, key, st)
				}),
		)
	}

	return stateMergeValues, nil, nil
}

func (opp *TransfersProcessor) Close() error {
	transfersProcessorPool.Put(opp)

	return nil
}

func PrepareSenderTotalAmounts(
	holder base.Address,
	required map[string]common.Big,
	getStateFunc base.GetStateFunc,
) (map[string]common.Big, error) {
	totalAmounts := map[string]common.Big{}

	for ca, rq := range required {
		g := tstate.NewStateKeyGenerator(ca)
		_, err := state.ExistsState(g.TokenBalance(holder.String()), fmt.Sprintf("token balance, %v", holder), getStateFunc)
		if err != nil {
			return nil, err
		}

		totalAmounts[g.TokenBalance(holder.String())] = rq
	}

	return totalAmounts, nil
}
