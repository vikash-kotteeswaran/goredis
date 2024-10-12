package core

type Command struct {
	Command  string
	Desc     string
	Executor func(Action)
}

func (command *Command) GetExecutor() func(Action) {
	return command.Executor
}
