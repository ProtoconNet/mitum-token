package token

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var TransfersFromItemHint = hint.MustNewHint("mitum-token-transfers-from-item-v0.0.1")

type TransfersFromItem struct {
	hint.BaseHinter
	contract base.Address
	receiver base.Address
	target   base.Address
	amount   common.Big
	currency types.CurrencyID
}

func NewTransfersFromItem(contract base.Address, receiver, target base.Address, amount common.Big, currency types.CurrencyID) TransfersFromItem {
	return TransfersFromItem{
		BaseHinter: hint.NewBaseHinter(TransfersFromItemHint),
		contract:   contract,
		receiver:   receiver,
		target:     target,
		amount:     amount,
		currency:   currency,
	}
}

func (it TransfersFromItem) IsValid([]byte) error {
	if err := it.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := util.CheckIsValiders(nil, false, it.contract, it.receiver, it.target); err != nil {
		return err
	}

	if it.contract.Equal(it.receiver) {
		return common.ErrSelfTarget.Wrap(errors.Errorf("receiver %v is same with contract account", it.receiver))
	}

	if it.contract.Equal(it.target) {
		return common.ErrSelfTarget.Wrap(errors.Errorf("target %v is same with contract account", it.target))
	}

	if it.receiver.Equal(it.target) {
		return common.ErrSelfTarget.Wrap(errors.Errorf("receiver %v is same with target", it.receiver))
	}

	if !it.amount.OverZero() {
		return common.ErrFactInvalid.Wrap(common.ErrValOOR.Wrap(errors.Errorf("transfer amount must be over zero, got %v", it.amount)))
	}

	return util.CheckIsValiders(nil, false,
		it.BaseHinter,
		it.contract,
		it.receiver,
		it.target,
		it.currency,
	)
}

func (it TransfersFromItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		it.receiver.Bytes(),
		it.target.Bytes(),
		it.amount.Bytes(),
		it.currency.Bytes(),
	)
}

func (it TransfersFromItem) Contract() base.Address {
	return it.contract
}

func (it TransfersFromItem) Receiver() base.Address {
	return it.receiver
}

func (it TransfersFromItem) Target() base.Address {
	return it.target
}

func (it TransfersFromItem) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 2)
	as[0] = it.receiver
	as[1] = it.target
	return as, nil
}

func (it TransfersFromItem) Amount() common.Big {
	return it.amount
}

func (it TransfersFromItem) Currency() types.CurrencyID {
	return it.currency
}
