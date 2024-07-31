package ssh

type ResourceQuery[T ResourceType] interface {
	Command(string) (string, error)
	ParseOutput(string) (*T, error)
}
