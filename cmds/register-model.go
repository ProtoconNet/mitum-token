package cmds

import (
	"context"

	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum-token/operation/token"
	"github.com/ProtoconNet/mitum-token/utils"

	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
)

type RegisterModelCommand struct {
	OperationCommand
	Symbol        TokenSymbolFlag      `arg:"" name:"symbol" help:"token symbol" required:"true"`
	Name          string               `arg:"" name:"name" help:"token name" required:"true"`
	InitialSupply currencycmds.BigFlag `arg:"" name:"initial-supply" help:"initial supply of token" required:"true"`
}

func (cmd *RegisterModelCommand) Run(pctx context.Context) error { // nolint:dupl
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

func (cmd *RegisterModelCommand) createOperation() (base.Operation, error) { // nolint:dupl}
	e := util.StringError(utils.ErrStringCreate("register-model operation"))

	fact := token.NewRegisterModelFact(
		[]byte(cmd.Token),
		cmd.sender, cmd.contract,
		cmd.Currency.CID, cmd.Symbol.Symbol,
		cmd.Name,
		cmd.InitialSupply.Big,
	)

	op := token.NewRegisterModel(fact)
	if err := op.Sign(cmd.Privatekey, cmd.NetworkID.NetworkID()); err != nil {
		return nil, e.Wrap(err)
	}

	return op, nil
}
