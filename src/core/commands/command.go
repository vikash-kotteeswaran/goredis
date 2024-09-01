package commands

type Command struct {
	command  string
	desc     string
	executor func(Action)
}

func (command *Command) GetExecutor() func(Action) {
	return command.executor
}
