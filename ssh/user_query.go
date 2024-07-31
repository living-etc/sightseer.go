package ssh

import (
	"strconv"
	"strings"
)

type UserQuery struct{}

func (query UserQuery) Command(platform string) (string, error) {
	var cmd string

	switch platform {
	default:
		cmd = "grep -e ^%v: /etc/passwd"
	}

	return cmd, nil
}

func (query UserQuery) ParseOutput(output string) (*User, error) {
	user := &User{}

	parts := strings.Split(output, ":")

	user.Username = parts[0]

	uid, err := strconv.Atoi(parts[2])
	if err != nil {
		user.Uid = -1
	} else {
		user.Uid = uid
	}

	gid, err := strconv.Atoi(parts[3])
	if err != nil {
		user.Gid = -1
	} else {
		user.Gid = gid
	}

	user.HomeDirectory = parts[5]
	user.Shell = parts[6]

	return user, nil
}
