package types

import (
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-token/utils"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var DesignHint = hint.MustNewHint("mitum-token-design-v0.0.1")

type Design struct {
	hint.BaseHinter
	tokenID types.CurrencyID
	symbol  string
	policy  Policy
}

func NewDesign(tokenID types.CurrencyID, symbol string, policy Policy) Design {
	return Design{
		BaseHinter: hint.NewBaseHinter(DesignHint),
		tokenID:    tokenID,
		symbol:     symbol,
		policy:     policy,
	}
}

func (d Design) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf(utils.ErrStringInvalid(d))

	if err := util.CheckIsValiders(nil, false,
		d.BaseHinter,
		d.tokenID,
		d.policy,
	); err != nil {
		return e.Wrap(err)
	}

	if d.symbol == "" {
		return e.Wrap(errors.Errorf("empty symbol"))
	}

	return nil
}

func (d Design) Bytes() []byte {
	return util.ConcatBytesSlice(
		d.tokenID.Bytes(),
		[]byte(d.symbol),
		d.policy.Bytes(),
	)
}

func (d Design) TokenID() types.CurrencyID {
	return d.tokenID
}

func (d Design) Symbol() string {
	return d.symbol
}

func (d Design) Policy() Policy {
	return d.policy
}
