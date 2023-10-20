package types

import (
	"github.com/ProtoconNet/mitum-token/utils"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var ApproveBoxHint = hint.MustNewHint("mitum-token-approve-box-v0.0.1")

type ApproveBox struct {
	hint.BaseHinter
	account  base.Address
	approved []ApproveInfo
}

func NewApproveBox(account base.Address, approved []ApproveInfo) ApproveBox {
	return ApproveBox{
		BaseHinter: hint.NewBaseHinter(ApproveBoxHint),
		account:    account,
		approved:   approved,
	}
}

func (a ApproveBox) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf(utils.ErrStringInvalid(a))

	if err := util.CheckIsValiders(nil, false,
		a.BaseHinter,
		a.account,
	); err != nil {
		return e.Wrap(err)
	}

	founds := map[string]struct{}{}
	for i := range a.approved {
		_, found := founds[a.approved[i].Account().String()]
		if found {
			return e.Wrap(errors.Errorf(utils.ErrStringDuplicate("approved", a.approved[i].Account().String())))
		} else {
			founds[a.approved[i].Account().String()] = struct{}{}
		}
	}

	return nil
}

func (a ApproveBox) Bytes() []byte {
	bs := make([][]byte, len(a.approved)+1)
	for i := range a.approved {
		bs[i] = a.approved[i].Bytes()
	}
	bs[len(a.approved)] = a.account.Bytes()

	return util.ConcatBytesSlice(bs...)
}

func (a ApproveBox) Account() base.Address {
	return a.account
}

func (a ApproveBox) Approved() []ApproveInfo {
	return a.approved
}

func (a ApproveBox) GetApproveInfo(ad base.Address) *ApproveInfo {
	for i := range a.approved {
		if ad.Equal(a.approved[i].Account()) {
			return &a.approved[i]
		}
	}
	return nil
}

func (a *ApproveBox) RemoveApproveInfo(ad base.Address) error {
	var nApproved []ApproveInfo
	for i := range a.approved {
		if ad.Equal(a.approved[i].Account()) {
			nApproved = append(nApproved, a.approved[:i]...)
			nApproved = append(nApproved, a.approved[i+1:]...)
			return nil
		}
		if i == len(a.approved)-1 {
			return errors.Errorf("not found approved, %s", ad)
		}
	}
	return nil
}

func (a *ApproveBox) SetApproveInfo(ap ApproveInfo) {
	var approved []ApproveInfo
	var count int
	for i := range a.approved {
		if ap.Account().Equal(a.approved[i].Account()) {
			approved = append(approved, ap)
		} else {
			approved = append(approved, a.approved[i])
			count = count + 1
		}
	}
	if count == len(a.approved) {
		approved = append(approved, ap)
	}

	a.approved = approved
	return
}
