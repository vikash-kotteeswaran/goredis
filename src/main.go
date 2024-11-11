package main

import (
	"bytes"
	"fmt"
	"goredis/src/config"
	"goredis/src/core"
	"os"
	"strconv"
	"syscall"
)

func main() {
	var connections map[int]core.Connection = map[int]core.Connection{}
	host, port, instSetupErr := core.SetupInstance()

	if instSetupErr != nil {
		fmt.Println("Error: " + instSetupErr.Error())
		os.Exit(1)
	}

	serverFd, socketCreateErr := core.SetupServer(host, port)
	if socketCreateErr != nil {
		fmt.Println("Error: " + socketCreateErr.Error())
		os.Exit(1)
	}

	fmt.Println("Server File Descriptor :: " + strconv.Itoa(serverFd) + " :: Created and Binded to Host and Port :: " + host + ":" + strconv.Itoa(port))

	listenErr := syscall.Listen(serverFd, config.CONCURRENCY_LIMIT)
	if listenErr != nil {
		fmt.Errorf("Error: " + "Failed to start listening from file descriptor")
		os.Exit(1)
	}

	fmt.Println("Server File Descriptor is Listening on :: " + host + ":" + strconv.Itoa(port))

	multiplexer, multiplexerErr := core.GetMultiplexer(config.CONCURRENCY_LIMIT)
	if multiplexerErr != nil {
		fmt.Println("Error: " + "Failed to create multiplexer for serving connections")
		os.Exit(1)
	}

	fmt.Println("Multiplexer Created")

	subscribeErr := multiplexer.Subscribe(serverFd)
	if subscribeErr != nil {
		fmt.Println("Error: " + "Failed to subscribe to multiplexer")
		os.Exit(1)
	}

	for !config.DOABORT {
		availableEvents, pollErr := multiplexer.Poll()
		if pollErr != nil {
			fmt.Println("Error: " + "Failed to fetch available file descriptors :: " + pollErr.Error())
			// os.Exit(1)
		}

		for _, availableEvent := range availableEvents {
			availableFd := availableEvent.Fd
			if availableFd == serverFd {
				connectedFd, _, connAccErr := syscall.Accept(serverFd)
				if connAccErr != nil {
					fmt.Println("Error: " + "Failed to accept connection")
					// os.Exit(1)
				}

				syscall.SetNonblock(connectedFd, true)
				multiplexer.Subscribe(connectedFd)

				connections[connectedFd] = core.Connection{Fd: connectedFd, Buffer: bytes.NewBuffer([]byte{})}
			} else {
				connection := connections[availableFd]
				core.Serve(connection)
				multiplexer.UnSubscribe(availableFd)
			}
		}
	}
}
