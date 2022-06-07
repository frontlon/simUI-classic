package utils

import (
	"strings"
)

var prefix = "t0Ic" //前缀

//创建rom唯一id
func CreateRomUniqId(t string, s int64) string {
	return prefix + "_" + t + "_" + ToString(s)
}

//分割rom唯一id
func SplitRomUniqId(uniqId string) (string, string) {
	if uniqId == "" {
		return "", ""
	}
	sli := strings.Split(uniqId, "_")
	return sli[1], sli[2] //返回 时间,大小
}

//检查是不是唯一ID
func HasRomUniqId(uniqId string) bool {
	if uniqId == "" {
		return false
	}
	sli := strings.Split(uniqId, "_")

	if sli[0] == prefix {
		return true
	}
	return false
}
