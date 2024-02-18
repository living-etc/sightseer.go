package kubernetes

import (
	"os/exec"
)

type RealCommandExecutor struct {
	kubeConfigPath string
}

func (executor RealCommandExecutor) executeCommand(command string) (string, error) {
	output, err := exec.Command(command).Output()

	return string(output), err
}
