package core

import (
	"bytes"
	"goredis/src/config"
	"io"
	"strconv"
	"syscall"
)

func ParseValue(conn *Connection) (interface{}, error) {
	var value interface{}
	var parseErr error

	bytesReadErr := readUntilCRLF(conn)
	if bytesReadErr != nil {
		return nil, bytesReadErr
	}

	firstByte, byteReadErr := conn.ReadByteFromBuffer()
	if byteReadErr != nil {
		if byteReadErr == io.EOF {
			return nil, nil
		}
		return nil, byteReadErr
	}

	switch firstByte {
	case '+':
		// For Simple String
		value, parseErr = ParseSimpleString(conn)
		break
	case '-':
		// For Simple Error
		value, parseErr = ParseSimpleString(conn)
		break
	case ':':
		value, parseErr = ParseInt64(conn)
		break
	case '*':
		value, parseErr = ParseArray(conn)
		break
	case '$':
		value, parseErr = ParseBulkString(conn)
		break
	default:
		break
	}

	if parseErr != nil {
		return nil, parseErr
	}

	return value, nil
}

func ParseSimpleString(conn *Connection) (interface{}, error) {
	simplestr, err := conn.ReadStringUntilFromBuffer('\r')
	if err != nil {
		return nil, err
	}
	// This removes '\r' that stays in simplestr
	slashRIdx := len(simplestr) - 1
	simplestr = simplestr[:slashRIdx]
	// This reads '\n' that remains after reading \r
	conn.ReadByteFromBuffer()

	return simplestr, nil
}

func ParseBulkString(conn *Connection) (interface{}, error) {
	strLen, lenParseErr := ParseInt64(conn)
	if lenParseErr != nil {
		return "", lenParseErr
	}

	strLenInt := int(strLen.(int64))

	if conn.Buffer.Len() < strLenInt {
		readUntilCRLF(conn)
	}

	readBulkStrBytes := make([]byte, strLenInt)
	_, readStrErr := conn.ReadFromBuffer(readBulkStrBytes)
	if readStrErr != nil {
		if readStrErr != io.EOF {
			return nil, readStrErr
		}
	}
	readBulkString := string(readBulkStrBytes)

	// This reads '\r' that remains after reading the bulk string
	conn.ReadByteFromBuffer()
	// This reads '\n' that remains after reading \r
	conn.ReadByteFromBuffer()

	return readBulkString, nil
}

func ParseInt64(conn *Connection) (interface{}, error) {
	intStr, valParseErr := ParseSimpleString(conn)
	if valParseErr != nil {
		return -1, valParseErr
	}
	intVal, intParseErr := strconv.ParseInt(intStr.(string), 10, 64)
	if intParseErr != nil {
		return -1, intParseErr
	}
	return intVal, nil
}

func ParseArray(conn *Connection) (interface{}, error) {
	arrLen, lenParseErr := ParseInt64(conn)
	if lenParseErr != nil {
		return nil, lenParseErr
	}
	arrLenInt := int(arrLen.(int64))

	array := make([]interface{}, arrLenInt)
	var arrEleParseErr error

	for idx := 0; idx < arrLenInt; idx++ {
		array[idx], arrEleParseErr = ParseValue(conn)
		if arrEleParseErr != nil {
			return nil, arrEleParseErr
		}
	}
	return array, nil
}

func readUntilCRLF(conn *Connection) error {
	tempBuf := make([]byte, config.CONNECTION_READ_BUF_SIZE)
	readAgainCount := 0
	for {
		if bytes.Contains(conn.Buffer.Bytes(), []byte{'\r', '\n'}) {
			break
		}

		nRead, readErr := conn.Read(tempBuf)
		if readErr != nil {
			if readErr == io.EOF {
				break
			}

			if readErr == syscall.EAGAIN {
				if readAgainCount >= config.CONNECTION_READ_RETRY {
					break
				}
				readAgainCount++
				continue
			}
			return readErr
		}

		if nRead > 0 {
			_, writeErr := conn.Buffer.Write(tempBuf[:nRead])
			if writeErr != nil {
				return writeErr
			}
		}

		if conn.Buffer.Len() > 0 {
			break
		}
	}

	return nil
}
