package core

import "fmt"

func ParseActions(conn *Connection) error {
	for {
		value, parseErr := ParseValue(conn)

		if parseErr != nil {
			return parseErr
		}

		if value == nil {
			break
		}

		action, commandParseErr := ParseCommandFromValue(value, conn)
		if commandParseErr != nil {
			return commandParseErr
		}
		conn.Actions = append(conn.Actions, action)
	}

	return nil
}

func ParseCommandFromValue(value interface{}, conn *Connection) (Action, error) {
	action := Action{}

	arrValue, isArr := value.([]interface{})

	if isArr {
		actionType := ATYPE.RETURN
		if conn.Meta.FromMaster {
			actionType = ATYPE.NO_RETURN
		}

		actionCmdName := arrValue[0].(string)
		action.Command = CommandMap[actionCmdName]
		action.Params = arrValue[1:]
		action.Connection = conn
		action.Store = &StoreObj
		action.Type = actionType
	}

	return action, nil
}

func PostActionProcess(action *Action) {
	if CurrInstance.Role == SLAVE_ROLE {
		CurrInstance.ReplOffset += action.Connection.Meta.BytesRead
	}

	if action.IsBroadCastable() {
		BroadCastToReplicas(action)
	}
}

func BroadCastToReplicas(action *Action) {
	commandToBeBroadCasted := action.GetCommand().Command
	commandParams, nParams := action.GetParams()
	actionToBeBroadCasted := make([]interface{}, nParams+1)
	for idx := range nParams + 1 {
		if idx == 0 {
			actionToBeBroadCasted[idx] = commandToBeBroadCasted
			continue
		}
		actionToBeBroadCasted[idx] = commandParams[idx-1]
	}

	for _, replica := range CurrInstance.Replicas {
		fmt.Println("Broadcasting to :: ", replica.Addr.AddressStr())
		replica.ReplOffset += action.Connection.Meta.BytesRead
		HitFromServer(actionToBeBroadCasted, &Address{Host: replica.Addr.Host, Port: replica.Addr.Port})
	}
}
