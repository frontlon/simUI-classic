package utils

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"math/rand"
	"regexp"
	"simUI/code/utils/pinyin"
	"strings"
	"time"
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

// 输入一个范围，生成范围随机数
func RandInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

// 输入一个长度，返回一个随机字符串
func RandStr(min, max int) string {
	rand.Seed(time.Now().UnixNano())

	length := rand.Intn(max-min+1) + min
	result := make([]byte, length)

	for i := 0; i < length; i++ {
		result[i] = byte(rand.Intn(26) + 97)
	}

	return string(result)
}

// 检查一个字符串是不是纯英文
func IsEnglish(str string) bool {
	reg := regexp.MustCompile(`[^a-zA-Z0-9!"# $%&'()*+,-./:;<=>?@[\\\]^_` + "`" + `{|}~]+`)
	enStr := reg.ReplaceAllString(str, "")
	return len(enStr) == len(str)
}

// 首字母大写
func ToTitleCase(str string) string {
	if str == "" {
		return ""
	}
	return strings.ToUpper(string(str[0])) + str[1:]
}
