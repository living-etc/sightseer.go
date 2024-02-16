package ssh

import (
	"testing"
)

func Test_serviceFromSystemctl(t *testing.T) {
	tests := []struct {
		systemctlOutput  string
		activeWant       string
		enabledWant      string
		loadedWant       string
		unitFileWant     string
		vendorPresetWant string
	}{
		{
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
			activeWant:       "active (running)",
			enabledWant:      "enabled",
			loadedWant:       "loaded",
			unitFileWant:     "/etc/systemd/system/kubelet.service",
			vendorPresetWant: "enabled",
		},
	}

	for _, tt := range tests {
		service, _ := serviceFromSystemctl(tt.systemctlOutput)

		if service.Active != tt.activeWant {
			t.Errorf("want %v, got %v", tt.activeWant, service.Active)
		}

		if service.Enabled != tt.enabledWant {
			t.Errorf("want %v, got %v", tt.enabledWant, service.Active)
		}

		if service.Loaded != tt.loadedWant {
			t.Errorf("want %v, got %v", tt.loadedWant, service.Loaded)
		}

		if service.UnitFile != tt.unitFileWant {
			t.Errorf("want %v, got %v", tt.unitFileWant, service.UnitFile)
		}

		if service.VendorPreset != tt.vendorPresetWant {
			t.Errorf("want %v, got %v", tt.vendorPresetWant, service.VendorPreset)
		}
	}
}
