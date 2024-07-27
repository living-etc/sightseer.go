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
		tests := TestCases.File(testPlatform)

		sshClient := VagrantSetup(testPlatform)

		for _, testcase := range tests {
			file, err := sshClient.File(testcase.path)

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
					if reflect.DeepEqual(file, testcase.fileWant) {
						t.Errorf("File failed:\nwant\t%v\ngot\t%v", testcase.fileWant, file)
					}
				}
			})
		}
	}
}

func TestService(t *testing.T) {
	for _, testPlatform := range testPlatforms {
		tests := TestCases.Service(testPlatform)

		sshClient := VagrantSetup(testPlatform)

		for _, testcase := range tests {
			service, err := sshClient.Service(testcase.serviceName)
			if err != nil {
				log.Fatalf("Error in %v: %v", testcase.name, err)
			}

			t.Run(testcase.serviceName, func(t *testing.T) {
				if !reflect.DeepEqual(service, testcase.serviceWant) {
					t.Errorf("Service failed:\nwant:\t%v\ngot\t\t%v", testcase.serviceWant, service)
				}
			})
		}
	}
}

func TestUser(t *testing.T) {
	for _, testPlatform := range testPlatforms {
		tests := TestCases.User(testPlatform)

		sshClient := VagrantSetup(testPlatform)

		for _, testcase := range tests {
			user, err := sshClient.User(testcase.username)
			if err != nil {
				t.Fatalf("Error in %v: %v", testcase.testName, err)
			}

			t.Run(testcase.testName, func(t *testing.T) {
				if !reflect.DeepEqual(user, testcase.userWant) {
					t.Fatalf("User failed:\nwant:\t%v\ngot:\t%v", testcase.userWant, user)
				}
			})
		}
	}
}

func TestSystemdTimer(t *testing.T) {
	for _, testPlatform := range testPlatforms {
		sshClient := VagrantSetup(testPlatform)

		tests := TestCases.SystemdTimer(testPlatform)

		for _, testcase := range tests {
			timer, err := sshClient.SystemdTimer(testcase.timerName)
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
					if !reflect.DeepEqual(timer, testcase.timerWant) {
						t.Fatalf("SystemdTimer failed:\nwant:\t%v\ngot:\t%v", testcase.timerWant, timer)
					}
				})
			}
		}
	}
}

func TestLinuxKernelParameter(t *testing.T) {
	for _, testPlatform := range testPlatforms {
		sshClient := VagrantSetup(testPlatform)

		tests := TestCases.LinuxKernelParameter(testPlatform)

		for _, testcase := range tests {
			linuxKernelParameter, err := sshClient.LinuxKernelParameter(testcase.parameterName)
			if err != nil {
				log.Fatalf("Error in %v: %v", testcase.name, err)
			}

			t.Run(testcase.name, func(t *testing.T) {
				if !reflect.DeepEqual(linuxKernelParameter, testcase.parameterWant) {
					t.Errorf(
						"LinuxKernelParameter failed:\nwant:\t%v\ngot:\t%v",
						testcase.parameterWant,
						linuxKernelParameter,
					)
				}
			})
		}
	}
}
