package types

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-token/utils"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type PolicyJSONMarshaler struct {
	hint.BaseHinter
	TotalSupply common.Big    `json:"total_supply"`
	ApproveList []ApproveInfo `json:"approve_list"`
}

func (p Policy) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(PolicyJSONMarshaler{
		BaseHinter:  p.BaseHinter,
		TotalSupply: p.totalSupply,
		ApproveList: p.approveList,
	})
}

type PolicyJSONUnmarshaler struct {
	Hint        hint.Hint       `json:"_hint"`
	TotalSupply string          `json:"total_supply"`
	ApproveList json.RawMessage `json:"approve_list"`
}

func (p *Policy) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError(utils.ErrStringDecodeJSON(*p))

	var u PolicyJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	return p.unmarshal(enc, u.Hint, u.TotalSupply, u.ApproveList)
}
