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
var PIC_EXTS = []string{"png", "jpg", "gif", "jpeg", "bmp"}; //支持的图片类型

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

	result := &exec.Cmd{}

	//这个写法牛不牛逼~但有更好的吗？有的话请告诉我。
	switch len(cmd) {
	case 0:
		result = exec.Command(exeFile)
	case 1:
		result = exec.Command(exeFile, cmd[0])
	case 2:
		result = exec.Command(exeFile, cmd[0], cmd[1])
	case 3:
		result = exec.Command(exeFile, cmd[0], cmd[1], cmd[2])
	case 4:
		result = exec.Command(exeFile, cmd[0], cmd[1], cmd[2], cmd[3])
	case 5:
		result = exec.Command(exeFile, cmd[0], cmd[1], cmd[2], cmd[3], cmd[4])
	case 6:
		result = exec.Command(exeFile, cmd[0], cmd[1], cmd[2], cmd[3], cmd[4], cmd[5])
	case 7:
		result = exec.Command(exeFile, cmd[0], cmd[1], cmd[2], cmd[3], cmd[4], cmd[5], cmd[6])
	case 8:
		result = exec.Command(exeFile, cmd[0], cmd[1], cmd[2], cmd[3], cmd[4], cmd[5], cmd[6], cmd[7])
	case 9:
		result = exec.Command(exeFile, cmd[0], cmd[1], cmd[2], cmd[3], cmd[4], cmd[5], cmd[6], cmd[7], cmd[8])
	case 10:
		result = exec.Command(exeFile, cmd[0], cmd[1], cmd[2], cmd[3], cmd[4], cmd[5], cmd[6], cmd[7], cmd[8], cmd[9])
	case 11:
		result = exec.Command(exeFile, cmd[0], cmd[1], cmd[2], cmd[3], cmd[4], cmd[5], cmd[6], cmd[7], cmd[8], cmd[9], cmd[10])
	case 12:
		result = exec.Command(exeFile, cmd[0], cmd[1], cmd[2], cmd[3], cmd[4], cmd[5], cmd[6], cmd[7], cmd[8], cmd[9], cmd[10], cmd[11])
	case 13:
		result = exec.Command(exeFile, cmd[0], cmd[1], cmd[2], cmd[3], cmd[4], cmd[5], cmd[6], cmd[7], cmd[8], cmd[9], cmd[10], cmd[11], cmd[12])
	case 14:
		result = exec.Command(exeFile, cmd[0], cmd[1], cmd[2], cmd[3], cmd[4], cmd[5], cmd[6], cmd[7], cmd[8], cmd[9], cmd[10], cmd[11], cmd[12], cmd[13])
	case 15:
		result = exec.Command(exeFile, cmd[0], cmd[1], cmd[2], cmd[3], cmd[4], cmd[5], cmd[6], cmd[7], cmd[8], cmd[9], cmd[10], cmd[11], cmd[12], cmd[13], cmd[14])
	case 16:
		result = exec.Command(exeFile, cmd[0], cmd[1], cmd[2], cmd[3], cmd[4], cmd[5], cmd[6], cmd[7], cmd[8], cmd[9], cmd[10], cmd[11], cmd[12], cmd[13], cmd[14], cmd[15])
	case 17:
		result = exec.Command(exeFile, cmd[0], cmd[1], cmd[2], cmd[3], cmd[4], cmd[5], cmd[6], cmd[7], cmd[8], cmd[9], cmd[10], cmd[11], cmd[12], cmd[13], cmd[14], cmd[15], cmd[16])
	case 18:
		result = exec.Command(exeFile, cmd[0], cmd[1], cmd[2], cmd[3], cmd[4], cmd[5], cmd[6], cmd[7], cmd[8], cmd[9], cmd[10], cmd[11], cmd[12], cmd[13], cmd[14], cmd[15], cmd[16], cmd[17])
	case 19:
		result = exec.Command(exeFile, cmd[0], cmd[1], cmd[2], cmd[3], cmd[4], cmd[5], cmd[6], cmd[7], cmd[8], cmd[9], cmd[10], cmd[11], cmd[12], cmd[13], cmd[14], cmd[15], cmd[16], cmd[17], cmd[18])
	case 20:
		result = exec.Command(exeFile, cmd[0], cmd[1], cmd[2], cmd[3], cmd[4], cmd[5], cmd[6], cmd[7], cmd[8], cmd[9], cmd[10], cmd[11], cmd[12], cmd[13], cmd[14], cmd[15], cmd[16], cmd[17], cmd[18], cmd[19])
	}

	if err := result.Start(); err != nil {
		return err
	}

	return nil
}
