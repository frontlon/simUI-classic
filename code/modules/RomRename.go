package modules

import (
	"VirtualNesGUI/code/config"
	"VirtualNesGUI/code/db"
	"VirtualNesGUI/code/utils"
	"fmt"
	"strings"
)

//rom重命名
func RomRename(setType uint8,id uint64,name string) error {

	rom, _ := (&db.Rom{}).GetById(id)

	if name == rom.Name {
		return nil
	}

	subRom, _ := (&db.Rom{}).GetSubRom(rom.Platform, rom.Name)
	platform := rom.Platform

	//alias
	iniCfg := &ini.File{}
	err := errors.New("")
	p := config.C.Platform[rom.Platform].Romlist
	//修改别名文件
	if setType == 1 {
		if p == "" || !utils.IsExist(p) {
			p = config.C.RootPath + config.C.Platform[platform].Name + ".ini"
			iniCfg = ini.Empty()
			config.C.Platform[platform].Romlist = p
		} else {
			iniCfg, err = ini.LoadSources(ini.LoadOptions{IgnoreInlineComment: true}, p)
			if err != nil {
				fmt.Println(err)
			}
		}

		//修改主rom
		k := utils.GetFileName(rom.RomPath)
		iniCfg.Section("Alias").NewKey(k, name)
		//修改子rom
		for _, v := range subRom {
			fileName := utils.GetFileName(v.RomPath)
			subName, _ := iniCfg.Section("Alias").GetKey(k) //修改主rom
			ns := strings.Replace(subName.String(), rom.Name+"__", name+"__", 1)
			iniCfg.Section("Alias").NewKey(fileName, ns)
		}

		if err := iniCfg.SaveTo(p); err != nil {
		}
	} else { //修改文件名
		//主rom
		if err := utils.Rename(rom.RomPath, name); err != nil {
			return err
		}

		//子rom
		for _, v := range subRom {
			fileName := utils.GetFileName(v.RomPath)
			newName := strings.Replace(fileName, rom.Name+"__", name+"__", 1)
			if err := utils.Rename(v.RomPath, newName); err != nil {
				return err
			}
		}

		//修改资源文件
		oldfileName := utils.GetFileName(rom.RomPath)

		if config.C.Platform[platform].PackingPath != "" {
			for _, ext := range config.PIC_EXTS {
				picpath := config.C.Platform[platform].PackingPath + "/" + oldfileName + ext
				if (utils.FileExists(picpath)) {
					if err := utils.Rename(picpath, name); err != nil {
						return err
					}
					break
				}
			}
		}

		if config.C.Platform[platform].SnapPath != "" {
			for _, ext := range config.PIC_EXTS {
				picpath := config.C.Platform[platform].SnapPath + "/" + oldfileName + ext
				if (utils.FileExists(picpath)) {
					if err := utils.Rename(picpath, name); err != nil {
						return err
					}
					break
				}
			}
		}

		if config.C.Platform[platform].ThumbPath != "" {
			for _, ext := range config.PIC_EXTS {
				picpath := config.C.Platform[platform].ThumbPath + "/" + oldfileName + ext
				if (utils.FileExists(picpath)) {
					if err := utils.Rename(picpath, name); err != nil {
						return err
					}
					break
				}
			}
		}

		if config.C.Platform[platform].PosterPath != "" {
			for _, ext := range config.PIC_EXTS {
				picpath := config.C.Platform[platform].PosterPath + "/" + oldfileName + ext
				if (utils.FileExists(picpath)) {
					if err := utils.Rename(picpath, name); err != nil {
						return err
					}
					break
				}
			}
		}

		if config.C.Platform[platform].DocPath != "" {
			for _, ext := range config.DOC_EXTS {
				picpath := config.C.Platform[platform].DocPath + "/" + oldfileName + ext
				if (utils.FileExists(picpath)) {
					if err := utils.Rename(picpath, name); err != nil {
						return err
					}
					break
				}
			}
		}
		if config.C.Platform[platform].StrategyPath != "" {
			for _, ext := range config.RUN_EXTS {
				picpath := config.C.Platform[platform].StrategyPath + "/" + oldfileName + ext
				if (utils.FileExists(picpath)) {
					if err := utils.Rename(picpath, name); err != nil {
						return err
					}
					break
				}
			}
		}

	}

	//更新数据库
	err = (&db.Rom{Id: id, Name: name, Pinyin: utils.TextToPinyin(name)}).UpdateName()
	if err != nil {
		return err
	}

	//更新配置
	err = (&db.Config{}).UpdateField("rename_type", setType)
	config.C.Default.RenameType = setType
	if err != nil {
		return err
	}
	return nil
}
