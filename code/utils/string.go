package utils

import (
	"crypto/md5"
	"fmt"
	"simUI/code/utils/pinyin"
)

/**
 * 字符串MD5
 **/
func Md5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	return fmt.Sprintf("%x", has)
}

/**
 * 字符转拼音
 **/
func TextToPinyin(str string) string {
	str, err := pinyin.New(str).Split("").Mode(pinyin.WithoutTone).Convert()
	if err != nil {
		// 错误处理
	}
	return str
}
