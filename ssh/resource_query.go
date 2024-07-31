package ssh

type ResourceQuery[T ResourceType] interface {
	Command(string) string
	ParseOutput(string) (*T, error)
}
