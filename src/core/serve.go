package core

import (
	"fmt"
)

func Serve(conn Connection) error {
	err := ParseActions(conn)
	if err != nil {
		return err
	}

	for _, action := range conn.Actions {
		action.Execute()
	}

	fmt.Println("Connection :: ", conn.GetConnAddr().String(), " has been served")

	return err
}
