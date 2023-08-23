package cmds

type TokenCommand struct {
	RegisterToken RegisterTokenCommand `cmd:"" name:"register-token" help:"register token to contract account"`
}
