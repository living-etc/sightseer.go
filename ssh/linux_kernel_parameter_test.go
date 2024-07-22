package ssh

import (
	"log"
	"testing"
)

func Test_linuxKernelParameterParseOutput(t *testing.T) {
	tests := []struct {
		name         string
		sysctlOutput string
		valueWant    string
	}{
		{
			name:         "Get a parameter",
			sysctlOutput: "vm.page_lock_unfairness = 5",
			valueWant:    "5",
		},
		{
			name:         "Spaces in value",
			sysctlOutput: "vm.lowmem_reserve_ratio = 256   256     32      0       0",
			valueWant:    "256   256     32      0       0",
		},
	}

	for _, testcase := range tests {
		linuxKernelParameter, err := parseOutput[LinuxKernelParameter, LinuxKernelParameterQuery](
			testcase.sysctlOutput,
		)
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
