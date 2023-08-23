package types

import (
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum-token/utils"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

func (p Policy) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":        p.Hint().String(),
			"total_supply": p.totalSupply,
			"approve_list": p.approveList,
		},
	)
}

type PolicyBSONUnmarshaler struct {
	Hint        string   `bson:"_hint"`
	TotalSupply string   `bson:"total_supply"`
	ApproveList bson.Raw `bson:"approve_list"`
}

func (p *Policy) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError(utils.ErrStringDecodeBSON(*p))

	var u PolicyBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e.Wrap(err)
	}

	return p.unmarshal(enc, ht, u.TotalSupply, u.ApproveList)
}
