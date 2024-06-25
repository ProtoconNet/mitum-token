package token

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-token/types"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type RegisterModelFactJSONMarshaler struct {
	TokenFactJSONMarshaler
	Symbol        types.TokenSymbol `json:"symbol"`
	Name          string            `json:"name"`
	Decimal       common.Big        `json:"decimal"`
	InitialSupply common.Big        `json:"initial_supply"`
}

func (fact RegisterModelFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(RegisterModelFactJSONMarshaler{
		TokenFactJSONMarshaler: fact.TokenFact.JSONMarshaler(),
		Symbol:                 fact.symbol,
		Name:                   fact.name,
		Decimal:                fact.decimal,
		InitialSupply:          fact.initialSupply,
	})
}

type RegisterModelFactJSONUnMarshaler struct {
	Symbol        string `json:"symbol"`
	Name          string `json:"name"`
	Decimal       string `json:"decimal"`
	InitialSupply string `json:"initial_supply"`
}

func (fact *RegisterModelFact) DecodeJSON(b []byte, enc encoder.Encoder) error {
	if err := fact.TokenFact.DecodeJSON(b, enc); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	var uf RegisterModelFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	if err := fact.unpack(enc, uf.Symbol, uf.Name, uf.Decimal, uf.InitialSupply); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	return nil
}
