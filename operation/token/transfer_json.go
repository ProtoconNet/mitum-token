package token

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-token/utils"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
)

type TransferFactJSONMarshaler struct {
	TokenFactJSONMarshaler
	Receiver base.Address `json:"receiver"`
	Amount   common.Big   `json:"amount"`
}

func (fact TransferFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(TransferFactJSONMarshaler{
		TokenFactJSONMarshaler: fact.TokenFact.JSONMarshaler(),
		Receiver:               fact.receiver,
		Amount:                 fact.amount,
	})
}

type TransferFactJSONUnMarshaler struct {
	Receiver string `json:"receiver"`
	Amount   string `json:"amount"`
}

func (fact *TransferFact) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError(utils.ErrStringDecodeJSON(*fact))

	if err := fact.TokenFact.DecodeJSON(b, enc); err != nil {
		return e.Wrap(err)
	}

	var uf TransferFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return e.Wrap(err)
	}

	return fact.unpack(enc,
		uf.Receiver,
		uf.Amount,
	)
}
