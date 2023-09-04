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
	TransferFromFactHint = hint.MustNewHint("mitum-token-transfer-from-operation-fact-v0.0.1")
	TransferFromHint     = hint.MustNewHint("mitum-token-transfer-from-operation-v0.0.1")
)

type TransferFromFact struct {
	TokenFact
	receiver base.Address
	target   base.Address
	amount   common.Big
}

func NewTransferFromFact(
	token []byte,
	sender, contract base.Address,
	tokenID, currency currencytypes.CurrencyID,
	receiver, target base.Address,
	amount common.Big,
) TransferFromFact {
	fact := TransferFromFact{
		TokenFact: NewTokenFact(
			base.NewBaseFact(TransferFromFactHint, token), sender, contract, tokenID, currency,
		),
		receiver: receiver,
		target:   target,
		amount:   amount,
	}
	fact.SetHash(fact.GenerateHash())
	return fact
}

func (fact TransferFromFact) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf(utils.ErrStringInvalid(fact))

	if err := fact.TokenFact.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	if err := fact.receiver.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	if err := fact.target.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	if fact.contract.Equal(fact.receiver) {
		return e.Wrap(errors.Errorf("contract address is same with receiver, %s", fact.receiver))
	}

	if fact.contract.Equal(fact.target) {
		return e.Wrap(errors.Errorf("contract address is same with target, %s", fact.target))
	}

	if fact.receiver.Equal(fact.target) {
		return e.Wrap(errors.Errorf("target is same with receiver, %s", fact.receiver))
	}

	if fact.sender.Equal(fact.target) {
		return e.Wrap(errors.Errorf("sender is same with target, %s", fact.target))
	}

	if !fact.amount.OverZero() {
		return e.Wrap(errors.Errorf("zero amount"))
	}

	return nil
}

func (fact TransferFromFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.TokenFact.Bytes(),
		fact.receiver.Bytes(),
		fact.target.Bytes(),
		fact.amount.Bytes(),
	)
}

func (fact TransferFromFact) Receiver() base.Address {
	return fact.receiver
}

func (fact TransferFromFact) Target() base.Address {
	return fact.target
}

func (fact TransferFromFact) Amount() common.Big {
	return fact.amount
}

type TransferFrom struct {
	TokenOperation
}

func NewTransferFrom(fact TransferFromFact) TransferFrom {
	return TransferFrom{TokenOperation: NewTokenOperation(TransferFromHint, fact)}
}