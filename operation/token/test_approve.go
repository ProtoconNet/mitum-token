package token

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/test"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type TestApproveProcessor struct {
	*test.BaseTestOperationProcessorNoItem[Approve]
}

func NewTestApproveProcessor(encs *encoder.Encoders) TestApproveProcessor {
	t := test.NewBaseTestOperationProcessorNoItem[Approve](encs)
	return TestApproveProcessor{BaseTestOperationProcessorNoItem: &t}
}

func (t *TestApproveProcessor) Create() *TestApproveProcessor {
	t.Opr, _ = NewApproveProcessor()(
		base.GenesisHeight,
		t.GetStateFunc,
		nil, nil,
	)
	return t
}

func (t *TestApproveProcessor) SetCurrency(
	cid string, am int64, addr base.Address, target []types.CurrencyID, instate bool,
) *TestApproveProcessor {
	t.BaseTestOperationProcessorNoItem.SetCurrency(cid, am, addr, target, instate)

	return t
}

func (t *TestApproveProcessor) SetAmount(
	am int64, cid types.CurrencyID, target []types.Amount,
) *TestApproveProcessor {
	t.BaseTestOperationProcessorNoItem.SetAmount(am, cid, target)

	return t
}

func (t *TestApproveProcessor) SetContractAccount(
	owner base.Address, priv string, amount int64, cid types.CurrencyID, target []test.Account, inState bool,
) *TestApproveProcessor {
	t.BaseTestOperationProcessorNoItem.SetContractAccount(owner, priv, amount, cid, target, inState)

	return t
}

func (t *TestApproveProcessor) SetAccount(
	priv string, amount int64, cid types.CurrencyID, target []test.Account, inState bool,
) *TestApproveProcessor {
	t.BaseTestOperationProcessorNoItem.SetAccount(priv, amount, cid, target, inState)

	return t
}

func (t *TestApproveProcessor) LoadOperation(fileName string,
) *TestApproveProcessor {
	t.BaseTestOperationProcessorNoItem.LoadOperation(fileName)

	return t
}

func (t *TestApproveProcessor) Print(fileName string,
) *TestApproveProcessor {
	t.BaseTestOperationProcessorNoItem.Print(fileName)

	return t
}

func (t *TestApproveProcessor) MakeOperation(
	sender base.Address, privatekey base.Privatekey, contract, approved base.Address, amount common.Big, currency types.CurrencyID,
) *TestApproveProcessor {
	op := NewApprove(
		NewApproveFact(
			[]byte("token"),
			sender,
			contract,
			currency,
			approved,
			amount,
		))
	_ = op.Sign(privatekey, t.NetworkID)
	t.Op = op

	return t
}

func (t *TestApproveProcessor) RunPreProcess() *TestApproveProcessor {
	t.BaseTestOperationProcessorNoItem.RunPreProcess()

	return t
}

func (t *TestApproveProcessor) RunProcess() *TestApproveProcessor {
	t.BaseTestOperationProcessorNoItem.RunProcess()

	return t
}

func (t *TestApproveProcessor) IsValid() *TestApproveProcessor {
	t.BaseTestOperationProcessorNoItem.IsValid()

	return t
}

func (t *TestApproveProcessor) Decode(fileName string) *TestApproveProcessor {
	t.BaseTestOperationProcessorNoItem.Decode(fileName)

	return t
}
