package compoments

import (
	"io"
	"net/http"
	"os"
	"simUI/code/config"
	"simUI/code/utils"
	"time"
)

//备份老图片
func BackupOldPic(p string, fileName string) error {
	//开始备份原图
	bakFolder := config.Cfg.CachePath + "thumb_bak/"

	//检测bak文件夹是否存在，不存在则创建bak目录
	folder := utils.FolderExists(bakFolder)
	if folder == false {

		if err := utils.CreateDir(bakFolder); err != nil {
			return err
		}
	}
	for _, ext := range config.PIC_EXTS {
		oldFileName := p + config.Cfg.Separator + fileName + ext //老图片文件名
		if utils.FileExists(oldFileName) {
			bakFileName := fileName + "_" + utils.ToString(time.Now().Unix()) + ext //生成备份文件名
			if err := utils.FileMove(oldFileName, bakFolder+bakFileName); err != nil {
				return err
			}
		}
	}
	return nil
}

//删除图片资源
func DeleteResPic(platformPath string, fileName string) error {
	//开始备份原图
	for _, ext := range config.PIC_EXTS {
		oldFileName := platformPath + config.Cfg.Separator + fileName + ext //图片文件名
		if utils.FileExists(oldFileName) {
			if err := utils.FileDelete(oldFileName); err != nil {
				return err
			}
		}
	}
	return nil
}

//下载展示图片
func DownloadRomThumbs(httpUrl string, localPath string) error {

	//下载文件
	response, err := http.Get(httpUrl)
	if err != nil {
		return err
	}

	f, err := os.Create(localPath)
	defer f.Close()

	if err != nil {
		return err
	}

	if _, err := io.Copy(f, response.Body); err != nil {
		return err
	}

	return nil
}
