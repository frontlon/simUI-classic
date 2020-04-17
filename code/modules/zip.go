package modules

import (
	"VirtualNesGUI/code/config"
	"VirtualNesGUI/code/utils"
	"archive/zip"
	"github.com/axgle/mahonia"
	"io"
	"os"
	"path/filepath"
	"strings"
)

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
		fpath := filepath.Join(config.C.UnzipPath, zipfileName, f.Name)

		fileExt := filepath.Ext(f.Name)

		//找到压缩包中的可执行文件
		if runPath == "" {
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
	err := os.RemoveAll(config.C.UnzipPath)
	if err != nil {
		return err
	}
	return nil
}

