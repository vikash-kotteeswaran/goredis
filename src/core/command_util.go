package core

import "fmt"

func WriteResponseToConn(conn *Connection, response []byte) {
	_, err := conn.Write(response)
	if err != nil {
		fmt.Println("Failed to perform write to connection :: ", err)
	}
}
