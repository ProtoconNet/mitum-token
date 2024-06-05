package state

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-token/utils"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type TokenBalanceStateValueJSONMarshaler struct {
	hint.BaseHinter
	Amount common.Big `json:"amount"`
}

func (s TokenBalanceStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(TokenBalanceStateValueJSONMarshaler{
		BaseHinter: s.BaseHinter,
		Amount:     s.Amount,
	})
}

type TokenBalanceStateValueJSONUnmarshaler struct {
	Amount string `json:"amount"`
}

func (s *TokenBalanceStateValue) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError(utils.ErrStringDecodeJSON(*s))

	var u TokenBalanceStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	big, err := common.NewBigFromString(u.Amount)
	if err != nil {
		return e.Wrap(err)
	}
	s.Amount = big

	return nil
}
