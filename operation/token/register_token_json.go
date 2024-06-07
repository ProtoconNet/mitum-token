package token

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-token/types"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type RegisterTokenFactJSONMarshaler struct {
	TokenFactJSONMarshaler
	Symbol        types.TokenSymbol `json:"symbol"`
	Name          string            `json:"name"`
	InitialSupply common.Big        `json:"initial_supply"`
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

func (fact *RegisterTokenFact) DecodeJSON(b []byte, enc encoder.Encoder) error {
	if err := fact.TokenFact.DecodeJSON(b, enc); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	var uf RegisterTokenFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	if err := fact.unpack(enc, uf.Symbol, uf.Name, uf.InitialSupply); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	return nil
}
