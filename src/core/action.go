package core

type Action struct {
	Command    Command
	Params     []interface{}
	Connection *Connection
	Store      *Store
}

type Actions []Action

func (action *Action) GetCommand() Command {
	return action.Command
}

func (action *Action) GetParams() ([]interface{}, int) {
	return action.Params, len(action.Params)
}

func (action *Action) GetStore() Store {
	return *action.Store
}

func (action *Action) GetActionConnection() Connection {
	return *action.Connection
}

func (action *Action) Execute() {
	command := action.GetCommand()
	commandExec := command.GetExecutor()
	if commandExec != nil {
		commandExec(*action)
	}
}
