package core

import "fmt"

func Serve(conn Connection) error {
	conn.SetConnAddr()
	conn.SetConnectionFrom()

	fmt.Println("\nServing Request from :: ", conn.GetConnAddr().AddressStr())

	err := ParseActions(&conn)
	if err != nil {
		conn.Write([]byte(UnParseString(err.Error(), true)))
	}

	for _, action := range conn.Actions {
		action.Execute()
	}

	if !conn.Meta.FromReplica && !conn.Meta.FromMaster {
		fmt.Println("Connection :: ", conn.GetConnAddr().AddressStr(), " has been served")
	} else if conn.Meta.FromMaster {
		fmt.Println("Request from Master :: ", conn.GetConnAddr().AddressStr(), " has been served")
	} else {
		fmt.Println("Request from Replica :: ", conn.GetConnAddr().AddressStr(), " has been served")
	}

	conn.Close()
	return err
}
