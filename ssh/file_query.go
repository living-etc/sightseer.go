package ssh

type FileQuery struct{}

func (query FileQuery) Regex() string {
	return `Access: \((?P<Mode>\d+)\/.*?\)\s+Uid: \(\s+(?P<OwnerId>\d+)\/\s+(?P<OwnerName>\w+)\)\s+Gid: \(\s+(?P<GroupId>\d+)\/\s+(?P<GroupName>\w+)\)`
}

func (query FileQuery) Command() string {
	return "stat %v"
}

func (query FileQuery) SetValues(file *File, values map[string]string) {
	file.OwnerName = values["OwnerName"]
	file.OwnerId = values["OwnerId"]
	file.GroupName = values["GroupName"]
	file.GroupId = values["GroupId"]
	file.Mode = values["Mode"]
}
