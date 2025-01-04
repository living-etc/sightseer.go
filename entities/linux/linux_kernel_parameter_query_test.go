package linux_test

import (
	"log"
	"testing"

	"github.com/living-etc/sightseer.go/entities/linux"
)

func Test_linuxKernelParameterParseOutput(t *testing.T) {
	tests := []struct {
		name          string
		sysctlOutput  string
		parameterWant *linux.LinuxKernelParameter
	}{
		{
			name:         "Get a parameter",
			sysctlOutput: "vm.page_lock_unfairness = 5",
			parameterWant: &linux.LinuxKernelParameter{
				Value: "5",
			},
		},
		{
			name:         "Spaces in value",
			sysctlOutput: "vm.lowmem_reserve_ratio = 256   256     32      0       0",
			parameterWant: &linux.LinuxKernelParameter{
				Value: "256   256     32      0       0",
			},
		},
	}

	for _, testcase := range tests {
		var linuxKernelParameterQuery linux.LinuxKernelParameterQuery
		linuxKernelParameter, err := linuxKernelParameterQuery.ParseOutput(
			testcase.sysctlOutput,
		)
		if err != nil {
			log.Fatalf("Error in %v: %v", testcase.name, err)
		}

		EvaluateStructTypesAreEqual(linuxKernelParameter, testcase.parameterWant, testcase.name, t)
	}
}
