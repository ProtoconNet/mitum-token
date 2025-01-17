package token

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (it *ApprovesItem) unpack(
	enc encoder.Encoder,
	ht hint.Hint,
	ca, apr, am, cid string,
) error {
	it.BaseHinter = hint.NewBaseHinter(ht)
	switch a, err := base.DecodeAddress(ca, enc); {
	case err != nil:
		return err
	default:
		it.contract = a
	}

	approved, err := base.DecodeAddress(apr, enc)
	if err != nil {
		return err
	}
	it.approved = approved

	if b, err := common.NewBigFromString(am); err != nil {
		return err
	} else {
		it.amount = b
	}
	it.currency = types.CurrencyID(cid)

	return nil
}
