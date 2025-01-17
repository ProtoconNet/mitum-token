package cmds

import (
	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum-token/operation/token"
	"github.com/ProtoconNet/mitum-token/state"
	"github.com/ProtoconNet/mitum-token/types"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/pkg/errors"
)

var Hinters []encoder.DecodeDetail
var SupportedProposalOperationFactHinters []encoder.DecodeDetail

var AddedHinters = []encoder.DecodeDetail{
	// revive:disable-next-line:line-length-limit
	{Hint: types.ApproveBoxHint, Instance: types.ApproveBox{}},
	{Hint: types.ApproveInfoHint, Instance: types.ApproveInfo{}},
	{Hint: types.PolicyHint, Instance: types.Policy{}},
	{Hint: types.DesignHint, Instance: types.Design{}},

	{Hint: state.DesignStateValueHint, Instance: state.DesignStateValue{}},
	{Hint: state.TokenBalanceStateValueHint, Instance: state.TokenBalanceStateValue{}},

	{Hint: token.RegisterModelHint, Instance: token.RegisterModel{}},
	{Hint: token.MintHint, Instance: token.Mint{}},
	{Hint: token.BurnHint, Instance: token.Burn{}},
	{Hint: token.ApproveHint, Instance: token.Approve{}},
	{Hint: token.ApprovesHint, Instance: token.Approves{}},
	{Hint: token.ApprovesItemHint, Instance: token.ApprovesItem{}},
	{Hint: token.TransferHint, Instance: token.Transfer{}},
	{Hint: token.TransfersHint, Instance: token.Transfers{}},
	{Hint: token.TransfersItemHint, Instance: token.TransfersItem{}},
	{Hint: token.TransferFromHint, Instance: token.TransferFrom{}},
	{Hint: token.TransfersFromHint, Instance: token.TransfersFrom{}},
	{Hint: token.TransfersFromItemHint, Instance: token.TransfersFromItem{}},
}

var AddedSupportedHinters = []encoder.DecodeDetail{
	{Hint: token.RegisterModelFactHint, Instance: token.RegisterModelFact{}},
	{Hint: token.MintFactHint, Instance: token.MintFact{}},
	{Hint: token.BurnFactHint, Instance: token.BurnFact{}},
	{Hint: token.ApproveFactHint, Instance: token.ApproveFact{}},
	{Hint: token.ApprovesFactHint, Instance: token.ApprovesFact{}},
	{Hint: token.TransferFactHint, Instance: token.TransferFact{}},
	{Hint: token.TransfersFactHint, Instance: token.TransfersFact{}},
	{Hint: token.TransferFromFactHint, Instance: token.TransferFromFact{}},
	{Hint: token.TransfersFromFactHint, Instance: token.TransfersFromFact{}},
}

func init() {
	Hinters = append(Hinters, currencycmds.Hinters...)
	Hinters = append(Hinters, AddedHinters...)

	SupportedProposalOperationFactHinters = append(SupportedProposalOperationFactHinters, currencycmds.SupportedProposalOperationFactHinters...)
	SupportedProposalOperationFactHinters = append(SupportedProposalOperationFactHinters, AddedSupportedHinters...)
}

func LoadHinters(encs *encoder.Encoders) error {
	for i := range Hinters {
		if err := encs.AddDetail(Hinters[i]); err != nil {
			return errors.Wrap(err, "add hinter to encoder")
		}
	}

	for i := range SupportedProposalOperationFactHinters {
		if err := encs.AddDetail(SupportedProposalOperationFactHinters[i]); err != nil {
			return errors.Wrap(err, "add supported proposal operation fact hinter to encoder")
		}
	}

	return nil
}
