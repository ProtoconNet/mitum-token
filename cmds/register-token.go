package cmds

import (
	"context"

	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum-token/operation/token"
	"github.com/ProtoconNet/mitum-token/utils"

	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
)

type RegisterTokenCommand struct {
	OperationCommand
	Symbol      currencycmds.CurrencyIDFlag `arg:"" name:"symbol" help:"token symbol" required:"true"`
	Name        string                      `arg:"" name:"name" help:"token name" required:"true"`
	TotalSupply currencycmds.BigFlag        `arg:"" name:"total-supply" help:"total supply of token" required:"true"`
}

func (cmd *RegisterTokenCommand) Run(pctx context.Context) error { // nolint:dupl
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

func (cmd *RegisterTokenCommand) createOperation() (base.Operation, error) { // nolint:dupl}
	e := util.StringError(utils.ErrStringCreate("register-token operation"))

	fact := token.NewRegisterTokenFact(
		[]byte(cmd.Token),
		cmd.sender, cmd.contract,
		cmd.Currency.CID, cmd.Symbol.CID,
		cmd.Name,
		cmd.TotalSupply.Big,
	)

	op := token.NewRegisterToken(fact)
	if err := op.Sign(cmd.Privatekey, cmd.NetworkID.NetworkID()); err != nil {
		return nil, e.Wrap(err)
	}

	return op, nil
}
