package ssh

type ServiceQuery struct{}

func (query ServiceQuery) Command() string {
	return "systemctl status %v --no-pager"
}

func (query ServiceQuery) Regex() string {
	return `Loaded: (?P<Loaded>\w+) \((?P<UnitFile>.*?); (?P<Enabled>\w+); .*: (?P<Preset>.*?)\)\s+Active: (?P<Active>.*? \(.+?\))`
}

func (query ServiceQuery) SetValues(service *Service, values map[string]string) {
	service.Active = values["Active"]
	service.Enabled = values["Enabled"]
	service.Loaded = values["Loaded"]
	service.UnitFile = values["UnitFile"]
	service.Preset = values["Preset"]
}
