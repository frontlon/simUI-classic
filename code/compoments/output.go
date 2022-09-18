package compoments

import (
	"archive/zip"
	"fmt"
	"github.com/go-ini/ini"
	"io"
	"os"
	"path/filepath"
	"simUI/code/config"
	"simUI/code/utils"
	"strings"
)

var zipMethod = 1 //0仅存储;1压缩
var TmpOutputIniPath = "./cache/config.ini" //ini文件路径
/**
 * 将文件添加到压缩文件中
 **/
func CompressZip(file *os.File, prefix string, zw *zip.Writer) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		if prefix != "" {
			prefix = prefix + config.Cfg.Separator + info.Name()
		} else {
			prefix = info.Name()
		}
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			f, err := os.Open(file.Name() + config.Cfg.Separator + fi.Name())
			defer f.Close()
			if err != nil {
				return err
			}
			err = CompressZip(f, prefix, zw)
			if err != nil {
				return err
			}
		}
	} else {

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		if prefix != "" {
			header.Name = prefix + config.Cfg.Separator + header.Name
		}
		if zipMethod == 1 {
			header.Method = zip.Deflate //压缩
		} else {
			header.Method = zip.Store //仅存储
		}

		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}

		_, err = io.Copy(writer, file)
		if err != nil {
			return err
		}
	}
	return nil
}

/**
 * 读取所有rom的展示图
 **/
func GetAllRomThumbsToMap(resPath string, names map[string]bool) (map[string][]string, error) {

	result := map[string][]string{}

	if err := filepath.Walk(resPath,
		func(p string, f os.FileInfo, err error) error {

			name := utils.GetFileName(p)
			if strings.Contains(name, "__") {
				nameArr := strings.Split(name, "__")
				name = nameArr[0]
			}

			if _, ok := result[name]; ok {
				result[name] = append(result[name], p)
			} else {
				result[name] = []string{p}
			}

			return nil
		}); err != nil {
		return nil, err
	}

	return result, nil
}

/**
 * 生成导出的ini文件
 **/
func CreateOutputIni(subGames map[string]string) error {
	cfg := ini.Empty()

	subgameSection, err := cfg.NewSection("subgame")
	if err != nil {
		fmt.Println("new mysql section failed:", err)
		return err
	}
	for k, v := range subGames {
		subgameSection.NewKey(k, v)
	}

	err = cfg.SaveTo(TmpOutputIniPath)
	if err != nil {
		fmt.Println("SaveTo failed: ", err)
		return err
	}
	return nil
}
