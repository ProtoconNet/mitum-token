package cmds

import (
	"context"
	"io"
	"os"

	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/pkg/errors"

	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/launch"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/logging"
	"github.com/ProtoconNet/mitum2/util/ps"
	"github.com/rs/zerolog"
)

type BaseCommand struct {
	Encoder  *jsonenc.Encoder
	Encoders *encoder.Encoders
	Log      *zerolog.Logger
	Out      io.Writer `kong:"-"`
}

func NewBaseCommand() *BaseCommand {
	return &BaseCommand{
		Out: os.Stdout,
	}
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

	return pctx, util.LoadFromContextOK(pctx,
		launch.EncodersContextKey, &cmd.Encoders,
		launch.EncoderContextKey, &cmd.Encoder,
	)
}

func PAddHinters(ctx context.Context) (context.Context, error) {
	e := util.StringError("failed to add hinters")

	var enc encoder.Encoder
	if err := util.LoadFromContextOK(ctx, launch.EncoderContextKey, &enc); err != nil {
		return ctx, e.Wrap(err)
	}
	var benc encoder.Encoder
	if err := util.LoadFromContextOK(ctx, currencycmds.BEncoderContextKey, &benc); err != nil {
		return ctx, e.Wrap(err)
	}

	if err := LoadHinters(enc); err != nil {
		return ctx, e.Wrap(err)
	}

	if err := LoadHinters(benc); err != nil {
		return ctx, e.Wrap(err)
	}

	return ctx, nil
}

type OperationCommand struct {
	BaseCommand
	currencycmds.OperationFlags
	Sender   currencycmds.AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	Contract currencycmds.AddressFlag    `arg:"" name:"contract" help:"contract address to register token" required:"true"`
	TokenID  currencycmds.CurrencyIDFlag `arg:"" name:"token-id" help:"token id" required:"true"`
	Currency currencycmds.CurrencyIDFlag `arg:"" name:"currency" help:"currency id" required:"true"`
	sender   base.Address
	contract base.Address
}

func NewOperationCommand() *OperationCommand {
	cmd := NewBaseCommand()
	return &OperationCommand{
		BaseCommand: *cmd,
	}
}

func (cmd *OperationCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	sender, err := cmd.Sender.Encode(enc)
	if err != nil {
		return errors.Wrapf(err, "invalid sender format, %q", cmd.Sender.String())
	}
	cmd.sender = sender

	contract, err := cmd.Contract.Encode(enc)
	if err != nil {
		return errors.Wrapf(err, "invalid contract account format, %q", cmd.Contract.String())
	}
	cmd.contract = contract

	return nil
}
