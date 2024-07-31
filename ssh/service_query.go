package ssh

import (
	"bufio"
	"strings"
)

type ServiceQuery struct{}

func (query ServiceQuery) Command(platform string) string {
	switch platform {
	default:
		return "systemctl show %v --no-pager"
	}
}

func (query ServiceQuery) ParseOutput(output string) (*Service, error) {
	service := &Service{}

	values := make(map[string]string)

	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		attribute := strings.SplitN(scanner.Text(), "=", 2)
		values[attribute[0]] = attribute[1]
	}

	service.Description = values["Description"]
	service.LoadState = values["LoadState"]
	service.UnitFileState = values["UnitFileState"]
	service.UnitFilePreset = values["UnitFilePreset"]
	service.ActiveState = values["ActiveState"]

	return service, nil
}
