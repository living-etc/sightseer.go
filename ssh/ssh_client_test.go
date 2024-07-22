package ssh

import (
	"log"
	"os"
	"testing"
	"time"
)

func VagrantSetup() *SshClient {
	privateKey, _ := os.ReadFile(".vagrant/machines/default/vmware_desktop/private_key")
	vmIP := "172.16.44.180"

	sshclient, err := NewSshClient(privateKey, vmIP, "vagrant")
	if err != nil {
		log.Fatalf("Error creating ssh client: %v", err)
	}

	return sshclient
}

func TestFile(t *testing.T) {
	tests := []struct {
		path      string
		owner     string
		error     bool
		erroutput string
	}{
		{
			path:      "/home/vagrant/.bashrc",
			owner:     "vagrant",
			error:     false,
			erroutput: "",
		},
		{
			path:      "/home/vagrant/.bashrc.doesnt.exist",
			owner:     "vagrant",
			error:     true,
			erroutput: `[Error]: Process exited with status 1 [Context]: stat: cannot statx '/home/vagrant/.bashrc.doesnt.exist': No such file or directory`,
		},
	}

	sshclient := VagrantSetup()

	for _, testcase := range tests {
		file, err := sshclient.File(testcase.path)

		t.Run(testcase.path+" owner", func(t *testing.T) {
			if testcase.error {
				errWant := testcase.erroutput
				errGot := err.Error()

				if errGot != errWant {
					t.Errorf("want %v, got %v", errWant, errGot)
				}
			} else {
				if file.OwnerName != "vagrant" {
					t.Errorf("want %v, got %v", "vagrant", file.OwnerName)
				}
			}
		})
	}
}

func TestService(t *testing.T) {
	tests := []struct {
		name        string
		activeWant  string
		serviceName string
	}{
		{
			name:        "SSH service is running",
			activeWant:  "active (running)",
			serviceName: "ssh",
		},
	}

	sshclient := VagrantSetup()

	for _, testcase := range tests {
		service, err := sshclient.Service("ssh")
		if err != nil {
			log.Fatalf("Error in %v: %v", testcase.name, err)
		}

		t.Run(testcase.serviceName+": active", func(t *testing.T) {
			if service.Active != testcase.activeWant {
				t.Errorf("want %v, got %v", testcase.activeWant, service.Active)
			}
		})
	}
}

func TestUser(t *testing.T) {
	tests := []struct {
		testName          string
		passwdOutput      string
		usernameWant      string
		uidWant           int
		gidWant           int
		homeDirectoryWant string
		shellWant         string
	}{
		{
			testName:          "User vagrant exists",
			passwdOutput:      `vagrant:x:1000:1000:vagrant:/home/vagrant:/bin/bash`,
			usernameWant:      "vagrant",
			uidWant:           1000,
			gidWant:           1000,
			homeDirectoryWant: "/home/vagrant",
			shellWant:         "/bin/bash",
		},
	}

	sshclient := VagrantSetup()

	for _, testcase := range tests {
		user, err := sshclient.User("vagrant")
		if err != nil {
			log.Fatalf("Error in %v: %v", testcase.testName, err)
		}

		t.Run(testcase.testName, func(t *testing.T) {
			if user.Username != testcase.usernameWant {
				t.Fatalf("want %v, got %v", testcase.usernameWant, err)
			}
		})
	}
}

func TestSystemdTimer(t *testing.T) {
	tests := []struct {
		name            string
		systemctlOutput string
		descriptionWant string
		loadedWant      bool
		unitFileWant    string
		enabledWant     bool
		presetWant      bool
		activeWant      string
		nextTriggerWant time.Time
		triggersWant    string
	}{
		{
			name: "Timer",
			systemctlOutput: `● logrotate.timer - Daily rotation of log files
     Loaded: loaded (/usr/lib/systemd/system/logrotate.timer; enabled; preset: enabled)
     Active: active (waiting) since Wed 2024-07-17 17:50:38 UTC; 24h ago
    Trigger: Fri 2024-07-19 01:02:03 UTC; 5h 28min left
   Triggers: ● logrotate.service
       Docs: man:logrotate(8)
             man:logrotate.conf(5)`,
			descriptionWant: "Daily rotation of log files",
			loadedWant:      true,
			unitFileWant:    "/usr/lib/systemd/system/logrotate.timer",
			enabledWant:     true,
			presetWant:      true,
			activeWant:      "active (waiting)",
			nextTriggerWant: time.Date(2024, time.July, 19, 01, 02, 03, 0, time.UTC),
			triggersWant:    "logrotate.service",
		},
	}

	sshclient := VagrantSetup()

	for _, testcase := range tests {
		timer, err := sshclient.SystemdTimer("logrotate.timer")
		if err != nil {
			log.Fatalf("Error in %v: %v", testcase.name, err)
		}

		t.Run(testcase.name, func(t *testing.T) {
			if timer.Description != testcase.descriptionWant {
				t.Fatalf("want %v, got %v", testcase.descriptionWant, timer.Description)
			}
		})
	}
}

func TestLinuxKernelParameter(t *testing.T) {
	tests := []struct {
		name          string
		parameterName string
		valueWant     string
	}{
		{
			name:          "Get a parameter",
			parameterName: "vm.page_lock_unfairness",
			valueWant:     "5",
		},
	}

	sshClient := VagrantSetup()

	for _, testcase := range tests {
		linuxKernelParameter, err := sshClient.LinuxKernelParameter(testcase.parameterName)
		if err != nil {
			log.Fatalf("Error in %v: %v", testcase.name, err)
		}

		t.Run(testcase.name, func(t *testing.T) {
			if linuxKernelParameter.Value != testcase.valueWant {
				t.Errorf("want %v, got %v", testcase.valueWant, linuxKernelParameter.Value)
			}
		})
	}
}
