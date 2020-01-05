package main

import (
	"VirtualNesGUI/code/utils"
	"archive/zip"
	"github.com/Lofanmi/pinyin-golang/pinyin"
	"github.com/axgle/mahonia"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
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
		for _, v := range romExt {
			if v == fileExt {
				runPath = fpath
				break
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

/*func UnzipRom(zipFile string, romExt []string) (string, error) {

	if strings.ToLower(filepath.Ext(zipFile)) != ".zip"{
		return zipFile,nil
	}

	zipReader, err := zip.OpenReader(zipFile)

	if err != nil {
		return "", err
	}
	defer zipReader.Close()
	romFile := &zip.File{}
	for _, f := range zipReader.File {
		fileExt := filepath.Ext(f.Name)

		//找到rom文件
		for _, v := range romExt {
			if v == fileExt {
				romFile = f
				break
			}
		}

		if romFile.Name != "" {
			break
		}
	}

	if romFile.Name == "" {
		return "", nil
	}

	//解决中文文件名乱码问题
	enc := mahonia.NewDecoder("gbk")
	romFile.Name = enc.ConvertString(romFile.Name)

	//拼接解压路径
	zipfileName := utils.GetFileName(zipFile)
	fpath := filepath.Join(Config.UnzipPath, zipfileName,romFile.Name)

	//如果文件存在，则无需解压了
	if utils.FileExists(fpath) {
		return fpath, nil
	}

	//开始解压
	if romFile.FileInfo().IsDir() {
		if err = utils.CreateDir(fpath);err != nil{
			return fpath,err
		}
	} else {
		if err = utils.CreateDir(filepath.Dir(fpath)); err != nil {
			return fpath, err
		}

		inFile, err := romFile.Open()
		if err != nil {
			return fpath, err
		}
		defer inFile.Close()

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, romFile.Mode())
		if err != nil {
			return fpath, err
		}
		defer outFile.Close()

		_, err = io.Copy(outFile, inFile)
		if err != nil {
			return fpath, err
		}
	}

	return fpath, nil
}*/

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

/*
 写日志
*/
func WriteLog(str string) {
	fileName := "log.txt"

	if !utils.IsExist(Config.CachePath) {
		if err := utils.CreateDir(Config.CachePath); err != nil {
			return
		}
	}

	f, _ := os.OpenFile(Config.CachePath+fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

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
