package core

import (
	"bytes"
	"fmt"
	"io"
	"syscall"
)

type Connection struct {
	io.ReadWriter
	Fd      int
	Buffer  *bytes.Buffer
	Actions Actions
	Meta    ConnMeta

	addr *Address
}

type ConnMeta struct {
	BytesRead   int
	FromReplica bool
	FromMaster  bool
}

func (conn *Connection) Read(readBytes []byte) (int, error) {
	nRead, readErr := syscall.Read(conn.Fd, readBytes)
	conn.Meta.BytesRead += nRead
	return nRead, readErr
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

func (conn *Connection) SetConnAddr() error {
	socketAddr, err := syscall.Getpeername(conn.Fd)
	if err != nil {
		err := fmt.Errorf("Encountered error while getting connection address :: ", err.Error())
		return err
	}
	conn.addr = &Address{}
	conn.addr.AbsorbSockAddr(&socketAddr)

	return nil
}

func (conn *Connection) GetConnAddr() *Address {
	if conn.addr == nil || (conn.addr == nil && conn.addr.Host == "" && conn.addr.Port < 0) {
		conn.SetConnAddr()
	}
	return conn.addr
}

func (conn *Connection) SetConnectionFrom() {
	conn.Meta.FromMaster = CurrInstance.MasterAddr.EqualsAddressHost(conn.addr)
	conn.Meta.FromReplica = CurrInstance.Replicas.ContainsAddressHost(conn.addr)
}

func (conn *Connection) Close() {
	syscall.Close(conn.Fd)
}
