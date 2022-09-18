package modules

import (
	"simUI/code/config"
	"simUI/code/db"
	"simUI/code/utils"
	"strings"
)

//快速创建一个平台
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
			CreateNewRomBaseFile(info.Rombase)
		}
	}

	return nil
}

//更新平台介绍
func UpdatePlatformDesc(platform uint32, desc string) error {

	desc = strings.Trim(desc, " ")

	//替换图片路径为相对路径
	desc = strings.ReplaceAll(desc, config.Cfg.RootPath, "")

	//更新到数据库
	err := (&db.Platform{Id: platform, Desc: desc}).UpdateDescById()
	if err != nil {
		return err
	}

	return nil
}

//更新平台缩略图类型
func UpdatePlatformThumb(platform uint32, thumb string) error {

	//更新到数据库
	err := (&db.Platform{Id: platform, Thumb: thumb}).UpdateThumbById()
	if err != nil {
		return err
	}

	return nil
}

//清空平台缩略图类型
func ClearPlatformThumb() error {

	//更新到数据库
	err := (&db.Platform{}).ClearAllThumb()
	if err != nil {
		return err
	}

	return nil
}

//更新平台缩略图类型
func UpdatePlatformThumbDirection(platform uint32, dir string) error {

	//更新到数据库
	err := (&db.Platform{Id: platform, ThumbDirection: dir}).UpdateThumbDirectionById()
	if err != nil {
		return err
	}

	return nil
}

//清空平台缩略图类型
func ClearPlatformThumbDirection() error {

	//更新到数据库
	err := (&db.Platform{}).ClearAllThumbDirection()
	if err != nil {
		return err
	}

	return nil
}
