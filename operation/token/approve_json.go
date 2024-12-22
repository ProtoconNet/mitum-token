package token

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/extras"
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
	if err := fact.TokenFact.DecodeJSON(b, enc); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	var uf ApproveFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	if err := fact.unpack(enc, uf.Approved, uf.Amount); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	return nil
}

type OperationMarshaler struct {
	common.BaseOperationJSONMarshaler
	extras.BaseOperationExtensionsJSONMarshaler
}

func (op Approve) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(OperationMarshaler{
		BaseOperationJSONMarshaler:           op.BaseOperation.JSONMarshaler(),
		BaseOperationExtensionsJSONMarshaler: op.BaseOperationExtensions.JSONMarshaler(),
	})
}

func (op *Approve) DecodeJSON(b []byte, enc encoder.Encoder) error {
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
