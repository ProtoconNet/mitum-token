package token

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/test"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
)

type TestTransferFromProcessor struct {
	*test.BaseTestOperationProcessorNoItem[TransferFrom]
}

func NewTestTransferFromProcessor(tp *test.TestProcessor) TestTransferFromProcessor {
	t := test.NewBaseTestOperationProcessorNoItem[TransferFrom](tp)
	return TestTransferFromProcessor{BaseTestOperationProcessorNoItem: &t}
}

func (t *TestTransferFromProcessor) Create() *TestTransferFromProcessor {
	t.Opr, _ = NewTransferFromProcessor()(
		base.GenesisHeight,
		t.GetStateFunc,
		nil, nil,
	)
	return t
}

func (t *TestTransferFromProcessor) SetCurrency(
	cid string, am int64, addr base.Address, target []types.CurrencyID, instate bool,
) *TestTransferFromProcessor {
	t.BaseTestOperationProcessorNoItem.SetCurrency(cid, am, addr, target, instate)

	return t
}

func (t *TestTransferFromProcessor) SetAmount(
	am int64, cid types.CurrencyID, target []types.Amount,
) *TestTransferFromProcessor {
	t.BaseTestOperationProcessorNoItem.SetAmount(am, cid, target)

	return t
}

func (t *TestTransferFromProcessor) SetContractAccount(
	owner base.Address, priv string, amount int64, cid types.CurrencyID, target []test.Account, inState bool,
) *TestTransferFromProcessor {
	t.BaseTestOperationProcessorNoItem.SetContractAccount(owner, priv, amount, cid, target, inState)

	return t
}

func (t *TestTransferFromProcessor) SetAccount(
	priv string, amount int64, cid types.CurrencyID, target []test.Account, inState bool,
) *TestTransferFromProcessor {
	t.BaseTestOperationProcessorNoItem.SetAccount(priv, amount, cid, target, inState)

	return t
}

func (t *TestTransferFromProcessor) LoadOperation(fileName string,
) *TestTransferFromProcessor {
	t.BaseTestOperationProcessorNoItem.LoadOperation(fileName)

	return t
}

func (t *TestTransferFromProcessor) Print(fileName string,
) *TestTransferFromProcessor {
	t.BaseTestOperationProcessorNoItem.Print(fileName)

	return t
}

func (t *TestTransferFromProcessor) MakeOperation(
	sender base.Address, privatekey base.Privatekey, contract, receiver, target base.Address, amount int64, currency types.CurrencyID,
) *TestTransferFromProcessor {
	op := NewTransferFrom(
		NewTransferFromFact(
			[]byte("token"),
			sender,
			contract,
			currency,
			receiver,
			target,
			common.NewBig(amount),
		))
	_ = op.Sign(privatekey, t.NetworkID)
	t.Op = op

	return t
}

func (t *TestTransferFromProcessor) RunPreProcess() *TestTransferFromProcessor {
	t.BaseTestOperationProcessorNoItem.RunPreProcess()

	return t
}

func (t *TestTransferFromProcessor) RunProcess() *TestTransferFromProcessor {
	t.BaseTestOperationProcessorNoItem.RunProcess()

	return t
}

func (t *TestTransferFromProcessor) IsValid() *TestTransferFromProcessor {
	t.BaseTestOperationProcessorNoItem.IsValid()

	return t
}

func (t *TestTransferFromProcessor) Decode(fileName string) *TestTransferFromProcessor {
	t.BaseTestOperationProcessorNoItem.Decode(fileName)

	return t
}
