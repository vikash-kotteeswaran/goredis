package core

import (
	"reflect"
)

func ParseActions(conn Connection) error {
	// value, valtype := ParseValue()
	// inval := value.(valtype)
	return nil
}

func Parse(conn Connection) (interface{}, reflect.Type) {
	// conn.Read()
	return string(""), reflect.TypeFor[int]()
}

func ParseSimpleString(buffer []byte) (interface{}, reflect.Type) {
	return string(""), reflect.TypeFor[int]()
}

func ParseBulkString() (interface{}, reflect.Type) {
	return string(""), reflect.TypeFor[int]()
}

func ParseLength() (interface{}, reflect.Type) {
	return string(""), reflect.TypeFor[int]()
}

func ParseArray() (interface{}, reflect.Type) {
	return string(""), reflect.TypeFor[int]()
}
