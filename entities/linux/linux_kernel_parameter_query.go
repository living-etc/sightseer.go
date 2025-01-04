package linux

import (
	"strings"
)

type LinuxKernelParameterQuery struct{}

func (query LinuxKernelParameterQuery) Command(platform string) (string, error) {
	var cmd string

	switch platform {
	default:
		cmd = "sudo sysctl -a | grep --color=none %v"
	}

	return cmd, nil
}

func (query LinuxKernelParameterQuery) ParseOutput(
	output string,
) (*LinuxKernelParameter, error) {
	linuxKernalParameter := &LinuxKernelParameter{}

	parts := strings.Split(output, " = ")

	linuxKernalParameter.Value = parts[1]

	return linuxKernalParameter, nil
}
