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
	ApprovesFactHint = hint.MustNewHint("mitum-token-approves-operation-fact-v0.0.1")
	ApprovesHint     = hint.MustNewHint("mitum-token-approves-operation-v0.0.1")
)

var MaxApprovesItems = 100

type ApprovesFact struct {
	base.BaseFact
	sender base.Address
	items  []ApprovesItem
}

func NewApprovesFact(
	token []byte,
	sender base.Address,
	items []ApprovesItem,
) ApprovesFact {
	fact := ApprovesFact{
		BaseFact: base.NewBaseFact(ApprovesFactHint, token),
		sender:   sender,
		items:    items,
	}
	fact.SetHash(fact.GenerateHash())
	return fact
}

func (fact ApprovesFact) IsValid(b []byte) error {
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

		if _, found := founds[item.contract.String()+"-"+item.approved.String()]; found {
			return common.ErrFactInvalid.Wrap(
				common.ErrDupVal.Wrap(errors.Errorf("contract account %v", item.contract)))
		}

		founds[item.Contract().String()] = struct{}{}
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}
	return nil
}

func (fact ApprovesFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact ApprovesFact) Bytes() []byte {
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

func (fact ApprovesFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact ApprovesFact) Sender() base.Address {
	return fact.sender
}

func (fact ApprovesFact) Items() []ApprovesItem {
	return fact.items
}

func (fact ApprovesFact) Addresses() ([]base.Address, error) {
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

func (fact ApprovesFact) FeeBase() map[types.CurrencyID][]common.Big {
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

func (fact ApprovesFact) FeePayer() base.Address {
	return fact.sender
}

func (fact ApprovesFact) FactUser() base.Address {
	return fact.sender
}

func (fact ApprovesFact) Signer() base.Address {
	return fact.sender
}

func (fact ApprovesFact) ActiveContract() []base.Address {
	var arr []base.Address
	for i := range fact.items {
		arr = append(arr, fact.items[i].contract)
	}
	return arr
}

type Approves struct {
	extras.ExtendedOperation
}

func NewApproves(fact ApprovesFact) Approves {
	return Approves{
		ExtendedOperation: extras.NewExtendedOperation(ApprovesHint, fact),
	}
}
