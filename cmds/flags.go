package cmds

import (
	"fmt"
	"github.com/ProtoconNet/mitum-token/types"
)

type TokenIDFlag struct {
	CID types.TokenID
}

func (v *TokenIDFlag) UnmarshalText(b []byte) error {
	cid := types.TokenID(string(b))
	if err := cid.IsValid(nil); err != nil {
		return fmt.Errorf("invalid token id, %q, %w", string(b), err)
	}
	v.CID = cid

	return nil
}

func (v *TokenIDFlag) String() string {
	return v.CID.String()
}
