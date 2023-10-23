package types

import (
	"bytes"
	"github.com/ProtoconNet/mitum2/base"
	"sort"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-token/utils"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var PolicyHint = hint.MustNewHint("mitum-token-policy-v0.0.1")

type Policy struct {
	hint.BaseHinter
	totalSupply common.Big
	approveList []ApproveBox
}

func NewPolicy(totalSupply common.Big, approveList []ApproveBox) Policy {
	return Policy{
		BaseHinter:  hint.NewBaseHinter(PolicyHint),
		totalSupply: totalSupply,
		approveList: approveList,
	}
}

func (p Policy) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf(utils.ErrStringInvalid(p))

	if err := p.BaseHinter.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	founds := map[string]struct{}{}
	for _, a := range p.approveList {
		if err := a.IsValid(nil); err != nil {
			return e.Wrap(err)
		}

		if _, ok := founds[a.account.String()]; ok {
			return e.Wrap(errors.Errorf(utils.ErrStringDuplicate("account", a.account.String())))
		}

		founds[a.account.String()] = struct{}{}
	}

	if !p.totalSupply.OverNil() {
		return e.Wrap(errors.Errorf("nil big"))
	}

	return nil
}

func (p Policy) Bytes() []byte {
	b := make([][]byte, len(p.approveList))
	for i, a := range p.approveList {
		b[i] = a.Bytes()
	}

	sort.Slice(b, func(i, j int) bool {
		return bytes.Compare(b[i], b[j]) < 1
	})

	return util.ConcatBytesSlice(
		p.totalSupply.Bytes(),
		util.ConcatBytesSlice(b...),
	)
}

func (p Policy) TotalSupply() common.Big {
	return p.totalSupply
}

func (p Policy) ApproveList() []ApproveBox {
	return p.approveList
}

func (p Policy) GetApproveBox(acc base.Address) *ApproveBox {
	var approvedBox ApproveBox
	idx := -1
	for i, apb := range p.approveList {
		if apb.Account().Equal(acc) {
			idx = i
			approvedBox = apb
			break
		}
	}
	if idx == -1 {
		return nil
	}
	return &approvedBox
}

func (p *Policy) MergeApproveBox(apb ApproveBox) {
	var approvedList = make([]ApproveBox, len(p.approveList))
	copy(approvedList, p.approveList)
	idx := -1
	for i, apb := range approvedList {
		if apb.Account().Equal(apb.Account()) {
			idx = i
			break
		}
	}
	if -1 < idx {
		approvedList[idx] = apb
	} else {
		approvedList = append(approvedList, apb)
	}
	p.approveList = approvedList
}

func (p *Policy) RemoveApproveBox(acc base.Address) {
	var approvedList []ApproveBox

	idx := -1
	for i, apb := range approvedList {
		if apb.Account().Equal(acc) {
			idx = i
			break
		}
	}
	if -1 < idx {
		approvedList = append(approvedList, p.approveList[:idx]...)
		approvedList = append(approvedList, p.approveList[idx+1:]...)
	}
	p.approveList = approvedList
}
