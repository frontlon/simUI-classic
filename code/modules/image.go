package modules

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"simUI/code/config"
	"simUI/code/utils"
)

var imageWidth = 320
var quality = 80

func CreateOptimizedImage(platform uint32, opt string) error {

	paths := config.GetResPath(platform)
	path := paths[opt]

	outputPath := config.Cfg.Platform[platform].OptimizedPath

	//先删除原文件
	utils.DeleteDir(outputPath)
	utils.CreateDir(outputPath)

	//读取文件总数
	files, _ := ioutil.ReadDir(path)
	fileCount := len(files)
	i := 0
	filepath.Walk(path, func(p string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() { /** 是否是目录 */
			return nil
		}

		/** 判断是不是图片 */
		format := utils.GetFileExt(p)

		outputPath := outputPath + config.Cfg.Separator + utils.GetFileNameAndExt(p)
		if p != "" {
			err := utils.ImageCompress(
				func() (io.Reader, error) {
					return os.Open(p)
				},
				func() (*os.File, error) {
					return os.Open(p)
				},
				outputPath,
				quality,
				imageWidth,
				format)

			if err != nil {
				return err
			}
		}

		if i%10 == 0 {
			utils.Loading("[2/3]已生成("+utils.ToString(i)+" / "+utils.ToString(fileCount)+")", config.Cfg.Platform[platform].Name)
		}

		i++
		return nil
	})

	//数据更新完成后，页面回调，更新页面DOM
	if _, err := utils.Window.Call("CB_createOptimizedCache"); err != nil {
		fmt.Print(err)
	}

	return nil
}
