package ssh

import (
	"bytes"
	"errors"
	"fmt"
	"log"

	"golang.org/x/crypto/ssh"
)

func check(err error, message string) {
	if err != nil {
		log.Fatalf("%v: %v", message, err)
	}
}

type SshClient struct {
	client   *ssh.Client
	executor CommandExecutor
}

type RealCommandExecutor struct {
	client *ssh.Client
}

type CommandExecutor interface {
	ExecuteCommand(command string) (string, error)
}

func NewSshClient(privateKey []byte, host string) (SshClient, error) {
	signer, err := ssh.ParsePrivateKey(privateKey)
	check(err, "Unable to parse private key file")

	config := &ssh.ClientConfig{
		User: "ubuntu",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", host+":22", config)

	sshClient := SshClient{
		client:   client,
		executor: RealCommandExecutor{client: client},
	}

	return sshClient, nil
}

func (sshClient SshClient) HasFile(filename string) bool {
	command := fmt.Sprintf("test -f %v", filename)
	_, err := sshClient.executor.ExecuteCommand(command)
	check(err, "Error occurred running command over SSH")

	if err != nil {
		return false
	}

	return true
}

func (sshclient SshClient) File(filename string) (File, error) {
	command := fmt.Sprintf("stat %v", filename)
	_, err := sshclient.executor.ExecuteCommand(command)
	check(err, "")

	return File{}, nil
}

func (sshClient SshClient) Hostname() string {
	output, err := sshClient.executor.ExecuteCommand("hostname -s")
	check(err, "Error occurred running command over SSH")
	return output
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

type Service struct {
	IsActive string
}

func (sshClient SshClient) Service(name string) Service {
	command := fmt.Sprintf("systemctl is-active %v", name)
	output, err := sshClient.executor.ExecuteCommand(command)
	check(err, "Error encountered checking service status")

	return Service{
		IsActive: output,
	}
}

func (sshClient SshClient) Command(command string) string {
	output, err := sshClient.executor.ExecuteCommand(command)
	check(err, "Error running command")
	return output
}
