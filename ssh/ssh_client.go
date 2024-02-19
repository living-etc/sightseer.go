package ssh

import (
	"fmt"
	"log"

	"golang.org/x/crypto/ssh"
)

type SshClient struct {
	host     string
	client   *ssh.Client
	executor CommandExecutor
}

func NewSshClient(privateKey []byte, host string) (SshClient, error) {
	signer, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		log.Fatal(err.Error())
	}

	config := &ssh.ClientConfig{
		User: "ubuntu",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", host+":22", config)

	sshClient := SshClient{
		host:     host,
		client:   client,
		executor: RealCommandExecutor{client: client},
	}

	return sshClient, nil
}

func (sshClient SshClient) fatalError(err error, message string) {
	if err != nil {
		log.Fatalf("[%v] %v: %v", sshClient.host, message, err)
	}
}

func (sshClient SshClient) HasFile(filename string) bool {
	command := fmt.Sprintf("test -f %v", filename)
	_, err := sshClient.executor.ExecuteCommand(command)
	sshClient.fatalError(err, "")

	if err != nil {
		return false
	}

	return true
}

func (sshclient SshClient) File(filename string) (File, error) {
	command := fmt.Sprintf("stat %v", filename)
	output, err := sshclient.executor.ExecuteCommand(command)
	if err != nil {
		return File{}, &CommandError{Context: output, Err: err.Error()}
	}

	file, err := fileFromStatOutput(output)

	return file, nil
}

func (sshClient SshClient) Hostname() string {
	output, err := sshClient.executor.ExecuteCommand("hostname -s")
	sshClient.fatalError(err, "")
	return output
}

func (sshClient SshClient) Service(name string) (Service, error) {
	command := fmt.Sprintf("systemctl status %v --no-pager", name)
	output, err := sshClient.executor.ExecuteCommand(command)

	service, err := serviceFromSystemctl(output)
	sshClient.fatalError(err, "")

	return service, nil
}

func (sshClient SshClient) Command(command string) string {
	output, err := sshClient.executor.ExecuteCommand(command)
	sshClient.fatalError(err, "")
	return output
}
