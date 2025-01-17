package state

import (
	"fmt"
	"strings"
)

var TokenPrefix = "token"

func StateKeyTokenPrefix(contract string) string {
	return fmt.Sprintf("%s:%s", TokenPrefix, contract)
}

type StateKeyGenerator struct {
	contract string
}

func NewStateKeyGenerator(contract string) StateKeyGenerator {
	return StateKeyGenerator{
		contract,
	}
}

func (g StateKeyGenerator) Design() string {
	return StateKeyDesign(g.contract)
}

func (g StateKeyGenerator) TokenBalance(address string) string {
	return StateKeyTokenBalance(g.contract, address)
}

func IsStateDesignKey(key string) bool {
	return strings.HasPrefix(key, TokenPrefix) && strings.HasSuffix(key, DesignSuffix)
}

func IsStateTokenBalanceKey(key string) bool {
	return strings.HasPrefix(key, TokenPrefix) && strings.HasSuffix(key, TokenBalanceSuffix)
}
