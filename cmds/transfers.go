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

type TransfersCommand struct {
	OperationCommand
	Receiver1 ccmds.AddressFlag `arg:"" name:"receiver" help:"token receiver" required:"true"`
	Receiver2 ccmds.AddressFlag `arg:"" name:"receiver" help:"token receiver" required:"true"`
	Amount    ccmds.BigFlag     `arg:"" name:"amount" help:"amount to transfer" required:"true"`
	receiver1 base.Address
	receiver2 base.Address
}

func (cmd *TransfersCommand) Run(pctx context.Context) error { // nolint:dupl
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

func (cmd *TransfersCommand) parseFlags() error {
	if err := cmd.OperationCommand.parseFlags(); err != nil {
		return err
	}

	receiver, err := cmd.Receiver1.Encode(cmd.Encoders.JSON())
	if err != nil {
		return errors.Wrapf(err, "invalid receiver format, %q", cmd.Receiver1.String())
	}
	cmd.receiver1 = receiver

	receiver, err = cmd.Receiver2.Encode(cmd.Encoders.JSON())
	if err != nil {
		return errors.Wrapf(err, "invalid receiver format, %q", cmd.Receiver2.String())
	}
	cmd.receiver2 = receiver

	return nil
}

func (cmd *TransfersCommand) createOperation() (base.Operation, error) { // nolint:dupl}
	e := util.StringError(utils.ErrStringCreate("transfer operation"))

	item1 := token.NewTransfersItem(cmd.contract,
		cmd.receiver1, cmd.Amount.Big, cmd.Currency.CID)

	item2 := token.NewTransfersItem(cmd.contract,
		cmd.receiver2, cmd.Amount.Big, cmd.Currency.CID)

	fact := token.NewTransfersFact(
		[]byte(cmd.Token), cmd.sender, []token.TransfersItem{item1, item2},
	)

	op := token.NewTransfers(fact)
	if err := op.Sign(cmd.Privatekey, cmd.NetworkID.NetworkID()); err != nil {
		return nil, e.Wrap(err)
	}

	return op, nil
}
