package state

import (
	"fmt"
	"github.com/pkg/errors"
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

func ParseStateKey(key string, Prefix string, expected int) ([]string, error) {
	parsedKey := strings.Split(key, ":")
	if parsedKey[0] != Prefix[:len(Prefix)-1] {
		return nil, errors.Errorf("State Key not include Prefix, %s", parsedKey)
	}
	if len(parsedKey) < expected {
		return nil, errors.Errorf("parsed State Key length under %v", expected)
	} else {
		return parsedKey, nil
	}
}

func IsStateDesignKey(key string) bool {
	return strings.HasPrefix(key, TokenPrefix) && strings.HasSuffix(key, DesignSuffix)
}

func IsStateTokenBalanceKey(key string) bool {
	return strings.HasPrefix(key, TokenPrefix) && strings.HasSuffix(key, TokenBalanceSuffix)
}
