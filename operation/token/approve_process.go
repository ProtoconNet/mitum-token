package token

import (
	"context"
	"fmt"
	"github.com/ProtoconNet/mitum-currency/v3/common"
	statecurrency "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	"sync"

	"github.com/ProtoconNet/mitum-token/types"
	"github.com/ProtoconNet/mitum-token/utils"

	currencystate "github.com/ProtoconNet/mitum-currency/v3/state"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
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

func NewApproveProcessor() currencytypes.GetNewProcessor {
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

	if err := currencystate.CheckExistsState(
		statecurrency.DesignStateKey(fact.Currency()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMCurrencyNF).Errorf("currency id, %v", fact.Currency())), nil
	}

	if _, _, aErr, cErr := currencystate.ExistsCAccount(
		fact.Sender(), "sender", true, false, getStateFunc); aErr != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", aErr)), nil
	} else if cErr != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMCAccountNA).
				Errorf("%v", cErr)), nil
	}

	_, _, aErr, cErr := currencystate.ExistsCAccount(
		fact.Contract(), "contract", true, true, getStateFunc)
	if aErr != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", aErr)), nil
	} else if cErr != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", cErr)), nil
	}

	if _, _, aErr, cErr := currencystate.ExistsCAccount(
		fact.Approved(), "approved", true, false, getStateFunc); aErr != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", aErr)), nil
	} else if cErr != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMCAccountNA).
				Errorf("%v", cErr)), nil
	}

	keyGenerator := state.NewStateKeyGenerator(fact.Contract())

	if st, err := currencystate.ExistsState(
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
	if err := currencystate.CheckExistsState(keyGenerator.TokenBalance(fact.Sender()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMStateNF).
				Errorf("token balance for sender %v in contract account %v", fact.Sender(), fact.Contract())), nil
	}

	if err := currencystate.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMSignInvalid).
				Errorf("%v", err)), nil
	}

	return ctx, nil, nil
}

func (opp *ApproveProcessor) Process(
	_ context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringError(ErrStringProcess(*opp))

	fact, ok := op.Fact().(ApproveFact)
	if !ok {
		return nil, nil, e.Wrap(errors.Errorf(utils.ErrStringTypeCast(ApproveFact{}, op.Fact())))
	}

	keyGenerator := state.NewStateKeyGenerator(fact.Contract())

	var sts []base.StateMergeValue

	v, baseErr, err := calculateCurrencyFee(fact.TokenFact, getStateFunc)
	if baseErr != nil || err != nil {
		return nil, baseErr, err
	}

	if len(v) > 0 {
		sts = append(sts, v...)
	}

	st, _ := currencystate.ExistsState(keyGenerator.Design(), "design", getStateFunc)
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
	de := types.NewDesign(design.Symbol(), design.Name(), policy)
	if err := de.IsValid(nil); err != nil {
		return nil, ErrInvalid(de, err), nil
	}
	sts = append(sts, currencystate.NewStateMergeValue(
		keyGenerator.Design(),
		state.NewDesignStateValue(de),
	))

	return sts, nil, nil
}

func (opp *ApproveProcessor) Close() error {
	approveProcessorPool.Put(opp)
	return nil
}
