package types

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-token/utils"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type DesignJSONMarshaler struct {
	hint.BaseHinter
	Symbol  TokenSymbol `json:"symbol"`
	Name    string      `json:"name"`
	Decimal string      `json:"decimal"`
	Policy  Policy      `json:"policy"`
}

func (d Design) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(DesignJSONMarshaler{
		BaseHinter: d.BaseHinter,
		Symbol:     d.symbol,
		Name:       d.name,
		Decimal:    d.decimal.String(),
		Policy:     d.policy,
	})
}

type DesignJSONUnmarshaler struct {
	Hint    hint.Hint       `json:"_hint"`
	Symbol  string          `json:"symbol"`
	Name    string          `json:"name"`
	Decimal string          `json:"decimal"`
	Policy  json.RawMessage `json:"policy"`
}

func (d *Design) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError(utils.ErrStringDecodeJSON(*d))

	var u DesignJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	return d.unpack(enc, u.Hint, u.Symbol, u.Name, u.Decimal, u.Policy)
}
