package core

import (
	"goredis/src/core/commands"
	"goredis/src/core/store"
	"net"
)

type Action struct {
	command    commands.Command
	params     []string
	connection *net.Conn
	store      *store.Store
}

type Actions []Action

func (action *Action) GetCommand() commands.Command {
	return action.command
}

func (action *Action) GetParams() ([]string, int) {
	return action.params, len(action.params)
}

func (action *Action) GetStore() store.Store {
	return *action.store
}

func (action *Action) GetActionConnection() net.Conn {
	return *action.connection
}

func (action *Action) Execute() {
	command := action.GetCommand()
	commandExec := command.GetExecutor()
	commandExec(*action)
}
