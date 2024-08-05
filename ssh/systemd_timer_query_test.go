package ssh_test

import (
	"log"
	"testing"
	"time"

	sightseer "github.com/living-etc/sightseer.go/ssh"
)

func Test_systemdTimerFromSystemctl(t *testing.T) {
	tests := []struct {
		name               string
		systemctlOutput    string
		descriptionWant    string
		loadStateWant      string
		unitFileWant       string
		unitFileStateWant  string
		unitFilePresetWant string
		activeStateWant    string
		nextTriggerWant    time.Time
		triggersWant       string
	}{
		{
			name: "logrotate.timer",
			systemctlOutput: `Unit=logrotate.service
TimersCalendar={ OnCalendar=*-*-* 00:00:00 ; next_elapse=Thu 2024-07-25 00:00:00 UTC }
OnClockChange=no
OnTimezoneChange=no
NextElapseUSecRealtime=Thu 2024-07-25 00:00:00 UTC
NextElapseUSecMonotonic=0
LastTriggerUSec=Wed 2024-07-24 15:32:51 UTC
LastTriggerUSecMonotonic=0
Result=success
AccuracyUSec=1h
RandomizedDelayUSec=0
FixedRandomDelay=no
Persistent=yes
WakeSystem=no
RemainAfterElapse=yes
Id=logrotate.timer
Names=logrotate.timer
Requires=-.mount sysinit.target
WantedBy=timers.target
Conflicts=shutdown.target
Before=shutdown.target logrotate.service timers.target
After=sysinit.target -.mount time-set.target time-sync.target
Triggers=logrotate.service
RequiresMountsFor=/var/lib/systemd/timers
Documentation="man:logrotate(8)" "man:logrotate.conf(5)"
Description=Daily rotation of log files
LoadState=loaded
ActiveState=active
FreezerState=running
SubState=waiting
FragmentPath=/usr/lib/systemd/system/logrotate.timer
UnitFileState=enabled
UnitFilePreset=enabled
StateChangeTimestamp=Wed 2024-07-24 15:32:51 UTC
StateChangeTimestampMonotonic=36766027
InactiveExitTimestamp=Wed 2024-07-24 16:32:18 UTC
InactiveExitTimestampMonotonic=2634405
ActiveEnterTimestamp=Wed 2024-07-24 16:32:18 UTC
ActiveEnterTimestampMonotonic=2634405
ActiveExitTimestampMonotonic=0
InactiveEnterTimestampMonotonic=0
CanStart=yes
CanStop=yes
CanReload=no
CanIsolate=no
CanClean=state
CanFreeze=no
StopWhenUnneeded=no
RefuseManualStart=no
RefuseManualStop=no
AllowIsolate=no
DefaultDependencies=yes
SurviveFinalKillSignal=no
OnSuccessJobMode=fail
OnFailureJobMode=replace
IgnoreOnIsolate=no
NeedDaemonReload=no
JobTimeoutUSec=infinity
JobRunningTimeoutUSec=infinity
JobTimeoutAction=none
ConditionResult=yes
AssertResult=yes
ConditionTimestamp=Wed 2024-07-24 16:32:18 UTC
ConditionTimestampMonotonic=2634278
AssertTimestamp=Wed 2024-07-24 16:32:18 UTC
AssertTimestampMonotonic=2634279
Transient=no
Perpetual=no
StartLimitIntervalUSec=10s
StartLimitBurst=5
StartLimitAction=none
FailureAction=none
SuccessAction=none
InvocationID=ae48cdee54aa4041b3069a351f794ec9
CollectMode=inactive
`,
			descriptionWant:    "Daily rotation of log files",
			loadStateWant:      "loaded",
			unitFileStateWant:  "enabled",
			unitFilePresetWant: "enabled",
			activeStateWant:    "active",
			nextTriggerWant:    time.Date(2024, time.July, 25, 00, 00, 00, 0, time.UTC),
			triggersWant:       "logrotate.service",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var systemdTimerQuery sightseer.SystemdTimerQuery
			timer, err := systemdTimerQuery.ParseOutput(tt.systemctlOutput)
			if err != nil {
				log.Fatalf("Error in %v: %v", tt.name, err)
			}

			if tt.descriptionWant != timer.Description {
				t.Fatalf("Description: want %v, got %v", tt.descriptionWant, timer.Description)
			}

			if tt.loadStateWant != timer.LoadState {
				t.Fatalf("LoadState: want %v, got %v", tt.loadStateWant, timer.LoadState)
			}

			if tt.unitFileStateWant != timer.UnitFileState {
				t.Fatalf(
					"UnitFileState: want %v, got %v",
					tt.unitFileStateWant,
					timer.UnitFileState,
				)
			}

			if tt.unitFilePresetWant != timer.UnitFilePreset {
				t.Fatalf(
					"UnitFilePreset: want %v, got %v",
					tt.unitFilePresetWant,
					timer.UnitFilePreset,
				)
			}

			if tt.activeStateWant != timer.ActiveState {
				t.Fatalf("ActiveState: want %v, got %v", tt.activeStateWant, timer.ActiveState)
			}

			if tt.nextTriggerWant != timer.NextTrigger {
				t.Fatalf("NextTrigger: want %v, got %v", tt.nextTriggerWant, timer.NextTrigger)
			}

			if tt.triggersWant != timer.Triggers {
				t.Fatalf("Triggers: want %v, got %v", tt.triggersWant, timer.Triggers)
			}
		})
	}
}
