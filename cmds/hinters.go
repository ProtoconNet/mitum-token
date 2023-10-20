package cmds

import (
	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum-token/operation/token"
	"github.com/ProtoconNet/mitum-token/state"
	"github.com/ProtoconNet/mitum-token/types"
	"github.com/ProtoconNet/mitum2/launch"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
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

	{Hint: token.RegisterTokenHint, Instance: token.RegisterToken{}},
	{Hint: token.MintHint, Instance: token.Mint{}},
	{Hint: token.BurnHint, Instance: token.Burn{}},
	{Hint: token.ApproveHint, Instance: token.Approve{}},
	{Hint: token.TransferHint, Instance: token.Transfer{}},
	{Hint: token.TransferFromHint, Instance: token.TransferFrom{}},
}

var AddedSupportedHinters = []encoder.DecodeDetail{
	{Hint: token.RegisterTokenFactHint, Instance: token.RegisterTokenFact{}},
	{Hint: token.MintFactHint, Instance: token.MintFact{}},
	{Hint: token.BurnFactHint, Instance: token.BurnFact{}},
	{Hint: token.ApproveFactHint, Instance: token.ApproveFact{}},
	{Hint: token.TransferFactHint, Instance: token.TransferFact{}},
	{Hint: token.TransferFromFactHint, Instance: token.TransferFromFact{}},
}

func init() {
	defaultLen := len(launch.Hinters)
	currencyExtendedLen := defaultLen + len(currencycmds.AddedHinters)
	allExtendedLen := currencyExtendedLen + len(AddedHinters)

	Hinters = make([]encoder.DecodeDetail, allExtendedLen)
	copy(Hinters, launch.Hinters)
	copy(Hinters[defaultLen:currencyExtendedLen], currencycmds.AddedHinters)
	copy(Hinters[currencyExtendedLen:], AddedHinters)

	defaultSupportedLen := len(launch.SupportedProposalOperationFactHinters)
	currencySupportedExtendedLen := defaultSupportedLen + len(currencycmds.AddedSupportedHinters)
	allSupportedExtendedLen := currencySupportedExtendedLen + len(AddedSupportedHinters)

	SupportedProposalOperationFactHinters = make(
		[]encoder.DecodeDetail,
		allSupportedExtendedLen)
	copy(SupportedProposalOperationFactHinters, launch.SupportedProposalOperationFactHinters)
	copy(SupportedProposalOperationFactHinters[defaultSupportedLen:currencySupportedExtendedLen], currencycmds.AddedSupportedHinters)
	copy(SupportedProposalOperationFactHinters[currencySupportedExtendedLen:], AddedSupportedHinters)
}

func LoadHinters(enc encoder.Encoder) error {
	e := util.StringError("failed to add to encoder")

	for _, hinter := range Hinters {
		if err := enc.Add(hinter); err != nil {
			return e.Wrap(err)
		}
	}

	for _, hinter := range SupportedProposalOperationFactHinters {
		if err := enc.Add(hinter); err != nil {
			return e.Wrap(err)
		}
	}

	return nil
}
