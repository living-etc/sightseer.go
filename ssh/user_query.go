package ssh

import "strconv"

type UserQuery struct{}

func (query UserQuery) Command() string {
	return "grep -e ^%v: /etc/passwd"
}

func (query UserQuery) Regex() string {
	return `(?P<Username>\w*?):(?P<Password>.*?):(?P<Uid>[0-9]*?):(?P<Gid>[0-9]*?):(?P<GECOS>.*?):(?P<HomeDirectory>.*?):(?P<Shell>.*?$)`
}

func (query UserQuery) SetValues(user *User, values map[string]string) {
	user.Username = values["Username"]
	user.Uid, _ = strconv.Atoi(values["Uid"])
	user.Gid, _ = strconv.Atoi(values["Gid"])
	user.HomeDirectory = values["HomeDirectory"]
	user.Shell = values["Shell"]
}
