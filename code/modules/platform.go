package modules

import (
	"simUI/code/config"
	"simUI/code/db"
	"simUI/code/utils"
)

// 快速创建一个平台
func CreatePlatform(id uint32) error {

	//读取平台信息
	info, err := (&db.Platform{}).GetById(id)
	if err != nil {
		utils.WriteLog(err.Error())
		return err
	}

	//创建rom目录
	if !utils.DirExists(config.Cfg.Platform[id].RomPath) {
		utils.CreateDir(config.Cfg.Platform[id].RomPath)
	}

	//创建资源目录
	dirList := config.GetResPath(id)
	for _, v := range dirList {
		if !utils.DirExists(v) {
			utils.CreateDir(v)
		}
	}

	if info.Rombase != "" {
		info.Rombase = config.Cfg.RootPath + info.Rombase
		if !utils.FileExists(info.Rombase) {
			CreateNewRomBaseFile(info.Rombase)
		}
	}

	return nil
}

// 更新平台缩略图类型
func UpdatePlatformFieldById(platform uint32, field, val string) error {
	//更新到数据库
	return (&db.Platform{Id: platform}).UpdateFieldById(field, val)
}

// 清空平台缩略图类型
func ClearAllPlatformAField(typ string) error {
	return (&db.Platform{}).ClearAllPlatformAField(typ)
}
