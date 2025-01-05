package ssh_test

import (
	"testing"

	"github.com/living-etc/sightseer.go/entities/linux"
)

func TestSshClient(t *testing.T) {
	EvaluateTestCases[linux.File, error]("File", t)
	EvaluateTestCases[linux.Service, error]("Service", t)
	EvaluateTestCases[linux.User, error]("User", t)
	EvaluateTestCases[linux.SystemdTimer, *linux.SystemdLoadError]("SystemdTimer", t)
	EvaluateTestCases[linux.LinuxKernelParameter, error]("LinuxKernelParameter", t)
	EvaluateTestCases[linux.Package, error]("Package", t)
}
