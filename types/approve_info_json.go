package types

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-token/utils"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type ApproveInfoJSONMarshaler struct {
	hint.BaseHinter
	Account  base.Address          `json:"account"`
	Approved map[string]common.Big `json:"approved"`
}

func (a ApproveInfo) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(ApproveInfoJSONMarshaler{
		BaseHinter: a.BaseHinter,
		Account:    a.account,
		Approved:   a.approved,
	})
}

type ApproveInfoJSONUnmarshaler struct {
	Hint     hint.Hint       `json:"_hint"`
	Account  string          `json:"account"`
	Approved json.RawMessage `json:"approved"`
}

func (a *ApproveInfo) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError(utils.ErrStringDecodeJSON(*a))

	var u ApproveInfoJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	return a.unmarshal(enc, u.Hint, u.Account, u.Approved)
}
