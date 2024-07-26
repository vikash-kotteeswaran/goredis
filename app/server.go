package main

import (
	"fmt"
	"goredis/app/core"
	"net"
	"os"
)

func main() {
	var listener net.Listener
	var connection net.Conn
	var err error

	listener, err = net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	connection, err = listener.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	core.Serve(connection)
}
