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
		return nil, ErrBaseOperationProcess("contract account cannot transfer token as approved", fact.Sender().String(), err), nil
	}

	if err := currencystate.CheckExistsState(extstate.StateKeyContractAccount(fact.Contract()), getStateFunc); err != nil {
		return nil, ErrBaseOperationProcess("contract", fact.Contract().String(), err), nil
	}

	if err := currencystate.CheckExistsState(currency.StateKeyCurrencyDesign(fact.Currency()), getStateFunc); err != nil {
		return nil, ErrStateNotFound("currency", fact.Currency().String(), err), nil
	}

	if err := currencystate.CheckExistsState(currency.StateKeyAccount(fact.Receiver()), getStateFunc); err != nil {
		return nil, ErrStateNotFound("receiver", fact.Receiver().String(), err), nil
	}

	if err := currencystate.CheckNotExistsState(extstate.StateKeyContractAccount(fact.Receiver()), getStateFunc); err != nil {
		return nil, ErrBaseOperationProcess("contract account cannot receive tokens", fact.Receiver().String(), err), nil
	}

	if err := currencystate.CheckExistsState(currency.StateKeyAccount(fact.Target()), getStateFunc); err != nil {
		return nil, ErrStateNotFound("target", fact.Target().String(), err), nil
	}

	if err := currencystate.CheckNotExistsState(extstate.StateKeyContractAccount(fact.Receiver()), getStateFunc); err != nil {
		return nil, ErrBaseOperationProcess("contract account cannot be target of transfer-from", fact.Receiver().String(), err), nil
	}

	g := state.NewStateKeyGenerator(fact.Contract(), fact.TokenID())

	st, err := currencystate.ExistsState(g.Design(), "key of design", getStateFunc)
	if err != nil {
		return nil, ErrStateNotFound("token design", utils.StringerChain(fact.Contract(), fact.TokenID()), err), nil
	}

	design, err := state.StateDesignValue(st)
	if err != nil {
		return nil, ErrStateNotFound("token design value", utils.StringerChain(fact.Contract(), fact.TokenID()), err), nil
	}

	al := design.Policy().ApproveList()

	idx := -1
	for i, ap := range al {
		if ap.Account().Equal(fact.Target()) {
			idx = i
			break
		}
	}

	if idx < 0 {
		return nil, ErrBaseOperationProcess(
			"sender is not approved account of target",
			utils.StringerChain(fact.Contract(), fact.TokenID(), fact.Sender(), fact.Target()),
			err,
		), nil
	}

	big, found := al[idx].Approved()[fact.Sender().String()]
	if !found {
		return nil, ErrBaseOperationProcess(
			"sender is not approved account of target",
			utils.StringerChain(fact.Contract(), fact.TokenID(), fact.Sender(), fact.Target()),
			err,
		), nil
	}

	if big.Compare(fact.Amount()) < 0 {
		return nil, ErrBaseOperationProcess(
			fmt.Sprintf("approved amount is less than amount to transfer, %s < %s", big, fact.Amount()),
			utils.StringerChain(fact.Contract(), fact.TokenID(), fact.Sender(), fact.Target()),
			err,
		), nil
	}

	st, err = currencystate.ExistsState(g.TokenBalance(fact.Target()), "key of token balance", getStateFunc)
	if err != nil {
		return nil, ErrStateNotFound("token balance", utils.StringerChain(fact.Contract(), fact.TokenID(), fact.Target()), err), nil
	}

	tb, err := state.StateTokenBalanceValue(st)
	if err != nil {
		return nil, ErrStateNotFound("token balance value", utils.StringerChain(fact.Contract(), fact.TokenID(), fact.Target()), err), nil
	}

	if tb.Compare(fact.Amount()) < 0 {
		return nil, ErrBaseOperationProcess(
			fmt.Sprintf("token balance is less than amount to transfer, %s < %s", tb, fact.Amount()),
			utils.StringerChain(fact.Contract(), fact.TokenID(), fact.Target()),
			err,
		), nil
	}

	if err := currencystate.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, ErrBaseOperationProcess("invalid signing", "", err), nil
	}

	return ctx, nil, nil
}

func (opp *TransferFromProcessor) Process(
	_ context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringError(ErrStringProcess(*opp))

	fact, ok := op.Fact().(TransferFromFact)
	if !ok {
		return nil, nil, e.Wrap(errors.Errorf(utils.ErrStringTypeCast(TransferFromFact{}, op.Fact())))
	}

	g := state.NewStateKeyGenerator(fact.Contract(), fact.TokenID())

	sts := make([]base.StateMergeValue, 4)

	v, baseErr, err := calculateCurrencyFee(fact.TokenFact, getStateFunc)
	if baseErr != nil || err != nil {
		return nil, baseErr, err
	}
	sts[0] = v

	st, err := currencystate.ExistsState(g.Design(), "key of design", getStateFunc)
	if err != nil {
		return nil, ErrStateNotFound("token design", utils.StringerChain(fact.Contract(), fact.TokenID()), err), nil
	}

	design, err := state.StateDesignValue(st)
	if err != nil {
		return nil, ErrStateNotFound("token design value", utils.StringerChain(fact.Contract(), fact.TokenID()), err), nil
	}

	al := design.Policy().ApproveList()

	idx := -1
	for i, ap := range al {
		if ap.Account().Equal(fact.Target()) {
			idx = i
			break
		}
	}

	if idx < 0 {
		return nil, ErrBaseOperationProcess(
			"sender is not approved account of target",
			utils.StringerChain(fact.Contract(), fact.TokenID(), fact.Sender(), fact.Target()),
			err,
		), nil
	}

	ap := al[idx].Approved()
	am := ap[fact.Sender().String()].Sub(fact.Amount())

	if am.IsZero() {
		delete(ap, fact.Sender().String())
	} else {
		ap[fact.Sender().String()] = am
	}

	al[idx] = types.NewApproveInfo(al[idx].Account(), ap)

	policy := types.NewPolicy(design.Policy().TotalSupply(), al)
	if err := policy.IsValid(nil); err != nil {
		return nil, ErrInvalid(policy, err), nil
	}

	design = types.NewDesign(design.TokenID(), design.Symbol(), policy)
	if err := design.IsValid(nil); err != nil {
		return nil, ErrInvalid(design, err), nil
	}

	sts[1] = currencystate.NewStateMergeValue(
		g.Design(),
		state.NewDesignStateValue(design),
	)

	st, err = currencystate.ExistsState(g.TokenBalance(fact.Sender()), "key of token balance", getStateFunc)
	if err != nil {
		return nil, ErrStateNotFound("token balance", utils.StringerChain(fact.Contract(), fact.TokenID(), fact.Target()), err), nil
	}

	sb, err := state.StateTokenBalanceValue(st)
	if err != nil {
		return nil, ErrStateNotFound("token balance value", utils.StringerChain(fact.Contract(), fact.TokenID(), fact.Target()), err), nil
	}

	sts[2] = currencystate.NewStateMergeValue(
		g.TokenBalance(fact.Target()),
		state.NewTokenBalanceStateValue(sb.Sub(fact.Amount())),
	)

	rb := common.ZeroBig
	switch st, found, err := getStateFunc(g.TokenBalance(fact.Receiver())); {
	case err != nil:
		return nil, ErrBaseOperationProcess("failed to check token balance", utils.StringerChain(fact.Contract(), fact.TokenID(), fact.Receiver()), err), nil
	case found:
		b, err := state.StateTokenBalanceValue(st)
		if err != nil {
			return nil, ErrBaseOperationProcess("failed to get token balance value from state", utils.StringerChain(fact.Contract(), fact.TokenID(), fact.Receiver()), err), nil
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