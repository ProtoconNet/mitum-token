package token

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/extras"
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
)

func (fact ApproveFact) MarshalBSON() ([]byte, error) {
	m := fact.TokenFact.marshalMap()
	m["approved"] = fact.approved
	m["amount"] = fact.amount

	return bsonenc.Marshal(m)
}

type ApproveFactBSONUnmarshaler struct {
	Approved string `bson:"approved"`
	Amount   string `bson:"amount"`
}

func (fact *ApproveFact) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	if err := fact.TokenFact.DecodeBSON(b, enc); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *fact)
	}

	var uf ApproveFactBSONUnmarshaler
	if err := bson.Unmarshal(b, &uf); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *fact)
	}

	if err := fact.unpack(enc, uf.Approved, uf.Amount); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *fact)
	}

	return nil
}

func (op *Approve) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	var ubo common.BaseOperation
	if err := ubo.DecodeBSON(b, enc); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *op)
	}

	op.BaseOperation = ubo

	var ueo extras.BaseOperationExtensions
	if err := ueo.DecodeBSON(b, enc); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *op)
	}

	op.BaseOperationExtensions = &ueo

	return nil
}
