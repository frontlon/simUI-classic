package controller

import (
	"VirtualNesGUI/code/utils"
	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
	"io"
	"os"
	"runtime"
	"strconv"
	"time"
)

//调用alert框
func ErrorMsg(w *window.Window, err string) *sciter.Value {

	if _, err := w.Call("errorBox", sciter.NewValue(err)); err != nil {
	}
	return sciter.NullValue();
}

//调用loading框
func Loading(w *window.Window, str string) *sciter.Value {
	if _, err := w.Call("startLoading", sciter.NewValue(str)); err != nil {
	}
	return sciter.NullValue();
}



/*
 写日志
*/
func WriteLog(str string) {
	fileName := "log.txt"

	cachePath := "./caches/"
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
