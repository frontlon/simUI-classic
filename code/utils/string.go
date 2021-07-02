package utils

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"simUI/code/utils/pinyin"
	"strings"
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
	str = strings.ToLower(str)
	str = strings.ReplaceAll(str, " ", "")
	return str
}

func Addslashes(str string) string {
	var buf bytes.Buffer
	for _, char := range str {
		switch char {
		case '\'', '"', '\\':
			buf.WriteRune('\\')
		}
		buf.WriteRune(char)
	}
	return buf.String()
}
