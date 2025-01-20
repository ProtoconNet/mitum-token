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
	sender  base.Address
	item    *ApprovesItem
	designs map[string]types.Design
}

func (opp *ApprovesItemProcessor) PreProcess(
	_ context.Context, _ base.Operation, getStateFunc base.GetStateFunc,
) error {
	e := util.StringError("preprocess ApprovesItemProcessor")

	if err := opp.item.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	if _, _, _, cErr := cstate.ExistsCAccount(opp.item.Approved(), "approved", true, false, getStateFunc); cErr != nil {
		return e.Wrap(common.ErrCAccountNA.Wrap(errors.Errorf("%v: approved %v is contract account", cErr, opp.item.Approved())))
	}

	keyGenerator := state.NewStateKeyGenerator(opp.item.Contract().String())

	if st, err := cstate.ExistsState(
		keyGenerator.Design(), "design", getStateFunc); err != nil {
		return e.Wrap(common.ErrServiceNF.Wrap(errors.Errorf("token design state for contract account %v",
			opp.item.Contract(),
		)))
	} else if design, err := state.StateDesignValue(st); err != nil {
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
	if err := cstate.CheckExistsState(keyGenerator.TokenBalance(opp.sender.String()), getStateFunc); err != nil {
		return e.Wrap(common.ErrStateNF.Wrap(errors.Errorf("token balance for sender %v in contract account %v", opp.sender, opp.item.Contract())))
	}

	return nil
}

func (opp *ApprovesItemProcessor) Process(
	_ context.Context, _ base.Operation, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, error) {
	e := util.StringError("preprocess ApprovesItemProcessor")

	var sts []base.StateMergeValue

	smv, err := cstate.CreateNotExistAccount(opp.item.Approved(), getStateFunc)
	if err != nil {
		return nil, e.Wrap(err)
	} else if smv != nil {
		sts = append(sts, smv)
	}

	design, _ := opp.designs[opp.item.Contract().String()]
	apb := design.Policy().GetApproveBox(opp.sender)
	if apb == nil {
		a := types.NewApproveBox(opp.sender, []types.ApproveInfo{types.NewApproveInfo(opp.item.Approved(), opp.item.Amount())})
		apb = &a
	} else {
		if opp.item.Amount().IsZero() {
			err := apb.RemoveApproveInfo(opp.item.Approved())
			if err != nil {
				return nil, e.Wrap(errors.Errorf("remove approved, %s: %w", opp.item.Approved().String(), err))
			}
		} else {
			apbInfo := apb.GetApproveInfo(opp.item.Approved())
			if apbInfo == nil {
				apb.SetApproveInfo(types.NewApproveInfo(opp.item.Approved(), opp.item.Amount()))
			} else {
				apb.SetApproveInfo(types.NewApproveInfo(opp.item.Approved(), apbInfo.Amount().Add(opp.item.Amount())))
			}

		}
	}

	policy := design.Policy()
	policy.MergeApproveBox(*apb)
	if err := policy.IsValid(nil); err != nil {
		return nil, ErrInvalid(policy, err)
	}
	de := types.NewDesign(design.Symbol(), design.Name(), design.Decimal(), policy)
	if err := de.IsValid(nil); err != nil {
		return nil, ErrInvalid(de, err)
	}
	opp.designs[opp.item.Contract().String()] = de

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

func NewApprovesProcessor() ctypes.GetNewProcessor {
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

	for i := range fact.items {
		tip := approvesItemProcessorPool.Get()
		t, ok := tip.(*ApprovesItemProcessor)
		if !ok {
			return ctx, base.NewBaseOperationProcessReasonError(
				common.ErrMPreProcess.Wrap(
					common.ErrMTypeMismatch).Errorf("expected %T, not %T", &ApprovesItemProcessor{}, tip)), nil

		}

		item := fact.items[i]
		t.sender = fact.Sender()
		t.item = &item

		if err := t.PreProcess(ctx, op, getStateFunc); err != nil {
			return ctx, base.NewBaseOperationProcessReasonError(
				common.ErrMPreProcess.
					Errorf("%v", err)), nil
		}
		t.Close()

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

	designs := make(map[string]types.Design)
	for i := range fact.items {
		keyGenerator := state.NewStateKeyGenerator(fact.items[i].Contract().String())
		st, _ := cstate.ExistsState(keyGenerator.Design(), "design", getStateFunc)
		design, _ := state.StateDesignValue(st)
		if _, found := designs[fact.items[i].contract.String()]; !found {
			designs[fact.items[i].contract.String()] = *design
		}
	}

	var stateMergeValues []base.StateMergeValue // nolint:prealloc
	for i := range fact.items {
		cip := approvesItemProcessorPool.Get()
		c, ok := cip.(*ApprovesItemProcessor)
		if !ok {
			return nil, base.NewBaseOperationProcessReasonError("expected %T, not %T", &ApprovesItemProcessor{}, cip), nil
		}

		item := fact.items[i]
		c.sender = fact.Sender()
		c.item = &item
		c.designs = designs

		s, err := c.Process(ctx, op, getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("process approves item: %w", err), nil
		}
		stateMergeValues = append(stateMergeValues, s...)

		c.Close()
	}

	for ca, de := range designs {
		g := state.NewStateKeyGenerator(ca)
		stateMergeValues = append(stateMergeValues, cstate.NewStateMergeValue(
			g.Design(),
			state.NewDesignStateValue(de),
		))
	}

	return stateMergeValues, nil, nil
}

func (opp *ApprovesProcessor) Close() error {
	approvesProcessorPool.Put(opp)

	return nil
}
