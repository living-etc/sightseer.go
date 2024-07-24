package ssh

type ResourceQuery[T ResourceType] interface {
	Command() string
	ParseOutput(string) (*T, error)
}
