package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

//复制文件
func FileCopy(src, dst string) (error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}

	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}

//移动文件
func FileMove(oldFile string,newFile string) (error) {
	if FileExists(oldFile) == true {
		err := os.Rename(oldFile,newFile)
		return err
	}
	return nil
}

/*
 重命名
*/
func Rename(oldpath string, filename string) error {
	newpath := filepath.Dir(oldpath) + "/" + filename + path.Ext(oldpath)
	return os.Rename(oldpath, newpath)
}

//删除文件
func FileDelete(src string) (error) {
	if FileExists(src) == true {
		err := os.Remove(src)
		return err
	}
	return nil
}

/**
 * 返回文件大小 + 单位
 */
func GetFileSizeString(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%.2fB", float64(size)/float64(1))
	} else if size < (1024 * 1024) {
		return fmt.Sprintf("%.2fKB", float64(size)/float64(1024))
	} else if size < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fMB", float64(size)/float64(1024*1024))
	} else if size < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fGB", float64(size)/float64(1024*1024*1024))
	} else {
		return fmt.Sprintf("%.2fEB", float64(size)/float64(1024*1024*1024*1024*1024))
	}
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

/**
 * 创建一个文件
 **/
func CreateFile(p string) error {
	if p == "" {
		return nil
	}
	f, err := os.Create(p)
	defer f.Close()
	if err != nil {
		return err
	}
	return nil
}

/**
 * 写文件（覆盖写）
 **/
func OverlayWriteFile(p string, t string) error {
	if p == "" {
		return nil
	}
	if err := ioutil.WriteFile(p, []byte(t), 0664); err != nil {
		return err
	}
	return nil
}
