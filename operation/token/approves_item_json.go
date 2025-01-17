package token

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type ApprovesItemJSONMarshaler struct {
	hint.BaseHinter
	Contract base.Address     `json:"contract"`
	Approved base.Address     `json:"approved"`
	Amount   string           `json:"amount"`
	Currency types.CurrencyID `json:"currency"`
}

func (it ApprovesItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(ApprovesItemJSONMarshaler{
		BaseHinter: it.BaseHinter,
		Contract:   it.contract,
		Approved:   it.approved,
		Amount:     it.Amount().String(),
		Currency:   it.currency,
	})
}

type ApprovesItemJSONUnmarshaler struct {
	Hint     hint.Hint `json:"_hint"`
	Contract string    `json:"contract"`
	Approved string    `json:"approved"`
	Amount   string    `json:"amount"`
	Currency string    `json:"currency"`
}

func (it *ApprovesItem) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var u ApprovesItemJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *it)
	}

	if err := it.unpack(enc, u.Hint, u.Contract, u.Approved, u.Amount, u.Currency); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *it)
	}

	return nil
}
