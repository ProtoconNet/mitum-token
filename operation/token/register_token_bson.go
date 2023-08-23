package token

import (
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum-token/utils"
	"github.com/ProtoconNet/mitum2/util"
)

func (fact RegisterTokenFact) MarshalBSON() ([]byte, error) {
	m := fact.TokenFact.marshalMap()

	m["symbol"] = fact.symbol
	m["total_supply"] = fact.totalSupply

	return bsonenc.Marshal(m)
}

type RegisterTokenFactBSONUnmarshaler struct {
	Symbol      string `bson:"symbol"`
	TotalSupply string `bson:"total_supply"`
}

func (fact *RegisterTokenFact) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError(utils.ErrStringDecodeBSON(*fact))

	if err := fact.TokenFact.DecodeBSON(b, enc); err != nil {
		return e.Wrap(err)
	}

	var uf RegisterTokenFactBSONUnmarshaler
	if err := bson.Unmarshal(b, &uf); err != nil {
		return e.Wrap(err)
	}

	return fact.unmarshal(enc,
		uf.Symbol,
		uf.TotalSupply,
	)
}
