package ssh

import (
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
	output, err := sshclient.executor.ExecuteCommand(command)
	file, _ := fileFromStatOutput(output)
	check(err, "")

	return file, nil
}

func (sshClient SshClient) Hostname() string {
	output, err := sshClient.executor.ExecuteCommand("hostname -s")
	check(err, "Error occurred running command over SSH")
	return output
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
