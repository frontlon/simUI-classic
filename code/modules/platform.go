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

	//创建rom目录
	if !utils.FolderExists(config.Cfg.Platform[id].RomPath) {
		utils.CreateDir(config.Cfg.Platform[id].RomPath)
	}

	//创建资源目录
	dirList := config.GetResPath(id)
	for _, v := range dirList {
		if !utils.FolderExists(v) {
			utils.CreateDir(v)
		}
	}

	if info.Rombase != "" {
		info.Rombase = config.Cfg.RootPath + info.Rombase
		if !utils.FileExists(info.Rombase) {
			CreateNewRomBaseFile(info.Rombase);
		}
	}

	return nil
}
