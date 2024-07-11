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

func (sshclient SshClient) File(filename string) (*File, error) {
	command := fmt.Sprintf("stat %v", filename)
	output, err := sshclient.executor.ExecuteCommand(command)
	if err != nil {
		return nil, err
	}

	file, err := fileFromStatOutput(output)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (sshClient SshClient) Hostname() (string, error) {
	output, err := sshClient.executor.ExecuteCommand("hostname -s")
	if err != nil {
		return "", err
	}
	return output, nil
}

func (sshClient SshClient) Service(name string) (*Service, error) {
	command := fmt.Sprintf("systemctl status %v --no-pager", name)

	output, err := sshClient.executor.ExecuteCommand(command)
	if err != nil {
		return nil, err
	}

	service, err := serviceFromSystemctl(output)
	if err != nil {
		return nil, err
	}

	return service, nil
}

func (sshClient SshClient) Command(command string) (string, error) {
	output, err := sshClient.executor.ExecuteCommand(command)
	if err != nil {
		return "", err
	}
	return output, nil
}
