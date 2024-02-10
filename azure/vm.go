package azure

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

	"golang.org/x/crypto/ssh"
)

type VM struct {
	PrivateIPAddress string
	PublicIPAddress  string
	DnsName          string
}

func check(err error, message string) {
	if err != nil {
		log.Fatalf("%v: %v", message, err)
	}
}

func (vm VM) ReachableOnPort(port int) (bool, error) {
	cmd := exec.Command("nc", vm.PublicIPAddress, strconv.Itoa(port), "-w 1")

	_, err := cmd.Output()
	if err != nil {
		return false, errors.New("non-zero exit code from nc: " + err.Error())
	}

	return true, nil
}

func (vm VM) newSSHSession() (*ssh.Client, error) {
	key, err := os.ReadFile("../keys/id_rsa")
	check(err, "Unable to read private key file")

	signer, err := ssh.ParsePrivateKey(key)
	check(err, "Unable to parse private key file")

	config := &ssh.ClientConfig{
		User: "ubuntu",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", vm.PublicIPAddress+":22", config)

	return client, err
}

func (vm VM) ConnectableOverSSH(publicKeyPath string) (bool, error) {
	client, err := vm.newSSHSession()
	check(err, "Unable to connect to "+vm.PublicIPAddress)
	defer client.Close()

	return true, nil
}

func (vm VM) commandOverSSH(command string) (string, error) {
	client, err := vm.newSSHSession()
	check(err, "Unable to connect to "+vm.PublicIPAddress)

	session, err := client.NewSession()
	check(err, "Unable to open SSH session")
	defer session.Close()

	cmd := fmt.Sprintf(command)

	var b bytes.Buffer
	session.Stdout = &b
	if err = session.Run(cmd); err != nil {
		return "", errors.New("Failed to run: " + err.Error())
	}

	return b.String(), nil
}

func (vm VM) Hostname() string {
	output, err := vm.commandOverSSH("echo -n $( hostname -s )")
	check(err, "Error occurred running command over SSH")
	return output
}

func (vm VM) HasFile(filename string) bool {
	command := fmt.Sprintf("test -f %v", filename)
	_, err := vm.commandOverSSH(command)
	check(err, "Error occurred running command over SSH")

	if err != nil {
		return false
	}

	return true
}
