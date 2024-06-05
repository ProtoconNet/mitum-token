package types

import (
	"regexp"

	"github.com/ProtoconNet/mitum2/util"
)

var (
	MinLengthTokenID = 3
	MaxLengthTokenID = 10
	ReValidTokenID   = regexp.MustCompile(`^[A-Z0-9][A-Z0-9_\.\!\$\*\@]*[A-Z0-9]$`)
	ReSpcecialChar   = regexp.MustCompile(`^[^\s:/?#\[\]@]*$`)
)

type TokenID string

func (cid TokenID) Bytes() []byte {
	return []byte(cid)
}

func (cid TokenID) String() string {
	return string(cid)
}

func (cid TokenID) IsValid([]byte) error {
	if l := len(cid); l < MinLengthTokenID || l > MaxLengthTokenID {
		return util.ErrInvalid.Errorf(
			"invalid length of token id, %d <= %d <= %d", MinLengthTokenID, l, MaxLengthTokenID)
	} else if !ReValidTokenID.Match([]byte(cid)) {
		return util.ErrInvalid.Errorf("wrong token id, %v", cid)
	}

	return nil
}
