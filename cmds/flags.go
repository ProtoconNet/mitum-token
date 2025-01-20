package cmds

import (
	"fmt"
	
	"github.com/ProtoconNet/mitum-token/types"
)

type TokenSymbolFlag struct {
	Symbol types.TokenSymbol
}

func (v *TokenSymbolFlag) UnmarshalText(b []byte) error {
	cid := types.TokenSymbol(string(b))
	if err := cid.IsValid(nil); err != nil {
		return fmt.Errorf("invalid token symbol, %q, %w", string(b), err)
	}
	v.Symbol = cid

	return nil
}

func (v *TokenSymbolFlag) String() string {
	return v.Symbol.String()
}
