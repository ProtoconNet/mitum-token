package token

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/extras"
	ctypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
	"github.com/pkg/errors"
)

var (
	TransferFactHint = hint.MustNewHint("mitum-token-transfer-operation-fact-v0.0.1")
	TransferHint     = hint.MustNewHint("mitum-token-transfer-operation-v0.0.1")
)

type TransferFact struct {
	TokenFact
	receiver base.Address
	amount   common.Big
}

func NewTransferFact(
	token []byte,
	sender, contract base.Address,
	currency ctypes.CurrencyID,
	receiver base.Address,
	amount common.Big,
) TransferFact {
	fact := TransferFact{
		TokenFact: NewTokenFact(
			base.NewBaseFact(TransferFactHint, token), sender, contract, currency,
		),
		receiver: receiver,
		amount:   amount,
	}
	fact.SetHash(fact.GenerateHash())
	return fact
}

func (fact TransferFact) IsValid(b []byte) error {
	if err := fact.TokenFact.IsValid(nil); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	if err := fact.receiver.IsValid(nil); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	if fact.sender.Equal(fact.receiver) {
		return common.ErrFactInvalid.Wrap(common.ErrSelfTarget.Wrap(errors.Errorf("sender %v is same with receiver", fact.sender)))
	}

	if fact.contract.Equal(fact.receiver) {
		return common.ErrFactInvalid.Wrap(common.ErrSelfTarget.Wrap(errors.Errorf("receiver %v is same with contract account", fact.receiver)))
	}

	if !fact.amount.OverZero() {
		return common.ErrFactInvalid.Wrap(common.ErrValOOR.Wrap(errors.Errorf("transfer amount must be over zero, got %v", fact.amount)))
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}
	return nil
}

func (fact TransferFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact TransferFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.TokenFact.Bytes(),
		fact.receiver.Bytes(),
		fact.amount.Bytes(),
	)
}

func (fact TransferFact) Receiver() base.Address {
	return fact.receiver
}

func (fact TransferFact) Amount() common.Big {
	return fact.amount
}

func (fact TransferFact) Addresses() ([]base.Address, error) {
	var as []base.Address

	as = append(as, fact.TokenFact.Sender())
	as = append(as, fact.TokenFact.Contract())
	as = append(as, fact.receiver)

	return as, nil
}

func (fact TransferFact) ActiveContract() []base.Address {
	return []base.Address{fact.contract}
}

type Transfer struct {
	extras.ExtendedOperation
}

func NewTransfer(fact TransferFact) Transfer {
	return Transfer{
		ExtendedOperation: extras.NewExtendedOperation(TransferHint, fact),
	}
}
