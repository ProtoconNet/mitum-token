package token

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

func (it TransfersFromItem) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":    it.Hint().String(),
			"contract": it.contract,
			"receiver": it.receiver,
			"target":   it.target,
			"amount":   it.amount,
			"currency": it.currency,
		},
	)
}

type TransfersFromItemBSONUnmarshaler struct {
	Hint     string `bson:"_hint"`
	Contract string `bson:"contract"`
	Receiver string `bson:"receiver"`
	Target   string `bson:"target"`
	Amount   string `bson:"amount"`
	Currency string `bson:"currency"`
}

func (it *TransfersFromItem) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	var u TransfersFromItemBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *it)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *it)
	}

	if err := it.unpack(enc, ht, u.Contract, u.Receiver, u.Target, u.Amount, u.Currency); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *it)
	}
	return nil
}
