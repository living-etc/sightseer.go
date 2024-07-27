package ssh

import "time"

type testCases struct{}

var TestCases testCases

type FileTestCase struct {
	testName string
	path     string
	fileWant *File
}

func (testCases) File(platform string) []FileTestCase {
	switch platform {
	case "ubuntu2404", "fedora40":
		return []FileTestCase{
			{
				testName: "File exists",
				path:     "/home/vagrant/.bashrc",
				fileWant: &File{
					Type:          "regular file",
					OwnerID:       1001,
					OwnerName:     "vagrant",
					GroupID:       1001,
					GroupName:     "vagrant",
					SizeBytes:     3771,
					Name:          ".bashrc",
					MountPoint:    "/",
					InodeNumber:   1835013,
					NoOfHardLinks: 1,
					Mode:          "644",
				},
			},
			{
				testName: "File doesn't exist",
				path:     "/home/vagrant/.bashrc.doesnt.exist",
				fileWant: nil,
			},
		}
	default:
		return []FileTestCase{}
	}
}

type ServiceTestCase struct {
	name        string
	serviceWant *Service
	serviceName string
}

func (testCases) Service(platform string) []ServiceTestCase {
	switch platform {
	case "ubuntu2404":
		return []ServiceTestCase{
			{
				name: "hello-world service is running",
				serviceWant: &Service{
					Description:    "Simple service for testing against",
					LoadState:      "loaded",
					UnitFileState:  "enabled",
					UnitFilePreset: "enabled",
					ActiveState:    "active",
				},
				serviceName: "hello-world.service",
			},
		}
	case "fedora40":
		return []ServiceTestCase{
			{
				name: "hello-world service is running",
				serviceWant: &Service{
					Description:    "Simple service for testing against",
					LoadState:      "loaded",
					UnitFileState:  "enabled",
					UnitFilePreset: "disabled",
					ActiveState:    "active",
				},
				serviceName: "hello-world.service",
			},
		}
	default:
		return []ServiceTestCase{}
	}
}

type UserTestCase struct {
	testName string
	username string
	userWant *User
}

func (testCases) User(platform string) []UserTestCase {
	switch platform {
	case "ubuntu2404", "fedora40":
		return []UserTestCase{
			{
				testName: "User vagrant exists",
				username: "vagrant",
				userWant: &User{
					Username:      "vagrant",
					Uid:           1000,
					Gid:           1000,
					HomeDirectory: "/home/vagrant",
					Shell:         "/bin/bash",
				},
			},
		}
	default:
		return []UserTestCase{}
	}
}

type SystemdTimerTestCase struct {
	testName  string
	timerName string
	timerWant *SystemdTimer
	errWant   error
}

func (testCases) SystemdTimer(platform string) []SystemdTimerTestCase {
	switch platform {
	case "ubuntu2404":
		return []SystemdTimerTestCase{
			{
				testName:  "Logrotate timer exists",
				timerName: "logrotate.timer",
				timerWant: &SystemdTimer{
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
		return []SystemdTimerTestCase{
			{
				testName:  "Logrotate timer does not exist",
				timerName: "logrotate.timer",
				timerWant: nil,
				errWant: &SystemdLoadError{
					UnitName:  "logrotate.timer",
					LoadState: "not-found",
					LoadError: `org.freedesktop.systemd1.NoSuchUnit "Unit logrotate.timer not found."`,
				},
			},
		}
	default:
		return []SystemdTimerTestCase{}
	}
}

type LinuxKernelParameterTestCase struct {
	name          string
	parameterName string
	parameterWant *LinuxKernelParameter
	errWant       error
}

func (testCases) LinuxKernelParameter(platform string) []LinuxKernelParameterTestCase {
	switch platform {
	case "ubuntu2404", "fedora40":
		return []LinuxKernelParameterTestCase{
			{
				name:          "Get a parameter",
				parameterName: "vm.page_lock_unfairness",
				parameterWant: &LinuxKernelParameter{
					Value: "5",
				},
			},
		}
	default:
		return []LinuxKernelParameterTestCase{}
	}
}
