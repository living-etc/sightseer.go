package ssh

import (
	"log"
	"testing"
)

func Test_serviceFromSystemctl(t *testing.T) {
	tests := []struct {
		name            string
		systemctlOutput string
		activeWant      string
		enabledWant     string
		loadedWant      string
		unitFileWant    string
		presetWant      string
	}{
		{
			name: "kubelet active",
			systemctlOutput: `● kubelet.service - Kubernetes Kubelet
     Loaded: loaded (/etc/systemd/system/kubelet.service; enabled; vendor preset: enabled)
     Active: active (running) since Thu 2024-02-15 22:48:43 UTC; 8min ago
       Docs: https://github.com/kubernetes/kubernetes
   Main PID: 532 (kubelet)
      Tasks: 11 (limit: 2263)
     Memory: 119.2M
        CPU: 4.143s
     CGroup: /system.slice/kubelet.service
             └─532 /usr/local/bin/kubelet --config=/var/lib/kubelet/kubelet-config.yaml --container-runtime=remote --container-runtime-endpoint=unix:///var/run/containerd/containerd.sock --image-pull-progress-dea…
`,
			activeWant:   "active (running)",
			enabledWant:  "enabled",
			loadedWant:   "loaded",
			unitFileWant: "/etc/systemd/system/kubelet.service",
			presetWant:   "enabled",
		},
		{
			name: "kube controller manager activating",
			systemctlOutput: `● kube-controller-manager.service - Kubernetes Controller Manager
     Loaded: loaded (/etc/systemd/system/kube-controller-manager.service; enabled; vendor preset: enabled)
     Active: activating (auto-restart) (Result: exit-code) since Sun 2024-02-18 18:20:04 UTC; 1s ago
       Docs: https://github.com/kubernetes/kubernetes
    Process: 7165 ExecStart=/usr/local/bin/kube-controller-manager \ (code=exited, status=1/FAILURE)
   Main PID: 7165 (code=exited, status=1/FAILURE)
        CPU: 103ms
`,
			activeWant:   "activating (auto-restart)",
			enabledWant:  "enabled",
			loadedWant:   "loaded",
			unitFileWant: "/etc/systemd/system/kube-controller-manager.service",
			presetWant:   "enabled",
		},
		{
			name: "ssh",
			systemctlOutput: `● ssh.service - OpenBSD Secure Shell server
     Loaded: loaded (/usr/lib/systemd/system/ssh.service; disabled; preset: enabled)
     Active: active (running) since Wed 2024-07-10 21:19:33 UTC; 8min ago
TriggeredBy: ● ssh.socket
       Docs: man:sshd(8)
             man:sshd_config(5)
   Main PID: 1275 (sshd)
      Tasks: 1 (limit: 2215)
     Memory: 4.5M (peak: 5.4M)
        CPU: 274ms
     CGroup: /system.slice/ssh.service
             └─1275 "sshd: /usr/sbin/sshd -D [listener] 0 of 10-100 startups"
`,
			activeWant:   "active (running)",
			enabledWant:  "disabled",
			loadedWant:   "loaded",
			unitFileWant: "/usr/lib/systemd/system/ssh.service",
			presetWant:   "enabled",
		},
	}

	for _, tt := range tests {
		service, err := ParseOutput[Service, ServiceQuery](tt.systemctlOutput)
		if err != nil {
			log.Fatal(err)
		}

		t.Run(tt.name, func(t *testing.T) {
			if service.Active != tt.activeWant {
				t.Errorf("want %v, got %v", tt.activeWant, service.Active)
			}

			if service.Enabled != tt.enabledWant {
				t.Errorf("want %v, got %v", tt.enabledWant, service.Enabled)
			}

			if service.Loaded != tt.loadedWant {
				t.Errorf("want %v, got %v", tt.loadedWant, service.Loaded)
			}

			if service.UnitFile != tt.unitFileWant {
				t.Errorf("want %v, got %v", tt.unitFileWant, service.UnitFile)
			}

			if service.Preset != tt.presetWant {
				t.Errorf("want %v, got %v", tt.presetWant, service.Preset)
			}
		})
	}
}
