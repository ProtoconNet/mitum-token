package token

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-token/utils"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
)

type RegisterTokenFactJSONMarshaler struct {
	TokenFactJSONMarshaler
	Symbol      string     `json:"symbol"`
	TotalSupply common.Big `json:"total_supply"`
}

func (fact RegisterTokenFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(RegisterTokenFactJSONMarshaler{
		TokenFactJSONMarshaler: fact.TokenFact.JSONMarshaler(),
		Symbol:                 fact.symbol,
		TotalSupply:            fact.totalSupply,
	})
}

type RegisterTokenFactJSONUnMarshaler struct {
	Symbol      string `json:"symbol"`
	TotalSupply string `json:"total_supply"`
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

	return fact.unmarshal(enc,
		uf.Symbol,
		uf.TotalSupply,
	)
}
