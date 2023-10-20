package token

import (
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum-token/utils"
	"github.com/ProtoconNet/mitum2/util"
)

func (fact TransferFromFact) MarshalBSON() ([]byte, error) {
	m := fact.TokenFact.marshalMap()

	m["receiver"] = fact.receiver
	m["target"] = fact.target
	m["amount"] = fact.amount

	return bsonenc.Marshal(m)
}

type TransferFromFactBSONUnmarshaler struct {
	Receiver string `bson:"receiver"`
	Target   string `bson:"target"`
	Amount   string `bson:"amount"`
}

func (fact *TransferFromFact) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError(utils.ErrStringDecodeBSON(*fact))

	if err := fact.TokenFact.DecodeBSON(b, enc); err != nil {
		return e.Wrap(err)
	}

	var uf TransferFromFactBSONUnmarshaler
	if err := bson.Unmarshal(b, &uf); err != nil {
		return e.Wrap(err)
	}

	return fact.unpack(enc,
		uf.Receiver,
		uf.Target,
		uf.Amount,
	)
}
