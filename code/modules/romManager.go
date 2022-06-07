package modules

import (
	"errors"
	"os"
	"path/filepath"
	"simUI/code/config"
	"simUI/code/db"
	"simUI/code/utils"
	"strings"
)

//查询重复rom
func CheckRomRepeat(platformId uint32) ([]map[string]interface{}, error) {

	romlist, _ := (&db.Rom{}).GetByPlatform(platformId)

	repeatList := map[int64][]map[string]interface{}{}
	for _, v := range romlist {
		//如果是相对路径，转换成绝对路径
		if !strings.Contains(v.RomPath, ":") {
			v.RomPath = config.Cfg.Platform[v.Platform].RomPath + config.Cfg.Separator + v.RomPath
		}
		f, err := os.Stat(v.RomPath)

		if err != nil {
			continue
		}

		rom := map[string]interface{}{}
		rom["id"] = v.Id
		rom["path"] = v.RomPath
		rom["name"] = v.Name
		rom["size"] = utils.ToString(f.Size())
		repeatList[f.Size()] = append(repeatList[f.Size()], rom)
	}

	result := []map[string]interface{}{}

	for _, v := range repeatList {
		if len(v) <= 1 {
			continue
		}
		for _, b := range v {
			result = append(result, b)
		}

	}
	return result, nil
}

//移动文件到其他目录
func MoveRomByFile(f string, p string) error {

	fileName := utils.GetFileNameAndExt(f)
	newPath := p + config.Cfg.Separator + fileName

	if err := utils.FileMove(f, newPath); err != nil {
		return err
	}

	return nil
}

//移动僵尸文件到其他目录
func MoveZombieByFile(f string, p string) error {

	fileName := utils.GetFileName(f)
	ext := utils.GetFileExt(f)
	oldPathArr := strings.Split(utils.GetFileAbsPath(f), config.Cfg.Separator)
	newPath := p + config.Cfg.Separator + fileName + "_" + oldPathArr[len(oldPathArr)-1] + ext
	if err := utils.FileMove(f, newPath); err != nil {
		return err
	}

	return nil
}

//查询无效资源
func CheckRomZombie(platformId uint32) ([]map[string]string, error) {

	romlist, _ := (&db.Rom{}).GetMasterRomByPlatform(platformId)

	notExistsList := []map[string]string{}
	existsMap := map[string]string{}
	//读取已存在rom
	for _, v := range romlist {
		name := utils.GetFileName(v.RomPath)
		existsMap[name] = ""
	}

	res := config.GetResPath(platformId)

	//先检查重复资料
	for k, path := range res {
		//攻略文件单独去处理
		if k == "files" || k == "audio" {
			continue
		}
		existsList := map[string][]string{}
		if path == "" {
			continue
		}
		if err := filepath.Walk(path,
			func(p string, f os.FileInfo, err error) error {
				if f == nil {
					return nil
				}
				if f.IsDir() == true {
					return nil
				}
				name := utils.GetFileName(p)
				if name == "" {
					return nil
				}
				//检查子游戏
				if strings.Contains(p, "__") {
					repeat := map[string]string{}
					repeat["path"] = p
					repeat["type"] = "3"
					notExistsList = append(notExistsList, repeat)
				} else if _, ok := existsMap[name]; !ok {
					//检查无效文件
					repeat := map[string]string{}
					repeat["path"] = p
					repeat["type"] = "1"
					notExistsList = append(notExistsList, repeat)
				} else {
					//检查重复文件
					existsList[name] = append(existsList[name], p)
					if len(existsList[name]) > 1 {
						repeat := map[string]string{}
						repeat["path"] = p
						repeat["type"] = "2"
						notExistsList = append(notExistsList, repeat)
					}
				}
				return nil
			}); err != nil {

		}
	}

	//处理攻略，攻略文件都是以__分割的文件名，且一个rom可存在多个攻略文件
	if err := filepath.Walk(res["files"],
		func(p string, f os.FileInfo, err error) error {
			if f == nil {
				return nil
			}
			if f.IsDir() == true {
				return nil
			}
			name := utils.GetFileName(p)
			if name == "" {
				return nil
			}
			//检查子游戏
			nameArr := strings.Split(name, "__")
			if _, ok := existsMap[nameArr[0]]; !ok {
				//检查无效文件
				repeat := map[string]string{}
				repeat["path"] = p
				repeat["type"] = "1"
				notExistsList = append(notExistsList, repeat)
			}
			return nil
		}); err != nil {
	}

	//处理音频，音频文件都是以__分割的文件名，且一个rom可存在多个音频文件
	if err := filepath.Walk(res["audio"],
		func(p string, f os.FileInfo, err error) error {
			if f == nil {
				return nil
			}
			if f.IsDir() == true {
				return nil
			}
			name := utils.GetFileName(p)
			if name == "" {
				return nil
			}
			//检查子游戏
			nameArr := strings.Split(name, "__")
			if _, ok := existsMap[nameArr[0]]; !ok {
				//检查无效文件
				repeat := map[string]string{}
				repeat["path"] = p
				repeat["type"] = "1"
				notExistsList = append(notExistsList, repeat)
			}
			return nil
		}); err != nil {
	}

	return notExistsList, nil
}

//绑定子游戏
func BindSubGame(pid uint64, sid uint64) (*db.Rom, error) {
	master, err := (&db.Rom{}).GetById(pid)
	if err != nil {
		return nil, errors.New(config.Cfg.Lang["MasterGameNotFound"])
	}

	slave, err := (&db.Rom{}).GetById(sid)
	if err != nil {
		return nil, errors.New(config.Cfg.Lang["SlaveGameNotFound"])
	}

	subRoms, err := (&db.Rom{}).GetSubRom(master.Platform,slave.FileMd5)
	if len(subRoms) > 0{
		return nil, errors.New(config.Cfg.Lang["SubGameExists"])
	}

	//更新数据库
	if (&db.Rom{}).UpdatePnameById(sid, master.FileMd5) != nil {
		return nil, err
	}
	if (&db.RomSubGame{}).UpdatePname(slave.Platform,slave.FileMd5,master.FileMd5) != nil {
		return nil, err
	}

	//整理返回值
	slave.Pname = master.FileMd5

	return slave, nil
}

//解绑子游戏
func UnBindSubGame(id uint64) (*db.Rom, error) {
	//读取子游戏数据
	vo, err := (&db.Rom{}).GetById(id)

	//更新数据库
	if (&db.Rom{}).UpdatePnameById(id, "") != nil {
		return nil, err
	}
	if (&db.RomSubGame{}).DeleteByFileMd5(vo.Platform,vo.FileMd5) != nil {
		return nil, err
	}

	//整理返回值
	vo.Pname = ""

	return vo, nil
}
