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

type TransferFromCommand struct {
	OperationCommand
	Receiver currencycmds.AddressFlag `arg:"" name:"receiver" help:"token receiver" required:"true"`
	Target   currencycmds.AddressFlag `arg:"" name:"target" help:"target approving" required:"true"`
	Amount   currencycmds.BigFlag     `arg:"" name:"amount" help:"amount to transfer" required:"true"`
	receiver base.Address
	target   base.Address
}

func NewTransferFromCommand() TransferFromCommand {
	cmd := NewOperationCommand()
	return TransferFromCommand{
		OperationCommand: *cmd,
	}
}

func (cmd *TransferFromCommand) Run(pctx context.Context) error { // nolint:dupl
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

func (cmd *TransferFromCommand) parseFlags() error {
	if err := cmd.OperationCommand.parseFlags(); err != nil {
		return err
	}

	receiver, err := cmd.Receiver.Encode(enc)
	if err != nil {
		return errors.Wrapf(err, "invalid receiver format, %q", cmd.Receiver.String())
	}
	cmd.receiver = receiver

	target, err := cmd.Target.Encode(enc)
	if err != nil {
		return errors.Wrapf(err, "invalid target format, %q", cmd.Target.String())
	}
	cmd.target = target

	return nil
}

func (cmd *TransferFromCommand) createOperation() (base.Operation, error) { // nolint:dupl}
	e := util.StringError(utils.ErrStringCreate("transfer-from operation"))

	fact := token.NewTransferFromFact(
		[]byte(cmd.Token),
		cmd.sender, cmd.contract,
		cmd.TokenID.CID, cmd.Currency.CID,
		cmd.receiver,
		cmd.target,
		cmd.Amount.Big,
	)

	op := token.NewTransferFrom(fact)
	if err := op.HashSign(cmd.Privatekey, cmd.NetworkID.NetworkID()); err != nil {
		return nil, e.Wrap(err)
	}

	return op, nil
}
