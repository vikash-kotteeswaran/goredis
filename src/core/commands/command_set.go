package commands

import (
	"fmt"
	"goredis/src/config"
	"goredis/src/configcore"
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

	if strings.ToLower(params[2]) == "px" {
		if nParam <= 3 {
			fmt.Println("Cannot perform command :: Params not correct :: ", params, nParam)
			return
		}
		var parseErr error = nil
		ttl, parseErr = strconv.ParseInt(params[3], 10, 64)
		if parseErr != nil {
			fmt.Println("Performing Set Action :: Incorrect ttl given :: ", params[3])
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

	response := []byte("+OK\r\n")
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

	key := params[0]
	value, getErr := store.Get(key)
	if getErr != nil {
		fmt.Println("Error while getting key value :: Key :: ", key)
		return
	}

	response := []byte(value.(string) + "\r\n")
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

	echoStr := params[0]
	response := []byte(echoStr + "\r\n")
	_, err := conn.Write(response)
	if err != nil {
		fmt.Println("Failed to perform command")
	}
}

func pingExec(action Action) {
	conn := action.GetActionConnection()
	response := []byte("+PONG\r\n")
	_, err := conn.Write(response)
	if err != nil {
		fmt.Println("Failed to perform command")
	}
}

func infoExec(action Action) {
	conn := action.GetActionConnection()
	params, nParam := action.GetParams()
	if nParam == 0 {
		params = []string{configcore.CURR_INST_INFO}
	}

	var infoType string = params[0]
	var info string = configcore.CurrInstance.GetInfo(infoType)
	response := []byte(info + "\r\n")
	_, err := conn.Write(response)
	if err != nil {
		fmt.Println("Failed to perform command")
	}
}

func abortExec(action Action) {
	config.DOABORT = true
}
