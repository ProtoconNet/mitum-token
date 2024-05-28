package token

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
)

func (fact BurnFact) MarshalBSON() ([]byte, error) {
	m := fact.TokenFact.marshalMap()

	m["target"] = fact.target
	m["amount"] = fact.amount

	return bsonenc.Marshal(m)
}

type BurnFactBSONUnmarshaler struct {
	Target string `bson:"target"`
	Amount string `bson:"amount"`
}

func (fact *BurnFact) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	if err := fact.TokenFact.DecodeBSON(b, enc); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *fact)
	}

	var uf BurnFactBSONUnmarshaler
	if err := bson.Unmarshal(b, &uf); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *fact)
	}

	if err := fact.unpack(enc, uf.Target, uf.Amount); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *fact)
	}

	return nil
}
