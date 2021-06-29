package modules

import (
	"encoding/json"
	"errors"
	"simUI/code/config"
	"simUI/code/db"
	"simUI/code/utils"
	"strings"
)

func GetStrategyFile(id uint64) ([]map[string]string, error) {
	vo, err := (&db.Rom{}).GetById(id)
	if err != nil {
		return nil, err
	}

	//搜索攻略文件
	exists, _ := utils.ScanDirByKeyword(config.Cfg.Platform[vo.Platform].FilesPath, vo.Name+"__")
	volist := []map[string]string{}
	for _, v := range exists {
		name := utils.GetFileName(v)
		namearr := strings.Split(name, "__")
		p := strings.Replace(v, config.Cfg.RootPath, "", 1)
		vo := make(map[string]string)
		vo["name"] = namearr[1]
		vo["path"] = p
		volist = append(volist, vo)
	}
	return volist, nil
}

/**
 * 上传文件
 **/
func UploadStrategyFile(id uint64, name string, p string) (string, error) {
	vo, _ := (&db.Rom{}).GetById(id)
	if config.Cfg.Platform[vo.Platform].FilesPath == "" {
		return "", errors.New(config.Cfg.Lang["FilesMenuCanNotBeEmpty"])
	}
	ext := utils.GetFileExt(p)
	fileName := utils.GetFileName(vo.RomPath)
	newPath := config.Cfg.Platform[vo.Platform].FilesPath + config.Cfg.Separator + fileName + "__" + name + ext
	if err := utils.FileCopy(p, newPath); err != nil {
		return "", err
	}
	relPath := strings.Replace(newPath, config.Cfg.RootPath, "", -1)
	return relPath, nil
}

/**
 * 更新数据
 **/
func UpdateStrategyFiles(id uint64, data string) error {
	vo, _ := (&db.Rom{}).GetById(id)
	if config.Cfg.Platform[vo.Platform].FilesPath == "" {
		return errors.New(config.Cfg.Lang["FilesMenuCanNotBeEmpty"])
	}

	//整理需要删除的文件
	d := []map[string]string{}
	json.Unmarshal([]byte(data), &d)
	newData := []string{}
	for _, v := range d {
		newData = append(newData, v["path"])
	}

	//读取已存在的文件
	exists, _ := utils.ScanDirByKeyword(config.Cfg.Platform[vo.Platform].FilesPath, vo.Name+"__")
	for _, v := range exists {
		rel := strings.Replace(v, config.Cfg.RootPath, "", 1)
		if !utils.InSliceString(rel, newData) {
			utils.FileDelete(v)
		}
	}

	return nil
}

//打开攻略文件
func OpenStrategyFiles(p string) error {

	if !utils.IsAbsPath(p) {
		p = config.Cfg.RootPath + p
	}

	if !utils.FileExists(p) {
		return errors.New(p + config.Cfg.Lang["FileNotFound"])
	}

	cmd:= []string{}
	cmd = append(cmd,p)

	if err := utils.RunGame("", cmd); err != nil {
		return err
	}

	return nil
}
