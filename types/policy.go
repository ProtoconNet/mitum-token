package types

import (
	"bytes"
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

func (p Policy) ApproveList() []ApproveInfo {
	return p.approveList
}
