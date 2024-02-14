package ssh

import (
	"bytes"
	"errors"
	"fmt"

	"golang.org/x/crypto/ssh"
)

type RealCommandExecutor struct {
	client *ssh.Client
}

func (executor RealCommandExecutor) ExecuteCommand(command string) (string, error) {
	session, err := executor.client.NewSession()
	check(err, "Unable to open SSH session")
	defer session.Close()

	cmd := fmt.Sprintf("echo -n $( %v )", command)

	var b bytes.Buffer
	session.Stdout = &b
	if err = session.Run(cmd); err != nil {
		return "", errors.New("Failed to run: " + err.Error())
	}

	return b.String(), nil
}
