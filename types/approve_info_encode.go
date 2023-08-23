package types

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-token/utils"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (a *ApproveInfo) unmarshal(enc encoder.Encoder, ht hint.Hint, ac string, bap []byte) error {
	e := util.StringError(utils.ErrStringUnmarshal(*a))

	a.BaseHinter = hint.NewBaseHinter(ht)

	switch ad, err := base.DecodeAddress(ac, enc); {
	case err != nil:
		return e.Wrap(err)
	default:
		a.account = ad
	}

	m, err := utils.DecodeMap(bap)
	if err != nil {
		return e.Wrap(err)
	}

	approved := make(map[string]common.Big)
	for k, v := range m {
		big, err := common.NewBigFromInterface(v)
		if err != nil {
			return e.Wrap(err)
		}

		approved[k] = big
	}
	a.approved = approved

	return nil
}
