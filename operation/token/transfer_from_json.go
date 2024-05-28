package token

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type TransferFromFactJSONMarshaler struct {
	TokenFactJSONMarshaler
	Receiver base.Address `json:"receiver"`
	Target   base.Address `json:"target"`
	Amount   common.Big   `json:"amount"`
}

func (fact TransferFromFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(TransferFromFactJSONMarshaler{
		TokenFactJSONMarshaler: fact.TokenFact.JSONMarshaler(),
		Receiver:               fact.receiver,
		Target:                 fact.target,
		Amount:                 fact.amount,
	})
}

type TransferFromFactJSONUnMarshaler struct {
	Receiver string `json:"receiver"`
	Target   string `json:"target"`
	Amount   string `json:"amount"`
}

func (fact *TransferFromFact) DecodeJSON(b []byte, enc encoder.Encoder) error {
	if err := fact.TokenFact.DecodeJSON(b, enc); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	var uf TransferFromFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	if err := fact.unpack(enc, uf.Receiver, uf.Target, uf.Amount); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	return nil
}
