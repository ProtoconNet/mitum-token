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

type TransfersFromCommand struct {
	OperationCommand
	Receiver ccmds.AddressFlag `arg:"" name:"receiver" help:"token receiver" required:"true"`
	Target1  ccmds.AddressFlag `arg:"" name:"target1" help:"target approving" required:"true"`
	Target2  ccmds.AddressFlag `arg:"" name:"target2" help:"target approving" required:"true"`
	Amount   ccmds.BigFlag     `arg:"" name:"amount" help:"amount to transfer" required:"true"`
	receiver base.Address
	target1  base.Address
	target2  base.Address
}

func (cmd *TransfersFromCommand) Run(pctx context.Context) error { // nolint:dupl
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

func (cmd *TransfersFromCommand) parseFlags() error {
	if err := cmd.OperationCommand.parseFlags(); err != nil {
		return err
	}

	receiver, err := cmd.Receiver.Encode(cmd.Encoders.JSON())
	if err != nil {
		return errors.Wrapf(err, "invalid receiver format, %q", cmd.Receiver.String())
	}
	cmd.receiver = receiver

	target, err := cmd.Target1.Encode(cmd.Encoders.JSON())
	if err != nil {
		return errors.Wrapf(err, "invalid target format, %q", cmd.Target1.String())
	}
	cmd.target1 = target

	target, err = cmd.Target2.Encode(cmd.Encoders.JSON())
	if err != nil {
		return errors.Wrapf(err, "invalid target format, %q", cmd.Target2.String())
	}
	cmd.target2 = target

	return nil
}

func (cmd *TransfersFromCommand) createOperation() (base.Operation, error) { // nolint:dupl}
	e := util.StringError(utils.ErrStringCreate("transfer-from operation"))

	item1 := token.NewTransfersFromItem(cmd.contract,
		cmd.receiver, cmd.target1, cmd.Amount.Big, cmd.Currency.CID)

	item2 := token.NewTransfersFromItem(cmd.contract,
		cmd.receiver, cmd.target2, cmd.Amount.Big, cmd.Currency.CID)

	fact := token.NewTransfersFromFact(
		[]byte(cmd.Token), cmd.sender, []token.TransfersFromItem{item1, item2},
	)

	op := token.NewTransfersFrom(fact)
	if err := op.Sign(cmd.Privatekey, cmd.NetworkID.NetworkID()); err != nil {
		return nil, e.Wrap(err)
	}

	return op, nil
}
