package token

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-token/utils"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
)

type BurnFactJSONMarshaler struct {
	TokenFactJSONMarshaler
	Target base.Address `json:"target"`
	Amount common.Big   `json:"amount"`
}

func (fact BurnFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(BurnFactJSONMarshaler{
		TokenFactJSONMarshaler: fact.TokenFact.JSONMarshaler(),
		Target:                 fact.target,
		Amount:                 fact.amount,
	})
}

type BurnFactJSONUnMarshaler struct {
	Target string `json:"target"`
	Amount string `json:"amount"`
}

func (fact *BurnFact) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError(utils.ErrStringDecodeJSON(*fact))

	if err := fact.TokenFact.DecodeJSON(b, enc); err != nil {
		return e.Wrap(err)
	}

	var uf BurnFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return e.Wrap(err)
	}

	return fact.unmarshal(enc,
		uf.Target,
		uf.Amount,
	)
}
