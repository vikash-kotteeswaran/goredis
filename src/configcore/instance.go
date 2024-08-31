package configcore

import (
	"flag"
	"fmt"
	"goredis/app/constants"
	"strconv"
	"strings"
)

// Info Type
const (
	REPL_INFO      string = "replication"
	CURR_INST_INFO string = "instance"
)

// Role Types
const (
	MASTER_ROLE = "master"
	SLAVE_ROLE  = "slave"
)

var replication Replication = Replication{Role: MASTER_ROLE}
var addr Address = Address{Host: constants.HOST, Port: constants.PORT}
var CurrInstance Instance = Instance{Name: "default", Addr: addr, Repl: replication}

type Address struct {
	Host string
	Port int
}

type Instance struct {
	Name string
	Addr Address
	Repl Replication
}

type Replication struct {
	Role          string
	ReplicaOfAddr Address
}

func (addr *Address) String() string {
	var addrInfo string
	addrInfo = ("host:" + addr.Host + "\n" +
		"port:" + strconv.Itoa(addr.Port))
	return addrInfo
}

func (instance *Instance) String() string {
	var instanceInfo string
	instanceInfo = ("name:" + instance.Name + "\n" +
		instance.Addr.String() + "\n" +
		instance.Repl.String())
	return instanceInfo
}

func (instance *Instance) GetInfo(infoType string) string {
	var instanceInfo string
	switch infoType {
	case REPL_INFO:
		instanceInfo = instance.Repl.String()
	case CURR_INST_INFO:
		instanceInfo = instance.String()
	default:
		instanceInfo = instance.String()
	}
	return instanceInfo
}

func (repl *Replication) String() string {
	var replInfo string
	replInfo = ("role:" + repl.Role + "\n" +
		"<replicaof>\n" + repl.ReplicaOfAddr.String() + "\n<\\replicaof>\n")
	return replInfo
}

func SetupInstance() (string, int, error) {
	var name string
	var host string
	var port int
	var replicaofAddr string

	flag.StringVar(&name, "name", "default", "name of the instance")
	flag.StringVar(&host, "host", constants.HOST, "host address for the instance")
	flag.IntVar(&port, "port", constants.PORT, "port address for the instance")
	flag.StringVar(&replicaofAddr, "replicaof", "", "replica of address for the instance")
	flag.Parse()

	CurrInstance.Name = name
	CurrInstance.Addr.Host = host
	CurrInstance.Addr.Port = port
	if replicaofAddr != "" {
		CurrInstance.Repl.Role = SLAVE_ROLE
		replicaofAddrSplit := strings.Split(replicaofAddr, ":")
		if len(replicaofAddrSplit) != 2 {
			fmt.Println("Invalid replicaof address")
			return "", -1, fmt.Errorf("Invalid replicaof address")
		}
		CurrInstance.Repl.ReplicaOfAddr.Host = replicaofAddrSplit[0]
		CurrInstance.Repl.ReplicaOfAddr.Port, _ = strconv.Atoi(replicaofAddrSplit[1])

		err := setupReplica(CurrInstance.Repl.ReplicaOfAddr)

		if err != nil {
			return "", -1, err
		}
	} else {
		CurrInstance.Repl.Role = MASTER_ROLE
	}

	return host, port, nil
}

func setupReplica(replicaOfAddr Address) error {

	return nil
}
