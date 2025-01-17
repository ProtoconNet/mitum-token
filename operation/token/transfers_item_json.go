package token

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type TransfersItemJSONMarshaler struct {
	hint.BaseHinter
	Contract base.Address     `json:"contract"`
	Receiver base.Address     `json:"receiver"`
	Amount   string           `json:"amount"`
	Currency types.CurrencyID `json:"currency"`
}

func (it TransfersItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(TransfersItemJSONMarshaler{
		BaseHinter: it.BaseHinter,
		Contract:   it.contract,
		Receiver:   it.receiver,
		Amount:     it.Amount().String(),
		Currency:   it.currency,
	})
}

type TransfersItemJSONUnmarshaler struct {
	Hint     hint.Hint `json:"_hint"`
	Contract string    `json:"contract"`
	Receiver string    `json:"receiver"`
	Amount   string    `json:"amount"`
	Currency string    `json:"currency"`
}

func (it *TransfersItem) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var u TransfersItemJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *it)
	}

	if err := it.unpack(enc, u.Hint, u.Contract, u.Receiver, u.Amount, u.Currency); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *it)
	}

	return nil
}
