package state

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
	"sync"
)

type TokenBalanceStateValueMerger struct {
	*common.BaseStateValueMerger
	existing TokenBalanceStateValue
	add      common.Big
	remove   common.Big
	sync.Mutex
}

func NewTokenBalanceStateValueMerger(height base.Height, key string, st base.State) *TokenBalanceStateValueMerger {
	nst := st
	if st == nil {
		nst = common.NewBaseState(base.NilHeight, key, nil, nil, nil)
	}

	s := &TokenBalanceStateValueMerger{
		BaseStateValueMerger: common.NewBaseStateValueMerger(height, nst.Key(), nst),
	}

	s.existing = NewTokenBalanceStateValue(common.ZeroBig)
	if nst.Value() != nil {
		s.existing = nst.Value().(TokenBalanceStateValue) //nolint:forcetypeassert //...
	}
	s.add = common.ZeroBig
	s.remove = common.ZeroBig

	return s
}

func (s *TokenBalanceStateValueMerger) Merge(value base.StateValue, ops util.Hash) error {
	s.Lock()
	defer s.Unlock()

	switch t := value.(type) {
	case AddTokenBalanceStateValue:
		s.add = s.add.Add(t.Amount)
	case DeductTokenBalanceStateValue:
		s.remove = s.remove.Add(t.Amount)
	default:
		return errors.Errorf("unsupported token balance state value, %T", value)
	}

	s.AddOperation(ops)

	return nil
}

func (s *TokenBalanceStateValueMerger) CloseValue() (base.State, error) {
	s.Lock()
	defer s.Unlock()

	newValue, err := s.closeValue()
	if err != nil {
		return nil, errors.WithMessage(err, "close TokenBalanceStateValueMerger")
	}

	s.BaseStateValueMerger.SetValue(newValue)

	return s.BaseStateValueMerger.CloseValue()
}

func (s *TokenBalanceStateValueMerger) closeValue() (base.StateValue, error) {
	existingAmount := s.existing.Amount

	if s.add.OverZero() {
		existingAmount = existingAmount.Add(s.add)
	}

	if s.remove.OverZero() {
		existingAmount = existingAmount.Sub(s.remove)
	}

	return NewTokenBalanceStateValue(
		existingAmount,
	), nil
}
