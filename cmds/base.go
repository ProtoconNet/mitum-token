package cmds

import (
	"context"
	"fmt"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/pkg/errors"
	"io"
	"os"

	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum2/launch"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/logging"
	"github.com/ProtoconNet/mitum2/util/ps"
	"github.com/rs/zerolog"
)

type BaseCommand struct {
	Encoder  encoder.Encoder   `kong:"-"`
	Encoders *encoder.Encoders `kong:"-"`
	Log      *zerolog.Logger   `kong:"-"`
	Out      io.Writer         `kong:"-"`
}

func (cmd *BaseCommand) prepare(pctx context.Context) (context.Context, error) {
	cmd.Out = os.Stdout
	pps := ps.NewPS("cmd")

	_ = pps.
		AddOK(launch.PNameEncoder, currencycmds.PEncoder, nil)

	_ = pps.POK(launch.PNameEncoder).
		PostAddOK(launch.PNameAddHinters, PAddHinters)

	var log *logging.Logging
	if err := util.LoadFromContextOK(pctx, launch.LoggingContextKey, &log); err != nil {
		return pctx, err
	}

	cmd.Log = log.Log()

	pctx, err := pps.Run(pctx) //revive:disable-line:modifies-parameter
	if err != nil {
		return pctx, err
	}

	if err := util.LoadFromContextOK(pctx, launch.EncodersContextKey, &cmd.Encoders); err != nil {
		return pctx, err
	}

	cmd.Encoder = cmd.Encoders.JSON()

	return pctx, nil
}

func (cmd *BaseCommand) print(f string, a ...interface{}) {
	_, _ = fmt.Fprintf(cmd.Out, f, a...)
	_, _ = fmt.Fprintln(cmd.Out)
}

func PAddHinters(pctx context.Context) (context.Context, error) {
	e := util.StringError("add hinters")

	var encs *encoder.Encoders
	if err := util.LoadFromContextOK(pctx, launch.EncodersContextKey, &encs); err != nil {
		return pctx, e.Wrap(err)
	}

	if err := LoadHinters(encs); err != nil {
		return pctx, e.Wrap(err)
	}

	return pctx, nil
}

type OperationCommand struct {
	BaseCommand
	currencycmds.OperationFlags
	Sender   currencycmds.AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	Contract currencycmds.AddressFlag    `arg:"" name:"contract" help:"contract address to register token" required:"true"`
	Currency currencycmds.CurrencyIDFlag `arg:"" name:"currency" help:"currency id" required:"true"`
	sender   base.Address
	contract base.Address
}

func (cmd *OperationCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	sender, err := cmd.Sender.Encode(cmd.Encoders.JSON())
	if err != nil {
		return errors.Wrapf(err, "invalid sender format, %q", cmd.Sender.String())
	}
	cmd.sender = sender

	contract, err := cmd.Contract.Encode(cmd.Encoders.JSON())
	if err != nil {
		return errors.Wrapf(err, "invalid contract account format, %q", cmd.Contract.String())
	}
	cmd.contract = contract

	return nil
}
