package utils

import (
	"fmt"
	"strconv"
)

/*
 传入任意类型，将数据转换为字符串
*/
func ToString(str interface{}) string {
	var val string
	switch t := str.(type) {
	case string:
		val = t
	case int:
		val = strconv.Itoa(t)
	case int8:
		val = strconv.Itoa(int(t))
	case int16:
		val = strconv.Itoa(int(t))
	case int32:
		val = strconv.Itoa(int(t))
	case int64:
		val = strconv.Itoa(int(t))
	case uint:
		val = strconv.Itoa(int(t))
	case uint8:
		val = strconv.Itoa(int(t))
	case uint16:
		val = strconv.Itoa(int(t))
	case uint32:
		val = strconv.Itoa(int(t))
	case uint64:
		val = strconv.Itoa(int(t))
	case float32:
		val = strconv.FormatFloat(float64(t), 'f', -1, 64)
	case float64:
		val = strconv.FormatFloat(t, 'f', -1, 64)
	case []uint8: //[]byte类型
		val = string(t)
	case bool:
		if t == true {
			val = "true"
		} else {
			val = "false"
		}
	default:
		val = t.(string)
	}
	return val
}

/**
 传入任意类型，将数据转换为整型
**/
func ToInt(str interface{}) int {
	var val int
	switch t := str.(type) {
	case int:
		val = int(t)
	case int8:
		val = int(t)
	case int16:
		val = int(t)
	case int32:
		val = int(t)
	case int64:
		val = int(t)
	case uint:
		val = int(t)
	case uint8:
		val = int(t)
	case uint16:
		val = int(t)
	case uint32:
		val = int(t)
	case uint64:
		val = int(t)
	case float32:
		val = int(t)
	case float64:
		val = int(t)
	case string:
		val, _ = strconv.Atoi(str.(string))
	case bool:
		if str.(bool) == true {
			val = 1
		} else {
			val = 0
		}
	default:
		fmt.Println("cvcvcvcvcv")

		val = 0
	}
	return val
}
