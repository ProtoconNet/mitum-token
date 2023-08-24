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
	tokenID, currency currencytypes.CurrencyID,
	receiver base.Address,
	amount common.Big,
) TransferFact {
	fact := TransferFact{
		TokenFact: NewTokenFact(
			base.NewBaseFact(TransferFactHint, token), sender, contract, tokenID, currency,
		),
		receiver: receiver,
		amount:   amount,
	}
	fact.SetHash(fact.GenerateHash())
	return fact
}

func (fact TransferFact) IsValid([]byte) error {
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

	return nil
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

type Transfer struct {
	TokenOperation
}

func NewTransfer(fact TransferFact) Transfer {
	return Transfer{TokenOperation: NewTokenOperation(TransferHint, fact)}
}
