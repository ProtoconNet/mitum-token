package state

import (
	"fmt"

	"github.com/ProtoconNet/mitum2/base"
)

var TokenPrefix = "token"

func StateKeyTokenPrefix(contract base.Address) string {
	return fmt.Sprintf("%s:%s", TokenPrefix, contract)
}

type StateKeyGenerator struct {
	contract base.Address
}

func NewStateKeyGenerator(contract base.Address) StateKeyGenerator {
	return StateKeyGenerator{
		contract,
	}
}

func (g StateKeyGenerator) Design() string {
	return StateKeyDesign(g.contract)
}

func (g StateKeyGenerator) TokenBalance(address base.Address) string {
	return StateKeyTokenBalance(g.contract, address)
}
