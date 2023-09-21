package state

import (
	"fmt"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-token/utils"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var (
	TokenBalanceStateValueHint = hint.MustNewHint("mitum-token-balance-state-value-v0.0.1")
	TokenBalanceSuffix         = ":tokenbalance"
)

type TokenBalanceStateValue struct {
	hint.BaseHinter
	amount common.Big
}

func NewTokenBalanceStateValue(amount common.Big) TokenBalanceStateValue {
	return TokenBalanceStateValue{
		BaseHinter: hint.NewBaseHinter(TokenBalanceStateValueHint),
		amount:     amount,
	}
}

func (s TokenBalanceStateValue) Hint() hint.Hint {
	return s.BaseHinter.Hint()
}

func (s TokenBalanceStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf(utils.ErrStringInvalid(s))

	if err := s.BaseHinter.IsValid(TokenBalanceStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if !s.amount.OverNil() {
		return e.Wrap(errors.Errorf("nil big"))
	}

	return nil
}

func (s TokenBalanceStateValue) HashBytes() []byte {
	return s.amount.Bytes()
}

func StateTokenBalanceValue(st base.State) (common.Big, error) {
	e := util.ErrNotFound.Errorf(ErrStringStateNotFound(st.Key()))

	v := st.Value()
	if v == nil {
		return common.NilBig, e.Wrap(errors.Errorf("nil value"))
	}

	s, ok := v.(TokenBalanceStateValue)
	if !ok {
		return common.NilBig, e.Wrap(errors.Errorf(utils.ErrStringTypeCast(TokenBalanceStateValue{}, v)))
	}

	return s.amount, nil
}

func StateKeyTokenBalance(contract base.Address, address base.Address) string {
	return fmt.Sprintf("%s:%s:%s", StateKeyTokenPrefix(contract), address, TokenBalanceSuffix)
}
