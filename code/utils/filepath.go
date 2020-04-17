package utils

import (
	"os"
	"path"
	"path/filepath"
	"strings"
)

/*
 从路径中读取文件名，不包含扩展名
*/
func GetFileName(p string) string {
	base := filepath.Base(p);
	str := strings.TrimSuffix(base, path.Ext(base))
	return str
}

/*
 从路径中读取文件名扩展名
*/
func GetFileExt(p string) string {
	return path.Ext(p)
}

/*
 重命名
*/
func Rename(oldpath string, filename string) error {
	newpath := filepath.Dir(oldpath) + "/" + filename + path.Ext(oldpath)
	return os.Rename(oldpath, newpath)
}

/*
 从完整路径中获取文件路径，不包含结尾  /
*/
func GetFilePath(p string) string {
	return filepath.Dir(p)
}

/*
 从路径中读取文件名+扩展名
*/
func GetFileNameAndExt(p string) string {
	return filepath.Base(p);
}

/**
 * 检测文件是否存在（文件夹也返回false）
 **/
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
	if p == ""{
		return false
	}
	ff, err := os.Stat(p)
	if err != nil {
		return false
	}
	if ff.IsDir() == false{
		return false
	}
	return true
}

/**
 * 判断文件夹或文件是否存在
 **/
func IsExist(f string) bool {
	if f == ""{
		return false
	}
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}


/**
 * 创建多层文件夹
 **/
func CreateDir(p string) error {
	err := os.MkdirAll(p, os.ModePerm)
	if err != nil {
		return err
	}
	if err := os.Chmod(p, os.ModePerm);err != nil{
		return err
	}
	return nil
}



