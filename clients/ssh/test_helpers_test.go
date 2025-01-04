package ssh_test

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	sightseer "github.com/living-etc/sightseer.go/ssh"
)

var testPlatforms = []string{
	"ubuntu2404",
	"fedora40",
}

type VagrantMachineConfig struct {
	ip   string
	name string
	port string
	user string
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

func EvaluateStructTypesAreEqual(got any, want any, testname string, t *testing.T) {
	diffs, _ := Diff.Structs(got, want)
	t.Run(testname, func(t *testing.T) {
		if len(diffs) > 0 {
			output := "Package mismatch"
			for _, diff := range diffs {
				output += fmt.Sprintf(
					"\n%v:\n\tgot\t%v\n\twant\t%v\n",
					diff.FieldName,
					diff.LhsValue,
					diff.RhsValue,
				)
			}
			t.Fatal(output)
		}
	})
}

func InitSshClient(platform string) *sightseer.SshClient {
	machine := testMachines[platform]

	privateKey, err := os.ReadFile(
		fmt.Sprintf(".vagrant/machines/%v/vmware_desktop/private_key", machine.name),
	)
	if err != nil {
		log.Fatalf("Error reading private key for %v", machine.name)
	}

	sshClient, err := sightseer.NewSshClient(
		privateKey,
		machine.ip,
		machine.port,
		machine.user,
		platform,
	)
	if err != nil {
		log.Fatalf("Error creating ssh client: %v", err)
	}

	return sshClient
}

func EvaluateTestCases[T sightseer.ResourceType, E error](resourceType string, t *testing.T) {
	for _, testPlatform := range testPlatforms {
		sshClient := InitSshClient(testPlatform)

		tests := TestCases.Get(resourceType, testPlatform)

		for _, testcase := range tests {
			testName := fmt.Sprintf("%v/%v/%v", resourceType, testPlatform, testcase.testName)

			methodValue := reflect.ValueOf(sshClient).MethodByName(resourceType)
			if !methodValue.IsValid() {
				log.Fatalf("Method '%v' not found on %v", resourceType, reflect.TypeOf(sshClient))
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
				EvaluateStructTypesAreEqual(
					resourceGot,
					testcase.resourceWant,
					testName,
					t,
				)
			}
		}
	}
}
