package cmds

type TokenCommand struct {
	RegisterToken RegisterTokenCommand `cmd:"" name:"register-token" help:"register token to contract account"`
	Mint          MintCommand          `cmd:"" name:"mint" help:"mint token to receiver"`
	Burn          BurnCommand          `cmd:"" name:"burn" help:"burn token of target"`
	Approve       ApproveCommand       `cmd:"" name:"approve" help:"approve token to approved account"`
	Transfer      TransferCommand      `cmd:"" name:"transfer" help:"transfer token to receiver"`
}
