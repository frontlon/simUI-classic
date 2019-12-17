package main

import (
	"VirtualNesGUI/code/db"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

var constSeparator = "__"                                    //rom子分隔符
var DOC_EXTS = []string{".txt", ".md", ".html", ".htm"}      //doc文档支持的扩展名
var PIC_EXTS = []string{".png", ".jpg", ".gif", ".jpeg", ".bmp"}; //支持的图片类型

type RomDetail struct {
	Info            *db.Rom
	DocContent      string
	StrategyContent string
	Sublist         []*db.Rom
}

/**
 * 读取游戏介绍文本
 **/
func getDocContent(f string) string {
	if f == "" {
		return ""
	}
	text, err := ioutil.ReadFile(f)
	content := ""
	if err != nil {
		return content
	}
	//enc := mahonia.NewDecoder("gbk")
	//content = enc.ConvertString(string(text))
	content = string(text)
	return content
}

/**
 * 运行游戏
 **/
func runGame(exeFile string, cmd []string) error {

	//更改程序运行目录
	if err := os.Chdir(filepath.Dir(exeFile)); err != nil {
		return err
	}

	result := exec.Command(exeFile, cmd...)

	if err := result.Start(); err != nil {
		return err
	}

	return nil
}
