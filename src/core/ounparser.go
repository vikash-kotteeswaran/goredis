package core

import (
	"fmt"
	"strconv"
)

const CLRF = "\r\n"

func UnParseValue(value interface{}, isSimpleStr bool) string {
	var unparsed string

	switch value.(type) {
	case string:
		unparsed = UnParseString(value.(string), isSimpleStr)
		break
	case int:
		unparsed = UnParseInt(int64(value.(int)))
		break
	case int64:
		unparsed = UnParseInt(value.(int64))
		break
	case []interface{}:
		unparsed = UnParseArray(value.([]interface{}))
		break
	default:
		fmt.Println("Unknown Value")
		break
	}

	return unparsed
}

func UnParseString(value string, isSimpleStr bool) string {
	unparsed := ""
	if isSimpleStr {
		unparsed = "+" + value + CLRF
	} else {
		unparsed = "$" + strconv.Itoa(len(value)) + CLRF + value + CLRF
	}
	return unparsed
}

func UnParseInt(value int64) string {
	unparsed := ":" + strconv.FormatInt(value, 10) + CLRF
	return unparsed
}

func UnParseArray(value []interface{}) string {
	unparsed := "*" + strconv.Itoa(len(value)) + CLRF

	for _, element := range value {
		unparsed += UnParseValue(element, false)
	}

	return unparsed
}
