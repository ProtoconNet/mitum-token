package token

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var ApprovesItemHint = hint.MustNewHint("mitum-token-approves-item-v0.0.1")

type ApprovesItem struct {
	hint.BaseHinter
	contract base.Address
	approved base.Address
	amount   common.Big
	currency types.CurrencyID
}

func NewApprovesItem(contract base.Address, approved base.Address, amount common.Big, currency types.CurrencyID) ApprovesItem {
	return ApprovesItem{
		BaseHinter: hint.NewBaseHinter(ApprovesItemHint),
		contract:   contract,
		approved:   approved,
		amount:     amount,
		currency:   currency,
	}
}

func (it ApprovesItem) IsValid([]byte) error {
	if err := it.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := util.CheckIsValiders(nil, false, it.contract, it.approved); err != nil {
		return err
	}

	if it.approved.Equal(it.contract) {
		return common.ErrSelfTarget.Wrap(errors.Errorf("approved %v is same with contract account", it.approved))
	}

	if !it.amount.OverZero() {
		return common.ErrFactInvalid.Wrap(common.ErrValOOR.Wrap(errors.Errorf("approved amount must be over zero, got %v", it.amount)))
	}

	return util.CheckIsValiders(nil, false,
		it.BaseHinter,
		it.contract,
		it.approved,
		it.currency,
	)
}

func (it ApprovesItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		it.approved.Bytes(),
		it.amount.Bytes(),
		it.currency.Bytes(),
	)
}

func (it ApprovesItem) Contract() base.Address {
	return it.contract
}

func (it ApprovesItem) Approved() base.Address {
	return it.approved
}

func (it ApprovesItem) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 1)
	as[0] = it.approved
	return as, nil
}

func (it ApprovesItem) Amount() common.Big {
	return it.amount
}

func (it ApprovesItem) Currency() types.CurrencyID {
	return it.currency
}
