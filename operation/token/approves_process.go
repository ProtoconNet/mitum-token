package token

import (
	"context"
	"sync"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/state"
	cstate "github.com/ProtoconNet/mitum-currency/v3/state"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	tstate "github.com/ProtoconNet/mitum-token/state"
	ttype "github.com/ProtoconNet/mitum-token/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var approvesItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(ApprovesItemProcessor)
	},
}

var approvesProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(ApprovesProcessor)
	},
}

func (Approves) Process(
	_ context.Context, _ base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	// NOTE Process is nil func
	return nil, nil, nil
}

type ApprovesItemProcessor struct {
	//h      util.Hash
	sender base.Address
	item   *ApprovesItem
}

func (opp *ApprovesItemProcessor) PreProcess(
	_ context.Context, _ base.Operation, getStateFunc base.GetStateFunc,
) error {
	e := util.StringError("preprocess ApprovesItemProcessor")

	if err := opp.item.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	if _, _, _, cErr := state.ExistsCAccount(opp.item.Approved(), "approved", true, false, getStateFunc); cErr != nil {
		return e.Wrap(common.ErrCAccountNA.Wrap(errors.Errorf("%v: approved %v is contract account", cErr, opp.item.Approved())))
	}

	keyGenerator := tstate.NewStateKeyGenerator(opp.item.Contract().String())

	if st, err := state.ExistsState(
		keyGenerator.Design(), "design", getStateFunc); err != nil {
		return e.Wrap(common.ErrServiceNF.Wrap(errors.Errorf("token design state for contract account %v",
			opp.item.Contract(),
		)))
	} else if design, err := tstate.StateDesignValue(st); err != nil {
		return e.Wrap(common.ErrServiceNF.Wrap(errors.Errorf("token design state value for contract account %v",
			opp.item.Contract(),
		)))
	} else if apb := design.Policy().GetApproveBox(opp.sender); apb == nil {
		if opp.item.Amount().IsZero() {
			return e.Wrap(common.ErrValueInvalid.Wrap(errors.Errorf("sender %v has not approved any accounts", opp.sender)))
		}
	} else if aprInfo := apb.GetApproveInfo(opp.item.Approved()); aprInfo == nil {
		if opp.item.Amount().IsZero() {
			return e.Wrap(common.ErrValueInvalid.Wrap(errors.Errorf("approved account %v has not been approved",
				opp.item.Approved())))
		}
	}
	if err := state.CheckExistsState(keyGenerator.TokenBalance(opp.sender.String()), getStateFunc); err != nil {
		return e.Wrap(common.ErrStateNF.Wrap(errors.Errorf("token balance for sender %v in contract account %v", opp.sender, opp.item.Contract())))
	}

	return nil
}

func (opp *ApprovesItemProcessor) Process(
	_ context.Context, _ base.Operation, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, error) {
	e := util.StringError("preprocess ApprovesItemProcessor")

	g := tstate.NewStateKeyGenerator(opp.item.Contract().String())

	var sts []base.StateMergeValue

	smv, err := state.CreateNotExistAccount(opp.item.Approved(), getStateFunc)
	if err != nil {
		return nil, e.Wrap(err)
	} else if smv != nil {
		sts = append(sts, smv)
	}

	st, _ := state.ExistsState(g.Design(), "design", getStateFunc)
	design, _ := tstate.StateDesignValue(st)
	apb := design.Policy().GetApproveBox(opp.sender)
	if apb == nil {
		a := ttype.NewApproveBox(opp.sender, []ttype.ApproveInfo{ttype.NewApproveInfo(opp.item.Approved(), opp.item.Amount())})
		apb = &a
	} else {
		if opp.item.Amount().IsZero() {
			err := apb.RemoveApproveInfo(opp.item.Approved())
			if err != nil {
				return nil, e.Wrap(errors.Errorf("remove approved, %s: %w", opp.item.Approved().String(), err))
			}
		} else {
			apb.SetApproveInfo(ttype.NewApproveInfo(opp.item.Approved(), opp.item.Amount()))
		}
	}

	policy := design.Policy()
	policy.MergeApproveBox(*apb)
	if err := policy.IsValid(nil); err != nil {
		return nil, ErrInvalid(policy, err)
	}
	de := ttype.NewDesign(design.Symbol(), design.Name(), design.Decimal(), policy)
	if err := de.IsValid(nil); err != nil {
		return nil, ErrInvalid(de, err)
	}
	sts = append(sts, cstate.NewStateMergeValue(
		g.Design(),
		tstate.NewDesignStateValue(de),
	))

	return sts, nil
}

func (opp *ApprovesItemProcessor) Close() {
	//opp.h = nil
	opp.item = nil

	approvesItemProcessorPool.Put(opp)
}

type ApprovesProcessor struct {
	*base.BaseOperationProcessor
}

func NewApprovesProcessor() types.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringError("create new ApprovesProcessor")

		nopp := approvesProcessorPool.Get()
		opp, ok := nopp.(*ApprovesProcessor)
		if !ok {
			return nil, e.Wrap(errors.Errorf("expected ApprovesProcessor, not %T", nopp))
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

func (opp *ApprovesProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	fact, ok := op.Fact().(ApprovesFact)
	if !ok {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMTypeMismatch).Errorf("expected %T, not %T", ApprovesFact{}, op.Fact()),
		), nil
	}

	var wg sync.WaitGroup
	errChan := make(chan *base.BaseOperationProcessReasonError, len(fact.items))
	for i := range fact.items {
		wg.Add(1)
		go func(item ApprovesItem) {
			defer wg.Done()
			tip := approvesItemProcessorPool.Get()
			t, ok := tip.(*ApprovesItemProcessor)
			if !ok {
				err := base.NewBaseOperationProcessReasonError(
					common.ErrMPreProcess.Wrap(
						common.ErrMTypeMismatch).Errorf("expected %T, not %T", &ApprovesItemProcessor{}, tip))
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

func (opp *ApprovesProcessor) Process( // nolint:dupl
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	fact, ok := op.Fact().(ApprovesFact)
	if !ok {
		return nil, base.NewBaseOperationProcessReasonError("expected %T, not %T", ApprovesFact{}, op.Fact()), nil
	}

	var stateMergeValues []base.StateMergeValue // nolint:prealloc
	var wg sync.WaitGroup
	var mu sync.Mutex
	errChan := make(chan *base.BaseOperationProcessReasonError, len(fact.items))
	for i := range fact.items {
		wg.Add(1)
		go func(item ApprovesItem) {
			defer wg.Done()
			cip := approvesItemProcessorPool.Get()
			c, ok := cip.(*ApprovesItemProcessor)
			if !ok {
				err := base.NewBaseOperationProcessReasonError("expected %T, not %T", &ApprovesItemProcessor{}, cip)
				errChan <- &err
				return
			}

			c.sender = fact.Sender()
			c.item = &item

			s, err := c.Process(ctx, op, getStateFunc)
			if err != nil {
				err := base.NewBaseOperationProcessReasonError("process approves item: %w", err)
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

	return stateMergeValues, nil, nil
}

func (opp *ApprovesProcessor) Close() error {
	approvesProcessorPool.Put(opp)

	return nil
}
