package types

import (
	"github.com/ProtoconNet/mitum-token/utils"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (a *ApproveBox) unpack(enc encoder.Encoder, ht hint.Hint, ac string, bap []byte) error {
	e := util.StringError(utils.ErrStringUnmarshal(*a))

	a.BaseHinter = hint.NewBaseHinter(ht)

	switch ad, err := base.DecodeAddress(ac, enc); {
	case err != nil:
		return e.Wrap(err)
	default:
		a.account = ad
	}

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
	a.approved = al

	return nil
}
