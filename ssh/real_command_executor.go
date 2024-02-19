package ssh

import (
	"log"
	"strings"

	"golang.org/x/crypto/ssh"
)

type RealCommandExecutor struct {
	client *ssh.Client
}

func (executor RealCommandExecutor) ExecuteCommand(command string) (string, error) {
	session, err := executor.client.NewSession()
	if err != nil {
		log.Fatalf("Error creating SSH session: %v", err.Error())
	}
	defer session.Close()

	output, err := session.CombinedOutput(command)
	o := strings.TrimSuffix(string(output), "\n")

	if err != nil {
		return o, err
	}

	return o, nil
}
