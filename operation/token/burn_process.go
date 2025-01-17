package token

import (
	"context"
	"fmt"
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"sync"

	"github.com/ProtoconNet/mitum-token/types"
	"github.com/ProtoconNet/mitum-token/utils"

	cstate "github.com/ProtoconNet/mitum-currency/v3/state"
	ctypes "github.com/ProtoconNet/mitum-currency/v3/types"
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

func NewBurnProcessor() ctypes.GetNewProcessor {
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
	fact, ok := op.Fact().(BurnFact)
	if !ok {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMTypeMismatch).
				Errorf("expected %T, not %T", BurnFact{}, op.Fact())), nil
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", err)), nil
	}

	if !fact.Sender().Equal(fact.Target()) {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMValueInvalid).
				Errorf("target %v is not token owner in contract account %v", fact.Target(), fact.Contract())), nil
	}

	g := state.NewStateKeyGenerator(fact.Contract().String())

	if err := cstate.CheckExistsState(g.Design(), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMServiceNF).
				Errorf("token design for contract account %v", fact.Contract())), nil
	}

	st, err := cstate.ExistsState(g.TokenBalance(fact.Target().String()), "token balance", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMStateNF).
				Errorf("token balance state of target %v in contract account %v", fact.Target(), fact.Contract())), nil
	}

	tb, err := state.StateTokenBalanceValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMStateValInvalid).
				Errorf("token balance state value of target %v in contract account %v", fact.Target(), fact.Contract())), nil
	}

	if tb.Compare(fact.Amount()) < 0 {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMValueInvalid).
				Errorf("token balance of target %v is less than amount to burn in contract account %v, %v < %v",
					fact.Target(), fact.Contract(), tb, fact.Amount())), nil
	}

	return ctx, nil, nil
}

func (opp *BurnProcessor) Process(
	_ context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	fact, _ := op.Fact().(BurnFact)

	g := state.NewStateKeyGenerator(fact.Contract().String())

	var sts []base.StateMergeValue

	st, _ := cstate.ExistsState(g.Design(), "design", getStateFunc)
	design, _ := state.StateDesignValue(st)

	policy := types.NewPolicy(
		design.Policy().TotalSupply().Sub(fact.Amount()),
		design.Policy().ApproveList(),
	)
	if err := policy.IsValid(nil); err != nil {
		return nil, ErrInvalid(policy, err), nil
	}

	de := types.NewDesign(design.Symbol(), design.Name(), design.Decimal(), policy)
	if err := de.IsValid(nil); err != nil {
		return nil, ErrInvalid(de, err), nil
	}

	sts = append(sts, cstate.NewStateMergeValue(
		g.Design(),
		state.NewDesignStateValue(de),
	))

	st, err := cstate.ExistsState(g.TokenBalance(fact.Target().String()), "token balance", getStateFunc)
	if err != nil {
		return nil, ErrBaseOperationProcess(err, "token balance not found, %s, %s", fact.Contract(), fact.Target()), nil
	}

	_, err = state.StateTokenBalanceValue(st)
	if err != nil {
		return nil, ErrBaseOperationProcess(err, "token balance value not found, %s, %s", fact.Contract(), fact.Target()), nil
	}

	sts = append(sts, common.NewBaseStateMergeValue(
		g.TokenBalance(fact.Target().String()),
		state.NewDeductTokenBalanceStateValue(fact.Amount()),
		func(height base.Height, st base.State) base.StateValueMerger {
			return state.NewTokenBalanceStateValueMerger(height, g.TokenBalance(fact.Target().String()), st)
		},
	))

	return sts, nil, nil
}

func (opp *BurnProcessor) Close() error {
	burnProcessorPool.Put(opp)
	return nil
}
