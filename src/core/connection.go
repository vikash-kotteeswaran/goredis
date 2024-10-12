package core

import (
	"bytes"
	"io"
	"net"
	"syscall"
)

type Connection struct {
	io.ReadWriter
	Fd      int
	Buffer  *bytes.Buffer
	Actions Actions
}

func (conn *Connection) Read(readBytes []byte) (int, error) {
	return syscall.Read(conn.Fd, readBytes)
}

func (conn *Connection) Write(writeBytes []byte) (int, error) {
	return syscall.Write(conn.Fd, writeBytes)
}

func (conn *Connection) WriteToBuffer(writeBytes []byte) (int, error) {
	return conn.Buffer.Write(writeBytes)
}

func (conn *Connection) ReadFromBuffer(readBytes []byte) (int, error) {
	return conn.Buffer.Read(readBytes)
}

func (conn *Connection) ReadByteFromBuffer() (byte, error) {
	return conn.Buffer.ReadByte()
}

func (conn *Connection) ReadStringUntilFromBuffer(byte byte) (string, error) {
	return conn.Buffer.ReadString(byte)
}

func (conn *Connection) GetConnAddr() net.IP {
	addr, _ := syscall.GetsockoptInet4Addr(conn.Fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR)
	ip := net.IPv4(addr[0], addr[1], addr[2], addr[3])
	return ip
}

func (conn *Connection) Close() {
	syscall.Close(conn.Fd)
}
