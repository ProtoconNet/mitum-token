package token

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/test"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	tokentypes "github.com/ProtoconNet/mitum-token/types"
	"github.com/ProtoconNet/mitum2/base"
)

type TestRegisterTokenProcessor struct {
	*test.BaseTestOperationProcessorNoItem[RegisterModel]
}

func NewTestRegisterTokenProcessor(tp *test.TestProcessor) TestRegisterTokenProcessor {
	t := test.NewBaseTestOperationProcessorNoItem[RegisterModel](tp)
	return TestRegisterTokenProcessor{BaseTestOperationProcessorNoItem: &t}
}

func (t *TestRegisterTokenProcessor) Create() *TestRegisterTokenProcessor {
	t.Opr, _ = NewRegisterModelProcessor()(
		base.GenesisHeight,
		t.GetStateFunc,
		nil, nil,
	)
	return t
}

func (t *TestRegisterTokenProcessor) SetCurrency(
	cid string, am int64, addr base.Address, target []types.CurrencyID, instate bool,
) *TestRegisterTokenProcessor {
	t.BaseTestOperationProcessorNoItem.SetCurrency(cid, am, addr, target, instate)

	return t
}

func (t *TestRegisterTokenProcessor) SetAmount(
	am int64, cid types.CurrencyID, target []types.Amount,
) *TestRegisterTokenProcessor {
	t.BaseTestOperationProcessorNoItem.SetAmount(am, cid, target)

	return t
}

func (t *TestRegisterTokenProcessor) SetContractAccount(
	owner base.Address, priv string, amount int64, cid types.CurrencyID, target []test.Account, inState bool,
) *TestRegisterTokenProcessor {
	t.BaseTestOperationProcessorNoItem.SetContractAccount(owner, priv, amount, cid, target, inState)

	return t
}

func (t *TestRegisterTokenProcessor) SetAccount(
	priv string, amount int64, cid types.CurrencyID, target []test.Account, inState bool,
) *TestRegisterTokenProcessor {
	t.BaseTestOperationProcessorNoItem.SetAccount(priv, amount, cid, target, inState)

	return t
}

func (t *TestRegisterTokenProcessor) LoadOperation(fileName string,
) *TestRegisterTokenProcessor {
	t.BaseTestOperationProcessorNoItem.LoadOperation(fileName)

	return t
}

func (t *TestRegisterTokenProcessor) Print(fileName string,
) *TestRegisterTokenProcessor {
	t.BaseTestOperationProcessorNoItem.Print(fileName)

	return t
}

func (t *TestRegisterTokenProcessor) MakeOperation(
	sender base.Address, privatekey base.Privatekey, contract base.Address,
	symbol, name string, initialSupply int64, currency types.CurrencyID,
) *TestRegisterTokenProcessor {
	op := NewRegisterModel(
		NewRegisterModelFact(
			[]byte("token"),
			sender,
			contract,
			currency,
			tokentypes.TokenSymbol(symbol),
			name,
			common.NewBig(initialSupply),
		))
	_ = op.Sign(privatekey, t.NetworkID)
	t.Op = op

	return t
}

func (t *TestRegisterTokenProcessor) RunPreProcess() *TestRegisterTokenProcessor {
	t.BaseTestOperationProcessorNoItem.RunPreProcess()

	return t
}

func (t *TestRegisterTokenProcessor) RunProcess() *TestRegisterTokenProcessor {
	t.BaseTestOperationProcessorNoItem.RunProcess()

	return t
}

func (t *TestRegisterTokenProcessor) IsValid() *TestRegisterTokenProcessor {
	t.BaseTestOperationProcessorNoItem.IsValid()

	return t
}

func (t *TestRegisterTokenProcessor) Decode(fileName string) *TestRegisterTokenProcessor {
	t.BaseTestOperationProcessorNoItem.Decode(fileName)

	return t
}
