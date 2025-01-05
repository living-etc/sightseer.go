package git

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	git "github.com/living-etc/sightseer.go/entities/git"
)

type GitClient struct {
	path string
}

func NewGitClient(path string) GitClient {
	return GitClient{path: path}
}

func (client *GitClient) IsValidRepo() error {
	cmdParts := []string{
		"git",
		"git-revparse",
		"--is-inside-working-tree",
	}
	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
	cmd.Dir = client.path

	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return errors.New(stderr.String())
	}
	return nil
}

func (client *GitClient) Commit(identifier string) (*git.Commit, error) {
	formatParts := []string{
		"{\"Hash\":\"%H\"",
		"\"AbbreviatedHash\":\"%h\"",
		"\"TreeHash\":\"%T\"",
		"\"AbbreviatedTreeHash\":\"%t\"",
		"\"ParentHashes\":\"%P\"",
		"\"AbbreviatedTree\":\"%p\"",
		"\"AuthorName\":\"%an\"",
		"\"AuthorNameMailmap\":\"%aN\"",
		"\"AuthorEmail\":\"%ae\"",
		"\"AuthorEmailMailmap\":\"%aE\"",
		"\"AuthorEmailLocalPart\":\"%al\"",
		"\"AuthorEmailLocalPartMailmap\":\"%aL\"",
		"\"AuthorDate\":\"%ad\"",
		"\"AuthorDateRFC2822\":\"%aD\"",
		"\"AuthorDateUNIX\":\"%at\"",
		"\"AuthorDateISO8601Like\":\"%ai\"",
		"\"AuthorDateISO8601Strict\":\"%aI\"",
		"\"AuthorDateShort\":\"%as\"",
		"\"AuthorDateHuman\":\"%ah\"",
		"\"CommitterName\":\"%cn\"",
		"\"CommitterNameMailmap\":\"%cN\"",
		"\"CommitterEmail\":\"%ce\"",
		"\"CommitterEmailMailmap\":\"%cE\"",
		"\"CommitterEmailLocalPart\":\"%cl\"",
		"\"CommitterEmailLocalPartMailmap\":\"%cL\"",
		"\"CommitterDate\":\"%cd\"",
		"\"CommitterDateRFC2822\":\"%cD\"",
		"\"CommitterDateUNIX\":\"%ct\"",
		"\"CommitterDateISO8601Like\":\"%ci\"",
		"\"CommitterDateISO8601Strict\":\"%cI\"",
		"\"CommitterDateShort\":\"%cs\"",
		"\"RefNames\":\"%d\"",
		"\"RefNamesUnwrapped\":\"%D\"",
		"\"Encoding\":\"%e\"",
		"\"Subject\":\"%s\"",
		"\"SubjectSanitised\":\"%f\"",
		"\"Body\":\"%b\"",
		"\"CommitNotes\":\"%N\"",
		"\"VerificationMessageRaw\":\"%GG\"",
		"\"SignerName\":\"%GS\"",
		"\"SigningKey\":\"%GK\"",
		"\"SigningKeyFingerprint\":\"%GF\"",
		"\"SigningKeyTrustLevel\":\"%GT\"",
		"\"ReflogSelector\":\"%gD\"",
		"\"ReflogSelectorShortened\":\"%gd\"",
		"\"ReflogIdentityName\":\"%gn\"",
		"\"ReflogIdentityNameMailmap\":\"%gN\"",
		"\"ReflogIdentityEmail\":\"%ge\"",
		"\"ReflogIdentityEmailMailmap\":\"%gE\"",
		"\"ReflogSubject\":\"%gs\"}",
	}

	cmdParts := []string{
		"git",
		"show",
		"--no-patch",
		fmt.Sprintf("--format=%s", strings.Join(formatParts, ",")),
		identifier,
	}
	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
	cmd.Dir = client.path

	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, errors.New(stderr.String())
	}

	commit := &git.Commit{}
	err = json.Unmarshal([]byte(stdout.String()), &commit)
	if err != nil {
		return nil, err
	}

	return commit, nil
}
