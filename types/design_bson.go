package types

import (
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum-token/utils"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

func (d Design) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":   d.Hint().String(),
			"symbol":  d.symbol,
			"name":    d.name,
			"decimal": d.decimal,
			"policy":  d.policy,
		},
	)
}

type DesignBSONUnmarshaler struct {
	Hint    string   `bson:"_hint"`
	Symbol  string   `bson:"symbol"`
	Name    string   `bson:"name"`
	Decimal string   `bson:"decimal"`
	Policy  bson.Raw `bson:"policy"`
}

func (d *Design) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError(utils.ErrStringDecodeBSON(*d))

	var u DesignBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e.Wrap(err)
	}

	return d.unpack(enc, ht, u.Symbol, u.Name, u.Decimal, u.Policy)
}
