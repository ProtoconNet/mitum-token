package types

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-token/utils"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type ApproveBoxJSONMarshaler struct {
	hint.BaseHinter
	Account  base.Address  `json:"account"`
	Approved []ApproveInfo `json:"approved"`
}

func (a ApproveBox) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(ApproveBoxJSONMarshaler{
		BaseHinter: a.BaseHinter,
		Account:    a.account,
		Approved:   a.approved,
	})
}

type ApproveBoxJSONUnmarshaler struct {
	Hint     hint.Hint       `json:"_hint"`
	Account  string          `json:"account"`
	Approved json.RawMessage `json:"approved"`
}

func (a *ApproveBox) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError(utils.ErrStringDecodeJSON(*a))

	var u ApproveBoxJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	return a.unpack(enc, u.Hint, u.Account, u.Approved)
}
