package ssh

import (
	"testing"
)

func TestSshClient(t *testing.T) {
	EvaluateTestCases[File, error]("File", t)
	EvaluateTestCases[Service, error]("Service", t)
	EvaluateTestCases[User, error]("User", t)
	EvaluateTestCases[SystemdTimer, *SystemdLoadError]("SystemdTimer", t)
	EvaluateTestCases[LinuxKernelParameter, error]("LinuxKernelParameter", t)
}
