package main

import (
	"github.com/Lofanmi/pinyin-golang/pinyin"
	"os"
	"path"
	"strings"
)
var dict = pinyin.NewDict()

/**
 * 去掉rom扩展名，从文件名中读取Rom名称
 **/
func GetFileName(filename string) string {
	str := strings.TrimSuffix(filename, path.Ext(filename))
	return strings.ToLower(str)
}

/**
 * 字符转拼音
 **/
func TextToPinyin(str string) string {
	return strings.Replace(dict.Sentence(str).ASCII()," ","",-1)
}

/**
 * 检测文件是否存在（文件夹也返回false）
 **/
func Exists(path string) bool {

	if path == "" {
		return false
	}

	finfo, err := os.Stat(path)
	isset := false
	if err != nil || finfo.IsDir() == true {
		isset = false
	} else {
		isset = true
	}
	return isset
}

/**
 * 检测文件夹是否存在（存在文件也返回false）
 **/
func ExistsFolder(path string) bool {
	ff, err := os.Stat(path)
	if err != nil {
		return false
	}
	if ff.IsDir() == false{
		return false
	}
	return true
}