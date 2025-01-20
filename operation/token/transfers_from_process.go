package token

import (
	"context"
	"sync"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	cstate "github.com/ProtoconNet/mitum-currency/v3/state"
	ctypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-token/state"
	"github.com/ProtoconNet/mitum-token/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var transfersFromItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(TransfersFromItemProcessor)
	},
}

var transfersFromProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(TransfersFromProcessor)
	},
}

func (TransfersFrom) Process(
	_ context.Context, _ base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	// NOTE Process is nil func
	return nil, nil, nil
}

type TransfersFromItemProcessor struct {
	sender  base.Address
	item    *TransfersFromItem
	designs map[string]types.Design
}

func (opp *TransfersFromItemProcessor) PreProcess(
	_ context.Context, _ base.Operation, getStateFunc base.GetStateFunc,
) error {
	e := util.StringError("preprocess TransfersFromItemProcessor")

	if err := opp.item.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	if _, _, _, cErr := cstate.ExistsCAccount(opp.item.Receiver(), "receiver", true, false, getStateFunc); cErr != nil {
		return e.Wrap(common.ErrCAccountNA.Wrap(errors.Errorf("%v: receiver %v is contract account", cErr, opp.item.Receiver())))
	}

	if _, _, aErr, cErr := cstate.ExistsCAccount(opp.item.Target(), "target", true, false, getStateFunc); aErr != nil {
		return e.Wrap(aErr)
	} else if cErr != nil {
		return e.Wrap(common.ErrCAccountNA.Wrap(errors.Errorf("%v", cErr)))
	}

	g := state.NewStateKeyGenerator(opp.item.Contract().String())

	design, _ := opp.designs[opp.item.Contract().String()]
	approveBoxList := design.Policy().ApproveList()

	idx := -1
	for i, apb := range approveBoxList {
		if apb.Account().Equal(opp.item.Target()) {
			idx = i
			break
		}
	}

	if idx < 0 {
		return e.Wrap(common.ErrAccountNAth.Wrap(
			errors.Errorf(
				"target %v has not approved any accounts in contract account %v", opp.item.Target(), opp.item.Contract())))
	}

	aprInfo := approveBoxList[idx].GetApproveInfo(opp.sender)
	if aprInfo == nil {
		return e.Wrap(common.ErrAccountNAth.Wrap(errors.Errorf(
			"sender %v has not been approved by target %v in contract account %v",
			opp.sender, opp.item.Target(), opp.item.Contract())))
	}

	if aprInfo.Amount().Compare(opp.item.Amount()) < 0 {
		return e.Wrap(common.ErrValueInvalid.Wrap(errors.Errorf(
			"approved amount of sender %v is less than amount to transfer in contract account %v, %v < %v",
			opp.sender, opp.item.Contract(), aprInfo.Amount(), opp.item.Amount())))
	}

	st, err := cstate.ExistsState(g.TokenBalance(opp.item.Target().String()), "token balance", getStateFunc)
	if err != nil {
		return e.Wrap(common.ErrStateNF.Wrap(errors.Errorf(
			"token balance of target %v in contract account %v", opp.item.Target(), opp.item.Contract())))
	}

	tb, err := state.StateTokenBalanceValue(st)
	if err != nil {
		return e.Wrap(common.ErrStateValInvalid.Wrap(errors.Errorf(
			"token balance of target %v in contract account %v", opp.item.Target(), opp.item.Contract())))
	}

	if tb.Compare(opp.item.Amount()) < 0 {
		return e.Wrap(common.ErrValueInvalid.Wrap(errors.Errorf(
			"token balance of target %v is less than amount to transfer-from in contract account %v, %v < %v",
			opp.item.Target(), opp.item.Contract(), tb, opp.item.Amount())))
	}

	return nil
}

func (opp *TransfersFromItemProcessor) Process(
	_ context.Context, _ base.Operation, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, error) {
	e := util.StringError("process TransfersFromItemProcessor")

	g := state.NewStateKeyGenerator(opp.item.Contract().String())
	var sts []base.StateMergeValue

	design, _ := opp.designs[opp.item.Contract().String()]
	approveBoxList := design.Policy().ApproveList()

	idx := -1
	for i, apb := range approveBoxList {
		if apb.Account().Equal(opp.item.Target()) {
			idx = i
			break
		}
	}

	if idx < 0 {
		return nil, e.Wrap(common.ErrAccountNAth.Wrap(
			errors.Errorf(
				"target %v has not approved any accounts in contract account %v during process operation", opp.item.Target(), opp.item.Contract())))
	}

	apb := approveBoxList[idx]
	am := apb.GetApproveInfo(opp.sender).Amount().Sub(opp.item.Amount())

	if am.IsZero() {
		err := apb.RemoveApproveInfo(opp.sender)
		if err != nil {
			return nil, e.Wrap(err)
		}
	} else {
		apb.SetApproveInfo(types.NewApproveInfo(opp.sender, am))
	}

	approveBoxList[idx] = apb

	receiver := opp.item.Receiver()
	amount := opp.item.Amount()
	smv, err := cstate.CreateNotExistAccount(receiver, getStateFunc)
	if err != nil {
		return nil, e.Wrap(err)
	} else if smv != nil {
		sts = append(sts, smv)
	}

	switch st, found, err := getStateFunc(g.TokenBalance(receiver.String())); {
	case err != nil:
		return nil, e.Wrap(err)
	case found:
		_, err := state.StateTokenBalanceValue(st)
		if err != nil {
			return nil, e.Wrap(err)
		}
	}

	sts = append(sts, common.NewBaseStateMergeValue(
		g.TokenBalance(receiver.String()),
		state.NewAddTokenBalanceStateValue(amount),
		func(height base.Height, st base.State) base.StateValueMerger {
			return state.NewTokenBalanceStateValueMerger(height, g.TokenBalance(receiver.String()), st)
		},
	))

	return sts, nil
}

func (opp *TransfersFromItemProcessor) Close() {
	opp.item = nil

	transfersFromItemProcessorPool.Put(opp)
}

type TransfersFromProcessor struct {
	*base.BaseOperationProcessor
}

func NewTransfersFromProcessor() ctypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringError("create new TransfersFromProcessor")

		nopp := transfersFromProcessorPool.Get()
		opp, ok := nopp.(*TransfersFromProcessor)
		if !ok {
			return nil, e.Wrap(errors.Errorf("expected TransfersFromProcessor, not %T", nopp))
		}

		b, err := base.NewBaseOperationProcessor(
			height, getStateFunc, newPreProcessConstraintFunc, newProcessConstraintFunc)
		if err != nil {
			return nil, e.Wrap(err)
		}

		opp.BaseOperationProcessor = b

		return opp, nil
	}
}

func (opp *TransfersFromProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	fact, ok := op.Fact().(TransfersFromFact)
	if !ok {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMTypeMismatch).Errorf("expected %T, not %T", TransfersFromFact{}, op.Fact()),
		), nil
	}

	requiredMap := make(map[string]map[string]common.Big)
	designs := make(map[string]types.Design)
	for i := range fact.Items() {
		required, found := requiredMap[fact.Items()[i].Target().String()]
		if !found {
			rq := make(map[string]common.Big)
			rq[fact.Items()[i].Contract().String()] = fact.Items()[i].Amount()
			requiredMap[fact.Items()[i].Target().String()] = rq
		} else {
			rq, found := required[fact.Items()[i].Contract().String()]
			if !found {
				required[fact.Items()[i].Contract().String()] = fact.Items()[i].Amount()
			} else {
				required[fact.Items()[i].Contract().String()] = rq.Add(fact.Items()[i].Amount())
			}
		}

		keyGenerator := state.NewStateKeyGenerator(fact.Items()[i].Contract().String())
		st, _ := cstate.ExistsState(keyGenerator.Design(), "design", getStateFunc)
		design, _ := state.StateDesignValue(st)
		if _, found := designs[fact.Items()[i].contract.String()]; !found {
			designs[fact.Items()[i].contract.String()] = *design
		}
	}

	for holder, required := range requiredMap {
		_, err := PrepareSenderTotalAmounts(holder, required, getStateFunc)
		if err != nil {
			return ctx, base.NewBaseOperationProcessReasonError(
				common.ErrMPreProcess.
					Errorf("%v", err)), nil
		}
	}

	for i := range fact.Items() {
		tip := transfersFromItemProcessorPool.Get()
		t, ok := tip.(*TransfersFromItemProcessor)
		if !ok {
			return ctx, base.NewBaseOperationProcessReasonError(
				common.ErrMPreProcess.Wrap(
					common.ErrMTypeMismatch).Errorf("expected %T, not %T", &TransfersFromItemProcessor{}, tip)), nil
		}

		item := fact.items[i]
		t.sender = fact.Sender()
		t.item = &item
		t.designs = designs

		if err := t.PreProcess(ctx, op, getStateFunc); err != nil {
			return ctx, base.NewBaseOperationProcessReasonError(
				common.ErrMPreProcess.
					Errorf("%v", err)), nil
		}
		t.Close()
	}

	return ctx, nil, nil
}

func (opp *TransfersFromProcessor) Process( // nolint:dupl
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	fact, ok := op.Fact().(TransfersFromFact)
	if !ok {
		return nil, base.NewBaseOperationProcessReasonError("expected %T, not %T", TransfersFromFact{}, op.Fact()), nil
	}

	requiredMap := make(map[string]map[string]common.Big)
	designs := make(map[string]types.Design)
	for i := range fact.Items() {
		required, found := requiredMap[fact.Items()[i].Target().String()]
		if !found {
			rq := make(map[string]common.Big)
			rq[fact.Items()[i].Contract().String()] = fact.Items()[i].Amount()
			requiredMap[fact.Items()[i].Target().String()] = rq
		} else {
			rq, found := required[fact.Items()[i].Contract().String()]
			if !found {
				required[fact.Items()[i].Contract().String()] = fact.Items()[i].Amount()
			} else {
				required[fact.Items()[i].Contract().String()] = rq.Add(fact.Items()[i].Amount())
			}
		}

		keyGenerator := state.NewStateKeyGenerator(fact.Items()[i].Contract().String())
		st, _ := cstate.ExistsState(keyGenerator.Design(), "design", getStateFunc)
		design, _ := state.StateDesignValue(st)
		if _, found := designs[fact.Items()[i].contract.String()]; !found {
			designs[fact.Items()[i].Contract().String()] = *design
		}
	}

	var stateMergeValues []base.StateMergeValue // nolint:prealloc
	for i := range fact.Items() {
		cip := transfersFromItemProcessorPool.Get()
		c, ok := cip.(*TransfersFromItemProcessor)
		if !ok {
			return nil, base.NewBaseOperationProcessReasonError("expected %T, not %T", &TransfersFromItemProcessor{}, cip), nil
		}

		item := fact.Items()[i]
		c.sender = fact.Sender()
		c.item = &item
		c.designs = designs

		s, err := c.Process(ctx, op, getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("process transfersFrom item: %w", err), nil
		}
		stateMergeValues = append(stateMergeValues, s...)
		c.Close()
	}

	for holder, required := range requiredMap {
		totalAmounts, _ := PrepareSenderTotalAmounts(holder, required, getStateFunc)

		for key, total := range totalAmounts {
			stateMergeValues = append(
				stateMergeValues,
				common.NewBaseStateMergeValue(
					key,
					state.NewDeductTokenBalanceStateValue(total),
					func(height base.Height, st base.State) base.StateValueMerger {
						return state.NewTokenBalanceStateValueMerger(height, key, st)
					}),
			)
		}
	}

	return stateMergeValues, nil, nil
}

func (opp *TransfersFromProcessor) Close() error {
	transfersFromProcessorPool.Put(opp)

	return nil
}
