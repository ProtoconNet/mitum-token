package token

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/extras"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
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

func (fact *BurnFact) DecodeJSON(b []byte, enc encoder.Encoder) error {
	if err := fact.TokenFact.DecodeJSON(b, enc); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	var uf BurnFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	if err := fact.unpack(enc, uf.Target, uf.Amount); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	return nil
}

func (op Burn) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(OperationMarshaler{
		BaseOperationJSONMarshaler:           op.BaseOperation.JSONMarshaler(),
		BaseOperationExtensionsJSONMarshaler: op.BaseOperationExtensions.JSONMarshaler(),
	})
}

func (op *Burn) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *op)
	}

	op.BaseOperation = ubo

	var ueo extras.BaseOperationExtensions
	if err := ueo.DecodeJSON(b, enc); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *op)
	}

	op.BaseOperationExtensions = &ueo

	return nil
}
