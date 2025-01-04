package ssh_test

import (
	"testing"

	sightseer "github.com/living-etc/sightseer.go/ssh"
)

func TestSshClient(t *testing.T) {
	EvaluateTestCases[sightseer.File, error]("File", t)
	EvaluateTestCases[sightseer.Service, error]("Service", t)
	EvaluateTestCases[sightseer.User, error]("User", t)
	EvaluateTestCases[sightseer.SystemdTimer, *sightseer.SystemdLoadError]("SystemdTimer", t)
	EvaluateTestCases[sightseer.LinuxKernelParameter, error]("LinuxKernelParameter", t)
	EvaluateTestCases[sightseer.Package, error]("Package", t)
}
