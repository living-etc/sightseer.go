package git

import (
	"os/exec"
	"strings"
)

type GitClient struct {
	path string
}

func NewGitClient(path string) GitClient {
	return GitClient{path: path}
}

type Ref struct {
	Subject string
}

func (client *GitClient) Ref(identifier string) (*Ref, error) {
	cmdParts := []string{"git", "show", "-s", "--format=%s", identifier}
	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
	cmd.Dir = client.path

	var out strings.Builder
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	return &Ref{
		Subject: strings.Trim(out.String(), "\n"),
	}, nil
}
