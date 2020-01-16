package controller

import (
	"VirtualNesGUI/code/db"
	"VirtualNesGUI/code/utils"
	"archive/zip"
	"github.com/axgle/mahonia"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var LAST_PROCESS = 0;                                                        //最后运行的模拟器进程
var constSeparator = "__"                                                    //rom子分隔符
var DOC_EXTS = []string{".txt", ".md"}                                       //doc文档支持的扩展名
var PIC_EXTS = []string{".png", ".jpg", ".gif", ".ico", ".jpeg", ".bmp"}     //支持的图片类型
var RUN_EXTS = []string{
	".html", ".htm", ".mht", ".mhtml", ".url",
	".pdf", ".chm", ".doc", ".docx", ".ppt", ".pptx", "xls", "xlsx", ".rtf",
	".exe", ".com", ".cmd", ".bat", ".lnk",
	".ico", ".png", ".jpg", ".gif", ".jpeg", ".bmp", ".mp4", ".avi", ".wmv"} //可直接运行的doc文档支持的扩展名

type RomDetail struct {
	Info            *db.Rom   //基础信息
	DocContent      string    //简介内容
	StrategyContent string    //攻略内容
	StrategyFile    string    //攻略文件
	Sublist         []*db.Rom //子游戏
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

	//保存进程id
	LAST_PROCESS = result.Process.Pid

	return nil
}

/**
 * 关闭游戏
 **/
func killGame() error {

	if LAST_PROCESS == 0 {
		return nil
	}
	c := exec.Command("taskkill.exe", "/T", "/PID", utils.ToString(LAST_PROCESS))
	err := c.Start()

	LAST_PROCESS = 0
	return err
}


/*
 zip解压
*/

func UnzipRom(zipFile string, romExt []string) (string, error) {

	if strings.ToLower(filepath.Ext(zipFile)) != ".zip" {
		return zipFile, nil
	}

	zipReader, err := zip.OpenReader(zipFile)
	if err != nil {
		return "", err
	}
	defer zipReader.Close()

	runPath := "" //压缩包中的执行文件

	for _, f := range zipReader.File {

		//解决中文文件名乱码问题
		enc := mahonia.NewDecoder("gbk")
		f.Name = enc.ConvertString(f.Name)

		//拼接解压路径
		zipfileName := utils.GetFileName(zipFile)
		fpath := filepath.Join(Config.UnzipPath, zipfileName, f.Name)

		fileExt := filepath.Ext(f.Name)

		//找到压缩包中的可执行文件
		if runPath == ""{
			for _, v := range romExt {
				if v == fileExt {
					runPath = fpath
					break
				}
			}
		}
		//如果文件存在，则无需解压了
		if utils.FileExists(fpath) {
			continue
		}

		//开始解压
		if f.FileInfo().IsDir() {
			if err = utils.CreateDir(fpath); err != nil {
				return "", err
			}
		} else {
			if err = utils.CreateDir(filepath.Dir(fpath)); err != nil {
				return "", err
			}

			inFile, err := f.Open()
			if err != nil {
				return "", err
			}
			defer inFile.Close()

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return "", err
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return "", err
			}

		}
	}
	return runPath, nil
}

/*
 清理解压缓存
*/
func ClearZipRom() error {
	err := os.RemoveAll(Config.UnzipPath)
	if err != nil {
		return err
	}
	return nil
}
