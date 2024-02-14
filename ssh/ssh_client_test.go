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
		Path  string
		Owner string
		Group string
		Mode  string
	}{
		{
			Path:  "/usr/local/bin/kube-proxy",
			Owner: "root",
			Group: "root",
			Mode:  "0755",
		},
	}

	mockExecutor := MockCommandExecutor{
		MockResponse: statOutput,
		MockError:    nil,
	}

	sshclient := SshClient{
		client:   nil,
		executor: &mockExecutor,
	}

	for _, testcase := range tests {
		file, _ := sshclient.File(testcase.Path)
		t.Run(testcase.Path, func(t *testing.T) {
		})

		_ = file
	}
}
