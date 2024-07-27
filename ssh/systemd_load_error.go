package ssh

import "fmt"

type SystemdLoadError struct {
	UnitName  string
	LoadState string
	LoadError string
}

func (err *SystemdLoadError) Error() string {
	return fmt.Sprintf(
		"[Unit]: %v - [LoadState]: %v - [LoadError]: %v",
		err.UnitName,
		err.LoadState,
		err.LoadError,
	)
}
