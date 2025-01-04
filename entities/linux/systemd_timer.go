package linux

import (
	"time"
)

type SystemdTimer struct {
	Id             string
	Description    string
	LoadState      string
	UnitFileState  string
	UnitFilePreset string
	ActiveState    string
	NextTrigger    time.Time
	Triggers       string
}
