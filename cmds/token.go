package cmds

type TokenCommand struct {
	RegisterToken RegisterTokenCommand `cmd:"" name:"register-token" help:"register token to contract account"`
	Mint          MintCommand          `cmd:"" name:"mint" help:"mint token to receiver"`
}
