package token

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-token/utils"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (fact *RegisterTokenFact) unpack(enc encoder.Encoder,
	symbol, name, ts string,
) error {
	e := util.StringError(utils.ErrStringUnPack(*fact))

	fact.symbol = currencytypes.CurrencyID(symbol)
	fact.name = name

	big, err := common.NewBigFromString(ts)
	if err != nil {
		return e.Wrap(err)
	}
	fact.initialSupply = big

	return nil
}
