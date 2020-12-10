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
	dir := filepath.Dir(p)
	return strings.Replace(dir, "./", "", 1)
}

/*
 从路径中读取文件名+扩展名
*/
func GetFileNameAndExt(p string) string {
	return filepath.Base(p);
}

/*
 路径转换为绝对路径
*/
func AbsPath(p string) string {
	if(!filepath.IsAbs(p)){
		p,_ = filepath.Abs(p)
	}

	return p
}

/*
 检查路径是否为绝对路径
*/
func IsAbsPath(p string) bool {
	return filepath.IsAbs(p)
}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
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
 * 判断文件夹或文件是否存在
 **/
func IsExist(f string) bool {
	if f == "" {
		return false
	}
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}


/**
 * 根据关键字扫描目录和子目录，查询出符合条件的文件名
 * 如果keyword不为空，则查询keyword关键字相关的程序
 **/
func ScanDirByKeyword(dir string, keyword string) ([]string, error) {

	files := []string{}
	keyword = strings.ToUpper(keyword) //忽略后缀匹配的大小写

	err := filepath.Walk(dir, func(filename string, fi os.FileInfo, err error) error { //遍历目录
		if err != nil { //忽略错误
			return err
		}

		if fi.IsDir() { // 忽略目录
			return nil
		}

		if strings.Contains(strings.ToUpper(fi.Name()), keyword) { //如果包含关键字
			files = append(files, filename)
		}

		return nil
	})

	return files, err

}
