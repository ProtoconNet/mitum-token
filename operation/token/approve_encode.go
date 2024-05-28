package token

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (fact *ApproveFact) unpack(enc encoder.Encoder,
	ap, am string,
) error {
	switch a, err := base.DecodeAddress(ap, enc); {
	case err != nil:
		return err
	default:
		fact.approved = a
	}

	big, err := common.NewBigFromString(am)
	if err != nil {
		return err
	}
	fact.amount = big

	return nil
}
