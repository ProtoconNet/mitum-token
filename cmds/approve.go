package cmds

import (
	"context"

	ccmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum-token/operation/token"
	"github.com/ProtoconNet/mitum-token/utils"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

type ApproveCommand struct {
	OperationCommand
	Approved ccmds.AddressFlag `arg:"" name:"approved" help:"approved account" required:"true"`
	Amount   ccmds.BigFlag     `arg:"" name:"amount" help:"amount to approve" required:"true"`
	approved base.Address
}

func (cmd *ApproveCommand) Run(pctx context.Context) error { // nolint:dupl
	if _, err := cmd.prepare(pctx); err != nil {
		return err
	}

	if err := cmd.parseFlags(); err != nil {
		return err
	}

	op, err := cmd.createOperation()
	if err != nil {
		return err
	}

	ccmds.PrettyPrint(cmd.Out, op)

	return nil
}

func (cmd *ApproveCommand) parseFlags() error {
	if err := cmd.OperationCommand.parseFlags(); err != nil {
		return err
	}

	approved, err := cmd.Approved.Encode(cmd.Encoders.JSON())
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
		cmd.Currency.CID,
		cmd.approved,
		cmd.Amount.Big,
	)

	op := token.NewApprove(fact)
	if err := op.Sign(cmd.Privatekey, cmd.NetworkID.NetworkID()); err != nil {
		return nil, e.Wrap(err)
	}

	return op, nil
}
