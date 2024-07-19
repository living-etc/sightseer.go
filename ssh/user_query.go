package ssh

import (
	"strconv"
)

type UserQuery struct{}

func (query UserQuery) Command() string {
	return "grep -e ^%v: /etc/passwd"
}

func (query UserQuery) Regex() string {
	return `(?P<Username>\w*?):(?P<Password>.*?):(?P<Uid>[0-9]*?):(?P<Gid>[0-9]*?):(?P<GECOS>.*?):(?P<HomeDirectory>.*?):(?P<Shell>.*?$)`
}

func (query UserQuery) SetValues(values map[string]string) (*User, error) {
	user := &User{}

	user.Username = values["Username"]

	uid, err := strconv.Atoi(values["Uid"])
	if err != nil {
		uid = -1
	}
	user.Uid = uid

	gid, err := strconv.Atoi(values["Gid"])
	if err != nil {
		gid = -1
	}
	user.Gid = gid

	user.HomeDirectory = values["HomeDirectory"]
	user.Shell = values["Shell"]

	return user, nil
}
