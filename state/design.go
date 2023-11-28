package state

import (
	"fmt"

	"github.com/ProtoconNet/mitum-token/types"
	"github.com/ProtoconNet/mitum-token/utils"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var (
	DesignStateValueHint = hint.MustNewHint("mitum-token-design-state-value-v0.0.1")
	DesignSuffix         = "design"
)

type DesignStateValue struct {
	hint.BaseHinter
	design types.Design
}

func NewDesignStateValue(design types.Design) DesignStateValue {
	return DesignStateValue{
		BaseHinter: hint.NewBaseHinter(DesignStateValueHint),
		design:     design,
	}
}

func (s DesignStateValue) Hint() hint.Hint {
	return s.BaseHinter.Hint()
}

func (s DesignStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf(utils.ErrStringInvalid(s))

	if err := s.BaseHinter.IsValid(DesignStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if err := s.design.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (s DesignStateValue) HashBytes() []byte {
	return s.design.Bytes()
}

func StateDesignValue(st base.State) (*types.Design, error) {
	e := util.ErrNotFound.Errorf(ErrStringStateNotFound(st.Key()))

	v := st.Value()
	if v == nil {
		return nil, e.Wrap(errors.Errorf("nil value"))
	}

	s, ok := v.(DesignStateValue)
	if !ok {
		return nil, e.Wrap(errors.Errorf(utils.ErrStringTypeCast(DesignStateValue{}, v)))
	}

	return &s.design, nil
}

func StateKeyDesign(contract base.Address) string {
	return fmt.Sprintf("%s:%s", StateKeyTokenPrefix(contract), DesignSuffix)
}
