package ssh

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/ssh"
)

type RealCommandExecutor struct {
	client *ssh.Client
}

func (executor RealCommandExecutor) ExecuteCommand(command string) (string, error) {
	session, err := executor.client.NewSession()
	check(err, "Unable to open SSH session")
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err = session.Run(command); err != nil {
		return "", errors.New(fmt.Sprintf(`Failed to run "%v": %v`, command, err.Error()))
	}

	output := strings.TrimSuffix(b.String(), "\n")

	return output, nil
}
