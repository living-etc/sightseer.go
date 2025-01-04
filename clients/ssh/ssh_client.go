package ssh

import (
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/crypto/ssh"

	"github.com/living-etc/sightseer.go/entities/linux"
)

type SshClient struct {
	host     string
	client   *ssh.Client
	platform string
}

func NewSshClient(
	privateKey []byte,
	host string,
	port string,
	user string,
	platform string,
) (*SshClient, error) {
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
		host:     host,
		client:   client,
		platform: platform,
	}

	return sshClient, nil
}

func (sshclient *SshClient) Service(name string) (*linux.Service, error) {
	return get(name, sshclient, linux.ServiceQuery{})
}

func (sshclient *SshClient) File(name string) (*linux.File, error) {
	return get(name, sshclient, linux.FileQuery{})
}

func (sshclient *SshClient) User(name string) (*linux.User, error) {
	return get(name, sshclient, linux.UserQuery{})
}

func (sshclient *SshClient) SystemdTimer(name string) (*linux.SystemdTimer, error) {
	return get(name, sshclient, linux.SystemdTimerQuery{})
}

func (sshclient *SshClient) LinuxKernelParameter(name string) (*linux.LinuxKernelParameter, error) {
	return get(name, sshclient, linux.LinuxKernelParameterQuery{})
}

func (sshclient *SshClient) Package(name string) (*linux.Package, error) {
	return get(name, sshclient, linux.PackageQuery{})
}

func get[T ResourceType, Q ResourceQuery[T]](
	identifier string,
	sshclient *SshClient,
	query Q,
) (*T, error) {
	cmdTemplate, err := query.Command(sshclient.platform)
	if err != nil {
		return nil, err
	}

	session, err := sshclient.client.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	var commandOutput string

	command := fmt.Sprintf(cmdTemplate, identifier)
	err = session.Run(command)
	if err != nil {
		commandOutput = stderr.String()
	} else {
		commandOutput = stdout.String()
	}

	commandOutput = strings.TrimSuffix(commandOutput, "\n")

	var q Q
	resource, err := q.ParseOutput(commandOutput)
	if err != nil {
		return nil, err
	}

	return resource, nil
}
