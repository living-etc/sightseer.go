package ssh

import (
	"fmt"
	"log"
	"os"
	"reflect"
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

func InitSshClient(platform string) *SshClient {
	machine := testMachines[platform]

	privateKey, err := os.ReadFile(
		fmt.Sprintf(".vagrant/machines/%v/vmware_desktop/private_key", machine.name),
	)
	if err != nil {
		log.Fatalf("Error reading private key for %v", machine.name)
	}

	sshClient, err := NewSshClient(privateKey, machine.ip, machine.port, machine.user, platform)
	if err != nil {
		log.Fatalf("Error creating ssh client: %v", err)
	}

	return sshClient
}

//func VagrantUp(t *testing.T) {
//	t.Log("[Vagrant Up] Launching vagrant VMs for testing")
//
//	cmd := exec.Command("vagrant", "up")
//	if err := cmd.Run(); err != nil {
//		log.Fatal(err)
//	}
//}
//
//func VagrantHalt(t *testing.T) {
//	t.Log("[Vagrant Halt] Halting vagrant VMs after testing")
//
//	cmd := exec.Command("vagrant", "halt")
//	if err := cmd.Run(); err != nil {
//		log.Fatal(err)
//	}
//}

func TestSshClient(t *testing.T) {
	// VagrantUp(t)
	// defer VagrantHalt(t)

	EvaluateTestCases[File, error]("File", t)
	EvaluateTestCases[Service, error]("Service", t)
	EvaluateTestCases[User, error]("User", t)
	EvaluateTestCases[SystemdTimer, *SystemdLoadError]("SystemdTimer", t)
	EvaluateTestCases[LinuxKernelParameter, error]("LinuxKernelParameter", t)
}

func EvaluateTestCases[T ResourceType, E error](resourceType string, t *testing.T) {
	for _, testPlatform := range testPlatforms {
		sshClient := InitSshClient(testPlatform)

		tests := TestCases.Get(resourceType, testPlatform)

		for _, testcase := range tests {
			testName := fmt.Sprintf("%v/%v/%v", resourceType, testPlatform, testcase.testName)

			t.Run(testName, func(t *testing.T) {
				methodValue := reflect.ValueOf(sshClient).MethodByName(resourceType)
				if !methodValue.IsValid() {
					log.Fatal("Method not found")
				}

				args := []reflect.Value{reflect.ValueOf(testcase.resourceIdentifier)}
				results := methodValue.Call(args)

				resourceGot := results[0].Interface().(*T)

				if testcase.errWant != nil {
					var err E
					if results[1].Interface() != nil {
						err = results[1].Interface().(E)
					}

					if err.Error() != testcase.errWant.Error() {
						t.Fatalf("\nwant:\t%v\ngot:\t%v", testcase.errWant, err)
					}
				}

				if testcase.resourceWant != nil {
					if !reflect.DeepEqual(resourceGot, testcase.resourceWant) {
						t.Fatalf("\nwant:\t%v\ngot:\t%v", testcase.resourceWant, resourceGot)
					}
				}
			})
		}
	}
}
