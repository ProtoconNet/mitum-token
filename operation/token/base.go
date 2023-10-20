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
	base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	sender, currency := fact.Sender(), fact.Currency()

	policy, err := state.ExistsCurrencyPolicy(currency, getStateFunc)
	if err != nil {
		return nil, ErrBaseOperationProcess(err, "currency policy not found, %s", currency.String()), nil
	}

	fee, err := policy.Feeer().Fee(common.ZeroBig)
	if err != nil {
		return nil, ErrBaseOperationProcess(err, "failed to check fee of currency, %s", currency.String()), nil
	}

	st, err := state.ExistsState(currencystate.StateKeyBalance(sender, currency), "key of currency balance", getStateFunc)
	if err != nil {
		return nil, ErrBaseOperationProcess(err, "currency balance not found, %s", utils.JoinStringers(sender, currency)), nil
	}
	sb := state.NewStateMergeValue(st.Key(), st.Value())

	switch b, err := currencystate.StateBalanceValue(st); {
	case err != nil:
		return nil, ErrBaseOperationProcess(err, "failed to get balance value, %s", utils.JoinStringers(sender, currency)), nil
	case b.Big().Compare(fee) < 0:
		return nil, ErrBaseOperationProcess(err, "not enough balance of sender, %s", utils.JoinStringers(sender, currency)), nil
	}

	v, ok := sb.Value().(currencystate.BalanceStateValue)
	if !ok {
		return nil, ErrBaseOperationProcess(nil, "expected %T, not %T", currencystate.BalanceStateValue{}, sb.Value()), nil
	}
	return state.NewStateMergeValue(sb.Key(), currencystate.NewBalanceStateValue(v.Amount.WithBig(v.Amount.Big().Sub(fee)))), nil, nil
}
