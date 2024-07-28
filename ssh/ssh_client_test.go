package ssh

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
)

type VagrantMachineConfig struct {
	ip   string
	name string
	port string
	user string
}

var testPlatforms = []string{
	"ubuntu2404",
	"fedora40",
}

var testMachines = map[string]VagrantMachineConfig{
	"ubuntu2404": {
		name: "ubuntu2404",
		ip:   "127.0.0.1",
		port: "2222",
		user: "vagrant",
	},
	"fedora40": {
		name: "fedora40",
		ip:   "127.0.0.1",
		port: "2200",
		user: "vagrant",
	},
}

func VagrantSetup(machineName string) *SshClient {
	machine := testMachines[machineName]

	privateKey, err := os.ReadFile(
		fmt.Sprintf(".vagrant/machines/%v/vmware_desktop/private_key", machine.name),
	)
	if err != nil {
		log.Fatalf("Error reading private key for %v", machine.name)
	}

	sshClient, err := NewSshClient(privateKey, machine.ip, machine.port, machine.user)
	if err != nil {
		log.Fatalf("Error creating ssh client: %v", err)
	}

	return sshClient
}

func TestFile(t *testing.T) {
	for _, testPlatform := range testPlatforms {
		sshClient := VagrantSetup(testPlatform)

		tests := TestCases.Get("File", testPlatform)

		for _, testcase := range tests {
			file, err := sshClient.File(testcase.resourceIdentifier)

			t.Run(testcase.testName, func(t *testing.T) {
				if err != nil {
					containsErrorString := strings.Contains(
						err.Error(),
						"No such file or directory",
					)

					if !containsErrorString {
						t.Errorf(
							"Error failed:\nwanted err to contain:\t%v\ngot:\t\t\t%v",
							"No such file or directory",
							err.Error(),
						)
					}
				} else {
					if reflect.DeepEqual(file, testcase.resourceWant) {
						t.Errorf("File failed:\nwant\t%v\ngot\t%v", testcase.resourceWant, file)
					}
				}
			})
		}
	}
}

func TestService(t *testing.T) {
	for _, testPlatform := range testPlatforms {
		sshClient := VagrantSetup(testPlatform)

		tests := TestCases.Get("Service", testPlatform)

		for _, testcase := range tests {
			service, err := sshClient.Service(testcase.resourceIdentifier)
			if err != nil {
				log.Fatalf("Error in %v: %v", testcase.resourceIdentifier, err)
			}

			t.Run(testcase.testName, func(t *testing.T) {
				if !reflect.DeepEqual(service, testcase.resourceWant) {
					t.Errorf(
						"Service failed:\nwant:\t%v\ngot\t\t%v",
						testcase.resourceWant,
						service,
					)
				}
			})
		}
	}
}

func TestUser(t *testing.T) {
	for _, testPlatform := range testPlatforms {
		sshClient := VagrantSetup(testPlatform)

		tests := TestCases.Get("User", testPlatform)

		for _, testcase := range tests {
			user, err := sshClient.User(testcase.resourceIdentifier)
			if err != nil {
				t.Fatalf("Error in %v: %v", testcase.testName, err)
			}

			t.Run(testcase.testName, func(t *testing.T) {
				if !reflect.DeepEqual(user, testcase.resourceWant) {
					t.Fatalf("User failed:\nwant:\t%v\ngot:\t%v", testcase.resourceWant, user)
				}
			})
		}
	}
}

func TestSystemdTimer(t *testing.T) {
	for _, testPlatform := range testPlatforms {
		sshClient := VagrantSetup(testPlatform)

		tests := TestCases.Get("SystemdTimer", testPlatform)

		for _, testcase := range tests {
			timer, err := sshClient.SystemdTimer(testcase.resourceIdentifier)
			if err != nil {
				if err.Error() != testcase.errWant.Error() {
					t.Fatalf(
						"SystemdTimer Error failed:\nwant:\t%v\ngot:\t%v",
						testcase.errWant,
						err,
					)
				}
			} else {
				t.Run(testcase.testName, func(t *testing.T) {
					if !reflect.DeepEqual(timer, testcase.resourceWant) {
						t.Fatalf("SystemdTimer failed:\nwant:\t%v\ngot:\t%v", testcase.resourceWant, timer)
					}
				})
			}
		}
	}
}

func TestLinuxKernelParameter(t *testing.T) {
	for _, testPlatform := range testPlatforms {
		sshClient := VagrantSetup(testPlatform)

		tests := TestCases.Get("LinuxKernelParameter", testPlatform)

		for _, testcase := range tests {
			linuxKernelParameter, err := sshClient.LinuxKernelParameter(testcase.resourceIdentifier)
			if err != nil {
				log.Fatalf("Error in %v: %v", testcase.resourceIdentifier, err)
			}

			t.Run(testcase.testName, func(t *testing.T) {
				if !reflect.DeepEqual(linuxKernelParameter, testcase.resourceWant) {
					t.Errorf(
						"LinuxKernelParameter failed:\nwant:\t%v\ngot:\t%v",
						testcase.resourceWant,
						linuxKernelParameter,
					)
				}
			})
		}
	}
}
