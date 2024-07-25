package core

import (
	"fmt"
	"io"
	"net"
	"strings"
)

const (
	PING = "PING"
)

type command struct {
	name string
}

func Serve(conn net.Conn) {
	request := readCommand(conn)
	commandParts := strings.Split(request, " ")
	command := commandParts[0]
	// parameters := commandParts[1:];
	// nParameters := len(parameters);

	switch command {
	case PING:
		pingAction(conn)
		break
	default:
		pingAction(conn)
		break
	}

	conn.Close()
}

func readCommand(conn net.Conn) string {
	var commandBytes []byte
	bufferLength := 512
	readBuf := make([]byte, bufferLength)
	doBreak := false

	for {
		nRead, err := conn.Read(readBuf)

		if nRead < bufferLength {
			err = io.EOF
		}

		if err != nil {
			if err != io.EOF {
				fmt.Println("Failed to read command")
			}
			doBreak = true
		}

		commandBytes = append(commandBytes, readBuf...)

		if doBreak {
			break
		}
	}

	return string(commandBytes)
}

// Actions for Command

func pingAction(conn net.Conn) {
	response := []byte("+PONG\r\n")
	_, err := conn.Write(response)
	if err != nil {
		fmt.Println("Failed to write command")
	}
}
