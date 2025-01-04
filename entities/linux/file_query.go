package linux

import (
	"bufio"
	"reflect"
	"strconv"
	"strings"
)

type FileQuery struct{}

func (query FileQuery) Command(platform string) (string, error) {
	var cmd string

	switch platform {
	default:
		cmd = `stat %v --printf="Type=%%F\nGroupID=%%g\nGroupName=%%G\nMode=%%a\nOwnerID=%%u\nOwnerName=%%U\nSizeBytes=%%s\nName=%%n\nMountPoint=%%m\nInodeNumber=%%i\nNoOfHardLinks=%%h\n"`
	}

	return cmd, nil
}

func (query FileQuery) ParseOutput(output string) (*File, error) {
	if strings.Contains(output, "No such file or directory") {
		return nil, &FileError{ErrorReason: "No such file or directory"}
	}
	file := File{}

	values := make(map[string]string)

	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		attribute := strings.SplitN(scanner.Text(), "=", 2)
		values[attribute[0]] = attribute[1]
	}

	file.Mode = values["Mode"]
	file.Type = values["Type"]
	file.Name = values["Name"]
	file.OwnerName = values["OwnerName"]
	file.GroupName = values["GroupName"]
	file.MountPoint = values["MountPoint"]

	intValueAttributes := []string{
		"OwnerID",
		"GroupID",
		"SizeBytes",
		"InodeNumber",
		"NoOfHardLinks",
	}

	for _, intValueAttribute := range intValueAttributes {
		value, err := strconv.Atoi(values[intValueAttribute])

		var valueToSet int
		if err != nil {
			valueToSet = -1
		} else {
			valueToSet = value
		}

		reflect.
			ValueOf(&file).
			Elem().
			FieldByName(intValueAttribute).
			SetInt(int64(valueToSet))
	}

	return &file, nil
}
