package linux

import (
	"bufio"
	"strings"
)

type PackageQuery struct{}

func (query PackageQuery) Command(platform string) (string, error) {
	var cmd string

	switch platform {
	case "ubuntu2404":
		cmd = "sudo dpkg -s %v"
	case "fedora40":
		cmd = "rpm -qi %v"
	}

	return cmd, nil
}

func (query PackageQuery) ParseOutput(output string) (*Package, error) {
	var pkg Package

	result := make(map[string]string)
	var currentKey string

	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") {
			result[currentKey] += "\n" + strings.TrimSpace(line)
		} else {
			kv := strings.SplitN(line, ":", 2)
			if len(kv) == 2 {
				result[kv[0]] = strings.TrimSpace(kv[1])
			} else {
				result[kv[0]] = ""
			}
			currentKey = kv[0]
		}
	}

	pkg.Name = result["Package"]
	pkg.Status = result["Status"]
	pkg.Priority = result["Priority"]
	pkg.Section = result["Section"]
	pkg.InstalledSize = result["Installed-Size"]
	pkg.Maintainer = result["Maintainer"]
	pkg.Architecture = result["Architecture"]
	pkg.MultiArch = result["Multi-Arch"]
	pkg.Source = result["Source"]
	pkg.Version = result["Version"]
	pkg.Replaces = result["Replaces"]
	pkg.Provides = result["Provides"]
	pkg.Depends = result["Depends"]
	pkg.PreDepends = result["Pre-Depends"]
	pkg.Recommends = result["Recommends"]
	pkg.Suggests = result["Suggests"]
	pkg.Conflicts = result["Conflicts"]
	pkg.Conffiles = result["Conffiles"]
	pkg.Description = result["Description"]
	pkg.Homepage = result["Homepage"]
	pkg.OriginalMaintainer = result["Original-Maintainer"]

	return &pkg, nil
}
