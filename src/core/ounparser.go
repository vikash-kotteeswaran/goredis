package core

import (
	"fmt"
	"strconv"
)

func UnParseValue(value interface{}, isSimpleStr bool, isRDB bool) string {
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
	case []byte:
		if isRDB {
			unparsed = UnParseRDB(value.([]byte))
		} else {
			fmt.Println("Unknown Value")
		}
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
		unparsed += UnParseValue(element, false, false)
	}

	return unparsed
}

func UnParseRDB(rdb []byte) string {
	return "$" + strconv.Itoa(len(rdb)) + CLRF + string(rdb)
}
