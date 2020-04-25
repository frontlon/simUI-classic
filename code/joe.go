package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func main() {

	rawFolder := "./raw"
	listFile := "./list.txt"

	if !FileExists(listFile) {
		fmt.Println("list.txt 文件不存在")
		return
	}

	//创建文件夹
	folder := time.Now().Format("2006-01-02 15.04.05")
	if !FolderExists(folder) {
		if err := CreateDir(folder); err != nil {
			fmt.Println("创建文件夹失败:", err)
			return
		}
	}

	fi, err := os.Open(listFile)
	if err != nil {
		fmt.Printf("打开文件失败: %s\n", err)
		return
	}

	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}

		base := filepath.Base(string(a));
		fname := strings.TrimSuffix(base, path.Ext(base))
		oldFile := rawFolder + "/" + fname + ".go"

		newFile := folder + "/" + fname + ".go"

		if err = os.Rename(oldFile, newFile); err != nil {
			fmt.Println("移动文件失败:", oldFile)
		}else{
			fmt.Println("移动文件成功:", oldFile)

		}
	}
}

func FileExists(path string) bool {

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
 * 检测路径是否存在
 **/
func FolderExists(p string) bool {
	if p == "" {
		return false
	}
	ff, err := os.Stat(p)
	if err != nil {
		return false
	}
	if ff.IsDir() == false {
		return false
	}
	return true
}

/**
 * 创建多层文件夹
 **/
func CreateDir(p string) error {
	err := os.MkdirAll(p, os.ModePerm)
	if err != nil {
		return err
	}
	if err := os.Chmod(p, os.ModePerm); err != nil {
		return err
	}
	return nil
}
