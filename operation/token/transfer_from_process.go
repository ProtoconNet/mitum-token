package token

import (
	"context"
	"fmt"
	"sync"

	"github.com/ProtoconNet/mitum-token/types"
	"github.com/ProtoconNet/mitum-token/utils"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencystate "github.com/ProtoconNet/mitum-currency/v3/state"
	"github.com/ProtoconNet/mitum-currency/v3/state/currency"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-token/state"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var transferFromProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(TransferFromProcessor)
	},
}

func (TransferFrom) Process(
	_ context.Context, _ base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type TransferFromProcessor struct {
	*base.BaseOperationProcessor
}

func NewTransferFromProcessor() currencytypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		t := TransferFromProcessor{}
		e := util.StringError(utils.ErrStringCreate(fmt.Sprintf("new %T", t)))

		nopp := transferFromProcessorPool.Get()
		opp, ok := nopp.(*TransferFromProcessor)
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

func (opp *TransferFromProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	fact, ok := op.Fact().(TransferFromFact)
	if !ok {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMTypeMismatch).
				Errorf("expected %T, not %T", TransferFromFact{}, op.Fact())), nil
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", err)), nil
	}

	_, err := currencystate.ExistsCurrencyPolicy(fact.Currency(), getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMCurrencyNF).Errorf("%v: %v", fact.Currency(), err)), nil
	}

	if _, _, aErr, cErr := currencystate.ExistsCAccount(fact.Sender(), "sender", true, false, getStateFunc); aErr != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", aErr)), nil
	} else if cErr != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMCAccountNA).
				Errorf("%v: sender account is contract account, %v", fact.Sender(), cErr)), nil
	}

	_, _, aErr, cErr := currencystate.ExistsCAccount(fact.Contract(), "contract", true, true, getStateFunc)
	if aErr != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", aErr)), nil
	} else if cErr != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", cErr)), nil
	}

	if _, _, aErr, cErr := currencystate.ExistsCAccount(fact.Receiver(), "receiver", true, false, getStateFunc); aErr != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", aErr)), nil
	} else if cErr != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMCAccountNA).
				Errorf("%v: receiver account is contract account, %v", fact.Receiver(), cErr)), nil
	}

	if _, _, aErr, cErr := currencystate.ExistsCAccount(fact.Target(), "target", true, false, getStateFunc); aErr != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", aErr)), nil
	} else if cErr != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMCAccountNA).
				Errorf("%v: target account is contract account, %v", fact.Target(), cErr)), nil
	}

	g := state.NewStateKeyGenerator(fact.Contract())

	st, err := currencystate.ExistsState(g.Design(), "design", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMServiceNF).Errorf("token design, %v",
				fact.Contract(),
			)), nil
	}

	design, err := state.StateDesignValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMStateValInvalid).
				Errorf("token design, %s, %s", fact.Contract(), fact.Target())), nil
	}

	approveBoxList := design.Policy().ApproveList()

	idx := -1
	for i, apb := range approveBoxList {
		if apb.Account().Equal(fact.Target()) {
			idx = i
			break
		}
	}

	if idx < 0 {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("sender has not approved %s, %s",
					fact.Contract(), fact.Sender())), nil
	}

	aprInfo := approveBoxList[idx].GetApproveInfo(fact.Sender())
	if aprInfo == nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("sender is not approved account of target, %v, %v, %v: %v",
					fact.Contract(), fact.Sender(), fact.Target(), err)), nil
	}

	if aprInfo.Amount().Compare(fact.Amount()) < 0 {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("approved amount is less than amount to transfer, %v < %v, %v, %v, %v:%v",
					aprInfo.Amount(), fact.Amount(), fact.Contract(), fact.Sender(), fact.Target(), err)), nil
	}

	st, err = currencystate.ExistsState(g.TokenBalance(fact.Target()), "token balance", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMStateNF).
				Errorf("token balance, %s, %s", fact.Contract(), fact.Target())), nil
	}

	tb, err := state.StateTokenBalanceValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMStateValInvalid).
				Errorf("token balance, %s, %s", fact.Contract(), fact.Target())), nil
	}

	if tb.Compare(fact.Amount()) < 0 {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMValueInvalid).
				Errorf("token balance is less than amount to transferfrom, %s < %s, %s, %s",
					tb, fact.Amount(), fact.Contract(), fact.Target())), nil
	}

	if err := currencystate.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMSignInvalid).
				Errorf("%v", err)), nil
	}

	return ctx, nil, nil
}

func (opp *TransferFromProcessor) Process(
	_ context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringError(ErrStringProcess(*opp))

	fact, _ := op.Fact().(TransferFromFact)

	g := state.NewStateKeyGenerator(fact.Contract())

	var sts []base.StateMergeValue

	v, baseErr, err := calculateCurrencyFee(fact.TokenFact, getStateFunc)
	if baseErr != nil || err != nil {
		return nil, baseErr, err
	}
	if len(v) > 0 {
		sts = append(sts, v...)
	}

	st, err := currencystate.ExistsState(g.Design(), "key of design", getStateFunc)
	if err != nil {
		return nil, ErrStateNotFound("token design", fact.Contract().String(), err), nil
	}

	design, err := state.StateDesignValue(st)
	if err != nil {
		return nil, ErrStateNotFound("token design value", fact.Contract().String(), err), nil
	}

	approveBoxList := design.Policy().ApproveList()

	idx := -1
	for i, apb := range approveBoxList {
		if apb.Account().Equal(fact.Target()) {
			idx = i
			break
		}
	}

	apb := approveBoxList[idx]
	am := apb.GetApproveInfo(fact.Sender()).Amount().Sub(fact.Amount())

	if am.IsZero() {
		err := apb.RemoveApproveInfo(fact.Sender())
		if err != nil {
			return nil, nil, e.Wrap(err)
		}
	} else {
		apb.SetApproveInfo(types.NewApproveInfo(fact.Sender(), am))
	}

	approveBoxList[idx] = apb

	policy := types.NewPolicy(design.Policy().TotalSupply(), approveBoxList)
	if err := policy.IsValid(nil); err != nil {
		return nil, ErrInvalid(policy, err), nil
	}

	de := types.NewDesign(design.Symbol(), design.Name(), policy)
	if err := de.IsValid(nil); err != nil {
		return nil, ErrInvalid(de, err), nil
	}

	sts = append(sts, currencystate.NewStateMergeValue(
		g.Design(),
		state.NewDesignStateValue(de),
	))

	st, err = currencystate.ExistsState(g.TokenBalance(fact.Target()), "key of token balance", getStateFunc)
	if err != nil {
		return nil, ErrStateNotFound("token balance", utils.JoinStringers(fact.Contract(), fact.Target()), err), nil
	}

	_, err = state.StateTokenBalanceValue(st)
	if err != nil {
		return nil, ErrStateNotFound("token balance value", utils.JoinStringers(fact.Contract(), fact.Target()), err), nil
	}

	sts = append(sts, common.NewBaseStateMergeValue(
		g.TokenBalance(fact.Target()),
		state.NewDeductTokenBalanceStateValue(fact.Amount()),
		func(height base.Height, st base.State) base.StateValueMerger {
			return state.NewTokenBalanceStateValueMerger(height, g.TokenBalance(fact.Target()), st)
		},
	))

	k := currency.StateKeyAccount(fact.Receiver())
	switch _, found, err := getStateFunc(k); {
	case err != nil:
		return nil, nil, e.Wrap(err)
	case !found:
		nilKys, err := currencytypes.NewNilAccountKeysFromAddress(fact.Receiver())
		if err != nil {
			return nil, ErrBaseOperationProcess(err, "failed to create AccountKeys instance for %v", fact.Receiver().String()), nil
		}
		acc, err := currencytypes.NewAccount(fact.Receiver(), nilKys)
		if err != nil {
			return nil, ErrBaseOperationProcess(err, "failed to create Account instance for %v", fact.Receiver().String()), nil
		}

		stv := currencystate.NewStateMergeValue(k, currency.NewAccountStateValue(acc))
		sts = append(sts, stv)
	}

	switch st, found, err := getStateFunc(g.TokenBalance(fact.Receiver())); {
	case err != nil:
		return nil, ErrBaseOperationProcess(err, "failed to check token balance, %s, %s", fact.Contract(), fact.Receiver()), nil
	case found:
		_, err := state.StateTokenBalanceValue(st)
		if err != nil {
			return nil, ErrBaseOperationProcess(err, "failed to get token balance value from state, %s, %s", fact.Contract(), fact.Receiver()), nil
		}
	}

	sts = append(sts, common.NewBaseStateMergeValue(
		g.TokenBalance(fact.Receiver()),
		state.NewAddTokenBalanceStateValue(fact.Amount()),
		func(height base.Height, st base.State) base.StateValueMerger {
			return state.NewTokenBalanceStateValueMerger(height, g.TokenBalance(fact.Receiver()), st)
		},
	))

	return sts, nil, nil
}

func (opp *TransferFromProcessor) Close() error {
	transferFromProcessorPool.Put(opp)
	return nil
}
