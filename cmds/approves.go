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

type ApprovesCommand struct {
	OperationCommand
	Approved currencycmds.AddressFlag `arg:"" name:"approved" help:"approved account" required:"true"`
	Amount   currencycmds.BigFlag     `arg:"" name:"amount" help:"amount to approve" required:"true"`
	approved base.Address
}

func (cmd *ApprovesCommand) Run(pctx context.Context) error { // nolint:dupl
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

	currencycmds.PrettyPrint(cmd.Out, op)

	return nil
}

func (cmd *ApprovesCommand) parseFlags() error {
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

func (cmd *ApprovesCommand) createOperation() (base.Operation, error) { // nolint:dupl}
	e := util.StringError(utils.ErrStringCreate("approves operation"))

	item := token.NewApprovesItem(cmd.contract,
		cmd.approved, cmd.Amount.Big, cmd.Currency.CID)

	fact := token.NewApprovesFact(
		[]byte(cmd.Token), cmd.sender, []token.ApprovesItem{item},
	)

	op := token.NewApproves(fact)
	if err := op.Sign(cmd.Privatekey, cmd.NetworkID.NetworkID()); err != nil {
		return nil, e.Wrap(err)
	}

	return op, nil
}
