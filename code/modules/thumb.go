package modules

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"simUI/code/config"
	"simUI/code/db"
	"simUI/code/utils"
	"time"
)

//下载展示图片
func DownloadRomThumbs(typeName string, id uint64, newPath string) (string, error) {

	rom := &db.Rom{
		Id: id,
	}

	//设定新的文件名
	vo, err := rom.GetById(id)
	if err != nil {
		return "", err
	}

	//下载文件
	response, err := http.Get(newPath)
	if err != nil {
		return "", err
	}

	//下载成功后，备份原文件
	res := config.GetResPath(vo.Platform)
	platformPath := res[typeName]

	if platformPath == "" {
		return "", errors.New(config.Cfg.Lang["NoSetThumbDir"])
	}

	//备份老图片
	if err := backupOldPic(platformPath, vo.RomPath); err != nil {
		return "", err
	}

	//生成新文件
	ext := utils.GetFileExt(newPath)
	if ext == ""{
		ext = ".jpg"
	}
	platformPathAbs, err := filepath.Abs(platformPath) //读取平台图片路径
	RomFileName := utils.GetFileName(vo.RomPath)
	newFileName := platformPathAbs + config.Cfg.Separator + RomFileName + ext //生成新文件的完整绝路路径地址

	f, err := os.Create(newFileName)
	defer f.Close()

	if err != nil {
		return "", err
	}
	if _, err := io.Copy(f, response.Body); err != nil {
		return "", err
	}

	return newFileName, nil
}

//编辑展示图片
func EditRomThumbs(typeName string, id uint64, picPath string) (string, error) {

	rom := &db.Rom{
		Id: id,
	}

	//设定新的文件名
	vo, err := rom.GetById(id)
	if err != nil {
		return "", err
	}

	//下载成功后，备份原文件
	res := config.GetResPath(vo.Platform)
	platformPath := res[typeName]

	if platformPath == "" {
		return "", errors.New(config.Cfg.Lang["NoSetThumbDir"])
	}

	//备份老图片
	if err := backupOldPic(platformPath, vo.RomPath); err != nil {
		return "", err
	}

	//生成新文件
	platformPathAbs, err := filepath.Abs(platformPath) //读取平台图片路径
	RomFileName := utils.GetFileName(vo.RomPath)
	newFileName := platformPathAbs + config.Cfg.Separator + RomFileName + utils.GetFileExt(picPath) //生成新文件的完整绝路路径地址
	//复制文件
	if err := utils.FileCopy(picPath, newFileName); err != nil {
		return "", errors.New(config.Cfg.Lang["ResFileNotFound"])
	}

	return newFileName, nil
}

//删除缩略图
func DeleteThumbs(typeName string, id uint64) error {

	rom := &db.Rom{
		Id: id,
	}

	//设定新的文件名
	vo, err := rom.GetById(id)
	if err != nil {
		return err
	}

	//下载成功后，备份原文件
	res := config.GetResPath(vo.Platform)
	platformPath := res[typeName]
	if platformPath == "" {
		return errors.New(config.Cfg.Lang["NoSetThumbDir"])
	}

	//备份老图片
	if err := backupOldPic(platformPath, vo.RomPath); err != nil {
		return err
	}

	return nil
}

//备份老图片
func backupOldPic(platformPath string, romPath string) error {
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
			err := os.Rename(oldFileName, bakFolder+bakFileName)                       //移动文件
			if err != nil {
				return err
			}
		}
	}
	return nil
}
