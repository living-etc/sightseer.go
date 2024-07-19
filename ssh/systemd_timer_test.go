package ssh

import (
	"log"
	"testing"
	"time"
)

func Test_systemdTimerFromSystemctl(t *testing.T) {
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
			name: "First test",
			systemctlOutput: `● logrotate.timer - Daily rotation of log files
     Loaded: loaded (/usr/lib/systemd/system/logrotate.timer; enabled; preset: enabled)
     Active: active (waiting) since Wed 2024-07-17 17:50:38 PST; 24h ago
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timer, err := parseOutput[SystemdTimer, SystemdTimerQuery](tt.systemctlOutput)
			if err != nil {
				log.Fatalf("Error in %v: %v", tt.name, err)
			}

			if tt.descriptionWant != timer.Description {
				t.Fatalf("want %v, got %v", tt.descriptionWant, timer.Description)
			}

			if tt.loadedWant != timer.Loaded {
				t.Fatalf("want %v, got %v", tt.loadedWant, timer.Loaded)
			}

			if tt.unitFileWant != timer.UnitFile {
				t.Fatalf("want %v, got %v", tt.unitFileWant, timer.UnitFile)
			}

			if tt.enabledWant != timer.Enabled {
				t.Fatalf("want %v, got %v", tt.enabledWant, timer.Enabled)
			}

			if tt.presetWant != timer.Preset {
				t.Fatalf("want %v, got %v", tt.presetWant, timer.Preset)
			}

			if tt.activeWant != timer.Active {
				t.Fatalf("want %v, got %v", tt.activeWant, timer.Active)
			}

			if tt.nextTriggerWant != timer.NextTrigger {
				t.Fatalf("want %v, got %v", tt.nextTriggerWant, timer.NextTrigger)
			}

			if tt.triggersWant != timer.Triggers {
				t.Fatalf("want %v, got %v", tt.triggersWant, timer.Triggers)
			}
		})
	}
}
