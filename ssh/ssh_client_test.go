package ssh

import (
	"log"
	"os"
	"testing"
)

func VagrantSetup() *SshClient {
	privateKey, _ := os.ReadFile(".vagrant/machines/default/vmware_desktop/private_key")
	vmIP := "172.16.44.180"

	sshclient, err := NewSshClient(privateKey, vmIP, "vagrant")
	if err != nil {
		log.Fatal(err)
	}

	return sshclient
}

func TestFile(t *testing.T) {
	tests := []struct {
		path      string
		owner     string
		error     bool
		erroutput string
	}{
		{
			path:      "/home/vagrant/.bashrc",
			owner:     "vagrant",
			error:     false,
			erroutput: "",
		},
		{
			path:      "/home/vagrant/.bashrc.doesnt.exist",
			owner:     "vagrant",
			error:     true,
			erroutput: `[Error]: Process exited with status 1 [Context]: stat: cannot statx '/home/vagrant/.bashrc.doesnt.exist': No such file or directory`,
		},
	}

	sshclient := VagrantSetup()

	for _, testcase := range tests {
		file, err := Get(testcase.path, sshclient, FileQuery{})

		t.Run(testcase.path+" owner", func(t *testing.T) {
			if testcase.error {
				errWant := testcase.erroutput
				errGot := err.Error()

				if errGot != errWant {
					t.Errorf("want %v, got %v", errWant, errGot)
				}
			} else {
				if file.OwnerName != "vagrant" {
					t.Errorf("want %v, got %v", "vagrant", file.OwnerName)
				}
			}
		})
	}
}

func TestService(t *testing.T) {
	tests := []struct {
		systemctlOutput string
		activeWant      string
		serviceName     string
	}{
		{
			activeWant:  "active (running)",
			serviceName: "ssh",
		},
	}

	sshclient := VagrantSetup()

	for _, testcase := range tests {
		service, err := Get("ssh", sshclient, ServiceQuery{})
		if err != nil {
			log.Fatal(err)
		}

		t.Run(testcase.serviceName+": active", func(t *testing.T) {
			if service.Active != testcase.activeWant {
				t.Errorf("want %v, got %v", testcase.activeWant, service.Active)
			}
		})
	}
}
