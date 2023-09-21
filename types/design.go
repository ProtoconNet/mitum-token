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
	symbol types.CurrencyID
	name   string
	policy Policy
}

func NewDesign(symbol types.CurrencyID, name string, policy Policy) Design {
	return Design{
		BaseHinter: hint.NewBaseHinter(DesignHint),
		symbol:     symbol,
		name:       name,
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

	return nil
}

func (d Design) Bytes() []byte {
	return util.ConcatBytesSlice(
		d.symbol.Bytes(),
		[]byte(d.name),
		d.policy.Bytes(),
	)
}

func (d Design) Symbol() types.CurrencyID {
	return d.symbol
}

func (d Design) Name() string {
	return d.name
}

func (d Design) Policy() Policy {
	return d.policy
}
