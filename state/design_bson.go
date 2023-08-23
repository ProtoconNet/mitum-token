package state

import (
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum-token/types"
	"github.com/ProtoconNet/mitum-token/utils"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

func (s DesignStateValue) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":  s.Hint().String(),
			"design": s.design,
		},
	)
}

type DesignStateValueBSONUnmarshaler struct {
	Hint   string   `bson:"_hint"`
	Design bson.Raw `bson:"design"`
}

func (s *DesignStateValue) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError(utils.ErrStringDecodeBSON(*s))

	var u DesignStateValueBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e.Wrap(err)
	}

	s.BaseHinter = hint.NewBaseHinter(ht)

	var design types.Design
	if err := design.DecodeBSON(u.Design, enc); err != nil {
		return e.Wrap(err)
	}
	s.design = design

	return nil
}
