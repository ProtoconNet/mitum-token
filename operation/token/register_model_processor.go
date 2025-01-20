package token

import (
	"context"
	"fmt"
	"sync"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	cstate "github.com/ProtoconNet/mitum-currency/v3/state"
	statee "github.com/ProtoconNet/mitum-currency/v3/state/extension"
	ctypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-token/state"
	"github.com/ProtoconNet/mitum-token/types"
	"github.com/ProtoconNet/mitum-token/utils"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var registerModelProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(RegisterModelProcessor)
	},
}

func (RegisterModel) Process(
	_ context.Context, _ base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type RegisterModelProcessor struct {
	*base.BaseOperationProcessor
}

func NewRegisterModelProcessor() ctypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		t := RegisterModelProcessor{}
		e := util.StringError(utils.ErrStringCreate(fmt.Sprintf("new %T", t)))

		nopp := registerModelProcessorPool.Get()
		opp, ok := nopp.(*RegisterModelProcessor)
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

func (opp *RegisterModelProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	fact, ok := op.Fact().(RegisterModelFact)
	if !ok {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMTypeMismatch).
				Errorf("expected %T, not %T", RegisterModelFact{}, op.Fact())), nil
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", err)), nil
	}

	g := state.NewStateKeyGenerator(fact.Contract().String())

	if found, _ := cstate.CheckNotExistsState(g.Design(), getStateFunc); found {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMServiceE).Errorf("token design for contract account %v", fact.Contract())), nil
	}

	return ctx, nil, nil
}

func (opp *RegisterModelProcessor) Process(
	_ context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringError(ErrStringProcess(*opp))

	fact, ok := op.Fact().(RegisterModelFact)
	if !ok {
		return nil, nil, e.Wrap(errors.Errorf(utils.ErrStringTypeCast(RegisterModelFact{}, op.Fact())))
	}

	g := state.NewStateKeyGenerator(fact.Contract().String())

	var sts []base.StateMergeValue

	policy := types.NewPolicy(fact.InitialSupply(), []types.ApproveBox{})
	if err := policy.IsValid(nil); err != nil {
		return nil, ErrInvalid(policy, err), nil
	}

	design := types.NewDesign(fact.Symbol(), fact.Name(), fact.Decimal(), policy)
	if err := design.IsValid(nil); err != nil {
		return nil, ErrInvalid(design, err), nil
	}

	sts = append(sts, cstate.NewStateMergeValue(
		g.Design(),
		state.NewDesignStateValue(design),
	))

	st, err := cstate.ExistsState(statee.StateKeyContractAccount(fact.Contract()), "contract account", getStateFunc)
	if err != nil {
		return nil, ErrStateNotFound("contract", fact.Contract().String(), err), nil
	}

	ca, err := statee.StateContractAccountValue(st)
	if err != nil {
		return nil, ErrStateNotFound("contract value", fact.Contract().String(), err), nil
	}
	nca := ca.SetIsActive(true)

	sts = append(sts, cstate.NewStateMergeValue(
		statee.StateKeyContractAccount(fact.Contract()),
		statee.NewContractAccountStateValue(nca),
	))

	if fact.InitialSupply().OverZero() {
		sts = append(sts, common.NewBaseStateMergeValue(
			g.TokenBalance(fact.Sender().String()),
			state.NewAddTokenBalanceStateValue(fact.InitialSupply()),
			func(height base.Height, st base.State) base.StateValueMerger {
				return state.NewTokenBalanceStateValueMerger(
					height,
					g.TokenBalance(fact.Sender().String()),
					st,
				)
			},
		))
	}

	return sts, nil, nil
}

func (opp *RegisterModelProcessor) Close() error {
	registerModelProcessorPool.Put(opp)
	return nil
}
