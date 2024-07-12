package ssh

import (
	"errors"
	"fmt"
	"regexp"

	"golang.org/x/crypto/ssh"
)

type (
	Resource     interface{}
	ResourceType interface {
		Resource
	}
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

type ResourceQuery[T ResourceType] interface {
	Regex() string
	Command() string
	SetValues(*T, map[string]string)
}

func Get[T ResourceType, Q ResourceQuery[T]](
	identifier string,
	sshclient *SshClient,
	query Q,
) (*T, error) {
	command := fmt.Sprintf(query.Command(), identifier)

	output, err := sshclient.executor.ExecuteCommand(command)
	if err != nil {
		return nil, err
	}

	resource, err := ParseOutput[T, Q](output)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

func ParseOutput[T ResourceType, Q ResourceQuery[T]](output string) (*T, error) {
	var q Q
	regex := q.Regex()

	re := regexp.MustCompile(regex)
	matches := re.FindStringSubmatch(output)
	if len(matches) == 0 {
		return nil, errors.New("failed to parse output")
	}

	values := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i != 0 && name != "" {
			values[name] = matches[i]
		}
	}

	result := new(T)
	q.SetValues(result, values)

	return result, nil
}

func (sshClient SshClient) Command(command string) (string, error) {
	output, err := sshClient.executor.ExecuteCommand(command)
	if err != nil {
		return "", err
	}
	return output, nil
}
