package ssh

import (
	"log"
	"testing"
)

func Test_userFromPasswd(t *testing.T) {
	tests := []struct {
		testName          string
		passwdOutput      string
		usernameWant      string
		uidWant           int
		gidWant           int
		homeDirectoryWant string
		shellWant         string
	}{
		{
			testName:          "User vagrant exists",
			passwdOutput:      `vagrant:x:1000:1000:vagrant:/home/vagrant:/bin/bash`,
			usernameWant:      "vagrant",
			uidWant:           1000,
			gidWant:           1000,
			homeDirectoryWant: "/home/vagrant",
			shellWant:         "/bin/bash",
		},
		{
			testName:          "empty entry",
			passwdOutput:      `::::::`,
			usernameWant:      "",
			uidWant:           0,
			gidWant:           0,
			homeDirectoryWant: "",
			shellWant:         "",
		},
		{
			testName:          "special characters in the entries",
			passwdOutput:      `dhcpcd:x:100:65534:DHCP Client Daemon,,,:/usr/lib/dhcpcd:/bin/false`,
			usernameWant:      "dhcpcd",
			uidWant:           100,
			gidWant:           65534,
			homeDirectoryWant: "/usr/lib/dhcpcd",
			shellWant:         "/bin/false",
		},
	}

	for _, tt := range tests {
		user, err := parseOutput[User, UserQuery](tt.passwdOutput)
		if err != nil {
			log.Fatal(err)
		}

		t.Run(tt.testName, func(t *testing.T) {
			if user.Username != tt.usernameWant {
				t.Errorf("want %v, got %v", tt.usernameWant, user.Username)
			}

			if user.Uid != tt.uidWant {
				t.Errorf("want %v, got %v", tt.uidWant, user.Uid)
			}

			if user.Gid != tt.gidWant {
				t.Errorf("want %v, got %v", tt.gidWant, user.Gid)
			}

			if user.HomeDirectory != tt.homeDirectoryWant {
				t.Errorf("want %v, got %v", tt.homeDirectoryWant, user.HomeDirectory)
			}

			if user.Shell != tt.shellWant {
				t.Errorf("want %v, got %v", tt.shellWant, user.Shell)
			}
		})
	}
}
