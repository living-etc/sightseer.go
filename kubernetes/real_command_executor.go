package kubernetes

import (
	"os/exec"
)

type RealCommandExecutor struct {
	kubeConfigPath string
}

func (executor RealCommandExecutor) executeCommand(binary string, args []string) (string, error) {
	cmd := exec.Command(binary, args...)

	output, err := cmd.CombinedOutput()

	return string(output), err
}
