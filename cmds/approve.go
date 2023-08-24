package cmds

import (
	"context"

	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum-token/operation/token"
	"github.com/ProtoconNet/mitum-token/utils"
	"github.com/pkg/errors"

	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
)

type ApproveCommand struct {
	OperationCommand
	Approved currencycmds.AddressFlag `arg:"" name:"approved" help:"approved account" required:"true"`
	Amount   currencycmds.BigFlag     `arg:"" name:"amount" help:"amount to approve" required:"true"`
	approved base.Address
}

func NewApproveCommand() ApproveCommand {
	cmd := NewOperationCommand()
	return ApproveCommand{
		OperationCommand: *cmd,
	}
}

func (cmd *ApproveCommand) Run(pctx context.Context) error { // nolint:dupl
	if _, err := cmd.prepare(pctx); err != nil {
		return err
	}

	encs = cmd.Encoders
	enc = cmd.Encoder

	if err := cmd.parseFlags(); err != nil {
		return err
	}

	op, err := cmd.createOperation()
	if err != nil {
		return err
	}

	currencycmds.PrettyPrint(cmd.Out, op)

	return nil
}

func (cmd *ApproveCommand) parseFlags() error {
	if err := cmd.OperationCommand.parseFlags(); err != nil {
		return err
	}

	approved, err := cmd.Approved.Encode(enc)
	if err != nil {
		return errors.Wrapf(err, "invalid approved format, %q", cmd.Approved.String())
	}
	cmd.approved = approved

	return nil
}

func (cmd *ApproveCommand) createOperation() (base.Operation, error) { // nolint:dupl}
	e := util.StringError(utils.ErrStringCreate("approve operation"))

	fact := token.NewApproveFact(
		[]byte(cmd.Token),
		cmd.sender, cmd.contract,
		cmd.TokenID.CID, cmd.Currency.CID,
		cmd.approved,
		cmd.Amount.Big,
	)

	op := token.NewApprove(fact)
	if err := op.HashSign(cmd.Privatekey, cmd.NetworkID.NetworkID()); err != nil {
		return nil, e.Wrap(err)
	}

	return op, nil
}
