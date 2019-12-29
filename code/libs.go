package main

import (
	"VirtualNesGUI/code/utils"
	"archive/zip"
	"github.com/Lofanmi/pinyin-golang/pinyin"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var dict = pinyin.NewDict()
var unzipDir = "%TEMP%"

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

	if filepath.Ext(zipFile) != "zip"{
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

	fpath := filepath.Join(unzipDir, romFile.Name)

	//如果文件存在，则无需解压了
	if utils.FileExists(fpath) {
		return fpath, nil
	}

	//开始解压
	if romFile.FileInfo().IsDir() {
		os.MkdirAll(fpath, os.ModePerm)
	} else {
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
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
}

/*
 zip解压
*/
func ClearZipRom() error {
	err := os.RemoveAll(unzipDir)
	if err != nil {
		return err
	}
	return nil
}
