package types

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-token/utils"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var ApproveInfoHint = hint.MustNewHint("mitum-token-approve-info-v0.0.1")

type ApproveInfo struct {
	hint.BaseHinter
	account base.Address
	amount  common.Big
}

func NewApproveInfo(account base.Address, amount common.Big) ApproveInfo {
	return ApproveInfo{
		BaseHinter: hint.NewBaseHinter(ApproveInfoHint),
		account:    account,
		amount:     amount,
	}
}

func (a ApproveInfo) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf(utils.ErrStringInvalid(a))

	if err := util.CheckIsValiders(nil, false,
		a.BaseHinter,
		a.account,
	); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (a ApproveInfo) Bytes() []byte {
	return util.ConcatBytesSlice(
		a.account.Bytes(),
		a.amount.Bytes(),
	)
}

func (a ApproveInfo) Account() base.Address {
	return a.account
}

func (a ApproveInfo) Amount() common.Big {
	return a.amount
}
