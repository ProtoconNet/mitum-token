package types

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-token/utils"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var DesignHint = hint.MustNewHint("mitum-token-design-v0.0.1")

type Design struct {
	hint.BaseHinter
	symbol  TokenSymbol
	name    string
	decimal common.Big
	policy  Policy
}

func NewDesign(symbol TokenSymbol, name string, decimal common.Big, policy Policy) Design {
	return Design{
		BaseHinter: hint.NewBaseHinter(DesignHint),
		symbol:     symbol,
		name:       name,
		decimal:    decimal,
		policy:     policy,
	}
}

func (d Design) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf(utils.ErrStringInvalid(d))

	if err := util.CheckIsValiders(nil, false,
		d.BaseHinter,
		d.symbol,
		d.policy,
	); err != nil {
		return e.Wrap(err)
	}

	if d.name == "" {
		return e.Wrap(errors.Errorf("empty symbol"))
	}
	if !d.decimal.OverNil() {
		return e.Wrap(errors.Errorf("decimal must be bigger than or equal to zero"))
	}

	return nil
}

func (d Design) Bytes() []byte {
	return util.ConcatBytesSlice(
		d.symbol.Bytes(),
		[]byte(d.name),
		d.decimal.Bytes(),
		d.policy.Bytes(),
	)
}

func (d Design) Symbol() TokenSymbol {
	return d.symbol
}

func (d Design) Name() string {
	return d.name
}

func (d Design) Decimal() common.Big {
	return d.decimal
}

func (d Design) Policy() Policy {
	return d.policy
}
