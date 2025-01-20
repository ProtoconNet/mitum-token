package main

import (
	"context"
	"fmt"
	"os"

	ccmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum-token/cmds"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/launch"
	launchcmd "github.com/ProtoconNet/mitum2/launch/cmd"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/logging"
	"github.com/alecthomas/kong"
	"github.com/pkg/errors"
)

var (
	Version   = "v0.0.1"
	BuildTime = "-"
	GitBranch = "master"
	GitCommit = "-"
)

//revive:disable:nested-structs
var CLI struct { //nolint:govet //...
	launch.BaseFlags
	Init      ccmds.INITCommand `cmd:"" help:"init node"`
	Run       cmds.RunCommand   `cmd:"" help:"run node"`
	Storage   cmds.Storage      `cmd:""`
	Operation struct {
		Currency ccmds.CurrencyCommand `cmd:"" help:"currency operation"`
		Suffrage ccmds.SuffrageCommand `cmd:"" help:"suffrage operation"`
		Token    cmds.TokenCommand     `cmd:"" help:"token operation"`
	} `cmd:"" help:"create operation"`
	Network struct {
		Client cmds.NetworkClientCommand `cmd:"" help:"network client"`
	} `cmd:"" help:"network"`
	Key struct {
		New     ccmds.KeyNewCommand      `cmd:"" help:"generate new key"`
		Address ccmds.KeyAddressCommand  `cmd:"" help:"generate address from key"`
		Load    launchcmd.KeyLoadCommand `cmd:"" help:"load key"`
		Sign    launchcmd.KeySignCommand `cmd:"" help:"sign"`
	} `cmd:"" help:"key"`
	Handover launchcmd.HandoverCommands `cmd:""`
	Version  struct{}                   `cmd:"" help:"version"`
}

var flagDefaults = kong.Vars{
	"log_out":                           "stderr",
	"log_format":                        "terminal",
	"log_level":                         "debug",
	"log_force_color":                   "false",
	"design_uri":                        launch.DefaultDesignURI,
	"create_account_threshold":          "100",
	"create_contract_account_threshold": "100",
	"suffrage_candidate_limiter_limit":  "77",
	"max_operation_in_proposal":         "99",
	"suffrage candidate lifespan":       "33",
	"max suffrage size":                 "33",
	"suffrage expel lifespan":           "44",
	"safe_threshold":                    base.SafeThreshold.String(),
	"network_id":                        "mitum",
}

func main() {
	kctx := kong.Parse(&CLI, flagDefaults)

	bi, err := util.ParseBuildInfo(Version, GitBranch, GitCommit, BuildTime)
	if err != nil {
		kctx.FatalIfErrorf(err)
	}

	if kctx.Command() == "version" {
		_, _ = fmt.Fprintln(os.Stdout, bi.String())

		return
	}

	pctx := util.ContextWithValues(context.Background(), map[util.ContextKey]interface{}{
		launch.VersionContextKey:     bi.Version,
		launch.FlagsContextKey:       CLI.BaseFlags,
		launch.KongContextContextKey: kctx,
	})

	pss := launch.DefaultMainPS()

	switch i, err := pss.Run(pctx); {
	case err != nil:
		kctx.FatalIfErrorf(err)
	default:
		pctx = i
		kctx = kong.Parse(&CLI, kong.BindTo(pctx, (*context.Context)(nil)), flagDefaults)
	}

	var log *logging.Logging
	if err := util.LoadFromContextOK(pctx, launch.LoggingContextKey, &log); err != nil {
		kctx.FatalIfErrorf(err)
	}

	log.Log().Debug().Interface("flags", os.Args).Msg("flags")
	log.Log().Debug().Interface("main_process", pss.Verbose()).Msg("processed")

	if err := func() error {
		defer log.Log().Debug().Msg("stopped")
		return errors.WithStack(kctx.Run(pctx))
	}(); err != nil {
		log.Log().Error().Err(err).Msg("stopped by error")
		kctx.FatalIfErrorf(err)
	}
}
