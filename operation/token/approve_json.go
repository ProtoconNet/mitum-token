package token

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-token/utils"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type ApproveFactJSONMarshaler struct {
	TokenFactJSONMarshaler
	Approved base.Address `json:"approved"`
	Amount   common.Big   `json:"amount"`
}

func (fact ApproveFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(ApproveFactJSONMarshaler{
		TokenFactJSONMarshaler: fact.TokenFact.JSONMarshaler(),
		Approved:               fact.approved,
		Amount:                 fact.amount,
	})
}

type ApproveFactJSONUnMarshaler struct {
	Approved string `json:"approved"`
	Amount   string `json:"amount"`
}

func (fact *ApproveFact) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError(utils.ErrStringDecodeJSON(*fact))

	if err := fact.TokenFact.DecodeJSON(b, enc); err != nil {
		return e.Wrap(err)
	}

	var uf ApproveFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return e.Wrap(err)
	}

	return fact.unpack(enc,
		uf.Approved,
		uf.Amount,
	)
}
