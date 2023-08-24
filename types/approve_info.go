package types

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-token/utils"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var ApproveInfoHint = hint.MustNewHint("mitum-token-approve-info-v0.0.1")

type ApproveInfo struct {
	hint.BaseHinter
	account  base.Address
	approved map[string]common.Big
}

func NewApproveInfo(account base.Address, approved map[string]common.Big) ApproveInfo {
	return ApproveInfo{
		BaseHinter: hint.NewBaseHinter(ApproveInfoHint),
		account:    account,
		approved:   approved,
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

	founds := map[string]struct{}{}
	for ac, big := range a.approved {
		if !big.OverZero() {
			return e.Wrap(errors.Errorf("zero big"))
		}

		if _, ok := founds[ac]; ok {
			return e.Wrap(errors.Errorf(utils.ErrStringDuplicate("approved", ac)))
		}

		founds[ac] = struct{}{}
	}

	return nil
}

func (a ApproveInfo) Bytes() []byte {
	b, _ := json.Marshal(a.approved)

	return util.ConcatBytesSlice(
		a.account.Bytes(),
		b,
	)
}

func (a ApproveInfo) Account() base.Address {
	return a.account
}

func (a ApproveInfo) Approved() map[string]common.Big {
	return a.approved
}
