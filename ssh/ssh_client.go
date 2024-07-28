package ssh

import (
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/crypto/ssh"
)

type ResourceType interface{}

type SshClient struct {
	host   string
	client *ssh.Client
}

func NewSshClient(privateKey []byte, host string, port string, user string) (*SshClient, error) {
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

	client, err := ssh.Dial("tcp", host+":"+port, config)
	if err != nil {
		return nil, err
	}

	sshClient := &SshClient{
		host:   host,
		client: client,
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

	output, err := sshclient.RunCommand(command)

	var q Q
	resource, err := q.ParseOutput(output)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

func (sshClient SshClient) RunCommand(command string) (string, error) {
	session, err := sshClient.client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	var output string

	err = session.Run(command)
	if err != nil {
		output = stderr.String()
	} else {
		output = stdout.String()
	}

	output = strings.TrimSuffix(output, "\n")

	return output, err
}
