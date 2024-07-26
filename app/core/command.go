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

func Serve(conn net.Conn) error {
	payload, err := readPayload(conn)
	requests := strings.Split(payload, "\n")

	for _, request := range requests {
		commandParts := strings.Split(request, " ")
		command := commandParts[0]
		// parameters := commandParts[1:];
		// nParameters := len(parameters);

		switch command {
		case PING:
			pingAction(conn)
			break
		default:
			break
		}
	}

	if err == io.EOF {
		conn.Close()
	}

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

// Actions for Command

func pingAction(conn net.Conn) {
	response := []byte("+PONG\r\n")
	_, err := conn.Write(response)
	if err != nil {
		fmt.Println("Failed to write command")
	}
}
