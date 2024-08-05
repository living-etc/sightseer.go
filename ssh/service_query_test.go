package ssh_test

import (
	"log"
	"testing"

	sightseer "github.com/living-etc/sightseer.go/ssh"
)

func TestParseOutput(t *testing.T) {
	tests := []struct {
		name            string
		systemctlOutput string
		serviceWant     *sightseer.Service
	}{
		{
			name: "hello-world service",
			systemctlOutput: `Type=simple
ExitType=main
Restart=always
RestartMode=normal
NotifyAccess=none
RestartUSec=1s
RestartSteps=0
RestartMaxDelayUSec=infinity
RestartUSecNext=1s
TimeoutStartUSec=1min 30s
TimeoutStopUSec=1min 30s
TimeoutAbortUSec=1min 30s
TimeoutStartFailureMode=terminate
TimeoutStopFailureMode=terminate
RuntimeMaxUSec=infinity
RuntimeRandomizedExtraUSec=0
WatchdogUSec=0
WatchdogTimestampMonotonic=0
RootDirectoryStartOnly=no
RemainAfterExit=no
GuessMainPID=yes
MainPID=2353
ControlPID=0
FileDescriptorStoreMax=0
NFileDescriptorStore=0
FileDescriptorStorePreserve=restart
StatusErrno=0
Result=success
ReloadResult=success
CleanResult=success
UID=0
GID=0
NRestarts=0
OOMPolicy=stop
ReloadSignal=1
ExecMainStartTimestamp=Wed 2024-07-24 19:36:15 UTC
ExecMainStartTimestampMonotonic=8562228
ExecMainExitTimestampMonotonic=0
ExecMainPID=2353
ExecMainCode=0
ExecMainStatus=0
ExecStart={ path=/bin/bash ; argv[]=/bin/bash -c while true; do sleep 1; done ; ignore_errors=no ; start_time=[Wed 2024-07-24 19:36:15 UTC] ; stop_time=[n/a] ; pid=2353 ; code=(null) ; status=0/0 }
ExecStartEx={ path=/bin/bash ; argv[]=/bin/bash -c while true; do sleep 1; done ; flags= ; start_time=[Wed 2024-07-24 19:36:15 UTC] ; stop_time=[n/a] ; pid=2353 ; code=(null) ; status=0/0 }
Slice=system.slice
ControlGroup=/system.slice/hello-world.service
ControlGroupId=4113
MemoryCurrent=589824
MemoryPeak=1232896
MemorySwapCurrent=0
MemorySwapPeak=0
MemoryZSwapCurrent=0
MemoryAvailable=1814130688
CPUUsageNSec=2464601000
TasksCurrent=2
IPIngressBytes=[no data]
IPIngressPackets=[no data]
IPEgressBytes=[no data]
IPEgressPackets=[no data]
IOReadBytes=[not set]
IOReadOperations=[not set]
IOWriteBytes=[not set]
IOWriteOperations=[not set]
Delegate=no
CPUAccounting=yes
CPUWeight=[not set]
StartupCPUWeight=[not set]
CPUShares=[not set]
StartupCPUShares=[not set]
CPUQuotaPerSecUSec=infinity
CPUQuotaPeriodUSec=infinity
IOAccounting=no
IOWeight=[not set]
StartupIOWeight=[not set]
BlockIOAccounting=no
BlockIOWeight=[not set]
StartupBlockIOWeight=[not set]
MemoryAccounting=yes
DefaultMemoryLow=0
DefaultStartupMemoryLow=0
DefaultMemoryMin=0
MemoryMin=0
MemoryLow=0
StartupMemoryLow=0
MemoryHigh=infinity
StartupMemoryHigh=infinity
MemoryMax=infinity
StartupMemoryMax=infinity
MemorySwapMax=infinity
StartupMemorySwapMax=infinity
MemoryZSwapMax=infinity
StartupMemoryZSwapMax=infinity
MemoryLimit=infinity
DevicePolicy=auto
TasksAccounting=yes
TasksMax=2214
IPAccounting=no
ManagedOOMSwap=auto
ManagedOOMMemoryPressure=auto
ManagedOOMMemoryPressureLimit=0
ManagedOOMPreference=none
MemoryPressureWatch=auto
MemoryPressureThresholdUSec=200ms
CoredumpReceive=no
UMask=0022
LimitCPU=infinity
LimitCPUSoft=infinity
LimitFSIZE=infinity
LimitFSIZESoft=infinity
LimitDATA=infinity
LimitDATASoft=infinity
LimitSTACK=infinity
LimitSTACKSoft=8388608
LimitCORE=infinity
LimitCORESoft=0
LimitRSS=infinity
LimitRSSSoft=infinity
LimitNOFILE=524288
LimitNOFILESoft=1024
LimitAS=infinity
LimitASSoft=infinity
LimitNPROC=7381
LimitNPROCSoft=7381
LimitMEMLOCK=8388608
LimitMEMLOCKSoft=8388608
LimitLOCKS=infinity
LimitLOCKSSoft=infinity
LimitSIGPENDING=7381
LimitSIGPENDINGSoft=7381
LimitMSGQUEUE=819200
LimitMSGQUEUESoft=819200
LimitNICE=0
LimitNICESoft=0
LimitRTPRIO=0
LimitRTPRIOSoft=0
LimitRTTIME=infinity
LimitRTTIMESoft=infinity
RootEphemeral=no
OOMScoreAdjust=0
CoredumpFilter=0x33
Nice=0
IOSchedulingClass=2
IOSchedulingPriority=4
CPUSchedulingPolicy=0
CPUSchedulingPriority=0
CPUAffinityFromNUMA=no
NUMAPolicy=n/a
TimerSlackNSec=50000
CPUSchedulingResetOnFork=no
NonBlocking=no
StandardInput=null
StandardOutput=journal
StandardError=inherit
TTYReset=no
TTYVHangup=no
TTYVTDisallocate=no
SyslogPriority=30
SyslogLevelPrefix=yes
SyslogLevel=6
SyslogFacility=3
LogLevelMax=-1
LogRateLimitIntervalUSec=0
LogRateLimitBurst=0
SecureBits=0
CapabilityBoundingSet=cap_chown cap_dac_override cap_dac_read_search cap_fowner cap_fsetid cap_kill cap_setgid cap_setuid cap_setpcap cap_linux_immutable cap_net_bind_service cap_net_broadcast cap_net_admin cap_net_raw cap_ipc_lock cap_ipc_owner cap_sys_module cap_sys_rawio cap_sys_chroot cap_sys_ptrace cap_sys_pacct cap_sys_admin cap_sys_boot cap_sys_nice cap_sys_resource cap_sys_time cap_sys_tty_config cap_mknod cap_lease cap_audit_write cap_audit_control cap_setfcap cap_mac_override cap_mac_admin cap_syslog cap_wake_alarm cap_block_suspend cap_audit_read cap_perfmon cap_bpf cap_checkpoint_restore
User=root
DynamicUser=no
SetLoginEnvironment=no
RemoveIPC=no
PrivateTmp=no
PrivateDevices=no
ProtectClock=no
ProtectKernelTunables=no
ProtectKernelModules=no
ProtectKernelLogs=no
ProtectControlGroups=no
PrivateNetwork=no
PrivateUsers=no
PrivateMounts=no
PrivateIPC=no
ProtectHome=no
ProtectSystem=no
SameProcessGroup=no
UtmpMode=init
IgnoreSIGPIPE=yes
NoNewPrivileges=no
SystemCallErrorNumber=2147483646
LockPersonality=no
RuntimeDirectoryPreserve=no
RuntimeDirectoryMode=0755
StateDirectoryMode=0755
CacheDirectoryMode=0755
LogsDirectoryMode=0755
ConfigurationDirectoryMode=0755
TimeoutCleanUSec=infinity
MemoryDenyWriteExecute=no
RestrictRealtime=no
RestrictSUIDSGID=no
RestrictNamespaces=no
MountAPIVFS=no
KeyringMode=private
ProtectProc=default
ProcSubset=all
ProtectHostname=no
MemoryKSM=no
RootImagePolicy=root=verity+signed+encrypted+unprotected+absent:usr=verity+signed+encrypted+unprotected+absent:home=encrypted+unprotected+absent:srv=encrypted+unprotected+absent:tmp=encrypted+unprotected+absent:var=encrypted+unprotected+absent
MountImagePolicy=root=verity+signed+encrypted+unprotected+absent:usr=verity+signed+encrypted+unprotected+absent:home=encrypted+unprotected+absent:srv=encrypted+unprotected+absent:tmp=encrypted+unprotected+absent:var=encrypted+unprotected+absent
ExtensionImagePolicy=root=verity+signed+encrypted+unprotected+absent:usr=verity+signed+encrypted+unprotected+absent:home=encrypted+unprotected+absent:srv=encrypted+unprotected+absent:tmp=encrypted+unprotected+absent:var=encrypted+unprotected+absent
KillMode=control-group
KillSignal=15
RestartKillSignal=15
FinalKillSignal=9
SendSIGKILL=yes
SendSIGHUP=no
WatchdogSignal=6
Id=hello-world.service
Names=hello-world.service
Requires=system.slice sysinit.target
WantedBy=multi-user.target
Conflicts=shutdown.target
Before=multi-user.target shutdown.target
After=systemd-journald.socket sysinit.target system.slice basic.target network.target
Description=Simple service for testing against
LoadState=loaded
ActiveState=active
FreezerState=running
SubState=running
FragmentPath=/etc/systemd/system/hello-world.service
UnitFileState=enabled
UnitFilePreset=enabled
StateChangeTimestamp=Wed 2024-07-24 19:36:15 UTC
StateChangeTimestampMonotonic=8562306
InactiveExitTimestamp=Wed 2024-07-24 19:36:15 UTC
InactiveExitTimestampMonotonic=8562306
ActiveEnterTimestamp=Wed 2024-07-24 19:36:15 UTC
ActiveEnterTimestampMonotonic=8562306
ActiveExitTimestampMonotonic=0
InactiveEnterTimestampMonotonic=0
CanStart=yes
CanStop=yes
CanReload=no
CanIsolate=no
CanFreeze=yes
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
ConditionTimestamp=Wed 2024-07-24 19:36:15 UTC
ConditionTimestampMonotonic=8558641
AssertTimestamp=Wed 2024-07-24 19:36:15 UTC
AssertTimestampMonotonic=8558642
Transient=no
Perpetual=no
StartLimitIntervalUSec=0
StartLimitBurst=5
StartLimitAction=none
FailureAction=none
SuccessAction=none
InvocationID=f5aedca5bccf4487ac3af35164dc5aae
CollectMode=inactive
`,
			serviceWant: &sightseer.Service{
				Description:    "Simple service for testing against",
				LoadState:      "loaded",
				UnitFileState:  "enabled",
				UnitFilePreset: "enabled",
				ActiveState:    "active",
			},
		},
	}

	for _, testcase := range tests {
		var serviceQuery sightseer.ServiceQuery
		service, err := serviceQuery.ParseOutput(testcase.systemctlOutput)
		if err != nil {
			log.Fatalf("Error in %v: %v", testcase.name, err)
		}

		EvaluateStructTypesAreEqual(service, testcase.serviceWant, testcase.name, t)
	}
}
