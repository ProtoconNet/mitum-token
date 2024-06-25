package token

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
)

func (fact RegisterModelFact) MarshalBSON() ([]byte, error) {
	m := fact.TokenFact.marshalMap()

	m["symbol"] = fact.symbol
	m["name"] = fact.name
	m["decimal"] = fact.decimal
	m["initial_supply"] = fact.initialSupply

	return bsonenc.Marshal(m)
}

type RegisterModelFactBSONUnmarshaler struct {
	Symbol        string `bson:"symbol"`
	Name          string `bson:"name"`
	Decimal       string `bson:"decimal"`
	InitialSupply string `bson:"initial_supply"`
}

func (fact *RegisterModelFact) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	if err := fact.TokenFact.DecodeBSON(b, enc); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *fact)
	}

	var uf RegisterModelFactBSONUnmarshaler
	if err := bson.Unmarshal(b, &uf); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *fact)
	}

	if err := fact.unpack(enc, uf.Symbol, uf.Name, uf.Decimal, uf.InitialSupply); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *fact)
	}

	return nil
}
