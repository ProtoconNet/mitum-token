package token

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-token/utils"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
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
	currency currencytypes.CurrencyID,
	approved base.Address,
	amount common.Big,
) ApproveFact {
	fact := ApproveFact{
		TokenFact: NewTokenFact(
			base.NewBaseFact(ApproveFactHint, token), sender, contract, currency,
		),
		approved: approved,
		amount:   amount,
	}
	fact.SetHash(fact.GenerateHash())
	return fact
}

func (fact ApproveFact) IsValid(b []byte) error {
	e := util.ErrInvalid.Errorf(utils.ErrStringInvalid(fact))

	if err := fact.TokenFact.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	if err := fact.approved.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	if fact.sender.Equal(fact.approved) {
		return e.Wrap(errors.Errorf("sender address is same with approved, %s", fact.approved))
	}

	if fact.contract.Equal(fact.approved) {
		return e.Wrap(errors.Errorf("contract address is same with approved, %s", fact.approved))
	}

	if !fact.amount.OverNil() {
		return e.Wrap(errors.Errorf("under zero"))
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return err
	}
	return nil
}

func (fact ApproveFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
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

func (fact ApproveFact) Addresses() ([]base.Address, error) {
	var as []base.Address

	as = append(as, fact.TokenFact.Sender())
	as = append(as, fact.TokenFact.Contract())
	as = append(as, fact.approved)

	return as, nil
}

type Approve struct {
	common.BaseOperation
}

func NewApprove(fact ApproveFact) Approve {
	return Approve{BaseOperation: common.NewBaseOperation(ApproveHint, fact)}
}
