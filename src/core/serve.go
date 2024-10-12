package core

import (
	"fmt"
)

func Serve(conn Connection) error {
	err := ParseActions(&conn)
	if err != nil {
		conn.Write([]byte(UnParseString(err.Error(), true)))
	}

	for _, action := range conn.Actions {
		action.Execute()
	}

	fmt.Println("Connection :: ", conn.GetConnAddr().String(), " has been served")

	conn.Close()
	return err
}
