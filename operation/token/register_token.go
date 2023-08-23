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
	RegisterTokenFactHint = hint.MustNewHint("mitum-token-register-token-operation-fact-v0.0.1")
	RegisterTokenHint     = hint.MustNewHint("mitum-token-register-token-operation-v0.0.1")
)

type RegisterTokenFact struct {
	TokenFact
	symbol      string
	totalSupply common.Big
}

func NewRegisterTokenFact(
	token []byte,
	sender, contract base.Address,
	tokenID, currency currencytypes.CurrencyID,
	symbol string,
	totalSupply common.Big,
) RegisterTokenFact {
	fact := RegisterTokenFact{
		TokenFact: NewTokenFact(
			base.NewBaseFact(RegisterTokenFactHint, token), sender, contract, tokenID, currency,
		),
		symbol:      symbol,
		totalSupply: totalSupply,
	}
	fact.SetHash(fact.GenerateHash())
	return fact
}

func (fact RegisterTokenFact) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf(utils.ErrStringInvalid(fact))

	if err := fact.TokenFact.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	if fact.symbol == "" {
		return e.Wrap(errors.Errorf("empty symbol"))
	}

	if !fact.totalSupply.OverNil() {
		return e.Wrap(errors.Errorf("nil big"))
	}

	return nil
}

func (fact RegisterTokenFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.TokenFact.Bytes(),
		[]byte(fact.symbol),
		fact.totalSupply.Bytes(),
	)
}

func (fact RegisterTokenFact) Symbol() string {
	return fact.symbol
}

func (fact RegisterTokenFact) TotalSupply() common.Big {
	return fact.totalSupply
}

type RegisterToken struct {
	TokenOperation
}

func NewRegisterToken(fact RegisterTokenFact) RegisterToken {
	return RegisterToken{TokenOperation: NewTokenOperation(RegisterTokenHint, fact)}
}
