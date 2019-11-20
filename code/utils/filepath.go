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
 从路径中读取文件名+扩展名
*/
func GetFullFileName(p string) string {
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
func PathExists(path string) bool {
	ff, err := os.Stat(path)
	if err != nil {
		return false
	}
	if ff.IsDir() == false{
		return false
	}
	return true
}