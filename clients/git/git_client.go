package git

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

type GitClient struct {
	path string
}

func NewGitClient(path string) GitClient {
	return GitClient{path: path}
}

type Commit struct {
	Hash                           string `json:"Hash"`
	AbbreviatedHash                string `json:"AbbreviatedHash"`
	TreeHash                       string `json:"TreeHash"`
	AbbreviatedTreeHash            string `json:"AbbreviatedTreeHash"`
	ParentHashes                   string `json:"ParentHashes"`
	AbbreviatedTree                string `json:"AbbreviatedTree"`
	AuthorName                     string `json:"AuthorName"`
	AuthorNameMailmap              string `json:"AuthorNameMailmap"`
	AuthorEmail                    string `json:"AuthorEmail"`
	AuthorEmailMailmap             string `json:"AuthorEmailMailmap"`
	AuthorEmailLocalPart           string `json:"AuthorEmailLocalPart"`
	AuthorEmailLocalPartMailmap    string `json:"AuthorEmailLocalPartMailmap"`
	AuthorDate                     string `json:"AuthorDate"`
	AuthorDateRFC2822              string `json:"AuthorDateRFC2822"`
	AuthorDateUNIX                 string `json:"AuthorDateUNIX"`
	AuthorDateISO8601Like          string `json:"AuthorDateISO8601Like"`
	AuthorDateISO8601Strict        string `json:"AuthorDateISO8601Strict"`
	AuthorDateShort                string `json:"AuthorDateShort"`
	CommitterName                  string `json:"CommitterName"`
	CommitterNameMailmap           string `json:"CommitterNameMailmap"`
	CommitterEmail                 string `json:"CommitterEmail"`
	CommitterEmailMailmap          string `json:"CommitterEmailMailmap"`
	CommitterEmailLocalPart        string `json:"CommitterEmailLocalPart"`
	CommitterEmailLocalPartMailmap string `json:"CommitterEmailLocalPartMailmap"`
	CommitterDate                  string `json:"CommitterDate"`
	CommitterDateRFC2822           string `json:"CommitterDateRFC2822"`
	CommitterDateUNIX              string `json:"CommitterDateUNIX"`
	CommitterDateISO8601Like       string `json:"CommitterDateISO8601Like"`
	CommitterDateISO8601Strict     string `json:"CommitterDateISO8601Strict"`
	CommitterDateShort             string `json:"CommitterDateShort"`
	RefNames                       string `json:"RefNames"`
	RefNamesUnwrapped              string `json:"RefNamesUnwrapped"`
	Encoding                       string `json:"Encoding"`
	Subject                        string `json:"Subject"`
	SubjectSanitised               string `json:"SubjectSanitised"`
	Body                           string `json:"Body"`
	CommitNotes                    string `json:"CommitNotes"`
	VerificationMessageRaw         string `json:"VerificationMessageRaw"`
	SignerName                     string `json:"SignerName"`
	SigningKey                     string `json:"SigningKey"`
	SigningKeyFingerprint          string `json:"SigningKeyFingerprint"`
	SigningKeyTrustLevel           string `json:"SigningKeyTrustLevel"`
	ReflogSelector                 string `json:"ReflogSelector"`
	ReflogSelectorShortened        string `json:"ReflogSelectorShortened"`
	ReflogIdentityName             string `json:"ReflogIdentityName"`
	ReflogIdentityNameMailmap      string `json:"ReflogIdentityNameMailmap"`
	ReflogIdentityEmail            string `json:"ReflogIdentityEmail"`
	ReflogIdentityEmailMailmap     string `json:"ReflogIdentityEmailMailmap"`
	ReflogSubject                  string `json:"ReflogSubject"`
}

func (client *GitClient) Commit(identifier string) (*Commit, error) {
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

	commit := &Commit{}
	err = json.Unmarshal([]byte(stdout.String()), &commit)
	if err != nil {
		return nil, err
	}

	return commit, nil
}
