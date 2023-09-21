package token

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-token/utils"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var (
	BurnFactHint = hint.MustNewHint("mitum-token-burn-operation-fact-v0.0.1")
	BurnHint     = hint.MustNewHint("mitum-token-burn-operation-v0.0.1")
)

type BurnFact struct {
	TokenFact
	target base.Address
	amount common.Big
}

func NewBurnFact(
	token []byte,
	sender, contract base.Address,
	currency currencytypes.CurrencyID,
	target base.Address,
	amount common.Big,
) BurnFact {
	fact := BurnFact{
		TokenFact: NewTokenFact(
			base.NewBaseFact(BurnFactHint, token), sender, contract, currency,
		),
		target: target,
		amount: amount,
	}
	fact.SetHash(fact.GenerateHash())
	return fact
}

func (fact BurnFact) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf(utils.ErrStringInvalid(fact))

	if err := fact.TokenFact.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	if err := fact.target.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	if fact.contract.Equal(fact.target) {
		return e.Wrap(errors.Errorf("contract address is same with target, %s", fact.target))
	}

	if !fact.amount.OverZero() {
		return e.Wrap(errors.Errorf("zero amount"))
	}

	return nil
}

func (fact BurnFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.TokenFact.Bytes(),
		fact.target.Bytes(),
		fact.amount.Bytes(),
	)
}

func (fact BurnFact) Target() base.Address {
	return fact.target
}

func (fact BurnFact) Amount() common.Big {
	return fact.amount
}

type Burn struct {
	common.BaseOperation
}

func NewBurn(fact BurnFact) Burn {
	return Burn{BaseOperation: common.NewBaseOperation(BurnHint, fact)}
}
