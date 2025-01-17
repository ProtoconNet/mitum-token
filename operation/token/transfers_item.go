package token

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var TransfersItemHint = hint.MustNewHint("mitum-token-transfers-item-v0.0.1")

type TransfersItem struct {
	hint.BaseHinter
	contract base.Address
	receiver base.Address
	amount   common.Big
	currency types.CurrencyID
}

func NewTransfersItem(contract base.Address, receiver base.Address, amount common.Big, currency types.CurrencyID) TransfersItem {
	return TransfersItem{
		BaseHinter: hint.NewBaseHinter(TransfersItemHint),
		contract:   contract,
		receiver:   receiver,
		amount:     amount,
		currency:   currency,
	}
}

func (it TransfersItem) IsValid([]byte) error {
	if err := it.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := util.CheckIsValiders(nil, false, it.contract, it.receiver); err != nil {
		return err
	}

	if it.receiver.Equal(it.contract) {
		return common.ErrSelfTarget.Wrap(errors.Errorf("receiver %v is same with contract account", it.receiver))
	}

	if !it.amount.OverZero() {
		return common.ErrFactInvalid.Wrap(common.ErrValOOR.Wrap(errors.Errorf("transfer amount must be over zero, got %v", it.amount)))
	}

	return util.CheckIsValiders(nil, false,
		it.BaseHinter,
		it.contract,
		it.receiver,
		it.currency,
	)
}

func (it TransfersItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		it.receiver.Bytes(),
		it.amount.Bytes(),
		it.currency.Bytes(),
	)
}

func (it TransfersItem) Contract() base.Address {
	return it.contract
}

func (it TransfersItem) Receiver() base.Address {
	return it.receiver
}

func (it TransfersItem) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 1)
	as[0] = it.receiver
	return as, nil
}

func (it TransfersItem) Amount() common.Big {
	return it.amount
}

func (it TransfersItem) Currency() types.CurrencyID {
	return it.currency
}
