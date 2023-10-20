package token

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-token/utils"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
)

type TokenFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Sender   base.Address             `json:"sender"`
	Contract base.Address             `json:"contract"`
	Currency currencytypes.CurrencyID `json:"currency"`
}

func (fact TokenFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(TokenFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Sender:                fact.sender,
		Contract:              fact.contract,
		Currency:              fact.currency,
	})
}

type TokenFactJSONUnmarshaler struct {
	base.BaseFactJSONUnmarshaler
	Sender   string `json:"sender"`
	Contract string `json:"contract"`
	Currency string `json:"currency"`
}

func (fact *TokenFact) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError(utils.ErrStringDecodeJSON(*fact))

	var uf TokenFactJSONUnmarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return e.Wrap(err)
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	return fact.unpack(enc,
		uf.Sender,
		uf.Contract,
		uf.Currency,
	)
}

func (fact TokenFact) JSONMarshaler() TokenFactJSONMarshaler {
	return TokenFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Sender:                fact.sender,
		Contract:              fact.contract,
		Currency:              fact.currency,
	}
}
