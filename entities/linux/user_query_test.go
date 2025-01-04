package linux_test

import (
	"log"
	"testing"

	"github.com/living-etc/sightseer.go/entities/linux"
)

func Test_userFromPasswd(t *testing.T) {
	tests := []struct {
		testName     string
		passwdOutput string
		userWant     *linux.User
	}{
		{
			testName:     "User vagrant exists",
			passwdOutput: `vagrant:x:1000:1000:vagrant:/home/vagrant:/bin/bash`,
			userWant: &linux.User{
				Username:      "vagrant",
				Uid:           1000,
				Gid:           1000,
				HomeDirectory: "/home/vagrant",
				Shell:         "/bin/bash",
			},
		},
		{
			testName:     "empty entry",
			passwdOutput: `::::::`,
			userWant: &linux.User{
				Username:      "",
				Uid:           -1,
				Gid:           -1,
				HomeDirectory: "",
				Shell:         "",
			},
		},
		{
			testName:     "special characters in the entries",
			passwdOutput: `dhcpcd:x:100:65534:DHCP Client Daemon,,,:/usr/lib/dhcpcd:/bin/false`,
			userWant: &linux.User{
				Username:      "dhcpcd",
				Uid:           100,
				Gid:           65534,
				HomeDirectory: "/usr/lib/dhcpcd",
				Shell:         "/bin/false",
			},
		},
	}

	for _, testcase := range tests {
		var userQuery linux.UserQuery
		user, err := userQuery.ParseOutput(testcase.passwdOutput)
		if err != nil {
			log.Fatalf("Error in %v: %v", testcase.testName, err)
		}

		EvaluateStructTypesAreEqual(user, testcase.userWant, testcase.testName, t)
	}
}
