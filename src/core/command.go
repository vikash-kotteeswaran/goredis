package core

type Command struct {
	Command  string
	Desc     string
	Type     CType
	Executor func(*Action)
}

type CType int

func (command *Command) GetExecutor() func(*Action) {
	return command.Executor
}
