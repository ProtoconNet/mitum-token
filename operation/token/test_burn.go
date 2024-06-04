package token

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/test"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
)

type TestBurnProcessor struct {
	*test.BaseTestOperationProcessorNoItem[Burn]
}

func NewTestBurnProcessor(tp *test.TestProcessor) TestBurnProcessor {
	t := test.NewBaseTestOperationProcessorNoItem[Burn](tp)
	return TestBurnProcessor{BaseTestOperationProcessorNoItem: &t}
}

func (t *TestBurnProcessor) Create() *TestBurnProcessor {
	t.Opr, _ = NewBurnProcessor()(
		base.GenesisHeight,
		t.GetStateFunc,
		nil, nil,
	)
	return t
}

func (t *TestBurnProcessor) SetCurrency(
	cid string, am int64, addr base.Address, target []types.CurrencyID, instate bool,
) *TestBurnProcessor {
	t.BaseTestOperationProcessorNoItem.SetCurrency(cid, am, addr, target, instate)

	return t
}

func (t *TestBurnProcessor) SetAmount(
	am int64, cid types.CurrencyID, target []types.Amount,
) *TestBurnProcessor {
	t.BaseTestOperationProcessorNoItem.SetAmount(am, cid, target)

	return t
}

func (t *TestBurnProcessor) SetContractAccount(
	owner base.Address, priv string, amount int64, cid types.CurrencyID, target []test.Account, inState bool,
) *TestBurnProcessor {
	t.BaseTestOperationProcessorNoItem.SetContractAccount(owner, priv, amount, cid, target, inState)

	return t
}

func (t *TestBurnProcessor) SetAccount(
	priv string, amount int64, cid types.CurrencyID, target []test.Account, inState bool,
) *TestBurnProcessor {
	t.BaseTestOperationProcessorNoItem.SetAccount(priv, amount, cid, target, inState)

	return t
}

func (t *TestBurnProcessor) LoadOperation(fileName string,
) *TestBurnProcessor {
	t.BaseTestOperationProcessorNoItem.LoadOperation(fileName)

	return t
}

func (t *TestBurnProcessor) Print(fileName string,
) *TestBurnProcessor {
	t.BaseTestOperationProcessorNoItem.Print(fileName)

	return t
}

func (t *TestBurnProcessor) MakeOperation(
	sender base.Address, privatekey base.Privatekey, contract, target base.Address, amount int64, currency types.CurrencyID,
) *TestBurnProcessor {
	op := NewBurn(
		NewBurnFact(
			[]byte("token"),
			sender,
			contract,
			currency,
			target,
			common.NewBig(amount),
		))
	_ = op.Sign(privatekey, t.NetworkID)
	t.Op = op

	return t
}

func (t *TestBurnProcessor) RunPreProcess() *TestBurnProcessor {
	t.BaseTestOperationProcessorNoItem.RunPreProcess()

	return t
}

func (t *TestBurnProcessor) RunProcess() *TestBurnProcessor {
	t.BaseTestOperationProcessorNoItem.RunProcess()

	return t
}

func (t *TestBurnProcessor) IsValid() *TestBurnProcessor {
	t.BaseTestOperationProcessorNoItem.IsValid()

	return t
}

func (t *TestBurnProcessor) Decode(fileName string) *TestBurnProcessor {
	t.BaseTestOperationProcessorNoItem.Decode(fileName)

	return t
}
