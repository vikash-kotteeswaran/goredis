package main

import (
	"flag"
	"fmt"
	"goredis/cli/config"
	"goredis/src/core"
	"os"
	"strconv"
)

func main() {
	var host string
	var port int

	flag.StringVar(&host, "host", config.HOST, "host address for the instance")
	flag.IntVar(&port, "port", config.PORT, "port address for the instance")
	flag.Parse()

	args := os.Args[1:]

	goredisParams := make([]interface{}, 0)
	for idx, _ := range args {
		prevParam, param := "", args[idx]
		if idx > 0 {
			prevParam = args[idx-1]
		}

		if param == "--host" || param == "-h" || param == "--port" || param == "-p" {
			continue
		}

		if prevParam == "--host" || param == "-h" || prevParam == "--port" || param == "-p" {
			continue
		}

		intParam, parseErr := strconv.ParseInt(param, 10, 64)
		if parseErr != nil {
			goredisParams = append(goredisParams, param)
		} else {
			goredisParams = append(goredisParams, intParam)
		}
	}

	response := core.HitFromServer(goredisParams, core.Address{Host: host, Port: port})

	fmt.Println(response)
}
