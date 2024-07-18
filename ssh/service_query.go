package ssh

type ServiceQuery struct{}

func (query ServiceQuery) Command() string {
	return "systemctl status %v --no-pager"
}

func (query ServiceQuery) Regex() string {
	return `Loaded: (?P<Loaded>\w+) \((?P<UnitFile>.*?); (?P<Enabled>\w+); .*: (?P<Preset>.*?)\)\s+Active: (?P<Active>.*? \(.+?\))`
}

func (query ServiceQuery) SetValues(values map[string]string) (*Service, error) {
	return &Service{
		Active:   values["Active"],
		Enabled:  values["Enabled"],
		Loaded:   values["Loaded"],
		UnitFile: values["UnitFile"],
		Preset:   values["Preset"],
	}, nil
}
