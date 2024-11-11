package core

import (
	"encoding/hex"
	"fmt"
	"goredis/src/config"
	"strconv"
	"strings"
)

func setKeyExec(action *Action) {
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

	response := []byte(UnParseValue("OK", true, false))
	if action.IsReturnable() {
		WriteResponseToConn(&conn, response)
	}
}

func getKeyExec(action *Action) {
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

	response := []byte(UnParseValue(value, false, false))
	if action.IsReturnable() {
		WriteResponseToConn(&conn, response)
	}
}

func echoExec(action *Action) {
	conn := action.GetActionConnection()
	params, nParam := action.GetParams()
	if nParam != 1 {
		fmt.Println("Cannot perform command :: Params not correct :: ", params, nParam)
		return
	}

	echoStr := params[0].(string)
	response := []byte(UnParseValue(echoStr, false, false))
	if action.IsReturnable() {
		WriteResponseToConn(&conn, response)
	}
}

func pingExec(action *Action) {
	conn := action.GetActionConnection()
	response := []byte(UnParseValue("PONG", true, false))
	if action.IsReturnable() {
		WriteResponseToConn(&conn, response)
	}
}

func replConfExec(action *Action) {
	var response []byte
	conn := action.GetActionConnection()
	params, nParam := action.GetParams()

	for paramIdx, param := range params {
		if param == "listening-port" && (paramIdx+1) < nParam {
			replAddr, replPort := Address{}, params[paramIdx+1].(int64)
			replAddr.Absorb(conn.GetConnAddr().AddressStr())
			replAddr.Host = "0.0.0.0"
			CurrInstance.AddReplica(replAddr.Host, int(replPort))
			fmt.Println("Replica added at :: " + replAddr.Host + ":" + strconv.Itoa(int(replPort)))
			response = []byte(UnParseValue("OK", true, false))
		} else if param == "GETACK" {
			response = []byte(UnParseValue([]interface{}{"REPLCONF", "ACK", CurrInstance.ReplOffset}, false, false))
		}
	}

	if action.IsReturnable() {
		WriteResponseToConn(&conn, response)
	}
}

func pSyncExec(action *Action) {
	conn := action.GetActionConnection()
	respStr := strings.Join([]string{"FULLRESYNC", CurrInstance.ReplId, strconv.Itoa(CurrInstance.ReplOffset)}, " ")
	response := []byte(UnParseValue(respStr, true, false))
	WriteResponseToConn(&conn, response)

	rdb, rdbDecErr := hex.DecodeString(EMPTY_RDB_HEX)
	if rdbDecErr != nil {
		fmt.Println("Failed to decode RDB File :: ", rdbDecErr.Error())
	}
	response = []byte(UnParseValue(rdb, false, true))
	WriteResponseToConn(&conn, response)
}

func infoExec(action *Action) {
	conn := action.GetActionConnection()
	params, nParam := action.GetParams()
	if nParam == 0 {
		params = make([]interface{}, 1)
		params[0] = CURR_INST_INFO
	}

	var infoType string = params[0].(string)
	var info string = CurrInstance.GetInfo(infoType)
	response := []byte(UnParseValue(info, false, false))
	if action.IsReturnable() {
		WriteResponseToConn(&conn, response)
	}
}

func abortExec(action *Action) {
	config.DOABORT = true
}
