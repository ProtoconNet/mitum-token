package token

import (
	"context"
	"fmt"
	"github.com/ProtoconNet/mitum-currency/v3/common"
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

func (opp *MintProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringError(ErrStringPreProcess(*opp))

	fact, ok := op.Fact().(MintFact)
	if !ok {
		return ctx, nil, e.Wrap(errors.Errorf(utils.ErrStringTypeCast(MintFact{}, op.Fact())))
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, ErrBaseOperationProcess(err, "invalid MintFact"), nil
	}

	if err := currencystate.CheckExistsState(currency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, ErrBaseOperationProcess(err, "sender not found, %s", fact.Sender().String()), nil
	}

	if err := currencystate.CheckNotExistsState(extstate.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, ErrBaseOperationProcess(err, "contract account cannot mint token, %s", fact.Sender().String()), nil
	}

	st, err := currencystate.ExistsState(extstate.StateKeyContractAccount(fact.Contract()), "key of contract account", getStateFunc)
	if err != nil {
		return nil, ErrBaseOperationProcess(err, "contract not found, %s", fact.Contract().String()), nil
	}

	ca, err := extstate.StateContractAccountValue(st)
	if err != nil {
		return nil, ErrBaseOperationProcess(err, "contract value not found, %s", fact.Contract().String()), nil
	}

	if !(ca.Owner().Equal(fact.sender) || ca.IsOperator(fact.Sender())) {
		return nil, ErrBaseOperationProcess(nil, "sender is neither the owner nor the operator of the target contract account, %q", fact.sender), nil
	}

	if err := currencystate.CheckExistsState(currency.StateKeyCurrencyDesign(fact.Currency()), getStateFunc); err != nil {
		return nil, ErrBaseOperationProcess(err, "currency not found, %s", fact.Currency().String()), nil
	}

	if err := currencystate.CheckExistsState(currency.StateKeyAccount(fact.Receiver()), getStateFunc); err != nil {
		return nil, ErrBaseOperationProcess(err, "receiver not found, %s", fact.Receiver().String()), nil
	}

	if err := currencystate.CheckNotExistsState(extstate.StateKeyContractAccount(fact.Receiver()), getStateFunc); err != nil {
		return nil, ErrBaseOperationProcess(err, "contract account cannot receive new tokens, %s", fact.Receiver().String()), nil
	}

	g := state.NewStateKeyGenerator(fact.Contract())

	if err := currencystate.CheckExistsState(g.Design(), getStateFunc); err != nil {
		return nil, ErrBaseOperationProcess(err, "token design not found, %s", fact.Contract().String()), nil
	}

	if err := currencystate.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, ErrBaseOperationProcess(err, "invalid signing"), nil
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
		return nil, ErrBaseOperationProcess(err, "token design not found, %s", fact.Contract().String()), nil
	}

	design, err := state.StateDesignValue(st)
	if err != nil {
		return nil, ErrBaseOperationProcess(err, "token design value not found, %s", fact.Contract().String()), nil
	}

	policy := types.NewPolicy(
		design.Policy().TotalSupply().Add(fact.Amount()),
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

	k := g.TokenBalance(fact.Receiver())
	switch st, found, err := getStateFunc(k); {
	case err != nil:
		return nil, ErrBaseOperationProcess(err, "failed to check token balance, %s, %s", fact.Contract(), fact.Receiver()), nil
	case found:
		_, err := state.StateTokenBalanceValue(st)
		if err != nil {
			return nil, ErrBaseOperationProcess(err, "failed to get token balance value from state, %s, %s", fact.Contract(), fact.Receiver()), nil
		}
	}

	sts = append(sts, common.NewBaseStateMergeValue(
		k,
		state.NewAddTokenBalanceStateValue(fact.Amount()),
		func(height base.Height, st base.State) base.StateValueMerger {
			return state.NewTokenBalanceStateValueMerger(
				height,
				k,
				st,
			)
		},
	))

	return sts, nil, nil
}

func (opp *MintProcessor) Close() error {
	mintProcessorPool.Put(opp)
	return nil
}
