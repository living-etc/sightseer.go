package ssh

import (
	"regexp"
)

type Service struct {
	Active       string
	Enabled      string
	Loaded       string
	UnitFile     string
	VendorPreset string
}

func serviceFromSystemctl(systemctlOutput string) (Service, error) {
	pattern := `Loaded: (?P<Loaded>\w+) \((?P<UnitFile>.*?); (?P<Enabled>\w+); vendor preset: (?P<VendorPreset>.*?)\)\s+Active: (?P<Active>.*? \(.+?\))`

	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(systemctlOutput)

	service := Service{
		Active:       matches[re.SubexpIndex("Active")],
		Enabled:      matches[re.SubexpIndex("Enabled")],
		Loaded:       matches[re.SubexpIndex("Loaded")],
		UnitFile:     matches[re.SubexpIndex("UnitFile")],
		VendorPreset: matches[re.SubexpIndex("VendorPreset")],
	}

	return service, nil
}
