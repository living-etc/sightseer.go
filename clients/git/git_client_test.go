package git_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"

	client "github.com/living-etc/sightseer.go/clients/git"
	git "github.com/living-etc/sightseer.go/entities/git"
)

func TestNewGitClient(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	repoDir := filepath.Join(wd + "/testdata/git-client-test-repo-for-sightseer")

	commitWant := &git.Commit{
		Hash:                           "108dd57ebe039833ad45f0e997fb68ad4fe49b6c",
		AbbreviatedHash:                "108dd57",
		TreeHash:                       "f93e3a1a1525fb5b91020da86e44810c87a2d7bc",
		AbbreviatedTreeHash:            "f93e3a1",
		ParentHashes:                   "",
		AbbreviatedTree:                "",
		AuthorName:                     "Sightseer Test User",
		AuthorNameMailmap:              "Sightseer Test User",
		AuthorEmail:                    "sightseer@example.com",
		AuthorEmailMailmap:             "sightseer@example.com",
		AuthorEmailLocalPart:           "sightseer",
		AuthorEmailLocalPartMailmap:    "sightseer",
		AuthorDate:                     "Sun Jan 5 16:47:26 2025 +0000",
		AuthorDateRFC2822:              "Sun, 5 Jan 2025 16:47:26 +0000",
		AuthorDateUNIX:                 "1736095646",
		AuthorDateISO8601Like:          "2025-01-05 16:47:26 +0000",
		AuthorDateISO8601Strict:        "2025-01-05T16:47:26+00:00",
		AuthorDateShort:                "2025-01-05",
		CommitterName:                  "Sightseer Test User",
		CommitterNameMailmap:           "Sightseer Test User",
		CommitterEmail:                 "sightseer@example.com",
		CommitterEmailMailmap:          "sightseer@example.com",
		CommitterEmailLocalPart:        "sightseer",
		CommitterEmailLocalPartMailmap: "sightseer",
		CommitterDate:                  "Sun Jan 5 16:47:26 2025 +0000",
		CommitterDateRFC2822:           "Sun, 5 Jan 2025 16:47:26 +0000",
		CommitterDateUNIX:              "1736095646",
		CommitterDateISO8601Like:       "2025-01-05 16:47:26 +0000",
		CommitterDateISO8601Strict:     "2025-01-05T16:47:26+00:00",
		CommitterDateShort:             "2025-01-05",
		RefNames:                       " (HEAD -> main, origin/main, origin/HEAD)",
		RefNamesUnwrapped:              "HEAD -> main, origin/main, origin/HEAD",
		Encoding:                       "",
		Subject:                        "Initial Commit",
		SubjectSanitised:               "Initial-Commit",
		Body:                           "",
		CommitNotes:                    "",
		VerificationMessageRaw:         "",
		SignerName:                     "",
		SigningKey:                     "",
		SigningKeyFingerprint:          "",
		SigningKeyTrustLevel:           "undefined",
		ReflogSelector:                 "",
		ReflogSelectorShortened:        "",
		ReflogIdentityName:             "",
		ReflogIdentityNameMailmap:      "",
		ReflogIdentityEmail:            "",
		ReflogIdentityEmailMailmap:     "",
		ReflogSubject:                  "",
	}

	t.Run("new client", func(t *testing.T) {
		gitClient := client.NewGitClient(repoDir)

		head, err := gitClient.Commit("main")
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(commitWant, head); diff != "" {
			t.Error(diff)
		}
	})

	t.Run("Check Valid Repo", func(t *testing.T) {
		gitClient := client.NewGitClient(repoDir)

		if err := gitClient.IsValidRepo(); err != nil {
			t.Fatalf("Not a valid git repo: %v", err)
		}
	})
}
