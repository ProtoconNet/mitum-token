package state

import (
	"encoding/json"
	"github.com/ProtoconNet/mitum2/util/encoder"

	"github.com/ProtoconNet/mitum-token/types"
	"github.com/ProtoconNet/mitum-token/utils"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type DesignStateValueJSONMarshaler struct {
	hint.BaseHinter
	Design types.Design `json:"design"`
}

func (s DesignStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(DesignStateValueJSONMarshaler{
		BaseHinter: s.BaseHinter,
		Design:     s.design,
	})
}

type DesignStateValueJSONUnmarshaler struct {
	Design json.RawMessage `json:"design"`
}

func (s *DesignStateValue) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError(utils.ErrStringDecodeJSON(*s))

	var u DesignStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	var design types.Design
	if err := design.DecodeJSON(u.Design, enc); err != nil {
		return e.Wrap(err)
	}
	s.design = design

	return nil
}
