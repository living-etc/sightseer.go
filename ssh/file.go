package ssh

type File struct {
	Type          string
	OwnerName     string
	OwnerID       int
	GroupName     string
	GroupID       int
	SizeBytes     int
	Name          string
	MountPoint    string
	InodeNumber   int
	NoOfHardLinks int
	Mode          string
}
