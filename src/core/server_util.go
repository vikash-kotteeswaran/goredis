package core

import (
	"errors"
	"fmt"
	"goredis/src/config"
	"io"
	"net"
	"os"
	"syscall"
)

func SetupServer(host string, port int) (int, error) {
	serverFd, socketCreateErr := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if socketCreateErr != nil {
		return -1, errors.New("Failed to create socket for server connections :: " + socketCreateErr.Error())
	}

	sockOptErr := syscall.SetsockoptInt(serverFd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	if sockOptErr != nil {
		return -1, errors.New("Failed to set socket options for server connections :: " + sockOptErr.Error())
	}

	nonBlockSockErr := syscall.SetNonblock(serverFd, true)
	if nonBlockSockErr != nil {
		return -1, errors.New("Failed to set socket as Non Blockings :: " + nonBlockSockErr.Error())
	}

	hostIP4 := net.ParseIP(host)
	hostIP4bytes := [4]byte{hostIP4[0], hostIP4[1], hostIP4[2], hostIP4[3]}

	bindErr := syscall.Bind(serverFd, &syscall.SockaddrInet4{Port: port, Addr: hostIP4bytes})
	if bindErr != nil {
		return -1, errors.New("Failed to bind host and port to socket :: " + bindErr.Error())
	}

	return serverFd, nil
}

func HitFromServer(params []interface{}, toServerAddr *Address) string {
	unparsed := UnParseValue(params, false, false)

	conn, connErr := net.Dial("tcp", toServerAddr.AddressStr())
	if connErr != nil {
		fmt.Println("Error: ", connErr.Error())
		return ""
	}

	conn.Write([]byte(unparsed))

	tempBuf := make([]byte, config.CONNECTION_READ_BUF_SIZE)
	response := ""
	for {
		nRead, readErr := conn.Read(tempBuf)
		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			fmt.Println("Error: ", readErr.Error())
			os.Exit(1)
		}
		response += string(tempBuf[:nRead])
	}

	// conn.Close()

	return response
}
