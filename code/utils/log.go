package utils

import (
	"io"
	"os"
	"runtime"
	"strconv"
	"time"
)

/*
 写日志
*/
func WriteLog(str string) {
	fileName := "log.txt"

	cachePath := "./cache/"
	if !IsExist(cachePath) {
		if err := CreateDir(cachePath); err != nil {
			return
		}
	}

	f, _ := os.OpenFile(cachePath+fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	c_date := time.Now().Format("2006-01-02 15:04:05")
	_, c_file, c_line, _ := runtime.Caller(1)

	content := c_date + "\t"               //日期
	content += c_file + "\t"               //文件
	content += strconv.Itoa(c_line) + "\t" //行号
	content += str + "\r\n"                //内容

	if _, err := io.WriteString(f, content); err != nil {
		return
	}

	defer f.Close()
	return
}
