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
	case "Package":
		return TestCases.Package(platform)
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

func (testCases) Package(platform string) []TestCase {
	switch platform {
	case "ubuntu2404":
		return []TestCase{
			{
				testName:           "Package is installed",
				resourceIdentifier: "openssh-server",
				resourceWant: &Package{
					Name:          "openssh-server",
					Status:        "install ok installed",
					Priority:      "optional",
					Section:       "net",
					InstalledSize: "2140",
					Maintainer:    "Ubuntu Developers <ubuntu-devel-discuss@lists.ubuntu.com>",
					Architecture:  "arm64",
					MultiArch:     "foreign",
					Source:        "openssh",
					Version:       "1:9.6p1-3ubuntu13",
					Replaces:      "openssh-client (<< 1:7.9p1-8), ssh, ssh-krb5",
					Provides:      "ssh-server",
					Depends:       "adduser, libpam-modules, libpam-runtime, lsb-base, openssh-client (= 1:9.6p1-3ubuntu13), openssh-sftp-server, procps, ucf, debconf (>= 0.5) | debconf-2.0, libaudit1 (>= 1:2.2.1), libc6 (>= 2.38), libcom-err2 (>= 1.43.9), libcrypt1 (>= 1:4.1.0), libgssapi-krb5-2 (>= 1.17), libkrb5-3 (>= 1.13~alpha1+dfsg), libpam0g (>= 0.99.7.1), libselinux1 (>= 3.1~), libssl3t64 (>= 3.0.13), libwrap0 (>= 7.6-4~), zlib1g (>= 1:1.1.4)",
					PreDepends:    "init-system-helpers (>= 1.54~)",
					Recommends:    "default-logind | logind | libpam-systemd, ncurses-term, xauth, ssh-import-id",
					Suggests:      "molly-guard, monkeysphere, ssh-askpass, ufw",
					Conflicts:     "sftp, ssh-socks, ssh2",
					Conffiles: `
/etc/default/ssh 500e3cf069fe9a7b9936108eb9d9c035
/etc/init.d/ssh 3649a6fe8c18ad1d5245fd91737de507
/etc/pam.d/sshd 8b4c7a12b031424b2a9946881da59812
/etc/ssh/moduli 366395e79244c54223455e5f83dafba3
/etc/ufw/applications.d/openssh-server 486b78d54b93cc9fdc950c1d52ff479e`,
					Description: `secure shell (SSH) server, for secure access from remote machines
This is the portable version of OpenSSH, a free implementation of
the Secure Shell protocol as specified by the IETF secsh working
group.
.
Ssh (Secure Shell) is a program for logging into a remote machine
and for executing commands on a remote machine.
It provides secure encrypted communications between two untrusted
hosts over an insecure network. X11 connections and arbitrary TCP/IP
ports can also be forwarded over the secure channel.
It can be used to provide applications with a secure communication
channel.
.
This package provides the sshd server.
.
In some countries it may be illegal to use any encryption at all
without a special permit.
.
sshd replaces the insecure rshd program, which is obsolete for most
purposes.`,
					Homepage:           "https://www.openssh.com/",
					OriginalMaintainer: "Debian OpenSSH Maintainers <debian-ssh@lists.debian.org>",
				},
			},
		}
	case "fedora40":
		return []TestCase{}
	default:
		return []TestCase{}
	}
}
