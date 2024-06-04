package token

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/test"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
)

type TestTransferProcessor struct {
	*test.BaseTestOperationProcessorNoItem[Transfer]
}

func NewTestTransferProcessor(tp *test.TestProcessor) TestTransferProcessor {
	t := test.NewBaseTestOperationProcessorNoItem[Transfer](tp)
	return TestTransferProcessor{BaseTestOperationProcessorNoItem: &t}
}

func (t *TestTransferProcessor) Create() *TestTransferProcessor {
	t.Opr, _ = NewTransferProcessor()(
		base.GenesisHeight,
		t.GetStateFunc,
		nil, nil,
	)
	return t
}

func (t *TestTransferProcessor) SetCurrency(
	cid string, am int64, addr base.Address, target []types.CurrencyID, instate bool,
) *TestTransferProcessor {
	t.BaseTestOperationProcessorNoItem.SetCurrency(cid, am, addr, target, instate)

	return t
}

func (t *TestTransferProcessor) SetAmount(
	am int64, cid types.CurrencyID, target []types.Amount,
) *TestTransferProcessor {
	t.BaseTestOperationProcessorNoItem.SetAmount(am, cid, target)

	return t
}

func (t *TestTransferProcessor) SetContractAccount(
	owner base.Address, priv string, amount int64, cid types.CurrencyID, target []test.Account, inState bool,
) *TestTransferProcessor {
	t.BaseTestOperationProcessorNoItem.SetContractAccount(owner, priv, amount, cid, target, inState)

	return t
}

func (t *TestTransferProcessor) SetAccount(
	priv string, amount int64, cid types.CurrencyID, target []test.Account, inState bool,
) *TestTransferProcessor {
	t.BaseTestOperationProcessorNoItem.SetAccount(priv, amount, cid, target, inState)

	return t
}

func (t *TestTransferProcessor) LoadOperation(fileName string,
) *TestTransferProcessor {
	t.BaseTestOperationProcessorNoItem.LoadOperation(fileName)

	return t
}

func (t *TestTransferProcessor) Print(fileName string,
) *TestTransferProcessor {
	t.BaseTestOperationProcessorNoItem.Print(fileName)

	return t
}

func (t *TestTransferProcessor) MakeOperation(
	sender base.Address, privatekey base.Privatekey, contract, receiver base.Address, amount int64, currency types.CurrencyID,
) *TestTransferProcessor {
	op := NewTransfer(
		NewTransferFact(
			[]byte("token"),
			sender,
			contract,
			currency,
			receiver,
			common.NewBig(amount),
		))
	_ = op.Sign(privatekey, t.NetworkID)
	t.Op = op

	return t
}

func (t *TestTransferProcessor) RunPreProcess() *TestTransferProcessor {
	t.BaseTestOperationProcessorNoItem.RunPreProcess()

	return t
}

func (t *TestTransferProcessor) RunProcess() *TestTransferProcessor {
	t.BaseTestOperationProcessorNoItem.RunProcess()

	return t
}

func (t *TestTransferProcessor) IsValid() *TestTransferProcessor {
	t.BaseTestOperationProcessorNoItem.IsValid()

	return t
}

func (t *TestTransferProcessor) Decode(fileName string) *TestTransferProcessor {
	t.BaseTestOperationProcessorNoItem.Decode(fileName)

	return t
}
