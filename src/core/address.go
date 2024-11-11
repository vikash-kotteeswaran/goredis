package core

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"syscall"
)

type Address struct {
	Host string
	Port int
}

func (address *Address) AbsorbSockAddr(sockAddr *syscall.Sockaddr) (bool, error) {
	sockAddrInet4 := (*sockAddr).(*syscall.SockaddrInet4)
	if sockAddrInet4 == nil {
		fmt.Println("Invalid replicaof address")
		return false, fmt.Errorf("Invalid replicaof address")
	}

	var addr string
	for addrIdx, addrEle := range sockAddrInet4.Addr {
		if addrIdx != 0 {
			addr += "."
		}
		addr += strconv.Itoa(int(addrEle))
	}

	address.Host = addr
	address.Port = sockAddrInet4.Port
	return true, nil
}

func (address *Address) Absorb(addrStr string) (bool, error) {
	addrSplit := strings.Split(addrStr, ":")
	if len(addrSplit) != 2 {
		fmt.Println("Invalid replicaof address")
		return false, fmt.Errorf("Invalid replicaof address")
	}

	address.Host = addrSplit[0]
	address.Port, _ = strconv.Atoi(addrSplit[1])

	return true, nil
}

func (addr *Address) String() string {
	var addrInfo string
	addrInfo = ("host:" + addr.Host + "\n" +
		"port:" + strconv.Itoa(addr.Port))
	return addrInfo
}

func (addr *Address) Equals(otherAddr string) bool {
	return addr.AddressStr() == otherAddr
}

func (addr *Address) EqualsHost(otherAddrHost string) bool {
	return addr.Host == otherAddrHost
}

func (addr *Address) EqualsNetAddr(insAddr net.Addr) bool {
	if insAddr == nil {
		return false
	}
	return addr.Equals(insAddr.String())
}

func (addr *Address) EqualsAddress(insAddr *Address) bool {
	if insAddr == nil {
		return false
	}
	return addr.Equals(insAddr.String())
}

func (addr *Address) EqualsAddressHost(insAddr *Address) bool {
	if insAddr == nil {
		return false
	}
	return addr.EqualsHost(insAddr.Host)
}

func (addr *Address) AddressStr() string {
	return addr.Host + ":" + strconv.Itoa(addr.Port)
}
