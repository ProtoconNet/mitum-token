package token

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
)

func (fact MintFact) MarshalBSON() ([]byte, error) {
	m := fact.TokenFact.marshalMap()

	m["receiver"] = fact.receiver
	m["amount"] = fact.amount

	return bsonenc.Marshal(m)
}

type MintFactBSONUnmarshaler struct {
	Receiver string `bson:"receiver"`
	Amount   string `bson:"amount"`
}

func (fact *MintFact) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	if err := fact.TokenFact.DecodeBSON(b, enc); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *fact)
	}

	var uf MintFactBSONUnmarshaler
	if err := bson.Unmarshal(b, &uf); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *fact)
	}

	if err := fact.unpack(enc, uf.Receiver, uf.Amount); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *fact)
	}

	return nil
}
