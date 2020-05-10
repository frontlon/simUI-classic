package utils

import (
	"encoding/base64"
	"strings"
)

//base64编码
func Base64Encode(s string) string {
	if s == "" {
		return ""
	}
	encodeString := base64.StdEncoding.EncodeToString([]byte(s))
	return strings.Replace(encodeString, "=", "_", -1)
}

//base64解码
func Base64Decode(s string) string {
	if s == "" {
		return ""
	}

	s = strings.Replace(s, "_", "=", -1)
	decodeBytes, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return ""
	}
	return string(decodeBytes)
}