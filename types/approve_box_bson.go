package types

import (
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum-token/utils"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

func (a ApproveBox) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":    a.Hint().String(),
			"account":  a.account,
			"approved": a.approved,
		},
	)
}

type ApproveBoxBSONUnmarshaler struct {
	Hint     string   `bson:"_hint"`
	Account  string   `bson:"account"`
	Approved bson.Raw `bson:"approved"`
}

func (a *ApproveBox) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError(utils.ErrStringDecodeBSON(*a))

	var u ApproveBoxBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e.Wrap(err)
	}

	return a.unpack(enc, ht, u.Account, u.Approved)
}
