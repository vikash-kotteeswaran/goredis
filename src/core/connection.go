package core

import (
	"net"
	"syscall"
)

type Connection struct {
	Fd      int
	Actions Actions
}

func (conn *Connection) Read(readBytes []byte) (int, error) {
	return syscall.Read(conn.Fd, readBytes)
}

func (conn *Connection) Write(writeBytes []byte) (int, error) {
	return syscall.Write(conn.Fd, writeBytes)
}

func (conn *Connection) GetConnAddr() net.IP {
	addr, _ := syscall.GetsockoptInet4Addr(conn.Fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR)
	ip := net.IPv4(addr[0], addr[1], addr[2], addr[3])
	return ip
}
