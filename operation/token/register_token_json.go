package token

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-token/utils"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
)

type RegisterTokenFactJSONMarshaler struct {
	TokenFactJSONMarshaler
	Symbol        currencytypes.CurrencyID `json:"symbol"`
	Name          string                   `json:"name"`
	InitialSupply common.Big               `json:"initial_supply"`
}

func (fact RegisterTokenFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(RegisterTokenFactJSONMarshaler{
		TokenFactJSONMarshaler: fact.TokenFact.JSONMarshaler(),
		Symbol:                 fact.symbol,
		Name:                   fact.name,
		InitialSupply:          fact.initialSupply,
	})
}

type RegisterTokenFactJSONUnMarshaler struct {
	Symbol        string `json:"symbol"`
	Name          string `json:"name"`
	InitialSupply string `json:"initial_supply"`
}

func (fact *RegisterTokenFact) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError(utils.ErrStringDecodeJSON(*fact))

	if err := fact.TokenFact.DecodeJSON(b, enc); err != nil {
		return e.Wrap(err)
	}

	var uf RegisterTokenFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return e.Wrap(err)
	}

	return fact.unpack(enc,
		uf.Symbol,
		uf.Name,
		uf.InitialSupply,
	)
}
