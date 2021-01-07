package modules

import (
	"simUI/code/config"
	"simUI/code/db"
	"simUI/code/utils"
)

//读取详情文件
func CreatePlatform(id uint32) error {

	//读取平台信息
	info, err := (&db.Platform{}).GetById(id)
	if err != nil {
		utils.WriteLog(err.Error())
		return err
	}

	dirList := []string{
		info.RomPath,
		info.ThumbPath,
		info.SnapPath,
		info.PosterPath,
		info.PackingPath,
		info.TitlePath,
		info.BackgroundPath,
		info.DocPath,
		info.StrategyPath,
		info.VideoPath,
	}

	for _, v := range dirList {
		path := config.Cfg.RootPath + v
		if (!utils.FolderExists(path)) {
			utils.CreateDir(path);
		}
	}

	info.Rombase = config.Cfg.RootPath + info.Rombase
	if (!utils.FileExists(info.Rombase)) {
		CreateNewRomBaseFile(info.Rombase);
	}

	return nil
}
