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

type BurnCommand struct {
	OperationCommand
	Target currencycmds.AddressFlag `arg:"" name:"target" help:"token target" required:"true"`
	Amount currencycmds.BigFlag     `arg:"" name:"amount" help:"amount to burn" required:"true"`
	target base.Address
}

func (cmd *BurnCommand) Run(pctx context.Context) error { // nolint:dupl
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

func (cmd *BurnCommand) parseFlags() error {
	if err := cmd.OperationCommand.parseFlags(); err != nil {
		return err
	}

	target, err := cmd.Target.Encode(enc)
	if err != nil {
		return errors.Wrapf(err, "invalid target format, %q", cmd.Target.String())
	}
	cmd.target = target

	return nil
}

func (cmd *BurnCommand) createOperation() (base.Operation, error) { // nolint:dupl}
	e := util.StringError(utils.ErrStringCreate("burn operation"))

	fact := token.NewBurnFact(
		[]byte(cmd.Token),
		cmd.sender, cmd.contract,
		cmd.Currency.CID,
		cmd.target,
		cmd.Amount.Big,
	)

	op := token.NewBurn(fact)
	if err := op.Sign(cmd.Privatekey, cmd.NetworkID.NetworkID()); err != nil {
		return nil, e.Wrap(err)
	}

	return op, nil
}
