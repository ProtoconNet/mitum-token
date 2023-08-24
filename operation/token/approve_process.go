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

var approveProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(ApproveProcessor)
	},
}

func (Approve) Process(
	_ context.Context, _ base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type ApproveProcessor struct {
	*base.BaseOperationProcessor
}

func NewApproveProcessor() currencytypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		t := ApproveProcessor{}
		e := util.StringError(utils.ErrStringCreate(fmt.Sprintf("new %T", t)))

		nopp := approveProcessorPool.Get()
		opp, ok := nopp.(*ApproveProcessor)
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

func (opp *ApproveProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringError(ErrStringPreProcess(*opp))

	fact, ok := op.Fact().(ApproveFact)
	if !ok {
		return ctx, nil, e.Wrap(errors.Errorf(utils.ErrStringTypeCast(ApproveFact{}, op.Fact())))
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e.Wrap(err)
	}

	if err := currencystate.CheckExistsState(currency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, ErrStateNotFound("sender", fact.Sender().String(), err), nil
	}

	if err := currencystate.CheckNotExistsState(extstate.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, ErrBaseOperationProcess("contract account cannot run approve-operation", fact.Sender().String(), err), nil
	}

	if err := currencystate.CheckExistsState(extstate.StateKeyContractAccount(fact.Contract()), getStateFunc); err != nil {
		return nil, ErrStateNotFound("contract", fact.Contract().String(), err), nil
	}

	if err := currencystate.CheckExistsState(currency.StateKeyCurrencyDesign(fact.Currency()), getStateFunc); err != nil {
		return nil, ErrStateNotFound("currency", fact.Currency().String(), err), nil
	}

	if err := currencystate.CheckExistsState(currency.StateKeyAccount(fact.Approved()), getStateFunc); err != nil {
		return nil, ErrStateNotFound("approved", fact.Approved().String(), err), nil
	}

	if err := currencystate.CheckNotExistsState(extstate.StateKeyContractAccount(fact.Approved()), getStateFunc); err != nil {
		return nil, ErrBaseOperationProcess("contract cannot become approved account", fact.Approved().String(), err), nil
	}

	g := state.NewStateKeyGenerator(fact.Contract(), fact.TokenID())

	if err := currencystate.CheckExistsState(g.Design(), getStateFunc); err != nil {
		return nil, ErrStateNotFound("token design", utils.StringerChain(fact.Contract(), fact.TokenID()), err), nil
	}

	if err := currencystate.CheckExistsState(g.TokenBalance(fact.Sender()), getStateFunc); err != nil {
		return nil, ErrStateNotFound("token balance", utils.StringerChain(fact.Contract(), fact.TokenID(), fact.Sender()), err), nil
	}

	if err := currencystate.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, ErrBaseOperationProcess("invalid signing", "", err), nil
	}

	return ctx, nil, nil
}

func (opp *ApproveProcessor) Process(
	_ context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringError(ErrStringProcess(*opp))

	fact, ok := op.Fact().(ApproveFact)
	if !ok {
		return nil, nil, e.Wrap(errors.Errorf(utils.ErrStringTypeCast(ApproveFact{}, op.Fact())))
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

	amount := fact.Amount()

	idx := -1
	for i, ap := range al {
		if ap.Account().Equal(fact.Sender()) {
			idx = i
			break
		}
	}

	if -1 < idx {
		for a, big := range al[idx].Approved() {
			if a == fact.Approved().String() {
				amount = amount.Add(big)
				break
			}
		}

		m := al[idx].Approved()
		m[fact.Approved().String()] = amount

		al[idx] = types.NewApproveInfo(al[idx].Account(), m)
	} else {
		m := map[string]common.Big{}
		m[fact.Approved().String()] = amount
		al = append(al, types.NewApproveInfo(fact.Sender(), m))
	}

	policy := types.NewPolicy(
		design.Policy().TotalSupply(),
		al,
	)
	if err := policy.IsValid(nil); err != nil {
		return nil, ErrInvalid(policy, err), nil
	}

	design = types.NewDesign(design.TokenID(), design.Symbol(), policy)
	if err := design.IsValid(nil); err != nil {
		return nil, ErrInvalid(design, err), nil
	}

	sts := make([]base.StateMergeValue, 2)

	v, baseErr, err := calculateCurrencyFee(fact.TokenFact, getStateFunc)
	if baseErr != nil || err != nil {
		return nil, baseErr, err
	}
	sts[0] = v

	sts[1] = currencystate.NewStateMergeValue(
		g.Design(),
		state.NewDesignStateValue(design),
	)

	return sts, nil, nil
}

func (opp *ApproveProcessor) Close() error {
	approveProcessorPool.Put(opp)
	return nil
}
