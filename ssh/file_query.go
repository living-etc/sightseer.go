package ssh

type FileQuery struct{}

func (query FileQuery) Command() string {
	return "stat %v"
}

func (query FileQuery) Regex() string {
	return `Access: \((?P<Mode>\d+)\/.*?\)\s+Uid: \(\s+(?P<OwnerId>\d+)\/\s+(?P<OwnerName>\w+)\)\s+Gid: \(\s+(?P<GroupId>\d+)\/\s+(?P<GroupName>\w+)\)`
}

func (query FileQuery) SetValues(values map[string]string) (*File, error) {
	return &File{
		OwnerName: values["OwnerName"],
		OwnerId:   values["OwnerId"],
		GroupName: values["GroupName"],
		GroupId:   values["GroupId"],
		Mode:      values["Mode"],
	}, nil
}
