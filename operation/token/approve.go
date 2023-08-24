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
	ApproveFactHint = hint.MustNewHint("mitum-token-approve-operation-fact-v0.0.1")
	ApproveHint     = hint.MustNewHint("mitum-token-approve-operation-v0.0.1")
)

type ApproveFact struct {
	TokenFact
	approved base.Address
	amount   common.Big
}

func NewApproveFact(
	token []byte,
	sender, contract base.Address,
	tokenID, currency currencytypes.CurrencyID,
	approved base.Address,
	amount common.Big,
) ApproveFact {
	fact := ApproveFact{
		TokenFact: NewTokenFact(
			base.NewBaseFact(ApproveFactHint, token), sender, contract, tokenID, currency,
		),
		approved: approved,
		amount:   amount,
	}
	fact.SetHash(fact.GenerateHash())
	return fact
}

func (fact ApproveFact) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf(utils.ErrStringInvalid(fact))

	if err := fact.TokenFact.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	if err := fact.approved.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	if fact.contract.Equal(fact.approved) {
		return e.Wrap(errors.Errorf("contract address is same with approved, %s", fact.approved))
	}

	if !fact.amount.OverZero() {
		return e.Wrap(errors.Errorf("zero amount"))
	}

	return nil
}

func (fact ApproveFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.TokenFact.Bytes(),
		fact.approved.Bytes(),
		fact.amount.Bytes(),
	)
}

func (fact ApproveFact) Approved() base.Address {
	return fact.approved
}

func (fact ApproveFact) Amount() common.Big {
	return fact.amount
}

type Approve struct {
	TokenOperation
}

func NewApprove(fact ApproveFact) Approve {
	return Approve{TokenOperation: NewTokenOperation(ApproveHint, fact)}
}
