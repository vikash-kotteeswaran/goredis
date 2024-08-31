package core

import (
	"fmt"
	"goredis/app/configcore"
	"goredis/app/constants"
	"io"
	"net"
	"strconv"
	"strings"
)

type command struct {
	name string
}

func Serve(conn net.Conn) error {
	payload, err := readPayload(conn)
	payload = strings.ReplaceAll(payload, "\n", "\\n")
	requests := strings.Split(payload, "\\n")

	for _, request := range requests {
		if request == "" || request == " " {
			continue
		}
		command, nParameters, parameters := parseCommandReq(request)

		switch command {
		case constants.SET:
			setKeyAction(conn, parameters, nParameters)
			break
		case constants.GET:
			getKeyAction(conn, parameters, nParameters)
			break
		case constants.PING:
			pingAction(conn)
			break
		case constants.ECHO:
			echoAction(conn, parameters, nParameters)
			break
		case constants.INFO:
			infoAction(conn, parameters, nParameters)
			break
		case constants.ABORT:
			abortAction()
			break
		default:
			break
		}
	}

	fmt.Println("Connection :: ", conn.RemoteAddr().String(), " has been served")

	return err
}

func readPayload(conn net.Conn) (string, error) {
	var readErr error
	var commandBytes []byte
	bufferLength := 512
	readBuf := make([]byte, bufferLength)
	endRead := false

	for {
		nRead, err := conn.Read(readBuf)

		if err != nil {
			if err != io.EOF {
				fmt.Println("Failed to read command :: ", err)
			}
			endRead = true
			readErr = err
		}

		commandBytes = append(commandBytes, readBuf[:nRead]...)

		if endRead {
			break
		}
	}

	return string(commandBytes), readErr
}

// Parse Command Request

func parseCommandReq(request string) (string, int, []string) {
	var command string
	var params []string
	var nParams int
	commandParts := []string{}

	insideQuotes := false
	addOnSpace := true

	partFormed := strings.Builder{}

	var reqPrevChar, reqChar rune

	for _, reqChar = range request {
		if addOnSpace && !insideQuotes && reqChar == rune(' ') {
			commandOrParam := partFormed.String()
			if commandOrParam != "" {
				commandParts = append(commandParts, partFormed.String())
				partFormed.Reset()
				continue
			}
		} else if reqChar == '\'' && reqPrevChar != rune('\\') {
			addOnSpace = insideQuotes
			insideQuotes = !insideQuotes
		}
		partFormed.WriteRune(reqChar)
		reqPrevChar = reqChar
	}

	// Adding the last parts of the request that cannot be handled in For Loop
	commandOrParam := partFormed.String()
	if commandOrParam != "" {
		commandParts = append(commandParts, partFormed.String())
	}

	nParams = len(commandParts) - 1

	if nParams > -1 {
		command = commandParts[0]
	}
	if nParams > 0 {
		params = commandParts[1:]
	}

	// First Element is the Command in itself and rest are its Parameters
	return command, nParams, params
}

// Actions for Command

func setKeyAction(conn net.Conn, params []string, nParam int) {
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
}

func getKeyAction(conn net.Conn, params []string, nParam int) {
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

func echoAction(conn net.Conn, params []string, nParam int) {
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

func pingAction(conn net.Conn) {
	response := []byte("+PONG\r\n")
	_, err := conn.Write(response)
	if err != nil {
		fmt.Println("Failed to perform command")
	}
}

func infoAction(conn net.Conn, params []string, nParam int) {
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

func abortAction() {
	constants.DOABORT = true
}
