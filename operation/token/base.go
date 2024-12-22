package token

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/types"
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
	if err := util.CheckIsValiders(nil, false,
		fact.BaseFact,
		fact.sender,
		fact.contract,
		fact.currency,
	); err != nil {
		return err
	}

	if fact.sender.Equal(fact.contract) {
		return common.ErrSelfTarget.Wrap(errors.Errorf("contract address is same with sender, %s", fact.sender))
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

func (fact TokenFact) FeeBase() map[types.CurrencyID][]common.Big {
	required := make(map[types.CurrencyID][]common.Big)
	required[fact.Currency()] = []common.Big{common.ZeroBig}

	return required
}

func (fact TokenFact) FeePayer() base.Address {
	return fact.sender
}

func (fact TokenFact) FactUser() base.Address {
	return fact.sender
}

func (fact TokenFact) Signer() base.Address {
	return fact.sender
}
