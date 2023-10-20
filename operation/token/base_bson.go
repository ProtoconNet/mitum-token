package token

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum-token/utils"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

func (fact TokenFact) marshalMap() map[string]interface{} {
	return map[string]interface{}{
		"_hint":    fact.Hint().String(),
		"sender":   fact.sender,
		"contract": fact.contract,
		"currency": fact.currency,
		"hash":     fact.BaseFact.Hash().String(),
		"token":    fact.BaseFact.Token(),
	}
}

func (fact TokenFact) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":    fact.Hint().String(),
			"sender":   fact.sender,
			"contract": fact.contract,
			"currency": fact.currency,
			"hash":     fact.BaseFact.Hash().String(),
			"token":    fact.BaseFact.Token(),
		},
	)
}

type TokenFactBSONUnmarshaler struct {
	Hint     string `bson:"_hint"`
	Sender   string `bson:"sender"`
	Contract string `bson:"contract"`
	Currency string `bson:"currency"`
}

func (fact *TokenFact) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError(utils.ErrStringDecodeBSON(*fact))

	var ubf common.BaseFactBSONUnmarshaler

	if err := enc.Unmarshal(b, &ubf); err != nil {
		return e.Wrap(err)
	}

	fact.BaseFact.SetHash(valuehash.NewBytesFromString(ubf.Hash))
	fact.BaseFact.SetToken(ubf.Token)

	var uf TokenFactBSONUnmarshaler
	if err := bson.Unmarshal(b, &uf); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(uf.Hint)
	if err != nil {
		return e.Wrap(err)
	}
	fact.BaseHinter = hint.NewBaseHinter(ht)

	return fact.unpack(enc,
		uf.Sender,
		uf.Contract,
		uf.Currency,
	)
}

//func (op TokenOperation) MarshalBSON() ([]byte, error) {
//	return bsonenc.Marshal(
//		bson.M{
//			"_hint": op.Hint().String(),
//			"hash":  op.Hash().String(),
//			"fact":  op.Fact(),
//			"signs": op.Signs(),
//		})
//}
//
//func (op *TokenOperation) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
//	e := util.StringError(utils.ErrStringDecodeBSON(*op))
//
//	var ubo common.BaseOperation
//	if err := ubo.DecodeBSON(b, enc); err != nil {
//		return e.Wrap(err)
//	}
//
//	op.BaseOperation = ubo
//
//	return nil
//}
