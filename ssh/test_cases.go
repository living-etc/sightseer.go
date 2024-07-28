package ssh

import "time"

type testCases struct{}

var TestCases testCases

type TestCase struct {
	testName           string
	resourceIdentifier string
	resourceWant       interface{}
	errWant            error
}

func (testCases) Get(resourceType string, platform string) []TestCase {
	switch resourceType {
	case "SystemdTimer":
		return TestCases.SystemdTimer(platform)
	case "LinuxKernelParameter":
		return TestCases.LinuxKernelParameter(platform)
	case "User":
		return TestCases.User(platform)
	case "Service":
		return TestCases.Service(platform)
	case "File":
		return TestCases.File(platform)
	default:
		return []TestCase{}
	}
}

type FileError struct {
	ErrorReason string
}

func (err FileError) Error() string { return err.ErrorReason }

func (testCases) File(platform string) []TestCase {
	switch platform {
	case "ubuntu2404":
		return []TestCase{
			{
				testName:           "File exists",
				resourceIdentifier: "/home/vagrant/.bashrc",
				resourceWant: &File{
					Type:          "regular file",
					OwnerID:       1000,
					OwnerName:     "vagrant",
					GroupID:       1000,
					GroupName:     "vagrant",
					SizeBytes:     3771,
					Name:          "/home/vagrant/.bashrc",
					MountPoint:    "/",
					InodeNumber:   1835013,
					NoOfHardLinks: 1,
					Mode:          "644",
				},
				errWant: nil,
			},
			{
				testName:           "File doesn't exist",
				resourceIdentifier: "/home/vagrant/.bashrc.doesnt.exist",
				resourceWant:       nil,
				errWant: &FileError{
					ErrorReason: "No such file or directory",
				},
			},
		}
	case "fedora40":
		return []TestCase{
			{
				testName:           "File exists",
				resourceIdentifier: "/home/vagrant/.bashrc",
				resourceWant: &File{
					Type:          "regular file",
					OwnerID:       1000,
					OwnerName:     "vagrant",
					GroupID:       1000,
					GroupName:     "vagrant",
					SizeBytes:     522,
					Name:          "/home/vagrant/.bashrc",
					MountPoint:    "/",
					InodeNumber:   17337356,
					NoOfHardLinks: 1,
					Mode:          "644",
				},
				errWant: nil,
			},
			{
				testName:           "File doesn't exist",
				resourceIdentifier: "/home/vagrant/.bashrc.doesnt.exist",
				resourceWant:       nil,
				errWant: &FileError{
					ErrorReason: "No such file or directory",
				},
			},
		}
	default:
		return []TestCase{}
	}
}

func (testCases) Service(platform string) []TestCase {
	switch platform {
	case "ubuntu2404":
		return []TestCase{
			{
				testName:           "hello-world service is running",
				resourceIdentifier: "hello-world.service",
				resourceWant: &Service{
					Description:    "Simple service for testing against",
					LoadState:      "loaded",
					UnitFileState:  "enabled",
					UnitFilePreset: "enabled",
					ActiveState:    "active",
				},
				errWant: nil,
			},
		}
	case "fedora40":
		return []TestCase{
			{
				testName:           "hello-world service is running",
				resourceIdentifier: "hello-world.service",
				resourceWant: &Service{
					Description:    "Simple service for testing against",
					LoadState:      "loaded",
					UnitFileState:  "enabled",
					UnitFilePreset: "disabled",
					ActiveState:    "active",
				},
				errWant: nil,
			},
		}
	default:
		return []TestCase{}
	}
}

func (testCases) User(platform string) []TestCase {
	switch platform {
	case "ubuntu2404", "fedora40":
		return []TestCase{
			{
				testName:           "User vagrant exists",
				resourceIdentifier: "vagrant",
				resourceWant: &User{
					Username:      "vagrant",
					Uid:           1000,
					Gid:           1000,
					HomeDirectory: "/home/vagrant",
					Shell:         "/bin/bash",
				},
				errWant: nil,
			},
		}
	default:
		return []TestCase{}
	}
}

func (testCases) SystemdTimer(platform string) []TestCase {
	switch platform {
	case "ubuntu2404":
		return []TestCase{
			{
				testName:           "Logrotate timer exists",
				resourceIdentifier: "logrotate.timer",
				resourceWant: &SystemdTimer{
					Id:             "logrotate.timer",
					Description:    "Daily rotation of log files",
					LoadState:      "loaded",
					UnitFileState:  "enabled",
					UnitFilePreset: "enabled",
					ActiveState:    "active",
					NextTrigger: time.Date(
						time.Now().Year(),
						time.Now().Month(),
						time.Now().Day()+1,
						0,
						0,
						0,
						0,
						time.UTC,
					),
					Triggers: "logrotate.service",
				},
				errWant: nil,
			},
		}
	case "fedora40":
		return []TestCase{
			{
				testName:           "Logrotate timer does not exist",
				resourceIdentifier: "logrotate.timer",
				resourceWant:       nil,
				errWant: &SystemdLoadError{
					UnitName:  "logrotate.timer",
					LoadState: "not-found",
					LoadError: `org.freedesktop.systemd1.NoSuchUnit "Unit logrotate.timer not found."`,
				},
			},
		}
	default:
		return []TestCase{}
	}
}

func (testCases) LinuxKernelParameter(platform string) []TestCase {
	switch platform {
	case "ubuntu2404", "fedora40":
		return []TestCase{
			{
				testName:           "Get a parameter",
				resourceIdentifier: "vm.page_lock_unfairness",
				resourceWant: &LinuxKernelParameter{
					Value: "5",
				},
				errWant: nil,
			},
		}
	default:
		return []TestCase{}
	}
}
