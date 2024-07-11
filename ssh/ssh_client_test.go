package ssh

import (
	"errors"
	"fmt"
	"log"
	"os"
	"testing"
)

const (
	statOutput = `  File: /usr/local/bin/kube-proxy
  Size: 43130880  	Blocks: 84240      IO Block: 4096   regular file
Device: 801h/2049d	Inode: 258242      Links: 1
Access: (0755/-rwxr-xr-x)  Uid: (    0/    root)   Gid: (    0/    root)
Access: 2024-02-13 18:28:13.260372630 +0000
Modify: 2024-02-13 18:28:13.196371385 +0000
Change: 2024-02-13 18:28:13.256372552 +0000
 Birth: 2024-02-13 18:28:11.628340885 +0000
`
)

type MockCommandExecutor struct {
	MockResponse string
	MockError    error
}

func (m MockCommandExecutor) ExecuteCommand(command string) (string, error) {
	return m.MockResponse, m.MockError
}

func TestFile(t *testing.T) {
	tests := []struct {
		path       string
		owner      string
		statoutput string
	}{
		{
			path:       "/usr/local/bin/kube-proxy",
			owner:      "root",
			statoutput: statOutput,
		},
	}

	for _, testcase := range tests {
		mockExecutor := MockCommandExecutor{
			MockResponse: testcase.statoutput,
			MockError:    nil,
		}

		sshclient := SshClient{
			client:   nil,
			executor: &mockExecutor,
		}

		file, _ := sshclient.File(testcase.path)
		t.Run(testcase.path+" owner", func(t *testing.T) {
			if file.OwnerName != "root" {
				t.Errorf("want %v, got %v", "root", file.OwnerName)
			}
		})
	}
}

func TestFile_Error(t *testing.T) {
	tests := []struct {
		path       string
		statoutput string
		erroutput  string
	}{
		{
			path:       "/opt/cni/bin/bandwidth",
			statoutput: `stat: cannot statx '/opt/cni/bin/bandwidth': Permission denied`,
			erroutput:  "Process exited with status 1",
		},
	}

	for _, testcase := range tests {
		mockExecutor := MockCommandExecutor{
			MockResponse: testcase.statoutput,
			MockError:    errors.New(testcase.erroutput),
		}

		sshclient := SshClient{
			client:   nil,
			executor: &mockExecutor,
		}

		_, err := sshclient.File(testcase.path)
		t.Run(testcase.path+" error", func(t *testing.T) {
			errWant := fmt.Sprintf(
				"[Error]: %v [Context]: %v",
				testcase.erroutput,
				testcase.statoutput,
			)
			errGot := err.Error()
			if errGot != errWant {
				t.Errorf("want %v, got %v", errWant, errGot)
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

	privateKey, _ := os.ReadFile(".vagrant/machines/default/vmware_desktop/private_key")
	vmIP := "172.16.44.180"

	sshclient, err := NewSshClient(privateKey, vmIP, "vagrant")
	if err != nil {
		log.Fatal(err)
	}

	for _, testcase := range tests {
		service, err := sshclient.Service("ssh")
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
