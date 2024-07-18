package ssh

import (
	"strconv"
	"strings"
	"time"
)

type SystemdTimerQuery struct{}

func (query SystemdTimerQuery) Command() string {
	return "systemctl status %v"
}

func (query SystemdTimerQuery) Regex() string {
	return `^.+ - (?P<Description>.*)\s+Loaded: (?P<Loaded>\w+) \((?P<UnitFile>.*?); (?P<Enabled>\w+); .*: (?P<Preset>.*?)\)\s+Active: (?P<Active>.*? \(.+?\))[\s\S]+Trigger: \w+ (?P<NextTriggerDate>.*) (?P<NextTriggerTime>.*) (?P<NextTriggerTimeZone>\w+);[\s\S]+Triggers: ● (?P<Triggers>.*)`
}

func (query SystemdTimerQuery) SetValues(values map[string]string) (*SystemdTimer, error) {
	systemdTimer := &SystemdTimer{}

	systemdTimer.Description = values["Description"]

	loadedString := values["Loaded"]
	if loadedString == "loaded" {
		systemdTimer.Loaded = true
	} else {
		systemdTimer.Loaded = false
	}

	systemdTimer.UnitFile = values["UnitFile"]

	enabledString := values["Enabled"]
	if enabledString == "enabled" {
		systemdTimer.Enabled = true
	} else {
		systemdTimer.Enabled = false
	}

	presetString := values["Preset"]
	if presetString == "enabled" {
		systemdTimer.Preset = true
	} else {
		systemdTimer.Preset = false
	}

	systemdTimer.Active = values["Active"]

	dateParts := strings.Split(values["NextTriggerDate"], "-")
	year, err := strconv.Atoi(dateParts[0])
	if err != nil {
		return nil, err
	}
	month, err := strconv.Atoi(dateParts[1])
	if err != nil {
		return nil, err
	}
	day, err := strconv.Atoi(dateParts[2])
	if err != nil {
		return nil, err
	}

	timeParts := strings.Split(values["NextTriggerTime"], ":")
	hours, err := strconv.Atoi(timeParts[0])
	if err != nil {
		return nil, err
	}
	minutes, err := strconv.Atoi(timeParts[1])
	if err != nil {
		return nil, err
	}
	seconds, err := strconv.Atoi(timeParts[2])
	if err != nil {
		return nil, err
	}

	systemdTimer.NextTrigger = time.Date(
		year,
		time.Month(month),
		day,
		hours,
		minutes,
		seconds,
		0,
		time.UTC,
	)

	systemdTimer.Triggers = values["Triggers"]

	return systemdTimer, nil
}
