package ssh

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

func VagrantSetup() []*SshClient {
	testMachines := []struct {
		ip   string
		name string
	}{
		{
			name: "ubuntu2404",
			ip:   "172.16.44.180",
		},
		//{
		//	name: "fedora40",
		//	ip:   "172.16.44.181",
		//},
	}

	var sshClients []*SshClient

	for _, machine := range testMachines {
		privateKey, err := os.ReadFile(
			fmt.Sprintf(".vagrant/machines/%v/vmware_desktop/private_key", machine.name),
		)
		if err != nil {
			log.Fatalf("Error reading private key for %v", machine.name)
		}

		sshclient, err := NewSshClient(privateKey, machine.ip, "vagrant")
		if err != nil {
			log.Fatalf("Error creating ssh client: %v", err)
		}

		sshClients = append(sshClients, sshclient)
	}

	return sshClients
}

func TestFile(t *testing.T) {
	tests := []struct {
		testName string
		path     string
		fileWant *File
	}{
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

	sshclients := VagrantSetup()

	for _, testcase := range tests {
		for _, sshclient := range sshclients {
			file, err := sshclient.File(testcase.path)

			t.Run(testcase.testName, func(t *testing.T) {
				if err != nil {
					containsErrorString := strings.Contains(
						err.Error(),
						"No such file or directory",
					)

					if !containsErrorString {
						t.Errorf(
							"Error failed:\nwanted err to contain:\t%v\ngot:\t\t\t%v",
							"No such file or directory",
							err.Error(),
						)
					}
				} else {
					if reflect.DeepEqual(file, testcase.fileWant) {
						t.Errorf("File failed:\nwant\t%v\ngot\t%v", testcase.fileWant, file)
					}
				}
			})
		}
	}
}

func TestService(t *testing.T) {
	tests := []struct {
		name        string
		serviceWant *Service
		serviceName string
	}{
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

	sshclients := VagrantSetup()

	for _, testcase := range tests {
		for _, sshclient := range sshclients {
			service, err := sshclient.Service(testcase.serviceName)
			if err != nil {
				log.Fatalf("Error in %v: %v", testcase.name, err)
			}

			t.Run(testcase.serviceName, func(t *testing.T) {
				if !reflect.DeepEqual(service, testcase.serviceWant) {
					t.Errorf("Service failed:\nwant:\t%v\ngot\t\t%v", testcase.serviceWant, service)
				}
			})
		}
	}
}

func TestUser(t *testing.T) {
	tests := []struct {
		testName string
		username string
		userWant *User
	}{
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

	sshclients := VagrantSetup()

	for _, testcase := range tests {
		for _, sshclient := range sshclients {
			user, err := sshclient.User("vagrant")
			if err != nil {
				log.Fatalf("Error in %v: %v", testcase.testName, err)
			}

			t.Run(testcase.testName, func(t *testing.T) {
				if !reflect.DeepEqual(user, testcase.userWant) {
					t.Fatalf("User failed:\nwant:\t%v\ngot:\t%v", testcase.userWant, user)
				}
			})
		}
	}
}

func TestSystemdTimer(t *testing.T) {
	tests := []struct {
		name      string
		timerWant *SystemdTimer
	}{
		{
			name: "Timer",
			timerWant: &SystemdTimer{
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
		},
	}

	sshclients := VagrantSetup()

	for _, testcase := range tests {
		for _, sshclient := range sshclients {
			timer, err := sshclient.SystemdTimer("logrotate.timer")
			if err != nil {
				log.Fatalf("Error in %v: %v", testcase.name, err)
			}

			t.Run(testcase.name, func(t *testing.T) {
				if !reflect.DeepEqual(timer, testcase.timerWant) {
					t.Fatalf("SystemdTimer failed:\nwant:\t%v\ngot:\t%v", testcase.timerWant, timer)
				}
			})
		}
	}
}

func TestLinuxKernelParameter(t *testing.T) {
	tests := []struct {
		name          string
		parameterName string
		parameterWant *LinuxKernelParameter
	}{
		{
			name:          "Get a parameter",
			parameterName: "vm.page_lock_unfairness",
			parameterWant: &LinuxKernelParameter{
				Value: "5",
			},
		},
	}

	sshClients := VagrantSetup()

	for _, testcase := range tests {
		for _, sshclient := range sshClients {
			linuxKernelParameter, err := sshclient.LinuxKernelParameter(testcase.parameterName)
			if err != nil {
				log.Fatalf("Error in %v: %v", testcase.name, err)
			}

			t.Run(testcase.name, func(t *testing.T) {
				if !reflect.DeepEqual(linuxKernelParameter, testcase.parameterWant) {
					t.Errorf(
						"LinuxKernelParameter failed:\nwant:\t%v\ngot:\t%v",
						testcase.parameterWant,
						linuxKernelParameter,
					)
				}
			})
		}
	}
}
