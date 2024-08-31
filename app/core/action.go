package core

import (
	"fmt"
	"net"
	"os"
)

func ListenAndServe(listenerPt *net.Listener, concurrSemaPt *(chan int)) {
	connection, err := (*listenerPt).Accept()
	defer connection.Close()

	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	fmt.Println("New connection :: ", connection.LocalAddr().String())

	Serve(connection)

	<-(*concurrSemaPt)
}
