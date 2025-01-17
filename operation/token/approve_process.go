package token

import (
	"context"
	"fmt"
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"sync"

	"github.com/ProtoconNet/mitum-token/types"
	"github.com/ProtoconNet/mitum-token/utils"

	cstate "github.com/ProtoconNet/mitum-currency/v3/state"
	ctypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-token/state"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var approveProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(ApproveProcessor)
	},
}

func (Approve) Process(
	_ context.Context, _ base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type ApproveProcessor struct {
	*base.BaseOperationProcessor
}

func NewApproveProcessor() ctypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		t := ApproveProcessor{}
		e := util.StringError(utils.ErrStringCreate(fmt.Sprintf("new %T", t)))

		nopp := approveProcessorPool.Get()
		opp, ok := nopp.(*ApproveProcessor)
		if !ok {
			return nil, e.Wrap(errors.Errorf("expected ApproveProcessor, not %T", nopp))
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

func (opp *ApproveProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	fact, ok := op.Fact().(ApproveFact)
	if !ok {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMTypeMismatch).
				Errorf("expected %T, not %T", ApproveFact{}, op.Fact())), nil
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", err)), nil
	}

	if _, _, _, cErr := cstate.ExistsCAccount(
		fact.Approved(), "approved", true, false, getStateFunc); cErr != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMCAccountNA).
				Errorf("%v: approved %v is contract account", cErr, fact.Approved())), nil
	}

	keyGenerator := state.NewStateKeyGenerator(fact.Contract().String())

	if st, err := cstate.ExistsState(
		keyGenerator.Design(), "design", getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMServiceNF).Errorf("token design state for contract account %v",
				fact.Contract(),
			)), nil
	} else if design, err := state.StateDesignValue(st); err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMServiceNF).Errorf("token design state value for contract account %v",
				fact.Contract(),
			)), nil
	} else if apb := design.Policy().GetApproveBox(fact.Sender()); apb == nil {
		if fact.Amount().IsZero() {
			return nil, base.NewBaseOperationProcessReasonError(
				common.ErrMPreProcess.Wrap(common.ErrMValueInvalid).
					Errorf("sender %v has not approved any accounts", fact.Sender())), nil
		}
	} else if aprInfo := apb.GetApproveInfo(fact.Approved()); aprInfo == nil {
		if fact.Amount().IsZero() {
			return nil, base.NewBaseOperationProcessReasonError(
				common.ErrMPreProcess.Wrap(common.ErrMValueInvalid).
					Errorf("approved account %v has not been approved",
						fact.Approved())), nil
		}
	}
	if err := cstate.CheckExistsState(keyGenerator.TokenBalance(fact.Sender().String()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMStateNF).
				Errorf("token balance for sender %v in contract account %v", fact.Sender(), fact.Contract())), nil
	}

	return ctx, nil, nil
}

func (opp *ApproveProcessor) Process(
	_ context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	fact, _ := op.Fact().(ApproveFact)

	keyGenerator := state.NewStateKeyGenerator(fact.Contract().String())

	var sts []base.StateMergeValue

	smv, err := cstate.CreateNotExistAccount(fact.Approved(), getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("%w", err), nil
	} else if smv != nil {
		sts = append(sts, smv)
	}

	st, _ := cstate.ExistsState(keyGenerator.Design(), "design", getStateFunc)
	design, _ := state.StateDesignValue(st)
	apb := design.Policy().GetApproveBox(fact.Sender())
	if apb == nil {
		a := types.NewApproveBox(fact.Sender(), []types.ApproveInfo{types.NewApproveInfo(fact.approved, fact.Amount())})
		apb = &a
	} else {
		if fact.Amount().IsZero() {
			err := apb.RemoveApproveInfo(fact.Approved())
			if err != nil {
				return nil, ErrBaseOperationProcess(err, "remove approved, %s", fact.Approved().String()), nil
			}
		} else {
			apb.SetApproveInfo(types.NewApproveInfo(fact.approved, fact.Amount()))
		}
	}

	policy := design.Policy()
	policy.MergeApproveBox(*apb)
	if err := policy.IsValid(nil); err != nil {
		return nil, ErrInvalid(policy, err), nil
	}
	de := types.NewDesign(design.Symbol(), design.Name(), design.Decimal(), policy)
	if err := de.IsValid(nil); err != nil {
		return nil, ErrInvalid(de, err), nil
	}
	sts = append(sts, cstate.NewStateMergeValue(
		keyGenerator.Design(),
		state.NewDesignStateValue(de),
	))

	return sts, nil, nil
}

func (opp *ApproveProcessor) Close() error {
	approveProcessorPool.Put(opp)
	return nil
}
