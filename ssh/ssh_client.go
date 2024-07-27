package ssh

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

type ResourceType interface{}

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

func (sshclient *SshClient) Service(name string) (*Service, error) {
	return get(name, sshclient, ServiceQuery{})
}

func (sshclient *SshClient) File(name string) (*File, error) {
	return get(name, sshclient, FileQuery{})
}

func (sshclient *SshClient) User(name string) (*User, error) {
	return get(name, sshclient, UserQuery{})
}

func (sshclient *SshClient) SystemdTimer(name string) (*SystemdTimer, error) {
	return get(name, sshclient, SystemdTimerQuery{})
}

func (sshclient *SshClient) LinuxKernelParameter(name string) (*LinuxKernelParameter, error) {
	return get(name, sshclient, LinuxKernelParameterQuery{})
}

func get[T ResourceType, Q ResourceQuery[T]](
	identifier string,
	sshclient *SshClient,
	query Q,
) (*T, error) {
	command := fmt.Sprintf(query.Command(), identifier)

	output, err := sshclient.executor.ExecuteCommand(command)
	if err != nil {
		return nil, err
	}

	var q Q
	resource, err := q.ParseOutput(output)
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
