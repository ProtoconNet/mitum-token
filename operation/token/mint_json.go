package token

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-token/utils"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
)

type MintFactJSONMarshaler struct {
	TokenFactJSONMarshaler
	Receiver base.Address `json:"receiver"`
	Amount   common.Big   `json:"amount"`
}

func (fact MintFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(MintFactJSONMarshaler{
		TokenFactJSONMarshaler: fact.TokenFact.JSONMarshaler(),
		Receiver:               fact.receiver,
		Amount:                 fact.amount,
	})
}

type MintFactJSONUnMarshaler struct {
	Receiver string `json:"receiver"`
	Amount   string `json:"amount"`
}

func (fact *MintFact) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError(utils.ErrStringDecodeJSON(*fact))

	if err := fact.TokenFact.DecodeJSON(b, enc); err != nil {
		return e.Wrap(err)
	}

	var uf MintFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return e.Wrap(err)
	}

	return fact.unmarshal(enc,
		uf.Receiver,
		uf.Amount,
	)
}
