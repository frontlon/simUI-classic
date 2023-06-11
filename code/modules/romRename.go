package modules

import (
	"errors"
	"simUI/code/config"
	"simUI/code/db"
	"simUI/code/utils"
	"strings"
)

//rom重命名
func RomRename(id uint64, newName string) error {

	//读取老信息
	rom, _ := (&db.Rom{}).GetById(id)
	if newName == rom.Name || newName == "" { //如果名称一样则不用修改
		return nil
	}

	err := errors.New("")

	//重命名文件
	if err = renameFile(newName, rom); err != nil {
		return err
	}

	//读取资料数据
	baseInfo := GetRomBaseById(rom.Platform, utils.GetFileName(rom.RomPath))
	baseName := newName
	if baseInfo != nil && baseInfo.Name != "" {
		baseName = baseInfo.Name
	}

	//写csv配置文件
	if baseInfo != nil {
		create := &RomBase{}
		create = baseInfo
		create.RomName = newName
		//写入配置文件
		if err := WriteRomBaseFile(rom.Platform, create); err != nil {
			return err
		}
	}

	//更新数据库
	p := rom.RomPath
	fpath := utils.GetFileAbsPath(rom.RomPath)
	fext := utils.GetFileExt(rom.RomPath)
	p = newName + fext
	if fpath != "." {
		p = fpath + "/" + newName + fext
	}

	infoMd5 := utils.GetRomMd5(baseName, p, rom.BaseType, rom.BaseYear, rom.BaseProducer, rom.BasePublisher, rom.BaseCountry, rom.BaseTranslate, rom.BaseVersion, rom.BaseNameEn, rom.BaseNameJp, rom.BaseOtherA, rom.BaseOtherB, rom.BaseOtherC, rom.BaseOtherD, rom.Score, rom.Size)

	err = (&db.Rom{
		Id:      id,
		Name:    baseName,
		RomPath: p,
		Pinyin:  utils.TextToPinyin(baseName),
		InfoMd5: infoMd5,
	}).UpdateName()
	if err != nil {
		return err
	}

	return nil
}

//rom批量重命名
func BatchRomRename(data []map[string]string) error {
	ids := []uint64{}
	create := map[string]map[string]string{}
	for _, v := range data {
		ids = append(ids, uint64(utils.ToInt(v["id"])))
		c := map[string]string{}
		c["id"] = v["id"]
		c["filename"] = v["filename"]
		create[c["filename"]] = c
	}
	//读取老信息
	volist, _ := (&db.Rom{}).GetByIds(ids)
	romlist := map[uint64]*db.Rom{}
	for _, v := range volist {
		romlist[v.Id] = v
		filename := utils.GetFileName(v.RomPath)
		//同名等于没改名
		if filename == create[filename]["filename"] {
			delete(create, filename)
		}
	}

	if len(create) == 0 {
		return nil
	}

	//开始遍历修改
	for _, v := range create {
		rom := romlist[uint64(utils.ToInt(v["id"]))]
		filename := v["filename"]

		err := errors.New("")

		if err = renameFile(filename, rom); err != nil {
			return err
		}

		//更新数据库
		fname := rom.RomPath
		fpath := utils.GetFileAbsPath(rom.RomPath)
		fext := utils.GetFileExt(rom.RomPath)
		fname = filename + fext
		if fpath != "." {
			fname = fpath + "/" + filename + fext
		}

		err = (&db.Rom{Id: uint64(utils.ToInt(v["id"])), Name: filename, RomPath: fname, Pinyin: utils.TextToPinyin(filename)}).UpdateName()
		if err != nil {
			return err
		}

	}

	return nil
}

//修改文件名
func renameFile(newName string, rom *db.Rom) error {
	platform := rom.Platform
	oldfileName := utils.GetFileName(rom.RomPath)

	resPaths := config.GetResPath(platform)
	resPaths["rom"] = config.Cfg.Platform[platform].RomPath

	//遍历资源目录
	for _, rpath := range resPaths {
		//读取相关资源文件
		files, _ := utils.ScanMasterSlaveFiles(rpath, oldfileName)
		for _, f := range files {
			fname := utils.GetFileName(f)
			newFilename := newName
			if strings.Contains(f, "__") {
				fileNameArr := strings.Split(fname, "__")
				newFilename = newFilename + "__" + fileNameArr[1]
			}
			//开始改名
			if err := utils.FileRename(f, newFilename); err != nil {
				return err
			}
		}

	}

	//改名音乐文件
	audioPath := resPaths["audio"] + config.Cfg.Separator + oldfileName
	if err := utils.FolderRename(audioPath, newName); err != nil {
		return err
	}

	return nil
}
