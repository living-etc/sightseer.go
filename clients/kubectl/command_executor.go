package kubectl

type CommandExecutor interface {
	ExecuteCommand(binary string, args []string) (string, error)
}
