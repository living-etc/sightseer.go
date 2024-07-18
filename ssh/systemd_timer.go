package ssh

import (
	"time"
)

type SystemdTimer struct {
	Description string
	Loaded      bool
	UnitFile    string
	Enabled     bool
	Preset      bool
	Active      string
	NextTrigger time.Time
	Triggers    string
}
