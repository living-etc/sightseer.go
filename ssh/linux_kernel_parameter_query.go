package ssh

type LinuxKernelParameterQuery struct{}

func (query LinuxKernelParameterQuery) Command() string {
	return "sudo sysctl -a | grep --color=none %v"
}

func (query LinuxKernelParameterQuery) Regex() string {
	return `= (?P<Value>.*)`
}

func (query LinuxKernelParameterQuery) SetValues(
	values map[string]string,
) (*LinuxKernelParameter, error) {
	linuxKernalParameter := &LinuxKernelParameter{}

	linuxKernalParameter.Value = values["Value"]

	return linuxKernalParameter, nil
}
