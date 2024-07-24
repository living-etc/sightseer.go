package ssh

import (
	"log"
	"reflect"
	"testing"
)

func Test_linuxKernelParameterParseOutput(t *testing.T) {
	tests := []struct {
		name          string
		sysctlOutput  string
		parameterWant *LinuxKernelParameter
	}{
		{
			name:         "Get a parameter",
			sysctlOutput: "vm.page_lock_unfairness = 5",
			parameterWant: &LinuxKernelParameter{
				Value: "5",
			},
		},
		{
			name:         "Spaces in value",
			sysctlOutput: "vm.lowmem_reserve_ratio = 256   256     32      0       0",
			parameterWant: &LinuxKernelParameter{
				Value: "256   256     32      0       0",
			},
		},
	}

	for _, testcase := range tests {
		var linuxKernelParameterQuery LinuxKernelParameterQuery
		linuxKernelParameter, err := linuxKernelParameterQuery.ParseOutput(
			testcase.sysctlOutput,
		)
		if err != nil {
			log.Fatalf("Error in %v: %v", testcase.name, err)
		}

		t.Run(testcase.name, func(t *testing.T) {
			if !reflect.DeepEqual(linuxKernelParameter, testcase.parameterWant) {
				t.Errorf(
					"LinuxKernelParameter failed:\nwant:\t%v\ngot:\t%v",
					testcase.parameterWant,
					linuxKernelParameter,
				)
			}
		})
	}
}
