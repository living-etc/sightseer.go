package ssh

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

type SshClient struct {
	host     string
	client   *ssh.Client
	executor CommandExecutor
}

func NewSshClient(privateKey []byte, host string, user string) (*SshClient, error) {
	signer, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", host+":22", config)
	if err != nil {
		return nil, err
	}

	sshClient := &SshClient{
		host:     host,
		client:   client,
		executor: RealCommandExecutor{client: client},
	}

	return sshClient, nil
}

type ResourceType interface {
	File | Service
}

func Get[T ResourceType](
	identifier string,
	sshclient *SshClient,
	commandString string,
	parser func(string) (*T, error),
) (*T, error) {
	command := fmt.Sprintf(commandString, identifier)

	output, err := sshclient.executor.ExecuteCommand(command)
	if err != nil {
		return nil, err
	}

	resource, err := parser(output)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

func (sshClient SshClient) Command(command string) (string, error) {
	output, err := sshClient.executor.ExecuteCommand(command)
	if err != nil {
		return "", err
	}
	return output, nil
}
