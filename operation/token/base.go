package token

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/state"
	currencystate "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-token/utils"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

type TokenFact struct {
	base.BaseFact
	sender   base.Address
	contract base.Address
	currency types.CurrencyID
}

func NewTokenFact(
	baseFact base.BaseFact,
	sender, contract base.Address,
	currency types.CurrencyID,
) TokenFact {
	return TokenFact{
		baseFact,
		sender,
		contract,
		currency,
	}
}

func (fact TokenFact) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf(utils.ErrStringInvalid(fact))

	if err := util.CheckIsValiders(nil, false,
		fact.BaseFact,
		fact.sender,
		fact.contract,
		fact.currency,
	); err != nil {
		return e.Wrap(err)
	}

	if fact.sender.Equal(fact.contract) {
		return e.Wrap(errors.Errorf("contract address is same with sender, %s", fact.sender))
	}

	return nil
}

func (fact TokenFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		fact.contract.Bytes(),
		fact.currency.Bytes(),
	)
}

func (fact TokenFact) Sender() base.Address {
	return fact.sender
}

func (fact TokenFact) Contract() base.Address {
	return fact.contract
}

func (fact TokenFact) Currency() types.CurrencyID {
	return fact.currency
}

func (fact TokenFact) Addresses() []base.Address {
	return []base.Address{fact.sender, fact.contract}
}

//type TokenOperation struct {
//	common.BaseOperation
//}
//
//func NewTokenOperation(ht hint.Hint, fact base.Fact) TokenOperation {
//	return TokenOperation{BaseOperation: common.NewBaseOperation(ht, fact)}
//}
//
//func (op *TokenOperation) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
//	return op.Sign(priv, networkID)
//}

func calculateCurrencyFee(fact TokenFact, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	//sender, currency := fact.Sender(), fact.Currency()
	var sts []base.StateMergeValue
	currencyPolicy, err := state.ExistsCurrencyPolicy(fact.Currency(), getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("currency not found, %q; %w", fact.Currency(), err), nil
	}

	fee, err := currencyPolicy.Feeer().Fee(common.ZeroBig)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			"failed to check fee of currency, %q; %w",
			fact.Currency(),
			err,
		), nil
	}

	senderBalSt, err := state.ExistsState(
		currencystate.StateKeyBalance(fact.Sender(), fact.Currency()),
		"key of sender balance",
		getStateFunc,
	)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			"sender balance not found, %q; %w",
			fact.Sender(),
			err,
		), nil
	}

	switch senderBal, err := currencystate.StateBalanceValue(senderBalSt); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError(
			"failed to get balance value, %q; %w",
			currencystate.StateKeyBalance(fact.Sender(), fact.Currency()),
			err,
		), nil
	case senderBal.Big().Compare(fee) < 0:
		return nil, base.NewBaseOperationProcessReasonError(
			"not enough balance of sender, %q",
			fact.Sender(),
		), nil
	}

	v, ok := senderBalSt.Value().(currencystate.BalanceStateValue)
	if !ok {
		return nil, base.NewBaseOperationProcessReasonError("expected BalanceStateValue, not %T", senderBalSt.Value()), nil
	}

	if currencyPolicy.Feeer().Receiver() != nil {
		if err := state.CheckExistsState(currencystate.StateKeyAccount(currencyPolicy.Feeer().Receiver()), getStateFunc); err != nil {
			return nil, nil, err
		} else if feeRcvrSt, found, err := getStateFunc(currencystate.StateKeyBalance(currencyPolicy.Feeer().Receiver(), fact.currency)); err != nil {
			return nil, nil, err
		} else if !found {
			return nil, nil, errors.Errorf("feeer receiver %s not found", currencyPolicy.Feeer().Receiver())
		} else if feeRcvrSt.Key() != senderBalSt.Key() {
			r, ok := feeRcvrSt.Value().(currencystate.BalanceStateValue)
			if !ok {
				return nil, nil, errors.Errorf("expected %T, not %T", currencystate.BalanceStateValue{}, feeRcvrSt.Value())
			}
			sts = append(sts, common.NewBaseStateMergeValue(
				feeRcvrSt.Key(),
				currencystate.NewAddBalanceStateValue(r.Amount.WithBig(fee)),
				func(height base.Height, st base.State) base.StateValueMerger {
					return currencystate.NewBalanceStateValueMerger(height, feeRcvrSt.Key(), fact.currency, st)
				},
			))

			sts = append(sts, common.NewBaseStateMergeValue(
				senderBalSt.Key(),
				currencystate.NewDeductBalanceStateValue(v.Amount.WithBig(fee)),
				func(height base.Height, st base.State) base.StateValueMerger {
					return currencystate.NewBalanceStateValueMerger(height, senderBalSt.Key(), fact.currency, st)
				},
			))
		}
	}
	return sts, nil, nil

	//policy, err := state.ExistsCurrencyPolicy(currency, getStateFunc)
	//if err != nil {
	//	return nil, ErrBaseOperationProcess(err, "currency policy not found, %s", currency.String()), nil
	//}
	//
	//fee, err := policy.Feeer().Fee(common.ZeroBig)
	//if err != nil {
	//	return nil, ErrBaseOperationProcess(err, "failed to check fee of currency, %s", currency.String()), nil
	//}
	//
	//st, err := state.ExistsState(currencystate.StateKeyBalance(sender, currency), "key of currency balance", getStateFunc)
	//if err != nil {
	//	return nil, ErrBaseOperationProcess(err, "currency balance not found, %s", utils.JoinStringers(sender, currency)), nil
	//}
	//sb := state.NewStateMergeValue(st.Key(), st.Value())
	//
	//switch b, err := currencystate.StateBalanceValue(st); {
	//case err != nil:
	//	return nil, ErrBaseOperationProcess(err, "failed to get balance value, %s", utils.JoinStringers(sender, currency)), nil
	//case b.Big().Compare(fee) < 0:
	//	return nil, ErrBaseOperationProcess(err, "not enough balance of sender, %s", utils.JoinStringers(sender, currency)), nil
	//}
	//
	//v, ok := sb.Value().(currencystate.BalanceStateValue)
	//if !ok {
	//	return nil, ErrBaseOperationProcess(nil, "expected %T, not %T", currencystate.BalanceStateValue{}, sb.Value()), nil
	//}
	//return state.NewStateMergeValue(sb.Key(), currencystate.NewBalanceStateValue(v.Amount.WithBig(v.Amount.Big().Sub(fee)))), nil, nil
}
