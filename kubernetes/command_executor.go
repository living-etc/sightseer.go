package kubernetes

type CommandExecutor interface {
	executeCommand(binary string, args []string) (string, error)
}
