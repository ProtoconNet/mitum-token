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

var mintProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(MintProcessor)
	},
}

func (Mint) Process(
	_ context.Context, _ base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type MintProcessor struct {
	*base.BaseOperationProcessor
}

func NewMintProcessor() currencytypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		t := MintProcessor{}
		e := util.StringError(utils.ErrStringCreate(fmt.Sprintf("new %T", t)))

		nopp := mintProcessorPool.Get()
		opp, ok := nopp.(*MintProcessor)
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

func (opp *MintProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringError(ErrStringPreProcess(*opp))

	fact, ok := op.Fact().(MintFact)
	if !ok {
		return ctx, nil, e.Wrap(errors.Errorf(utils.ErrStringTypeCast(MintFact{}, op.Fact())))
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e.Wrap(err)
	}

	if err := currencystate.CheckExistsState(currency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, ErrStateNotFound("sender", fact.Sender().String(), err), nil
	}

	if err := currencystate.CheckNotExistsState(extstate.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, ErrBaseOperationProcess("contract account cannot mint token", fact.Sender().String(), err), nil
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

	if err := currencystate.CheckExistsState(currency.StateKeyCurrencyDesign(fact.Currency()), getStateFunc); err != nil {
		return nil, ErrStateNotFound("currency", fact.Currency().String(), err), nil
	}

	if err := currencystate.CheckExistsState(currency.StateKeyAccount(fact.Receiver()), getStateFunc); err != nil {
		return nil, ErrStateNotFound("receiver", fact.Receiver().String(), err), nil
	}

	if err := currencystate.CheckNotExistsState(extstate.StateKeyContractAccount(fact.Receiver()), getStateFunc); err != nil {
		return nil, ErrBaseOperationProcess("contract account cannot receive new tokens", fact.Receiver().String(), err), nil
	}

	g := state.NewStateKeyGenerator(fact.Contract(), fact.TokenID())

	if err := currencystate.CheckExistsState(g.Design(), getStateFunc); err != nil {
		return nil, ErrStateNotFound("token design", utils.StringerChain(fact.Contract(), fact.TokenID()), err), nil
	}

	if err := currencystate.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, ErrBaseOperationProcess("invalid signing", "", err), nil
	}

	return ctx, nil, nil
}

func (opp *MintProcessor) Process(
	_ context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringError(ErrStringProcess(*opp))

	fact, ok := op.Fact().(MintFact)
	if !ok {
		return nil, nil, e.Wrap(errors.Errorf(utils.ErrStringTypeCast(MintFact{}, op.Fact())))
	}

	g := state.NewStateKeyGenerator(fact.Contract(), fact.TokenID())

	sts := make([]base.StateMergeValue, 3)

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

	policy := types.NewPolicy(
		design.Policy().TotalSupply().Add(fact.Amount()),
		design.Policy().ApproveList(),
	)
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

	sts[2] = currencystate.NewStateMergeValue(
		g.TokenBalance(fact.Receiver()),
		state.NewTokenBalanceStateValue(fact.Amount()),
	)

	return sts, nil, nil
}

func (opp *MintProcessor) Close() error {
	mintProcessorPool.Put(opp)
	return nil
}
