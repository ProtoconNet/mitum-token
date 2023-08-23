package state

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum-token/utils"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

func (s TokenBalanceStateValue) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":  s.Hint().String(),
			"amount": s.amount,
		},
	)
}

type TokenBalanceStateValueBSONUnmarshaler struct {
	Hint   string `bson:"_hint"`
	Amount string `bson:"amount"`
}

func (s *TokenBalanceStateValue) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError(utils.ErrStringDecodeBSON(*s))

	var u TokenBalanceStateValueBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e.Wrap(err)
	}

	s.BaseHinter = hint.NewBaseHinter(ht)

	big, err := common.NewBigFromString(u.Amount)
	if err != nil {
		return e.Wrap(err)
	}
	s.amount = big

	return nil
}
