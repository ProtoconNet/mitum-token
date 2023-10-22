package token

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-token/utils"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (fact *TransferFromFact) unpack(enc encoder.Encoder,
	ra, ta, am string,
) error {
	e := util.StringError(utils.ErrStringUnPack(*fact))

	switch a, err := base.DecodeAddress(ra, enc); {
	case err != nil:
		return e.Wrap(err)
	default:
		fact.receiver = a
	}

	switch a, err := base.DecodeAddress(ta, enc); {
	case err != nil:
		return e.Wrap(err)
	default:
		fact.target = a
	}

	big, err := common.NewBigFromString(am)
	if err != nil {
		return e.Wrap(err)
	}
	fact.amount = big

	return nil
}
