package git_test

import (
	"os/exec"
	"testing"

	client "github.com/living-etc/sightseer.go/clients/git"
)

func TestNewGitClient(t *testing.T) {
	repoDir := t.TempDir()

	setupCommands := [][]string{
		{"git", "init"},
		{"touch", "README.md"},
		{"git", "add", "README.md"},
		{"git", "commit", "-m", "Initial Commit"},
	}

	for _, setupCommand := range setupCommands {
		cmd := exec.Command(setupCommand[0], setupCommand[1:]...)
		cmd.Dir = repoDir

		err := cmd.Run()
		if err != nil {
			t.Fatal(err)
		}
	}

	t.Run("new client", func(t *testing.T) {
		gitClient := client.NewGitClient(repoDir)

		head, _ := gitClient.Ref("main")

		if head.Subject != "Initial Commit" {
			t.Errorf("want %v, got %v,", "Initial Commit", head.Subject)
		}
	})
}
