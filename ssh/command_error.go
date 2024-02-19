package ssh

import (
	"fmt"
)

type CommandError struct {
	Context string
	Err     string
}

func (err *CommandError) Error() string {
	return fmt.Sprintf("[Error]: %v [Context]: %v", err.Err, err.Context)
}
