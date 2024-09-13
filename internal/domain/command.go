package domain

type CommandType string

const (
	CommandGet    CommandType = "GET"
	CommandSet    CommandType = "SET"
	CommandDelete CommandType = "Delete"
)

func (t CommandType) Valid() bool {
	return t == CommandGet || t == CommandSet || t == CommandDelete
}

type Command struct {
	Type  CommandType
	Key   string
	Value string
}
