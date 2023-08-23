package types

import (
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-token/utils"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

func (d *Design) unmarshal(enc encoder.Encoder, ht hint.Hint, t, s string, bp []byte) error {
	e := util.StringError(utils.ErrStringUnmarshal(*d))

	d.BaseHinter = hint.NewBaseHinter(ht)
	d.tokenID = types.CurrencyID(t)
	d.symbol = s

	if hinter, err := enc.Decode(bp); err != nil {
		return e.Wrap(err)
	} else if p, ok := hinter.(Policy); !ok {
		return e.Wrap(errors.Errorf("expected %T, not %T", Policy{}, hinter))
	} else {
		d.policy = p
	}

	return nil
}
