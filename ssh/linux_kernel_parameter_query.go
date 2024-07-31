package ssh

import (
	"strings"
)

type LinuxKernelParameterQuery struct{}

func (query LinuxKernelParameterQuery) Command(platform string) string {
	switch platform {
	default:
		return "sudo sysctl -a | grep --color=none %v"
	}
}

func (query LinuxKernelParameterQuery) ParseOutput(
	output string,
) (*LinuxKernelParameter, error) {
	linuxKernalParameter := &LinuxKernelParameter{}

	parts := strings.Split(output, " = ")

	linuxKernalParameter.Value = parts[1]

	return linuxKernalParameter, nil
}
