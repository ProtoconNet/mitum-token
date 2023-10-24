package token

import (
	"context"
	"fmt"
	"sync"

	"github.com/ProtoconNet/mitum-token/types"
	"github.com/ProtoconNet/mitum-token/utils"

	currencystate "github.com/ProtoconNet/mitum-currency/v3/state"
	"github.com/ProtoconNet/mitum-currency/v3/state/currency"
	extstate "github.com/ProtoconNet/mitum-currency/v3/state/extension"
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
			return nil, e.Wrap(errors.Errorf(utils.ErrStringTypeCast(&t, nopp)))
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
	e := util.StringError(ErrStringPreProcess(*opp))

	fact, ok := op.Fact().(ApproveFact)
	if !ok {
		return ctx, nil, e.Wrap(errors.Errorf(utils.ErrStringTypeCast(ApproveFact{}, op.Fact())))
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e.Wrap(err)
	}

	if err := currencystate.CheckExistsState(currency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, ErrBaseOperationProcess(err, "sender account not found, %s", fact.Sender().String()), nil
	}

	if err := currencystate.CheckNotExistsState(extstate.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, ErrBaseOperationProcess(err, "contract account cannot run approve-operation, %s", fact.Sender().String()), nil
	}

	if err := currencystate.CheckExistsState(extstate.StateKeyContractAccount(fact.Contract()), getStateFunc); err != nil {
		return nil, ErrBaseOperationProcess(err, "contract account not found, %s", fact.Contract().String()), nil
	}

	if err := currencystate.CheckExistsState(currency.StateKeyCurrencyDesign(fact.Currency()), getStateFunc); err != nil {
		return nil, ErrBaseOperationProcess(err, "currency not found, %s", fact.Currency().String()), nil
	}

	if err := currencystate.CheckExistsState(currency.StateKeyAccount(fact.Approved()), getStateFunc); err != nil {
		return nil, ErrBaseOperationProcess(err, "approved account not found, %s", fact.Approved().String()), nil
	}

	if err := currencystate.CheckNotExistsState(extstate.StateKeyContractAccount(fact.Approved()), getStateFunc); err != nil {
		return nil, ErrBaseOperationProcess(err, "contract account cannot become approved account, %s", fact.Approved().String()), nil
	}

	keyGenerator := state.NewStateKeyGenerator(fact.Contract())

	if st, err := currencystate.ExistsState(keyGenerator.Design(), "key of design", getStateFunc); err != nil {
		return nil, ErrBaseOperationProcess(err, "token design not found, %s", fact.Contract().String()), nil
	} else if design, err := state.StateDesignValue(st); err != nil {
		return nil, ErrBaseOperationProcess(err, "token design value not found, %s", fact.Contract().String()), nil
	} else if apb := design.Policy().GetApproveBox(fact.Sender()); apb == nil {
		if fact.Amount().IsZero() {
			return nil, ErrBaseOperationProcess(err, "sender account has approved no accounts, %s", fact.Sender().String()), nil
		}
	} else if aprInfo := apb.GetApproveInfo(fact.Approved()); aprInfo == nil {
		if fact.Amount().IsZero() {
			return nil, ErrBaseOperationProcess(err, "approved account has not been approved, %s", fact.Approved().String()), nil
		}
	}
	if err := currencystate.CheckExistsState(keyGenerator.TokenBalance(fact.Sender()), getStateFunc); err != nil {
		return nil, ErrBaseOperationProcess(err, "token balance not found, %s", utils.JoinStringers(fact.Contract(), fact.Sender())), nil
	}

	if err := currencystate.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, ErrBaseOperationProcess(err, "invalid signing"), nil
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
	sts = append(sts, v)

	st, _ := currencystate.ExistsState(keyGenerator.Design(), "key of design", getStateFunc)
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
	if len(sts) != 2 {
		return nil, ErrBaseOperationProcess(nil, "insufficient state generated"), nil
	}

	return sts, nil, nil
}

func (opp *ApproveProcessor) Close() error {
	approveProcessorPool.Put(opp)
	return nil
}
