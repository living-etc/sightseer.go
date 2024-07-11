package ssh

import (
	"strings"

	"golang.org/x/crypto/ssh"
)

type RealCommandExecutor struct {
	client *ssh.Client
}

func (executor RealCommandExecutor) ExecuteCommand(command string) (string, error) {
	session, err := executor.client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	output, err := session.CombinedOutput(command)
	o := strings.TrimSuffix(string(output), "\n")

	if err != nil {
		return o, err
	}

	return o, nil
}
