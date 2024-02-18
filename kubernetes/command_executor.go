package kubernetes

type CommandExecutor interface {
	executeCommand(command string) (string, error)
}
