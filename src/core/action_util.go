package core

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
		actionCmdName := arrValue[0].(string)
		action.Command = CommandMap[actionCmdName]
		action.Params = arrValue[1:]
		action.Connection = conn
		action.Store = &StoreObj
	}

	return action, nil
}
