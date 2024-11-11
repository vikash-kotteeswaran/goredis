package core

import "fmt"

type Action struct {
	Type       AType
	Command    Command
	Params     []interface{}
	Connection *Connection
	Store      *Store
}

type Actions []Action

type AType int

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

	if commandExec != nil && action.IsExecutable() {
		commandExec(action)
	} else {
		fmt.Println("Command is not executable for the connection")
	}

	PostActionProcess(action)
}

// Action Config Helper

func (action *Action) IsBroadCastable() bool {
	return CurrInstance.HasReplicas() && action.Command.Type&CTYPE.CWRITE > 0 //&& !CurrInstance.Replicas.ContainsIp(action.Connection.GetConnAddr())
}

func (action *Action) IsReturnable() bool {
	return action.Type&ATYPE.RETURN > 0 || action.Command.Type&CTYPE.CINFO > 0 || action.Command.Type&CTYPE.CREPL > 0 || action.Command.Type&CTYPE.CCRITICAL > 0
}

func (action *Action) IsExecutable() bool {
	return (action.Command.Type&CTYPE.CREPL > 0 || action.Command.Type&CTYPE.CINFO > 0) && action.Connection.Meta.FromReplica || !action.Connection.Meta.FromReplica
}
