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

var registerTokenProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(RegisterTokenProcessor)
	},
}

func (RegisterToken) Process(
	_ context.Context, _ base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type RegisterTokenProcessor struct {
	*base.BaseOperationProcessor
}

func NewRegisterTokenProcessor() currencytypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		t := RegisterTokenProcessor{}
		e := util.StringError(utils.ErrStringCreate(fmt.Sprintf("new %T", t)))

		nopp := registerTokenProcessorPool.Get()
		opp, ok := nopp.(*RegisterTokenProcessor)
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

func (opp *RegisterTokenProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringError(ErrStringPreProcess(*opp))

	fact, ok := op.Fact().(RegisterTokenFact)
	if !ok {
		return ctx, nil, e.Wrap(errors.Errorf(utils.ErrStringTypeCast(RegisterTokenFact{}, op.Fact())))
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e.Wrap(err)
	}

	if err := currencystate.CheckExistsState(currency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, ErrStateNotFound("sender", fact.Sender().String(), err), nil
	}

	if err := currencystate.CheckNotExistsState(extstate.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, ErrBaseOperationProcess("contract account cannot register token", fact.Sender().String(), err), nil
	}

	st, err := currencystate.ExistsState(extstate.StateKeyContractAccount(fact.Contract()), "key of contract account", getStateFunc)
	if err != nil {
		return nil, ErrStateNotFound("contract", fact.Contract().String(), err), nil
	}

	ca, err := extstate.StateContractAccountValue(st)
	if err != nil {
		return nil, ErrStateNotFound("contract value", fact.Contract().String(), err), nil
	}

	if !ca.Owner().Equal(fact.Sender()) {
		return nil, ErrBaseOperationProcess("not contract account owner", fact.Sender().String(), nil), nil
	}

	if ca.IsActive() {
		return nil, ErrBaseOperationProcess("a design is already registered", fact.Contract().String(), nil), nil
	}

	if err := currencystate.CheckExistsState(currency.StateKeyCurrencyDesign(fact.Currency()), getStateFunc); err != nil {
		return nil, ErrStateNotFound("currency", fact.Currency().String(), err), nil
	}

	g := state.NewStateKeyGenerator(fact.Contract())

	if err := currencystate.CheckNotExistsState(g.Design(), getStateFunc); err != nil {
		return nil, ErrStateAlreadyExists("token design", fact.Contract().String(), err), nil
	}

	if fact.TotalSupply().OverZero() {
		if err := currencystate.CheckNotExistsState(g.TokenBalance(ca.Owner()), getStateFunc); err != nil {
			return nil, ErrStateAlreadyExists("token balance", utils.JoinStringers(fact.Contract(), ca.Owner()), err), nil
		}
	}

	if err := currencystate.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, ErrBaseOperationProcess("invalid signing", "", err), nil
	}

	return ctx, nil, nil
}

func (opp *RegisterTokenProcessor) Process(
	_ context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringError(ErrStringProcess(*opp))

	fact, ok := op.Fact().(RegisterTokenFact)
	if !ok {
		return nil, nil, e.Wrap(errors.Errorf(utils.ErrStringTypeCast(RegisterTokenFact{}, op.Fact())))
	}

	g := state.NewStateKeyGenerator(fact.Contract())

	sts := make([]base.StateMergeValue, 3, 4)

	v, baseErr, err := calculateCurrencyFee(fact.TokenFact, getStateFunc)
	if baseErr != nil || err != nil {
		return nil, baseErr, err
	}
	sts[0] = v

	policy := types.NewPolicy(fact.TotalSupply(), []types.ApproveInfo{})
	if err := policy.IsValid(nil); err != nil {
		return nil, ErrInvalid(policy, err), nil
	}

	design := types.NewDesign(fact.Symbol(), fact.Name(), policy)
	if err := design.IsValid(nil); err != nil {
		return nil, ErrInvalid(design, err), nil
	}

	sts[1] = currencystate.NewStateMergeValue(
		g.Design(),
		state.NewDesignStateValue(design),
	)

	st, err := currencystate.ExistsState(extstate.StateKeyContractAccount(fact.Contract()), "key of contract account", getStateFunc)
	if err != nil {
		return nil, ErrStateNotFound("contract", fact.Contract().String(), err), nil
	}

	ca, err := extstate.StateContractAccountValue(st)
	if err != nil {
		return nil, ErrStateNotFound("contract value", fact.Contract().String(), err), nil
	}
	ca.SetIsActive(true)

	sts[2] = currencystate.NewStateMergeValue(
		extstate.StateKeyContractAccount(fact.Contract()),
		extstate.NewContractAccountStateValue(ca),
	)

	if fact.TotalSupply().OverZero() {
		sts = append(sts, currencystate.NewStateMergeValue(
			g.TokenBalance(fact.Contract()),
			state.NewTokenBalanceStateValue(fact.TotalSupply()),
		))
	}

	return sts, nil, nil
}

func (opp *RegisterTokenProcessor) Close() error {
	registerTokenProcessorPool.Put(opp)
	return nil
}
