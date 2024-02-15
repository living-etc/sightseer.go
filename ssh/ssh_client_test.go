package ssh

import (
	"testing"
)

const statOutput = `  File: /usr/local/bin/kube-proxy
  Size: 43130880  	Blocks: 84240      IO Block: 4096   regular file
Device: 801h/2049d	Inode: 258242      Links: 1
Access: (0755/-rwxr-xr-x)  Uid: (    0/    root)   Gid: (    0/    root)
Access: 2024-02-13 18:28:13.260372630 +0000
Modify: 2024-02-13 18:28:13.196371385 +0000
Change: 2024-02-13 18:28:13.256372552 +0000
 Birth: 2024-02-13 18:28:11.628340885 +0000
`

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
		group      string
		mode       string
		statoutput string
	}{
		{
			path:       "/usr/local/bin/kube-proxy",
			owner:      "root",
			group:      "root",
			mode:       "0755",
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
