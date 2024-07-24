package ssh

import (
	"time"
)

type SystemdTimer struct {
	Description    string
	LoadState      string
	UnitFileState  string
	UnitFilePreset string
	ActiveState    string
	NextTrigger    time.Time
	Triggers       string
}
