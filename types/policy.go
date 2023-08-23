package types

import (
	"bytes"
	"sort"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var PolicyHint = hint.MustNewHint("mitum-token-policy-v0.0.1")

type Policy struct {
	hint.BaseHinter
	totalSupply common.Big
	approveList []ApproveInfo
}

func NewPolicy(totalSupply common.Big, approveList []ApproveInfo) Policy {
	return Policy{
		BaseHinter:  hint.NewBaseHinter(PolicyHint),
		totalSupply: totalSupply,
		approveList: approveList,
	}
}

func (p Policy) IsValid([]byte) error {
	if err := p.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	founds := map[string]struct{}{}
	for _, a := range p.approveList {
		if err := a.IsValid(nil); err != nil {
			return err
		}

		if _, ok := founds[a.account.String()]; ok {
			return util.ErrInvalid.Errorf("duplicate account found, %s", a.account)
		}

		founds[a.account.String()] = struct{}{}
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

func (p Policy) ApproveList() []ApproveInfo {
	return p.approveList
}
