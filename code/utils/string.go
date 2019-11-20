package utils

import (
	"crypto/md5"
	"fmt"
)

/**
 * 字符串MD5
 **/
func Md5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	return fmt.Sprintf("%x", has)
}