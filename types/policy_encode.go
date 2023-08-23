package types

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-token/utils"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (p *Policy) unmarshal(enc encoder.Encoder, ht hint.Hint, ts string, bap []byte) error {
	e := util.StringError(utils.ErrStringUnmarshal(*p))

	p.BaseHinter = hint.NewBaseHinter(ht)

	big, err := common.NewBigFromString(ts)
	if err != nil {
		return e.Wrap(err)
	}
	p.totalSupply = big

	hap, err := enc.DecodeSlice(bap)
	if err != nil {
		return e.Wrap(err)
	}

	al := make([]ApproveInfo, len(hap))
	for i, h := range hap {
		ap, ok := h.(ApproveInfo)
		if !ok {
			return e.Wrap(util.ErrInvalid.Errorf("expected %T, not %T", ApproveInfo{}, h))
		}

		al[i] = ap
	}
	p.approveList = al

	return nil
}
