package linux

type FileError struct {
	ErrorReason string
}

func (err FileError) Error() string { return err.ErrorReason }
