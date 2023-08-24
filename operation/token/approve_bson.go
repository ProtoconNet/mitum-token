package token

import (
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum-token/utils"
	"github.com/ProtoconNet/mitum2/util"
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
	e := util.StringError(utils.ErrStringDecodeBSON(*fact))

	if err := fact.TokenFact.DecodeBSON(b, enc); err != nil {
		return e.Wrap(err)
	}

	var uf ApproveFactBSONUnmarshaler
	if err := bson.Unmarshal(b, &uf); err != nil {
		return e.Wrap(err)
	}

	return fact.unmarshal(enc,
		uf.Approved,
		uf.Amount,
	)
}
