package main

import (
	"VirtualNesGUI/code/utils"
	"github.com/Lofanmi/pinyin-golang/pinyin"
	"os"
	"strings"
)

var dict = pinyin.NewDict()

/**
 * 字符转拼音
 **/
func TextToPinyin(str string) string {
	return strings.Replace(dict.Sentence(str).ASCII(), " ", "", -1)
}


/**
 * 读取文件唯一标识
 **/
func GetFileUniqId(f os.FileInfo) string {
	str := f.Name() + utils.ToString(f.Size()) + f.ModTime().String()
	return utils.Md5(str)
}
