package main

import (
	"fmt"
	"goredis/app/configcore"
	"goredis/app/constants"
	"goredis/app/core"
	"net"
	"os"
	"strconv"
	"syscall"
)

func main() {
	var connections map[int]int = map[int]int{}
	host, port, instSetupErr := configcore.SetupInstance()

	if instSetupErr != nil {
		fmt.Println("Error: " + instSetupErr.Error())
		os.Exit(1)
	}

	serverFd, socketCreateErr := core.SetupServer(host, port)
	if socketCreateErr != nil {
		fmt.Println("Error: " + socketCreateErr.Error())
		os.Exit(1)
	}

	fmt.Println("Server File Descriptor :: " + strconv.Itoa(serverFd) + ":: Created and Binded to Host and Port :: " + host + ":" + strconv.Itoa(port))

	listenErr := syscall.Listen(serverFd, constants.CONCURRENCY_LIMIT)
	if listenErr != nil {
		fmt.Errorf("Error: " + "Failed to start listening from file descriptor")
		os.Exit(1)
	}

	fmt.Println("Server File Descriptor is Listening on :: " + host + ":" + strconv.Itoa(port))

	multiplexer, multiplexerErr := core.GetMultiplexer(constants.CONCURRENCY_LIMIT)
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

	for !constants.DOABORT {
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

				connections[connectedFd] = connectedFd
			} else {
				connectedFd := connections[availableFd]
				osFile := os.NewFile(uintptr(connectedFd), "clientConnectionFile")
				conn, connErr := net.FileConn(osFile)
				if connErr != nil {
					fmt.Println("Error: " + strconv.Itoa(connectedFd) + " :: " + connErr.Error())
					break
				}
				core.Serve(conn)
				multiplexer.UnSubscribe(connectedFd)
				osFile.Close()
				conn.Close()
			}
		}
	}

	// listener, err = net.Listen("tcp", "0.0.0.0:6379")
	// if err != nil {
	// 	fmt.Println("Failed to bind to port 6379")
	// 	os.Exit(1)
	// }

	// for {
	// 	concurrencySema <- 1
	// 	go core.ListenAndServe(&listener, &concurrencySema)
	// }
}
