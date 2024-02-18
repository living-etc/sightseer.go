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
	regexpes := []string{
		`Active: (?P<Active>.*? \(.+?\))`,
		`Loaded: (?P<Loaded>\w+) \((?P<UnitFile>.*?); (?P<Enabled>\w+); vendor preset: (?P<VendorPreset>.*?)\)`,
	}

	result := make(map[string]string)
	for _, pattern := range regexpes {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(systemctlOutput)

		for i, name := range re.SubexpNames() {
			if i > 0 {
				result[name] = matches[i]
			}
		}
	}

	service := Service{
		Active:       result["Active"],
		Enabled:      result["Enabled"],
		Loaded:       result["Loaded"],
		UnitFile:     result["UnitFile"],
		VendorPreset: result["VendorPreset"],
	}

	return service, nil
}
