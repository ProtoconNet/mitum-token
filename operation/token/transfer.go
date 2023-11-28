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
	currency currencytypes.CurrencyID,
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
	e := util.ErrInvalid.Errorf(utils.ErrStringInvalid(fact))

	if err := fact.TokenFact.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	if err := fact.receiver.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	if fact.contract.Equal(fact.receiver) {
		return e.Wrap(errors.Errorf("contract address is same with receiver, %s", fact.receiver))
	}

	if !fact.amount.OverZero() {
		return e.Wrap(errors.Errorf("zero amount"))
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return err
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

type Transfer struct {
	common.BaseOperation
}

func NewTransfer(fact TransferFact) Transfer {
	return Transfer{BaseOperation: common.NewBaseOperation(TransferHint, fact)}
}
