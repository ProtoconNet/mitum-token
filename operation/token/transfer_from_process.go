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
	extstate "github.com/ProtoconNet/mitum-currency/v3/state/extension"
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
			return nil, e.Wrap(errors.Errorf(utils.ErrStringTypeCast(t, nopp)))
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
	e := util.StringError(ErrStringPreProcess(*opp))

	fact, ok := op.Fact().(TransferFromFact)
	if !ok {
		return ctx, nil, e.Wrap(errors.Errorf(utils.ErrStringTypeCast(TransferFromFact{}, op.Fact())))
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e.Wrap(err)
	}

	if err := currencystate.CheckExistsState(currency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, ErrStateNotFound("sender", fact.Sender().String(), err), nil
	}

	if err := currencystate.CheckNotExistsState(extstate.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, ErrBaseOperationProcess(err, "contract account cannot transfer token as approved, %s", fact.Sender().String()), nil
	}

	if err := currencystate.CheckExistsState(extstate.StateKeyContractAccount(fact.Contract()), getStateFunc); err != nil {
		return nil, ErrBaseOperationProcess(err, "contract not found, %s", fact.Contract().String()), nil
	}

	if err := currencystate.CheckExistsState(currency.StateKeyCurrencyDesign(fact.Currency()), getStateFunc); err != nil {
		return nil, ErrStateNotFound("currency", fact.Currency().String(), err), nil
	}

	if err := currencystate.CheckExistsState(currency.StateKeyAccount(fact.Receiver()), getStateFunc); err != nil {
		return nil, ErrStateNotFound("receiver", fact.Receiver().String(), err), nil
	}

	if err := currencystate.CheckNotExistsState(extstate.StateKeyContractAccount(fact.Receiver()), getStateFunc); err != nil {
		return nil, ErrBaseOperationProcess(err, "contract account cannot receive tokens, %s", fact.Receiver().String()), nil
	}

	if err := currencystate.CheckExistsState(currency.StateKeyAccount(fact.Target()), getStateFunc); err != nil {
		return nil, ErrStateNotFound("target", fact.Target().String(), err), nil
	}

	if err := currencystate.CheckNotExistsState(extstate.StateKeyContractAccount(fact.Receiver()), getStateFunc); err != nil {
		return nil, ErrBaseOperationProcess(err, "contract account cannot be target of transfer-from, %s", fact.Receiver().String()), nil
	}

	g := state.NewStateKeyGenerator(fact.Contract())

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

	if idx < 0 {
		return nil, ErrBaseOperationProcess(
			err,
			"sender has not approved %s, %s",
			fact.Contract(), fact.Sender(),
		), nil
	}

	aprInfo := approveBoxList[idx].GetApproveInfo(fact.Sender())
	if aprInfo == nil {
		return nil, ErrBaseOperationProcess(
			err,
			"sender is not approved account of target, %s, %s, %s",
			fact.Contract(), fact.Sender(), fact.Target(),
		), nil
	}

	if aprInfo.Amount().Compare(fact.Amount()) < 0 {
		return nil, ErrBaseOperationProcess(
			err,
			"approved amount is less than amount to transfer, %s < %s, %s, %s, %s",
			aprInfo.Amount(), fact.Amount(), fact.Contract(), fact.Sender(), fact.Target(),
		), nil
	}

	st, err = currencystate.ExistsState(g.TokenBalance(fact.Target()), "key of token balance", getStateFunc)
	if err != nil {
		return nil, ErrStateNotFound("token balance", utils.JoinStringers(fact.Contract(), fact.Target()), err), nil
	}

	tb, err := state.StateTokenBalanceValue(st)
	if err != nil {
		return nil, ErrStateNotFound("token balance value", utils.JoinStringers(fact.Contract(), fact.Target()), err), nil
	}

	if tb.Compare(fact.Amount()) < 0 {
		return nil, ErrBaseOperationProcess(
			err,
			"token balance is less than amount to transfer, %s < %s, %s, %s",
			tb, fact.Amount(), fact.Contract(), fact.Target(),
		), nil
	}

	if err := currencystate.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, ErrBaseOperationProcess(err, "invalid signing"), nil
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

	sts := make([]base.StateMergeValue, 4)

	v, baseErr, err := calculateCurrencyFee(fact.TokenFact, getStateFunc)
	if baseErr != nil || err != nil {
		return nil, baseErr, err
	}
	sts[0] = v

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

	sts[1] = currencystate.NewStateMergeValue(
		g.Design(),
		state.NewDesignStateValue(de),
	)

	st, err = currencystate.ExistsState(g.TokenBalance(fact.Sender()), "key of token balance", getStateFunc)
	if err != nil {
		return nil, ErrStateNotFound("token balance", utils.JoinStringers(fact.Contract(), fact.Target()), err), nil
	}

	sb, err := state.StateTokenBalanceValue(st)
	if err != nil {
		return nil, ErrStateNotFound("token balance value", utils.JoinStringers(fact.Contract(), fact.Target()), err), nil
	}

	sts[2] = currencystate.NewStateMergeValue(
		g.TokenBalance(fact.Target()),
		state.NewTokenBalanceStateValue(sb.Sub(fact.Amount())),
	)

	rb := common.ZeroBig
	switch st, found, err := getStateFunc(g.TokenBalance(fact.Receiver())); {
	case err != nil:
		return nil, ErrBaseOperationProcess(err, "failed to check token balance, %s, %s", fact.Contract(), fact.Receiver()), nil
	case found:
		b, err := state.StateTokenBalanceValue(st)
		if err != nil {
			return nil, ErrBaseOperationProcess(err, "failed to get token balance value from state, %s, %s", fact.Contract(), fact.Receiver()), nil
		}
		rb = b
	}

	sts[3] = currencystate.NewStateMergeValue(
		g.TokenBalance(fact.Receiver()),
		state.NewTokenBalanceStateValue(rb.Add(fact.Amount())),
	)

	return sts, nil, nil
}

func (opp *TransferFromProcessor) Close() error {
	transferFromProcessorPool.Put(opp)
	return nil
}
