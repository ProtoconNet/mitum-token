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

func (opp *BurnProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringError(ErrStringPreProcess(*opp))

	fact, ok := op.Fact().(BurnFact)
	if !ok {
		return ctx, nil, e.Wrap(errors.Errorf(utils.ErrStringTypeCast(BurnFact{}, op.Fact())))
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, ErrBaseOperationProcess(err, "invalid BurnFact"), nil
	}

	if err := currencystate.CheckExistsState(currency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, ErrBaseOperationProcess(err, "sender not found, %s", fact.Sender().String()), nil
	}

	if err := currencystate.CheckNotExistsState(extstate.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, ErrBaseOperationProcess(err, "contract account cannot burn token, %s", fact.Sender().String()), nil
	}

	st, err := currencystate.ExistsState(extstate.StateKeyContractAccount(fact.Contract()), "key of contract account", getStateFunc)
	if err != nil {
		return nil, ErrBaseOperationProcess(err, "contract not found, %s", fact.Contract().String()), nil
	}

	if !fact.Sender().Equal(fact.Target()) {
		return nil, ErrBaseOperationProcess(nil, "not token owner, %s", fact.Sender().String()), nil
	}

	if err := currencystate.CheckExistsState(currency.StateKeyCurrencyDesign(fact.Currency()), getStateFunc); err != nil {
		return nil, ErrBaseOperationProcess(err, "currency not found, %s", fact.Currency().String()), nil
	}

	if err := currencystate.CheckExistsState(currency.StateKeyAccount(fact.Target()), getStateFunc); err != nil {
		return nil, ErrBaseOperationProcess(err, "target not found, %s", fact.Target().String()), nil
	}

	if err := currencystate.CheckNotExistsState(extstate.StateKeyContractAccount(fact.Target()), getStateFunc); err != nil {
		return nil, ErrBaseOperationProcess(err, "burning tokens of contract accounts is impossible, %s", fact.Target().String()), nil
	}

	g := state.NewStateKeyGenerator(fact.Contract())

	if err := currencystate.CheckExistsState(g.Design(), getStateFunc); err != nil {
		return nil, ErrBaseOperationProcess(err, "token design not found, %s", fact.Contract().String()), nil
	}

	st, err = currencystate.ExistsState(g.TokenBalance(fact.Target()), "key of token balance", getStateFunc)
	if err != nil {
		return nil, ErrBaseOperationProcess(err, "token balance not found, %s, %s", fact.Contract(), fact.Target()), nil
	}

	tb, err := state.StateTokenBalanceValue(st)
	if err != nil {
		return nil, ErrBaseOperationProcess(err, "token balance value not found, %s, %s", fact.Contract(), fact.Target()), nil
	}

	if tb.Compare(fact.Amount()) < 0 {
		return nil, ErrBaseOperationProcess(err,
			"token balance is less than amount to burn, %s < %s, %s, %s",
			tb, fact.Amount(), fact.Contract(), fact.Target(),
		), nil
	}

	if err := currencystate.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, ErrBaseOperationProcess(err, "invalid signing", ""), nil
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
		return nil, ErrBaseOperationProcess(err, "token design", utils.JoinStringers(fact.Contract())), nil
	}

	design, err := state.StateDesignValue(st)
	if err != nil {
		return nil, ErrBaseOperationProcess(err, "token design value", utils.JoinStringers(fact.Contract())), nil
	}

	policy := types.NewPolicy(
		design.Policy().TotalSupply().Sub(fact.Amount()),
		design.Policy().ApproveList(),
	)
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
		return nil, ErrBaseOperationProcess(err, "token balance not found, %s, %s", fact.Contract(), fact.Target()), nil
	}

	tb, err := state.StateTokenBalanceValue(st)
	if err != nil {
		return nil, ErrBaseOperationProcess(err, "token balance value not found, %s, %s", fact.Contract(), fact.Target()), nil
	}

	sts = append(sts, currencystate.NewStateMergeValue(
		g.TokenBalance(fact.Target()),
		state.NewTokenBalanceStateValue(tb.Sub(fact.Amount())),
	))

	return sts, nil, nil
}

func (opp *BurnProcessor) Close() error {
	burnProcessorPool.Put(opp)
	return nil
}
