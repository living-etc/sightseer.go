package ssh

type ResourceQuery[T ResourceType] interface {
	Regex() string
	Command() string
	SetValues(*T, map[string]string)
}
