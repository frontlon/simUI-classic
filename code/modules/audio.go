package modules

import (
	"encoding/json"
	"errors"
	"simUI/code/config"
	"simUI/code/db"
	"simUI/code/utils"
	"strings"
)

func GetAudioList(id uint64) ([]map[string]string, error) {
	vo, err := (&db.Rom{}).GetById(id)
	name := utils.GetFileName(vo.RomPath)
	if err != nil {
		return nil, err
	}

	//搜索音频文件
	p := config.Cfg.Platform[vo.Platform].AudioPath + config.Cfg.Separator + name
	exists, _ := utils.ScanDir(p)
	volist := []map[string]string{}
	for _, v := range exists {
		ext := utils.GetFileExt(v)
		if !utils.InSliceString(ext,config.AUDIO_EXTS){
			continue
		}

		vo := make(map[string]string)
		vo["name"] = utils.GetFileName(v)
		vo["path"] = strings.Replace(v, config.Cfg.RootPath, "", 1)
		volist = append(volist, vo)
	}
	return volist, nil
}

/**
 * 上传文件
 **/
func UploadAudioFile(id uint64, name string, p string) (string, error) {
	vo, _ := (&db.Rom{}).GetById(id)
	if config.Cfg.Platform[vo.Platform].AudioPath == "" {
		return "", errors.New(config.Cfg.Lang["AudioMenuCanNotBeEmpty"])
	}
	ext := utils.GetFileExt(p)
	fileName := utils.GetFileName(vo.RomPath)
	newPath := config.Cfg.Platform[vo.Platform].AudioPath + config.Cfg.Separator + fileName + config.Cfg.Separator + name + ext

	rel := strings.Replace(newPath, config.Cfg.RootPath, "", 1)
	if rel == p {
		return p, nil
	}
	
	//创建目录
	_ = utils.CreateDir(utils.GetFilePath(newPath))

	//复制文件
	if err := utils.FileCopy(p, newPath); err != nil {
		return "", err
	}
	relPath := strings.Replace(newPath, config.Cfg.RootPath, "", -1)
	return relPath, nil
}

/**
 * 更新数据
 **/
func UpdateAudio(id uint64, data string) error {
	vo, _ := (&db.Rom{}).GetById(id)
	if config.Cfg.Platform[vo.Platform].FilesPath == "" {
		return errors.New(config.Cfg.Lang["AudioMenuCanNotBeEmpty"])
	}

	//整理需要删除的文件
	d := []map[string]string{}
	json.Unmarshal([]byte(data), &d)
	newData := []string{}
	for _, v := range d {
		newData = append(newData, v["path"])
	}

	//读取已存在的文件
	p := config.Cfg.Platform[vo.Platform].AudioPath + config.Cfg.Separator + vo.Name
	exists, _ := utils.ScanDir(p)
	for _, v := range exists {
		rel := strings.Replace(v, config.Cfg.RootPath, "", 1)
		if !utils.InSliceString(rel, newData) {
			utils.FileDelete(v)
		}
	}

	return nil
}
