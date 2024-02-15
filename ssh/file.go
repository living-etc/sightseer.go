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

func fileFromStatOutput(statoutput string) (File, error) {
	accessPermissionsRegex := `Access: \((?P<mode>\d+)\/.*?\)\s+Uid: \(\s+(?P<uid>\d+)\/\s+(?P<uname>\w+)\)\s+Gid: \(\s+(?P<gid>\d+)\/\s+(?P<gname>\w+)\)`

	re := regexp.MustCompile(accessPermissionsRegex)
	matches := re.FindStringSubmatch(statoutput)
	result := make(map[string]string)

	for i, name := range re.SubexpNames() {
		if i > 0 {
			result[name] = matches[i]
		}
	}

	return File{
		OwnerName: result["uname"],
		OwnerId:   result["uid"],
		GroupName: result["gname"],
		GroupId:   result["gid"],
		Mode:      result["mode"],
	}, nil
}
