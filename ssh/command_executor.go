package ssh

type CommandExecutor interface {
	ExecuteCommand(command string) (string, error)
}
