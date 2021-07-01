package modules

import (
	"errors"
	"simUI/code/config"
	"simUI/code/db"
	"simUI/code/utils"
	"strings"
)

//rom重命名
func RomRename(id uint64, filename string) error {

	//读取老信息
	rom, _ := (&db.Rom{}).GetById(id)
	if filename == rom.Name || filename == "" { //如果名称一样则不用修改
		return nil
	}

	//读取子游戏
	subRom, _ := (&db.Rom{}).GetSubRom(rom.Platform, rom.Name)

	err := errors.New("")

	if err = renameFile(filename, rom, subRom); err != nil {
		return err
	}

	//更新数据库
	fname := rom.RomPath
	fpath := utils.GetFilePath(rom.RomPath)
	fext := utils.GetFileExt(rom.RomPath)
	fname = filename + fext
	if fpath != "." {
		fname = fpath + "/" + filename + fext
	}

	err = (&db.Rom{Id: id, Name: filename, RomPath: fname, Pinyin: utils.TextToPinyin(filename)}).UpdateName()
	if err != nil {
		return err
	}

	return nil
}

//修改文件名
func renameFile(name string, rom *db.Rom, subRom []*db.Rom) error {
	platform := rom.Platform

	//主rom
	rompath := config.Cfg.Platform[platform].RomPath + "/" + rom.RomPath

	if utils.IsAbsPath(rom.RomPath) {
		rompath = rom.RomPath
	}
	if err := utils.Rename(rompath, name); err != nil {
		return err
	}

	//子rom
	for _, v := range subRom {
		fileName := utils.GetFileName(v.RomPath)
		newName := strings.Replace(fileName, rom.Name+"__", name+"__", 1)
		rompath := config.Cfg.Platform[platform].RomPath + "/" + v.RomPath
		if utils.IsAbsPath(v.RomPath) {
			rompath = v.RomPath
		}
		if err := utils.Rename(rompath, newName); err != nil {
			return err
		}
	}

	//修改资源文件
	oldfileName := utils.GetFileName(rom.RomPath)
	resExts := config.GetResExts()
	for resName, path := range config.GetResPath(platform) {
		if path != "" {
			for _, ext := range resExts[resName] {
				picpath := path + "/" + oldfileName + ext
				if utils.FileExists(picpath) {
					if err := utils.Rename(picpath, name); err != nil {
						return err
					}
					break
				}
			}
		}
	}
	//修改攻略文件
	masterName := utils.GetFileName(rom.RomPath)
	files, _ := utils.ScanDirByKeyword(config.Cfg.Platform[rom.Platform].FilesPath, masterName+"__")
	for _, f := range files {
		fArr := strings.Split(f, "__")
		fName := fArr[len(fArr)-1]
		fArr = strings.Split(fName, ".")
		fArr = utils.SliceDeleteLast(fArr)
		fName = strings.Join(fArr, ".")
		newName := name + "__" + fName
		if err := utils.Rename(f, newName); err != nil {
			return err
		}
	}
	return nil
}
