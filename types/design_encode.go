package types

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-token/utils"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

func (d *Design) unpack(enc encoder.Encoder, ht hint.Hint, symbol, name, decimal string, bp []byte) error {
	e := util.StringError(utils.ErrStringUnPack(*d))

	d.BaseHinter = hint.NewBaseHinter(ht)
	d.symbol = TokenSymbol(symbol)
	d.name = name
	if dc, err := common.NewBigFromString(decimal); err != nil {
		return err
	} else {
		d.decimal = dc
	}

	if hinter, err := enc.Decode(bp); err != nil {
		return e.Wrap(err)
	} else if p, ok := hinter.(Policy); !ok {
		return e.Wrap(errors.Errorf(utils.ErrStringTypeCast(Policy{}, hinter)))
	} else {
		d.policy = p
	}

	return nil
}
