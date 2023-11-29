package state

import (
	"fmt"
	"strings"

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

func IsStateDesignKey(key string) bool {
	return strings.HasPrefix(key, TokenPrefix) && strings.HasSuffix(key, DesignSuffix)
}

func IsStateTokenBalanceKey(key string) bool {
	return strings.HasPrefix(key, TokenPrefix) && strings.HasSuffix(key, TokenBalanceSuffix)
}
