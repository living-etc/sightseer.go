package ssh

type TypeMismatchError struct {
	Text string
}

func (err *TypeMismatchError) Error() string {
	return err.Text
}
