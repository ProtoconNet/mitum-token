package state

import (
	"fmt"

	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
)

var TokenPrefix = "token"

func StateKeyTokenPrefix(contract base.Address, tokenID types.CurrencyID) string {
	return fmt.Sprintf("%s:%s:%s", TokenPrefix, contract, tokenID)
}

type StateKeyGenerator struct {
	contract base.Address
	tokenID  types.CurrencyID
}

func NewStateKeyGenerator(contract base.Address, tokenID types.CurrencyID) StateKeyGenerator {
	return StateKeyGenerator{
		contract,
		tokenID,
	}
}

func (g StateKeyGenerator) Design() string {
	return StateKeyDesign(g.contract, g.tokenID)
}

func (g StateKeyGenerator) TokenBalance(address base.Address) string {
	return StateKeyTokenBalance(g.contract, g.tokenID, address)
}
