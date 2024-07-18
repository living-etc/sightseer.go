package ssh

type ResourceQuery[T ResourceType] interface {
	Regex() string
	Command() string
	SetValues(map[string]string) (*T, error)
}
