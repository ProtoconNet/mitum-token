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

var burnProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(BurnProcessor)
	},
}

func (Burn) Process(
	_ context.Context, _ base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type BurnProcessor struct {
	*base.BaseOperationProcessor
}

func NewBurnProcessor() currencytypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		t := BurnProcessor{}
		e := util.StringError(utils.ErrStringCreate(fmt.Sprintf("new %T", t)))

		nopp := burnProcessorPool.Get()
		opp, ok := nopp.(*BurnProcessor)
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

func (opp *BurnProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringError(ErrStringPreProcess(*opp))

	fact, ok := op.Fact().(BurnFact)
	if !ok {
		return ctx, nil, e.Wrap(errors.Errorf(utils.ErrStringTypeCast(BurnFact{}, op.Fact())))
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e.Wrap(err)
	}

	if err := currencystate.CheckExistsState(currency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, ErrStateNotFound("sender", fact.Sender().String(), err), nil
	}

	if err := currencystate.CheckNotExistsState(extstate.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, ErrBaseOperationProcess("contract account cannot burn token", fact.Sender().String(), err), nil
	}

	st, err := currencystate.ExistsState(extstate.StateKeyContractAccount(fact.Contract()), "key of contract account", getStateFunc)
	if err != nil {
		return nil, ErrStateNotFound("contract", fact.Contract().String(), err), nil
	}

	if !fact.Sender().Equal(fact.Target()) {
		ca, err := extstate.StateContractAccountValue(st)
		if err != nil {
			return nil, ErrStateNotFound("contract value", fact.Contract().String(), err), nil
		}

		if !ca.Owner().Equal(fact.Sender()) {
			return nil, ErrBaseOperationProcess("not contract account owner, neither token owner", fact.Sender().String(), nil), nil
		}
	}

	if err := currencystate.CheckExistsState(currency.StateKeyCurrencyDesign(fact.Currency()), getStateFunc); err != nil {
		return nil, ErrStateNotFound("currency", fact.Currency().String(), err), nil
	}

	if err := currencystate.CheckExistsState(currency.StateKeyAccount(fact.Target()), getStateFunc); err != nil {
		return nil, ErrStateNotFound("target", fact.Target().String(), err), nil
	}

	if err := currencystate.CheckNotExistsState(extstate.StateKeyContractAccount(fact.Target()), getStateFunc); err != nil {
		return nil, ErrBaseOperationProcess("burning tokens of contract accounts is impossible", fact.Target().String(), err), nil
	}

	g := state.NewStateKeyGenerator(fact.Contract())

	if err := currencystate.CheckExistsState(g.Design(), getStateFunc); err != nil {
		return nil, ErrStateNotFound("token design", fact.Contract().String(), err), nil
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
			fmt.Sprintf("token balance is less than amount to burn, %s < %s", tb, fact.Amount()),
			utils.JoinStringers(fact.Contract(), fact.Target()),
			err,
		), nil
	}

	if err := currencystate.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, ErrBaseOperationProcess("invalid signing", "", err), nil
	}

	return ctx, nil, nil
}

func (opp *BurnProcessor) Process(
	_ context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringError(ErrStringProcess(*opp))

	fact, ok := op.Fact().(BurnFact)
	if !ok {
		return nil, nil, e.Wrap(errors.Errorf(utils.ErrStringTypeCast(BurnFact{}, op.Fact())))
	}

	g := state.NewStateKeyGenerator(fact.Contract())

	sts := make([]base.StateMergeValue, 3)

	v, baseErr, err := calculateCurrencyFee(fact.TokenFact, getStateFunc)
	if baseErr != nil || err != nil {
		return nil, baseErr, err
	}
	sts[0] = v

	st, err := currencystate.ExistsState(g.Design(), "key of design", getStateFunc)
	if err != nil {
		return nil, ErrStateNotFound("token design", utils.JoinStringers(fact.Contract()), err), nil
	}

	design, err := state.StateDesignValue(st)
	if err != nil {
		return nil, ErrStateNotFound("token design value", utils.JoinStringers(fact.Contract()), err), nil
	}

	policy := types.NewPolicy(
		design.Policy().TotalSupply().Sub(fact.Amount()),
		design.Policy().ApproveList(),
	)
	if err := policy.IsValid(nil); err != nil {
		return nil, ErrInvalid(policy, err), nil
	}

	design = types.NewDesign(design.Symbol(), design.Name(), policy)
	if err := design.IsValid(nil); err != nil {
		return nil, ErrInvalid(design, err), nil
	}

	sts[1] = currencystate.NewStateMergeValue(
		g.Design(),
		state.NewDesignStateValue(design),
	)

	st, err = currencystate.ExistsState(g.TokenBalance(fact.Target()), "key of token balance", getStateFunc)
	if err != nil {
		return nil, ErrStateNotFound("token balance", utils.JoinStringers(fact.Contract(), fact.Target()), err), nil
	}

	tb, err := state.StateTokenBalanceValue(st)
	if err != nil {
		return nil, ErrStateNotFound("token balance value", utils.JoinStringers(fact.Contract(), fact.Target()), err), nil
	}

	sts[2] = currencystate.NewStateMergeValue(
		g.TokenBalance(fact.Target()),
		state.NewTokenBalanceStateValue(tb.Sub(fact.Amount())),
	)

	return sts, nil, nil
}

func (opp *BurnProcessor) Close() error {
	burnProcessorPool.Put(opp)
	return nil
}
