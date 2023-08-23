package types

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-token/utils"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type DesignJSONMarshaler struct {
	hint.BaseHinter
	TokenID types.CurrencyID `json:"token_id"`
	Symbol  string           `json:"symbol"`
	Policy  Policy           `json:"policy"`
}

func (d Design) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(DesignJSONMarshaler{
		BaseHinter: d.BaseHinter,
		TokenID:    d.tokenID,
		Symbol:     d.symbol,
	})
}

type DesignJSONUnmarshaler struct {
	Hint    hint.Hint       `json:"_hint"`
	TokenID string          `json:"token_id"`
	Symbol  string          `json:"symbol"`
	Policy  json.RawMessage `json:"policy"`
}

func (d *Design) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError(utils.ErrStringDecodeJSON(*d))

	var u DesignJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	return d.unmarshal(enc, u.Hint, u.TokenID, u.Symbol, u.Policy)
}
