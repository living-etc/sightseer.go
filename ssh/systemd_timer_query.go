package ssh

import (
	"bufio"
	"errors"
	"regexp"
	"strings"
	"time"
)

type SystemdTimerQuery struct{}

func (query SystemdTimerQuery) Command() string {
	return "systemctl show %v"
}

func (query SystemdTimerQuery) ParseOutput(output string) (*SystemdTimer, error) {
	systemdTimer := &SystemdTimer{}

	values := make(map[string]string)

	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		attribute := strings.SplitN(scanner.Text(), "=", 2)
		values[attribute[0]] = attribute[1]
	}

	if values["LoadState"] == "not-found" {
		return nil, &SystemdLoadError{
			UnitName:  values["Id"],
			LoadState: values["LoadState"],
			LoadError: values["LoadError"],
		}
	}

	systemdTimer.Id = values["Id"]
	systemdTimer.Description = values["Description"]
	systemdTimer.LoadState = values["LoadState"]
	systemdTimer.UnitFileState = values["UnitFileState"]
	systemdTimer.UnitFilePreset = values["UnitFilePreset"]
	systemdTimer.ActiveState = values["ActiveState"]
	systemdTimer.Triggers = values["Triggers"]

	re := regexp.MustCompile(`next_elapse=(\w+ \d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} \w+)`)
	match := re.FindStringSubmatch(values["TimersCalendar"])
	if len(match) == 0 {
		return nil, errors.New("No timestamp found in the input string")
	}

	nextTriggerTime, err := time.Parse("Mon 2006-01-02 15:04:05 MST", match[1])
	if err != nil {
		return nil, err
	}
	systemdTimer.NextTrigger = nextTriggerTime

	return systemdTimer, nil
}
