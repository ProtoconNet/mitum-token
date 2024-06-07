package types

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/pkg/errors"
	"regexp"
)

var (
	MinLengthTokenSymbol = 3
	MaxLengthTokenSymbol = 10
	ReValidTokenSymbol   = regexp.MustCompile(`^[A-Z0-9][A-Z0-9_\.\!\$\*\@]*[A-Z0-9]$`)
	ReSpcecialChar       = regexp.MustCompile(`^[^\s:/?#\[\]@]*$`)
)

type TokenSymbol string

func (ts TokenSymbol) Bytes() []byte {
	return []byte(ts)
}

func (ts TokenSymbol) String() string {
	return string(ts)
}

func (ts TokenSymbol) IsValid([]byte) error {
	if l := len(ts); l < MinLengthTokenSymbol || l > MaxLengthTokenSymbol {
		return common.ErrValOOR.Wrap(errors.Errorf(
			"invalid length of token symbol, %d <= %d <= %d", MinLengthTokenSymbol, l, MaxLengthTokenSymbol))
	} else if !ReValidTokenSymbol.Match([]byte(ts)) {
		return common.ErrValueInvalid.Wrap(errors.Errorf("wrong token symbol, %v", ts))
	}

	return nil
}
