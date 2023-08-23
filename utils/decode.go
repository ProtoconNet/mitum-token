package utils

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum2/util"
)

func DecodeMap(b []byte) (map[string]interface{}, error) {
	var m map[string]json.RawMessage
	if err := util.UnmarshalJSON(b, &m); err != nil {
		return nil, err
	}

	if len(m) < 1 {
		return nil, nil
	}

	r := map[string]interface{}{}
	for k, v := range m {
		var s interface{}
		if err := json.Unmarshal(v, &s); err != nil {
			return nil, err
		}
		r[k] = s
	}

	return r, nil
}
