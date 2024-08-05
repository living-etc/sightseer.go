package ssh_test

import (
	"log"
	"testing"

	sightseer "github.com/living-etc/sightseer.go/ssh"
)

func Test_fileFromStatOutput(t *testing.T) {
	tests := []struct {
		name       string
		statoutput string
		fileWant   *sightseer.File
	}{
		{
			name: "Regular file",
			statoutput: `Type=regular file
GroupID=1001
GroupName=vagrant
Permissions=644
OwnerID=1001
OwnerName=vagrant
SizeBytes=3771
Name=.bashrc
MountPoint=/
InodeNumber=1835013
NoOfHardLinks=1
Mode=644
`,
			fileWant: &sightseer.File{
				Type:          "regular file",
				OwnerID:       1001,
				OwnerName:     "vagrant",
				GroupID:       1001,
				GroupName:     "vagrant",
				SizeBytes:     3771,
				Name:          ".bashrc",
				MountPoint:    "/",
				InodeNumber:   1835013,
				NoOfHardLinks: 1,
				Mode:          "644",
			},
		},
	}

	t.Run("test file from stat", func(t *testing.T) {
		for _, tt := range tests {
			var fileQuery sightseer.FileQuery
			file, err := fileQuery.ParseOutput(tt.statoutput)
			if err != nil {
				log.Fatalf("Error in %v: %v", tt.name, err)
			}

			EvaluateStructTypesAreEqual(file, tt.fileWant, tt.name, t)
		}
	})
}
