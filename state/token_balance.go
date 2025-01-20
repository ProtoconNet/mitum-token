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
	TokenBalanceSuffix         = "tokenbalance"
)

type TokenBalanceStateValue struct {
	hint.BaseHinter
	Amount common.Big
}

func NewTokenBalanceStateValue(amount common.Big) TokenBalanceStateValue {
	return TokenBalanceStateValue{
		BaseHinter: hint.NewBaseHinter(TokenBalanceStateValueHint),
		Amount:     amount,
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

	if !s.Amount.OverNil() {
		return e.Wrap(errors.Errorf("nil big"))
	}

	return nil
}

func (s TokenBalanceStateValue) HashBytes() []byte {
	return s.Amount.Bytes()
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

	return s.Amount, nil
}

type AddTokenBalanceStateValue struct {
	Amount common.Big
}

func NewAddTokenBalanceStateValue(amount common.Big) AddTokenBalanceStateValue {
	return AddTokenBalanceStateValue{
		Amount: amount,
	}
}

func (b AddTokenBalanceStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid AddTokenBalanceStateValue")

	if err := util.CheckIsValiders(nil, false, b.Amount); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (b AddTokenBalanceStateValue) HashBytes() []byte {
	return b.Amount.Bytes()
}

type DeductTokenBalanceStateValue struct {
	Amount common.Big
}

func NewDeductTokenBalanceStateValue(amount common.Big) DeductTokenBalanceStateValue {
	return DeductTokenBalanceStateValue{
		Amount: amount,
	}
}

func (b DeductTokenBalanceStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid DeductTokenBalanceStateValue")

	if err := util.CheckIsValiders(nil, false, b.Amount); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (b DeductTokenBalanceStateValue) HashBytes() []byte {
	return b.Amount.Bytes()
}

func StateKeyTokenBalance(contract string, address string) string {
	return fmt.Sprintf("%s:%s:%s", StateKeyTokenPrefix(contract), address, TokenBalanceSuffix)
}
