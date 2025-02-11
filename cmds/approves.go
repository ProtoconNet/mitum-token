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

type ApprovesCommand struct {
	OperationCommand
	Approved1 ccmds.AddressFlag `arg:"" name:"approved" help:"approved account" required:"true"`
	Approved2 ccmds.AddressFlag `arg:"" name:"approved" help:"approved account" required:"true"`
	Amount    ccmds.BigFlag     `arg:"" name:"amount" help:"amount to approve" required:"true"`
	approved1 base.Address
	approved2 base.Address
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

	ccmds.PrettyPrint(cmd.Out, op)

	return nil
}

func (cmd *ApprovesCommand) parseFlags() error {
	if err := cmd.OperationCommand.parseFlags(); err != nil {
		return err
	}

	approved, err := cmd.Approved1.Encode(cmd.Encoders.JSON())
	if err != nil {
		return errors.Wrapf(err, "invalid approved format, %q", cmd.Approved1.String())
	}
	cmd.approved1 = approved

	approved, err = cmd.Approved2.Encode(cmd.Encoders.JSON())
	if err != nil {
		return errors.Wrapf(err, "invalid approved format, %q", cmd.Approved2.String())
	}
	cmd.approved2 = approved

	return nil
}

func (cmd *ApprovesCommand) createOperation() (base.Operation, error) { // nolint:dupl}
	e := util.StringError(utils.ErrStringCreate("approves operation"))

	item1 := token.NewApprovesItem(cmd.contract,
		cmd.approved1, cmd.Amount.Big, cmd.Currency.CID)

	item2 := token.NewApprovesItem(cmd.contract,
		cmd.approved2, cmd.Amount.Big, cmd.Currency.CID)

	fact := token.NewApprovesFact(
		[]byte(cmd.Token), cmd.sender, []token.ApprovesItem{item1, item2},
	)

	op := token.NewApproves(fact)
	if err := op.Sign(cmd.Privatekey, cmd.NetworkID.NetworkID()); err != nil {
		return nil, e.Wrap(err)
	}

	return op, nil
}
