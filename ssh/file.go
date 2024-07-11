package ssh

import (
	"regexp"
)

type File struct {
	OwnerName string
	OwnerId   string
	GroupName string
	GroupId   string
	Mode      string
}

func fileFromStatOutput(statoutput string) (*File, error) {
	accessPermissionsRegex := `Access: \((?P<Mode>\d+)\/.*?\)\s+Uid: \(\s+(?P<Uid>\d+)\/\s+(?P<Uname>\w+)\)\s+Gid: \(\s+(?P<Gid>\d+)\/\s+(?P<Gname>\w+)\)`

	re := regexp.MustCompile(accessPermissionsRegex)
	matches := re.FindStringSubmatch(statoutput)

	file := &File{
		OwnerName: matches[re.SubexpIndex("Uname")],
		OwnerId:   matches[re.SubexpIndex("Uid")],
		GroupName: matches[re.SubexpIndex("Gname")],
		GroupId:   matches[re.SubexpIndex("Gid")],
		Mode:      matches[re.SubexpIndex("Mode")],
	}

	return file, nil
}
