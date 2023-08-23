package token

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-token/utils"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (fact *RegisterTokenFact) unmarshal(enc encoder.Encoder,
	s, ts string,
) error {
	e := util.StringError(utils.ErrStringUnmarshal(*fact))

	fact.symbol = s

	big, err := common.NewBigFromString(ts)
	if err != nil {
		return e.Wrap(err)
	}
	fact.totalSupply = big

	return nil
}
