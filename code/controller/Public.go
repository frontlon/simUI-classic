package controller

import (
	"VirtualNesGUI/code/config"
	"VirtualNesGUI/code/utils"
	"github.com/sciter-sdk/go-sciter"
	"io"
	"os"
	"runtime"
	"strconv"
	"time"
)

//调用alert框
func ErrorMsg(err string) *sciter.Value {

	if _, err := config.Window.Call("errorBox", sciter.NewValue(err)); err != nil {
	}
	return sciter.NullValue();
}

//调用loading框
func Loading(str string,platform string) *sciter.Value {
	if _, err := config.Window.Call("startLoading", sciter.NewValue(str),sciter.NewValue(platform)); err != nil {
	}
	return sciter.NullValue();
}



/*
 写日志
*/
func WriteLog(str string) {
	fileName := "log.txt"

	cachePath := "./cache/"
	if !utils.IsExist(cachePath) {
		if err := utils.CreateDir(cachePath); err != nil {
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
