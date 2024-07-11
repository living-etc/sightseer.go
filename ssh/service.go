package ssh

import (
	"regexp"
)

type Service struct {
	Active   string
	Enabled  string
	Loaded   string
	UnitFile string
	Preset   string
}

func serviceFromSystemctl(systemctlOutput string) (Service, error) {
	pattern := `Loaded: (?P<Loaded>\w+) \((?P<UnitFile>.*?); (?P<Enabled>\w+); .*: (?P<Preset>.*?)\)\s+Active: (?P<Active>.*? \(.+?\))`

	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(systemctlOutput)

	service := Service{
		Active:   matches[re.SubexpIndex("Active")],
		Enabled:  matches[re.SubexpIndex("Enabled")],
		Loaded:   matches[re.SubexpIndex("Loaded")],
		UnitFile: matches[re.SubexpIndex("UnitFile")],
		Preset:   matches[re.SubexpIndex("Preset")],
	}

	return service, nil
}
