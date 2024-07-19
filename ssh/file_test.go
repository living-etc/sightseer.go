package ssh

import (
	"log"
	"testing"
)

func Test_fileFromStatOutput(t *testing.T) {
	tests := []struct {
		name          string
		statoutput    string
		ownerNameWant string
		ownerIdWant   string
		groupNameWant string
		groupIdWant   string
		modeWant      string
	}{
		{
			name: "Kube proxy binary",
			statoutput: `  File: /usr/local/bin/kube-proxy
  Size: 43130880  	Blocks: 84240      IO Block: 4096   regular file
Device: 801h/2049d	Inode: 258242      Links: 1
Access: (0755/-rwxr-xr-x)  Uid: (    0/    root)   Gid: (    0/    root)
Access: 2024-02-13 18:28:13.260372630 +0000
Modify: 2024-02-13 18:28:13.196371385 +0000
Change: 2024-02-13 18:28:13.256372552 +0000
 Birth: 2024-02-13 18:28:11.628340885 +0000
`,
			ownerNameWant: "root",
			ownerIdWant:   "0",
			groupNameWant: "root",
			groupIdWant:   "0",
			modeWant:      "0755",
		},
		{
			name: "admin.kubeconfig",
			statoutput: `  File: admin.kubeconfig
  Size: 6261      	Blocks: 16         IO Block: 4096   regular file
Device: 801h/2049d	Inode: 258140      Links: 1
Access: (0644/-rw-r--r--)  Uid: ( 1000/  ubuntu)   Gid: ( 1000/  ubuntu)
Access: 2024-02-15 14:05:30.314397696 +0000
Modify: 2024-02-15 14:05:29.518391511 +0000
Change: 2024-02-15 14:05:30.314397696 +0000
 Birth: 2024-02-15 14:05:29.490391294 +0000
`,
			ownerNameWant: "ubuntu",
			ownerIdWant:   "1000",
			groupNameWant: "ubuntu",
			groupIdWant:   "1000",
			modeWant:      "0644",
		},
	}

	t.Run("test file from stat", func(t *testing.T) {
		for _, tt := range tests {
			file, err := parseOutput[File, FileQuery](tt.statoutput)
			if err != nil {
				log.Fatalf("Error in %v: %v", tt.name, err)
			}

			ownerWant := tt.ownerNameWant
			ownerGot := file.OwnerName
			if ownerGot != ownerWant {
				t.Errorf("want %v, got %v", ownerWant, ownerGot)
			}

			groupWant := tt.groupNameWant
			groupGot := file.GroupName
			if groupGot != groupWant {
				t.Errorf("want %v, got %v", groupWant, groupGot)
			}

			modeWant := tt.modeWant
			modeGot := file.Mode
			if modeGot != modeWant {
				t.Errorf("want %v, got %v", modeWant, modeGot)
			}
		}
	})
}
