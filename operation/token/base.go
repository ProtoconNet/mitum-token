package token

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/state"
	currencystate "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-token/utils"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
	"github.com/pkg/errors"
)

type TokenFact struct {
	base.BaseFact
	sender   base.Address
	contract base.Address
	tokenID  types.CurrencyID
	currency types.CurrencyID
}

func NewTokenFact(
	baseFact base.BaseFact,
	sender, contract base.Address,
	tokenID, currency types.CurrencyID,
) TokenFact {
	return TokenFact{
		baseFact,
		sender,
		contract,
		tokenID,
		currency,
	}
}

func (fact TokenFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact TokenFact) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf(utils.ErrStringInvalid(fact))

	if err := util.CheckIsValiders(nil, false,
		fact.BaseFact,
		fact.sender,
		fact.contract,
		fact.tokenID,
		fact.currency,
	); err != nil {
		return e.Wrap(err)
	}

	if fact.sender.Equal(fact.contract) {
		return e.Wrap(errors.Errorf("sender is same with contract account, %s", fact.sender))
	}

	return nil
}

func (fact TokenFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		fact.contract.Bytes(),
		fact.tokenID.Bytes(),
		fact.currency.Bytes(),
	)
}

func (fact TokenFact) Sender() base.Address {
	return fact.sender
}

func (fact TokenFact) Contract() base.Address {
	return fact.contract
}

func (fact TokenFact) TokenID() types.CurrencyID {
	return fact.tokenID
}

func (fact TokenFact) Currency() types.CurrencyID {
	return fact.currency
}

func (fact TokenFact) Addresses() []base.Address {
	return []base.Address{fact.sender, fact.contract}
}

type TokenOperation struct {
	common.BaseOperation
}

func NewTokenOperation(ht hint.Hint, fact base.Fact) TokenOperation {
	return TokenOperation{BaseOperation: common.NewBaseOperation(ht, fact)}
}

func (op *TokenOperation) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	return op.Sign(priv, networkID)
}

func calculateCurrencyFee(fact TokenFact, getStateFunc base.GetStateFunc) (
	base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	sender, currency := fact.Sender(), fact.Currency()

	policy, err := state.ExistsCurrencyPolicy(currency, getStateFunc)
	if err != nil {
		return nil, ErrStateNotFound("currency policy", currency.String(), err), nil
	}

	fee, err := policy.Feeer().Fee(common.ZeroBig)
	if err != nil {
		return nil, ErrBaseOperationProcess("failed to check fee of currency", currency.String(), err), nil
	}

	st, err := state.ExistsState(currencystate.StateKeyBalance(sender, currency), "key of currency balance", getStateFunc)
	if err != nil {
		return nil, ErrStateNotFound("currency balance", utils.StringerChain(sender, currency), err), nil
	}
	sb := state.NewStateMergeValue(st.Key(), st.Value())

	switch b, err := currencystate.StateBalanceValue(st); {
	case err != nil:
		return nil, ErrBaseOperationProcess("failed to get balance value", utils.StringerChain(sender, currency), err), nil
	case b.Big().Compare(fee) < 0:
		return nil, ErrBaseOperationProcess("not enough balance of sender", utils.StringerChain(sender, currency), err), nil
	}

	v, ok := sb.Value().(currencystate.BalanceStateValue)
	if !ok {
		return nil, ErrBaseOperationProcess(utils.ErrStringTypeCast(currencystate.BalanceStateValue{}, sb.Value()), "", nil), nil
	}
	return state.NewStateMergeValue(sb.Key(), currencystate.NewBalanceStateValue(v.Amount.WithBig(v.Amount.Big().Sub(fee)))), nil, nil
}
