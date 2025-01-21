package token

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/extras"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
	"github.com/pkg/errors"
)

var (
	TransfersFactHint = hint.MustNewHint("mitum-token-transfers-operation-fact-v0.0.1")
	TransfersHint     = hint.MustNewHint("mitum-token-transfers-operation-v0.0.1")
)

var MaxTransfersItems = 100

type TransfersFact struct {
	base.BaseFact
	sender base.Address
	items  []TransfersItem
}

func NewTransfersFact(
	token []byte,
	sender base.Address,
	items []TransfersItem,
) TransfersFact {
	fact := TransfersFact{
		BaseFact: base.NewBaseFact(TransfersFactHint, token),
		sender:   sender,
		items:    items,
	}
	fact.SetHash(fact.GenerateHash())
	return fact
}

func (fact TransfersFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	if l := len(fact.items); l < 1 {
		return common.ErrFactInvalid.Wrap(common.ErrArrayLen.Wrap(errors.Errorf("empty items for TransfersFact")))
	} else if l > int(MaxTransfersItems) {
		return common.ErrFactInvalid.Wrap(
			common.ErrArrayLen.Wrap(errors.Errorf("items over allowed, %d > %d", l, MaxTransfersItems)))
	}

	if err := util.CheckIsValiders(nil, false,
		fact.BaseFact,
		fact.sender,
	); err != nil {
		return err
	}

	founds := map[string]struct{}{}
	for _, item := range fact.items {
		if err := item.IsValid(nil); err != nil {
			return common.ErrFactInvalid.Wrap(err)
		}

		if fact.sender.Equal(item.contract) {
			return common.ErrFactInvalid.Wrap(
				common.ErrSelfTarget.Wrap(errors.Errorf("sender %v is same with contract account", fact.sender)))
		}

		if fact.sender.Equal(item.receiver) {
			return common.ErrFactInvalid.Wrap(common.ErrSelfTarget.Wrap(errors.Errorf("sender %v is same with receiver", fact.sender)))
		}

		if item.contract.Equal(item.receiver) {
			return common.ErrFactInvalid.Wrap(common.ErrSelfTarget.Wrap(errors.Errorf("receiver %v is same with contract account", item.receiver)))
		}

		if !item.amount.OverZero() {
			return common.ErrFactInvalid.Wrap(common.ErrValOOR.Wrap(errors.Errorf("transfer amount must be over zero, got %v", item.amount)))
		}

		if _, found := founds[item.contract.String()+"-"+item.receiver.String()]; found {
			return common.ErrFactInvalid.Wrap(
				common.ErrDupVal.Wrap(
					errors.Errorf(
						"receiver account %v in contract account %v", item.receiver, item.contract)))
		}

		founds[item.contract.String()+"-"+item.receiver.String()] = struct{}{}
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}
	return nil
}

func (fact TransfersFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact TransfersFact) Bytes() []byte {
	is := make([][]byte, len(fact.items))
	for i := range fact.items {
		is[i] = fact.items[i].Bytes()
	}

	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		util.ConcatBytesSlice(is...),
	)
}

func (fact TransfersFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact TransfersFact) Sender() base.Address {
	return fact.sender
}

func (fact TransfersFact) Items() []TransfersItem {
	return fact.items
}

func (fact TransfersFact) Addresses() ([]base.Address, error) {
	var as []base.Address

	for i := range fact.items {
		if ads, err := fact.items[i].Addresses(); err != nil {
			return nil, err
		} else {
			as = append(as, ads...)
		}
	}

	as = append(as, fact.Sender())

	return as, nil
}

func (fact TransfersFact) FeeBase() map[types.CurrencyID][]common.Big {
	required := make(map[types.CurrencyID][]common.Big)

	for i := range fact.items {
		zeroBig := common.ZeroBig
		cid := fact.items[i].Currency()
		var amsTemp []common.Big
		if ams, found := required[cid]; found {
			ams = append(ams, zeroBig)
			required[cid] = ams
		} else {
			amsTemp = append(amsTemp, zeroBig)
			required[cid] = amsTemp
		}
	}

	return required
}

func (fact TransfersFact) FeePayer() base.Address {
	return fact.sender
}

func (fact TransfersFact) FactUser() base.Address {
	return fact.sender
}

func (fact TransfersFact) Signer() base.Address {
	return fact.sender
}

func (fact TransfersFact) ActiveContract() []base.Address {
	var arr []base.Address
	for i := range fact.items {
		arr = append(arr, fact.items[i].contract)
	}
	return arr
}

type Transfers struct {
	extras.ExtendedOperation
}

func NewTransfers(fact TransfersFact) Transfers {
	return Transfers{
		ExtendedOperation: extras.NewExtendedOperation(TransfersHint, fact),
	}
}
