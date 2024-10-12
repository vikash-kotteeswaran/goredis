package core

import (
	"fmt"
	"goredis/src/config"
	"strconv"
	"strings"
)

func setKeyExec(action Action) {
	conn := action.GetActionConnection()
	store := action.GetStore()
	params, nParam := action.GetParams()
	if nParam < 2 {
		fmt.Println("Cannot perform command :: Params not correct :: ", params, nParam)
		return
	}

	key := params[0]
	val := params[1]
	var ttl int64 = -1

	for idx, param := range params {
		paramStr, isStr := param.(string)
		if isStr && strings.ToLower(paramStr) == "px" {
			if idx+1 >= nParam {
				fmt.Println("Cannot perform command :: Params not correct :: ", params, nParam)
			}
			ttl = params[idx+1].(int64)
		}
	}
	inserted, err := store.Set(key, val, ttl)

	if !inserted {
		fmt.Println("Did not enter key :: Key :: ", key)
		return
	}

	if err != nil {
		fmt.Println("Error while entering key :: Key :: ", key)
		return
	}

	response := []byte(UnParseValue("OK", true))
	conn.Write(response)
}

func getKeyExec(action Action) {
	conn := action.GetActionConnection()
	store := action.GetStore()
	params, nParam := action.GetParams()

	if nParam != 1 {
		fmt.Println("Cannot perform command :: Params not correct :: ", params, nParam)
		return
	}

	key := params[0].(string)
	value, getErr := store.Get(key)
	if getErr != nil {
		fmt.Println("Error while getting key value :: Key :: ", key)
		return
	}

	response := []byte(UnParseValue(value, false))
	_, writeErr := conn.Write(response)
	if writeErr != nil {
		fmt.Println("Failed to perform command")
	}
}

func echoExec(action Action) {
	conn := action.GetActionConnection()
	params, nParam := action.GetParams()
	if nParam != 1 {
		fmt.Println("Cannot perform command :: Params not correct :: ", params, nParam)
		return
	}

	echoStr := params[0].(string)
	response := []byte(UnParseValue(echoStr, false))
	_, err := conn.Write(response)
	if err != nil {
		fmt.Println("Failed to perform command")
	}
}

func pingExec(action Action) {
	conn := action.GetActionConnection()
	response := []byte(UnParseValue("PONG", true))
	_, err := conn.Write(response)
	if err != nil {
		fmt.Println("Failed to perform command")
	}
}

func replConfExec(action Action) {
	conn := action.GetActionConnection()
	response := []byte(UnParseValue("OK", true))
	_, err := conn.Write(response)
	if err != nil {
		fmt.Println("Failed to perform command")
	}
}

func pSyncExec(action Action) {
	conn := action.GetActionConnection()
	respStr := strings.Join([]string{"FULLRESYNC", CurrInstance.ReplId, strconv.Itoa(CurrInstance.ReplOffset)}, " ")
	response := []byte(UnParseValue(respStr, true))
	_, err := conn.Write(response)
	if err != nil {
		fmt.Println("Failed to perform command")
	}
}

func infoExec(action Action) {
	conn := action.GetActionConnection()
	params, nParam := action.GetParams()
	if nParam == 0 {
		params = make([]interface{}, 1)
		params[0] = CURR_INST_INFO
	}

	var infoType string = params[0].(string)
	var info string = CurrInstance.GetInfo(infoType)
	response := []byte(UnParseValue(info, false))
	_, err := conn.Write(response)
	if err != nil {
		fmt.Println("Failed to perform command")
	}
}

func abortExec(action Action) {
	config.DOABORT = true
}
