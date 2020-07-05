package modules

import (
	"errors"
	"simUI/code/config"
	"simUI/code/db"
	"simUI/code/utils"
	"strings"
)

//rom重命名
func RomRename(id uint64, name string) error {

	//读取老信息
	rom, _ := (&db.Rom{}).GetById(id)
	if name == rom.Name || name == "" { //如果名称一样则不用修改
		return nil
	}

	//读取子游戏
	subRom, _ := (&db.Rom{}).GetSubRom(rom.Platform, rom.Name)

	err := errors.New("")

	if err = renameFile(name, rom, subRom); err != nil {
		return err
	}
	
	//更新数据库
	fname := rom.RomPath
	fpath := utils.GetFilePath(rom.RomPath)
	fext := utils.GetFileExt(rom.RomPath)
	fname = name + fext
	if fpath != "." {
		fname = fpath + "/" + name + fext
	}

	err = (&db.Rom{Id: id, Name: name, RomPath: fname, Pinyin: utils.TextToPinyin(name)}).UpdateName()
	if err != nil {
		return err
	}

	return nil
}

//修改别名文件
/*func renameAlias(name string, rom *db.Rom, subRom []*db.Rom) error {
	platform := rom.Platform
	iniCfg := &ini.File{}
	err := errors.New("")
	p := config.Cfg.Platform[rom.Platform].Romlist
	if p == "" || !utils.IsExist(p) {
		p = config.Cfg.RootPath + config.Cfg.Platform[platform].Name + ".ini"
		iniCfg = ini.Empty()
		config.Cfg.Platform[platform].Romlist = p
		//更新数据库字段
		if err = (&db.Platform{Id: platform}).UpdateFieldById("romlist", p); err != nil {
			return err
		}

	} else {
		iniCfg, err = ini.LoadSources(ini.LoadOptions{IgnoreInlineComment: true}, p)
		if err != nil {
			return err
		}
	}

	//修改主rom
	k := utils.GetFileName(rom.RomPath)
	iniCfg.Section("Alias").NewKey(k, name)
	//修改子rom
	for _, v := range subRom {
		fileName := utils.GetFileName(v.RomPath)
		subName, _ := iniCfg.Section("Alias").GetKey(k)
		ns := ""
		if subName.String() == "" {
			ns = name+"__" + fileName
		} else {
			ns = name+"__" + v.Name
		}
		iniCfg.Section("Alias").NewKey(fileName, ns)
	}

	if err := iniCfg.SaveTo(p); err != nil {
		return err
	}
	return nil
}*/

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

	if config.Cfg.Platform[platform].PackingPath != "" {
		for _, ext := range config.PIC_EXTS {
			picpath := config.Cfg.Platform[platform].PackingPath + "/" + oldfileName + ext
			if (utils.FileExists(picpath)) {
				if err := utils.Rename(picpath, name); err != nil {
					return err
				}
				break
			}
		}
	}

	if config.Cfg.Platform[platform].SnapPath != "" {
		for _, ext := range config.PIC_EXTS {
			picpath := config.Cfg.Platform[platform].SnapPath + "/" + oldfileName + ext
			if (utils.FileExists(picpath)) {
				if err := utils.Rename(picpath, name); err != nil {
					return err
				}
				break
			}
		}
	}

	if config.Cfg.Platform[platform].ThumbPath != "" {
		for _, ext := range config.PIC_EXTS {
			picpath := config.Cfg.Platform[platform].ThumbPath + "/" + oldfileName + ext
			if (utils.FileExists(picpath)) {
				if err := utils.Rename(picpath, name); err != nil {
					return err
				}
				break
			}
		}
	}

	if config.Cfg.Platform[platform].PosterPath != "" {
		for _, ext := range config.PIC_EXTS {
			picpath := config.Cfg.Platform[platform].PosterPath + "/" + oldfileName + ext
			if (utils.FileExists(picpath)) {
				if err := utils.Rename(picpath, name); err != nil {
					return err
				}
				break
			}
		}
	}

	if config.Cfg.Platform[platform].DocPath != "" {
		for _, ext := range config.DOC_EXTS {
			picpath := config.Cfg.Platform[platform].DocPath + "/" + oldfileName + ext
			if (utils.FileExists(picpath)) {
				if err := utils.Rename(picpath, name); err != nil {
					return err
				}
				break
			}
		}
	}
	if config.Cfg.Platform[platform].StrategyPath != "" {
		for _, ext := range config.RUN_EXTS {
			picpath := config.Cfg.Platform[platform].StrategyPath + "/" + oldfileName + ext
			if (utils.FileExists(picpath)) {
				if err := utils.Rename(picpath, name); err != nil {
					return err
				}
				break
			}
		}
	}
	return nil
}
