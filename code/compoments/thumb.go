package compoments

import (
	"simUI/code/config"
	"simUI/code/utils"
	"time"
)

//备份老图片
func BackupOldPic(platformPath string, romPath string) error {
	//开始备份原图
	bakFolder := config.Cfg.CachePath + "thumb_bak/"
	RomFileName := utils.GetFileName(romPath)

	//检测bak文件夹是否存在，不存在则创建bak目录
	folder := utils.FolderExists(bakFolder)
	if folder == false {

		if err := utils.CreateDir(bakFolder); err != nil {
			return err
		}
	}
	for _, ext := range config.PIC_EXTS {
		oldFileName := platformPath + config.Cfg.Separator + RomFileName + ext //老图片文件名
		if utils.FileExists(oldFileName) {
			bakFileName := RomFileName + "_" + utils.ToString(time.Now().Unix()) + ext //生成备份文件名
			if err := utils.FileMove(oldFileName, bakFolder+bakFileName); err != nil {
				return err
			}
		}
	}
	return nil
}

//删除图片资源
func DeleteResPic(platformPath string, romPath string) error {
	//开始备份原图
	RomFileName := utils.GetFileName(romPath)
	for _, ext := range config.PIC_EXTS {
		oldFileName := platformPath + config.Cfg.Separator + RomFileName + ext //图片文件名
		if utils.FileExists(oldFileName) {
			if err := utils.FileDelete(oldFileName); err != nil {
				return err
			}
		}
	}
	return nil
}
