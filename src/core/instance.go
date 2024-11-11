package core

import (
	"flag"
	"fmt"
	"goredis/src/config"
	"net"
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

var addr Address = Address{Host: config.HOST, Port: config.PORT}
var CurrInstance Instance = Instance{Name: "default", Role: MASTER_ROLE, Addr: addr}

type Instance struct {
	Name       string
	Role       string
	Addr       Address
	ReplId     string
	ReplOffset int
	Replicas   Instances
	MasterAddr Address
}

type Instances []Instance

func (instance *Instance) String() string {
	var instanceInfo string
	instanceInfo = ("name:" + instance.Name + "\n" +
		instance.Addr.String() + "\n" +
		"master_repl_id:" + instance.ReplId + "\n" +
		"master_repl_offset:" + strconv.Itoa(instance.ReplOffset) + "\n" +
		instance.Replicas.String())
	return instanceInfo
}

func (instance *Instance) GetInfo(infoType string) string {
	var instanceInfo string
	switch infoType {
	case REPL_INFO:
		instanceInfo = instance.Replicas.String()
	case CURR_INST_INFO:
		instanceInfo = instance.String()
	default:
		instanceInfo = instance.String()
	}
	return instanceInfo
}

func (ins *Instance) HasReplicas() bool {
	return len(ins.Replicas) > 0
}

func (ins *Instance) AddReplica(newReplHost string, newReplPort int) bool {
	var repl Instance

	repl.Name = "replica" + strconv.Itoa(CurrInstance.ReplOffset)
	repl.Addr.Host = newReplHost
	repl.Addr.Port = newReplPort
	repl.Role = SLAVE_ROLE
	repl.ReplId = "23jk4b4k36bk45jb6k3b4ib123k5bjk23btibd"
	repl.ReplOffset = 0
	repl.MasterAddr = CurrInstance.Addr

	ins.Replicas = append(ins.Replicas, repl)

	return true
}

func (repls *Instances) String() string {
	var replsInfo string

	for _, repl := range *repls {
		replsInfo = repl.String()
	}
	return replsInfo
}

func (repls *Instances) Contains(insAddr string) bool {
	isContained := false
	for _, repl := range *repls {
		if repl.Addr.AddressStr() == insAddr {
			isContained = true
			break
		}
	}

	return isContained
}

func (repls *Instances) ContainsHost(insAddrHost string) bool {
	isContained := false
	for _, repl := range *repls {
		if repl.Addr.Host == insAddrHost {
			isContained = true
			break
		}
	}

	return isContained
}

func (repls *Instances) ContainsNetAddr(insAddr net.Addr) bool {
	if insAddr == nil {
		return false
	}
	return repls.Contains(insAddr.String())
}

func (repls *Instances) ContainsAddress(insAddr *Address) bool {
	if insAddr == nil {
		return false
	}
	return repls.Contains(insAddr.String())
}

func (repls *Instances) ContainsAddressHost(insAddr *Address) bool {
	if insAddr == nil {
		return false
	}
	return repls.ContainsHost(insAddr.Host)
}

func SetupInstance() (string, int, error) {
	var name string
	var host string
	var port int
	var replicaofAddr string

	flag.StringVar(&name, "name", "default", "name of the instance")
	flag.StringVar(&host, "host", config.HOST, "host address for the instance")
	flag.IntVar(&port, "port", config.PORT, "port address for the instance")
	flag.StringVar(&replicaofAddr, "replicaof", "", "replica of address for the instance")
	flag.Parse()

	CurrInstance.Name = name
	CurrInstance.Addr.Host = host
	CurrInstance.Addr.Port = port
	if replicaofAddr != "" {
		CurrInstance.Role = SLAVE_ROLE
		CurrInstance.ReplId = "23jk4b4k36bk45jb6k3b4ib123k5bjk23btibd"
		CurrInstance.ReplOffset = 0
		absorbed, absorbErr := CurrInstance.MasterAddr.Absorb(replicaofAddr)

		if !absorbed && absorbErr != nil {
			fmt.Println("Invalid replicaof address")
			return "", -1, absorbErr
		}

		err := setupReplica(&CurrInstance.Addr, &CurrInstance.MasterAddr)

		if err != nil {
			return "", -1, err
		}
	} else {
		CurrInstance.ReplId = "23jk4b4k36bk45jb6k3b4ib123k5bjk23btibd"
		CurrInstance.ReplOffset = 0
		CurrInstance.Role = MASTER_ROLE
	}

	return host, port, nil
}

func setupReplica(replicaAddr *Address, replicaOfAddr *Address) error {
	fmt.Println("Setting up Replica at " + replicaAddr.AddressStr())

	// PING master
	pingErr := pingMaster(replicaOfAddr)
	if pingErr != nil {
		return pingErr
	}

	// REPLCONF with master
	replConfErr := replConfigureWithMaster(replicaAddr, replicaOfAddr)
	if replConfErr != nil {
		return replConfErr
	}

	// PSYNC with master
	pSyncErr := psyncWithMaster(replicaOfAddr)
	if pSyncErr != nil {
		return pSyncErr
	}

	fmt.Println("Replica has been setup Successfull at " + replicaAddr.AddressStr())

	fmt.Println()

	return nil
}

func pingMaster(replicaOfAddr *Address) error {
	fmt.Println("Pinging to Master at " + replicaOfAddr.AddressStr())

	respStr := HitFromServer([]interface{}{"PING"}, replicaOfAddr)
	if respStr != "+PONG\r\n" {
		return fmt.Errorf("ReplicaOf Address PING failed :: PING reponse :: ", respStr)
	}

	fmt.Println("Ping to Master Successfull")

	return nil
}

func replConfigureWithMaster(replicaAddr *Address, replicaOfAddr *Address) error {
	fmt.Println("Replica Configuring with Master")

	respStr := HitFromServer([]interface{}{"REPLCONF", "listening-host", replicaAddr.Host, "listening-port", replicaAddr.Port}, replicaOfAddr)
	if respStr != "+OK\r\n" {
		return fmt.Errorf("ReplicaOf Address REPLCONF listening-port failed :: REPLCONF reponse :: ", respStr)
	}

	respStr = HitFromServer([]interface{}{"REPLCONF", "capa", "eof", "capa", "psync2"}, replicaOfAddr)
	if respStr != "+OK\r\n" {
		return fmt.Errorf("ReplicaOf Address REPLCONF capa failed :: REPLCONF reponse :: ", respStr)
	}

	fmt.Println("Replica Configuring with Master Successfull")

	return nil
}

func psyncWithMaster(replicaOfAddr *Address) error {
	fmt.Println("Psyncing with Master at " + replicaOfAddr.AddressStr())

	respStr := HitFromServer([]interface{}{"PSYNC", "?", "-1"}, replicaOfAddr)
	if !strings.Contains(respStr, "FULLRESYNC") {
		return fmt.Errorf("ReplicaOf Address PSYNC failed :: PSYNC reponse :: ", respStr)
	}

	fmt.Println("Psyncing with Master Successfull")

	return nil
}
