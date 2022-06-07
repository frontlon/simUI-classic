package modules

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"simUI/code/compoments"
	"simUI/code/config"
	"simUI/code/db"
	"simUI/code/utils"
	"simUI/code/utils/go-sciter"
	"strings"
	"time"
)

/**
 * 创建缓存入口
 **/
func CreateRomCache(getPlatform uint32) error {

	//检查更新一个平台还是所有平台
	PlatformList := map[uint32]*db.Platform{}
	if getPlatform == 0 { //所有平台
		PlatformList = config.Cfg.Platform
	} else { //一个平台
		if _, ok := config.Cfg.Platform[getPlatform]; ok {
			PlatformList[getPlatform] = config.Cfg.Platform[getPlatform]
		}
	}

	//先检查平台，将不存在的平台数据先干掉
	ClearPlatform()

	//开始重建缓存
	for platform, _ := range PlatformList {

		//第一步：更新文件时间，并获取文件id列表
		fileMap, _ := updateFileDate(platform)

		//更新数据库file_md5
		updateRomFileMd5(platform, fileMap)

		//第二步：读取rom缓存
		romlist, err := getRomData(platform)
		if err != nil {
			utils.WriteLog(err.Error())
			continue
		}

		//第三步：更新rom数据
		romlist = updateRomDB(platform, fileMap, romlist)

		//第四步：读取rom目录
		menu, _ := createMenuList(platform)

		//第五步：更新menu数据
		updateMenuDB(platform, menu)

		//第六步：更新filter数据
		updateFilterDB(platform, romlist)

		//第六步：清理rom_config和rom_sub
		updateRomConfig(platform, romlist)
	}

	//收缩数据库
	db.Vacuum()

	//数据更新完成后，页面回调，更新页面DOM
	if _, err := utils.Window.Call("CB_createCache", sciter.NewValue("")); err != nil {
		fmt.Println(err)
	}
	return nil

}

/**
 * 创建缓存
 **/
func getRomData(platform uint32) ([]*db.Rom, error) {

	romlist := []*db.Rom{}
	fileMd5List := map[string]string{}
	RomPath := config.Cfg.Platform[platform].RomPath //rom文件路径

	//获取扩展名并转换成map
	RomExt := strings.Split(config.Cfg.Platform[platform].RomExts, ",") //rom扩展名
	RomExtMap := map[string]bool{}
	for _, v := range RomExt {
		RomExtMap[v] = true
	}

	//csv游戏信息
	BaseInfo, err := GetRomBaseList(platform)

	//rom配置信息（喜爱，隐藏）
	romConfig, err := (&db.RomSetting{}).GetByPlatformToMap(platform)

	//读取子游戏关系
	romSubGame, err := (&db.RomSubGame{}).GetByPlatformToMap(platform)

	//从数据库中读取file_md5值
	//dbData, _ := (&db.Rom{}).GetFileMd5ByPlatform(platform)

	if err != nil {
		return nil, errors.New(config.Cfg.Lang["CsvFormatError"] + err.Error())
	}

	//进入循环，遍历文件
	if err := filepath.Walk(RomPath,
		func(p string, f os.FileInfo, err error) error {

			if f == nil {
				return err
			}

			romExt := strings.ToLower(path.Ext(p)) //获取文件后缀

			//转换为相对路径
			p = strings.Replace(p, RomPath+config.Cfg.Separator, "", -1)

			//整理目录格式，并转换为数组
			newpath := strings.Replace(RomPath, "/", config.Cfg.Separator, -1)
			newpath = strings.Replace(RomPath, "\\", config.Cfg.Separator, -1)
			newpath = strings.Replace(p, newpath, "", -1)
			if len(newpath) > 0 && newpath[0:1] == "/" {
				newpath = strings.Replace(newpath, "/", "", 1)
			}
			subpath := strings.Split(newpath, config.Cfg.Separator)
			romExists := false //rom是否存在

			if _, ok := RomExtMap[romExt]; ok {
				romExists = true //rom存在
			}

			//如果该文件是游戏rom，并是否存在
			if f.IsDir() == false && romExists == true {
				romName := utils.GetFileName(f.Name())
				title := romName

				fileSize := utils.GetFileSizeString(f.Size())

				//fileMd5 := ""
				//if _, ok := dbData[p]; ok {
				//	fileMd5 = dbData[p]
				//} else {
				fileMd5 := utils.CreateRomUniqId(utils.ToString(f.ModTime().UnixNano()), f.Size())
				//}

				//先读取基础数据，如果没有基础数据，则读取别名
				baseName := ""
				base := &RomBase{}
				if _, ok := BaseInfo[romName]; ok {
					base = BaseInfo[romName]
					if BaseInfo[romName].Name != "" {
						baseName = BaseInfo[romName].Name
					}
				}

				if baseName != "" {
					title = baseName
				}

				infoMd5 := GetRomMd5(title, p, base.Type, base.Year, base.Producer, base.Publisher, base.Country, base.Translate, base.Version, base.NameEN, base.NameJP, base.OtherA, base.OtherB, base.OtherC, base.OtherD, fileSize)

				//如果游戏名称存在分隔符，说明是子游戏
				menu := ConstMenuRootKey //无目录，读取默认参数
				//定义目录，如果有子目录，则记录子目录名称
				if len(subpath) > 1 {
					menu = subpath[0]
				}

				//读取rom大小
				//如果游戏名称存在分隔符，说明是老版本子游戏
				if strings.Contains(title, ConstSeparator) {

					//拆分文件名
					sub := strings.Split(title, ConstSeparator)
					//去掉扩展名，生成标题
					sinfo := &db.Rom{
						Name:          sub[1],
						Pname:         sub[0],
						RomPath:       p,
						Menu:          menu,
						Platform:      platform,
						Star:          0,
						Hide:          0,
						Size:          fileSize,
						Pinyin:        utils.TextToPinyin(sub[1]),
						InfoMd5:       infoMd5,
						FileMd5:       fileMd5,
						SimConf:       "{}",
						BaseType:      base.Type,
						BaseYear:      base.Year,
						BaseProducer:  base.Producer,
						BasePublisher: base.Publisher,
						BaseCountry:   base.Country,
						BaseTranslate: base.Translate,
						BaseVersion:   base.Version,
						BaseNameEn:    base.NameEN,
						BaseNameJp:    base.NameJP,
						BaseOtherA:    base.OtherA,
						BaseOtherB:    base.OtherB,
						BaseOtherC:    base.OtherC,
						BaseOtherD:    base.OtherD,
					}
					romlist = append(romlist, sinfo)
				} else { //不是子游戏

					star := uint8(0)
					hide := uint8(0)
					runNum := uint64(0)
					runLastTime := int64(0)
					simId := uint32(0)
					simConf := ""
					if romConfig[fileMd5] != nil {
						star = romConfig[fileMd5].Star
						hide = romConfig[fileMd5].Hide
						runNum = romConfig[fileMd5].RunNum
						runLastTime = romConfig[fileMd5].RunLasttime
						simId = romConfig[fileMd5].SimId
						simConf = romConfig[fileMd5].SimConf
					}

					fileMd5List[title] = fileMd5
					fileMd5List[romName] = fileMd5

					pname := ""
					if _, ok := romSubGame[fileMd5]; ok {
						pname = romSubGame[fileMd5]
					}

					rinfo := &db.Rom{
						Menu:          menu,
						Name:          title,
						Platform:      platform,
						Pname:         pname,
						RomPath:       p,
						Star:          star,
						Hide:          hide,
						Size:          fileSize,
						Pinyin:        utils.TextToPinyin(title),
						InfoMd5:       infoMd5,
						FileMd5:       fileMd5,
						SimId:         simId,
						SimConf:       simConf,
						BaseType:      base.Type,
						BaseYear:      base.Year,
						BaseProducer:  base.Producer,
						BasePublisher: base.Publisher,
						BaseCountry:   base.Country,
						BaseTranslate: base.Translate,
						BaseVersion:   base.Version,
						BaseNameEn:    base.NameEN,
						BaseNameJp:    base.NameJP,
						BaseOtherA:    base.OtherA,
						BaseOtherB:    base.OtherB,
						BaseOtherC:    base.OtherC,
						BaseOtherD:    base.OtherD,
						RunNum:        runNum,
						RunLasttime:   runLastTime,
					}

					romlist = append(romlist, rinfo)
				}

			}
			return nil
		}); err != nil {
		return nil, errors.New(config.Cfg.Lang["RomMenuCanNotBeExists"])
	}

	//更新老版本子游戏 __pname
	if len(romlist) > 0 {
		for k, rom := range romlist {
			if rom.Pname != "" && !utils.HasRomUniqId(rom.Pname) {
				if _, ok := fileMd5List[rom.Pname]; ok {
					romlist[k].Pname = fileMd5List[rom.Pname]
				} else {
					romlist[k].Pname = ""
				}
			}
		}
	}
	return romlist, nil
}

/*
更新文件日期
*/
func updateFileDate(platform uint32) (map[string]string, error) {

	create := map[string][]string{}
	fileMap := map[string]string{}
	RomPath := config.Cfg.Platform[platform].RomPath //rom文件路径

	//获取扩展名并转换成map
	RomExt := strings.Split(config.Cfg.Platform[platform].RomExts, ",") //rom扩展名
	RomExtMap := map[string]bool{}
	for _, v := range RomExt {
		RomExtMap[v] = true
	}

	//进入循环，遍历文件
	if err := filepath.Walk(RomPath,
		func(p string, f os.FileInfo, err error) error {

			if f == nil {
				return err
			}

			romExt := strings.ToLower(path.Ext(p)) //获取文件后缀
			romExists := false                     //rom是否存在

			if _, ok := RomExtMap[romExt]; ok {
				romExists = true //rom存在
			}

			//如果是目录，或不是rom，则跳过
			if f.IsDir() == true || romExists == false {
				return nil
			}

			//检查修改时间是否完整
			fileMd5 := utils.CreateRomUniqId(utils.ToString(f.ModTime().UnixNano()), f.Size())

			c := create[fileMd5]
			if _, ok := create[fileMd5]; ok {
				c = append(c, p)
			} else {
				c = []string{p}
			}
			create[fileMd5] = c

			//转换为相对路径
			fp := strings.Replace(p, RomPath+config.Cfg.Separator, "", -1)

			fileMap[fp] = fileMd5

			return nil
		}); err != nil {
		return map[string]string{}, errors.New(config.Cfg.Lang["RomMenuCanNotBeExists"])
	}

	//开始检查是否有重复的md5，如果有重复的，则修改文件时间
	if len(create) > 0 {
		i := 1
		for _, v := range create {
			if len(v) > 1 {
				for _, p := range v {
					stat, _ := os.Stat(p)
					modTime := utils.ToString(stat.ModTime().UnixNano())

					t := ""
					s := 0
					e := 0
					if modTime == "0" || len(modTime) < 10 {
						t = utils.ToString(time.Now().UnixNano())
						s = utils.ToInt(t[:10]) + i
						e = utils.ToInt(t[11:])
					} else {
						t = utils.ToString(stat.ModTime().UnixNano())
						n := utils.ToString(time.Now().UnixNano())
						s = utils.ToInt(t[:10]) + i
						e = utils.ToInt(n[11:])
					}
					tm := time.Unix(int64(s), int64(e))

					//修改文件时间
					os.Chtimes(p, tm, tm)

					//获取md5
					fp := strings.Replace(p, RomPath+config.Cfg.Separator, "", -1)
					fileMap[fp] = utils.CreateRomUniqId(utils.ToString(tm.UnixNano()), stat.Size())
				}
			}
		}
	}

	return fileMap, nil
}

//创建菜单列表
func createMenuList(platform uint32) (map[string]*db.Menu, error) {

	menuList := map[string]*db.Menu{}

	FileInfo, err := ioutil.ReadDir(config.Cfg.Platform[platform].RomPath)
	if err != nil {
		return menuList, err
	}
	for _, v := range FileInfo {
		if v.IsDir() == true {
			menuList[v.Name()] = &db.Menu{
				Platform: platform,
				Name:     v.Name(),
				Pinyin:   utils.TextToPinyin(v.Name()),
			}
		}
	}
	return menuList, nil

}

/**
 * 删除不存在平台的缓存数据
 **/
func ClearPlatform() {
	pfs := []string{}
	for k, _ := range config.Cfg.Platform {
		pfs = append(pfs, utils.ToString(k))
	}

	//清空不存在的平台（rom表）
	if err := (&db.Rom{}).ClearByNotPlatform(pfs); err != nil {
		fmt.Println(err)
	}

	//清空不存在的平台（menu表）
	if err := (&db.Menu{}).ClearByNotPlatform(pfs); err != nil {
		fmt.Println(err)
	}

	//清空不存在的平台（rom_subpage表）
	if err := (&db.RomSubGame{}).ClearByNotPlatform(pfs); err != nil {
		fmt.Println(err)
	}

	//清空不存在的平台（rom_setting表）
	if err := (&db.RomSetting{}).ClearByNotPlatform(pfs); err != nil {
		fmt.Println(err)
	}
}

/**
 * 更新rom file_md5
 **/
func updateRomFileMd5(platform uint32, fileMap map[string]string) {

	romlist, _ := (&db.Rom{}).GetByPlatform(platform)

	//更新变更的rom唯一id
	dbMap := map[string]string{}
	for _, v := range romlist {
		dbMap[v.RomPath] = v.FileMd5
	}

	//更新rom唯一id
	replaceFileMd5 := []map[string]string{}
	if len(fileMap) > 0 && len(romlist) > 0 {
		for _, v := range romlist {

			if _, ok := fileMap[v.RomPath]; !ok {
				continue
			}
			if v.FileMd5 != fileMap[v.RomPath] {
				fmt.Println(v.FileMd5, fileMap[v.RomPath])
				m := map[string]string{}
				m["romPath"] = v.RomPath
				m["oldMd5"] = v.FileMd5
				m["newMd5"] = fileMap[v.RomPath]
				replaceFileMd5 = append(replaceFileMd5, m)
			}
		}

		//更新数据库
		(&db.Rom{}).BatchUpdateFileMd5(platform, replaceFileMd5)
		(&db.RomSubGame{}).BatchUpdateFileMd5(platform, replaceFileMd5)
		(&db.RomSetting{}).BatchUpdateFileMd5(platform, replaceFileMd5)
	}
}

/**
 * 更新rom cache
 **/
func updateRomDB(platform uint32, fileMap map[string]string, romlist []*db.Rom) []*db.Rom {

	romlistInfoMd5 := []string{} //磁盘文件
	romlistFileMd5 := []string{} //磁盘文件
	for _, v := range romlist {  //从romlist列表中抽出两个md5
		romlistInfoMd5 = append(romlistInfoMd5, v.InfoMd5)
		romlistFileMd5 = append(romlistFileMd5, v.FileMd5)
	}

	//数据库中抽出两个md5
	DbFileMd5, DbInfoMd5, _ := (&db.Rom{}).GetMd5ByPlatform(platform)
	addFileUniq := utils.SliceDiff(romlistFileMd5, DbFileMd5) //新增的
	subFileUniq := utils.SliceDiff(DbFileMd5, romlistFileMd5) //删除的

	//整理出要添加的数据体
	addData := []*db.Rom{}
	updateData := []*db.Rom{}

	for _, v := range romlist {
		if utils.InSliceString(v.FileMd5, addFileUniq) {
			addData = append(addData, v) //添加的数据
		} else if !utils.InSliceString(v.InfoMd5, DbInfoMd5) {
			updateData = append(updateData, v) //更新的数据
		}
	}

	//在已有数据中查找info_md5不一致的数据，就是修改的数据
	updateIssetData := []*db.Rom{}
	updateFileMd5 := []string{}
	for _, v := range updateData {

		if !utils.InSliceString(v.InfoMd5, DbInfoMd5) {
			updateFileMd5 = append(updateFileMd5, v.FileMd5)
			updateIssetData = append(updateIssetData, v) //添加的数据
		}
	}

	//删除重复数据
	(&db.Rom{}).DeleteRepeat(platform)

	//删除不存在的rom
	(&db.Rom{}).DeleteByMd5(platform, subFileUniq)

	//写入新rom
	(&db.Rom{}).BatchAdd(addData, 1)

	//更新现有rom
	(&db.Rom{}).DeleteByMd5(platform, updateFileMd5)
	(&db.Rom{}).BatchAdd(updateIssetData, 1)

	return romlist
}

/**
 * 更新menu cache
 **/
func updateMenuDB(platform uint32, menumap map[string]*db.Menu) error {

	//磁盘中目录列表
	diskMenus := []string{}
	if len(menumap) > 0 {
		for k, _ := range menumap {
			if k == ConstMenuRootKey {
				continue
			}
			diskMenus = append(diskMenus, k)
		}
	}

	//数据库中目录列表
	dbNames, err := (&db.Menu{}).GetAllNamesByPlatform(platform)
	if err != nil {
		return err
	}

	add := utils.SliceDiff(diskMenus, dbNames)
	sub := utils.SliceDiff(dbNames, diskMenus)

	//删除当前平台下不存在的菜单
	if err := (&db.Menu{}).DeleteNotExists(platform, sub); err != nil {
	}

	//取出需要写入数据库的rom数据。
	saveMenulist := []*db.Menu{}
	if len(add) > 0 {
		for _, v := range add {
			saveMenulist = append(saveMenulist, menumap[v])
		}
	}

	//保存数据到数据库cate表
	if len(saveMenulist) > 0 {
		for _, v := range saveMenulist {
			if err := v.Add(); err != nil {
			}
		}

	}

	//这些变量较大，写入完成后清理变量
	saveMenulist = []*db.Menu{}
	menumap = map[string]*db.Menu{}

	return nil
}

//读取路径Md5
func GetRomMd5(par ...string) string {
	str := strings.Join(par, ",")
	return utils.Md5(str)
}

/**
 * 更新filter cache
 **/
func updateFilterDB(platform uint32, romlist []*db.Rom) {

	dbf, _ := (&db.Filter{}).GetByPlatform(platform)
	romFilters, dbFilter := compoments.FilterFactory(romlist, dbf)

	//增加过滤器
	if len(romFilters) > 0 {
		addFilter := []string{}
		for t, v := range romFilters {
			if _, ok := dbFilter[t]; !ok {
				//如果数据库没数据，则添加全部
				addFilter = v
			} else {
				addFilter = utils.SliceDiff(romFilters[t], dbFilter[t])
			}

			//开始写入数据库
			if len(addFilter) > 0 {
				create := []*db.Filter{}
				for _, name := range addFilter {
					f := &db.Filter{
						Platform: platform,
						Type:     t,
						Name:     name,
					}
					create = append(create, f)
				}
				(&db.Filter{}).BatchAdd(create)
			}

		}
	}

	//删除过滤器
	if len(dbFilter) > 0 {
		subFilter := []string{}
		for t, v := range dbFilter {
			if _, ok := romFilters[t]; !ok {
				//如果数据库没数据，则添加全部
				subFilter = v
			} else {
				subFilter = utils.SliceDiff(dbFilter[t], romFilters[t])
			}

			if len(subFilter) > 0 {
				(&db.Filter{}).DeleteByFileNames(platform, t, subFilter)
			}
		}
	}

}

func updateRomConfig(platform uint32, romlist []*db.Rom) {

	if len(romlist) == 0 {
		return
	}

	md5s := []string{}
	for _, v := range romlist {
		md5s = append(md5s, v.FileMd5)
	}

	RomSettingDbList, _ := (&db.RomSetting{}).GetFileMd5ByPlatform(platform)
	RomSettingDiff := utils.SliceDiff(RomSettingDbList, md5s)
	(&db.RomSetting{}).DeleteByFileMd5s(platform, RomSettingDiff)

	RomSubGameDbList, _ := (&db.RomSubGame{}).GetFileMd5ByPlatform(platform)
	RomSubGameDiff := utils.SliceDiff(RomSubGameDbList, md5s)
	(&db.RomSubGame{}).DeleteByFileMd5s(platform, RomSubGameDiff)
}

/**
 * 清空缓存
 */
/**
 * 创建缓存入口
 **/
func TruncateRomCache(getPlatform uint32) error {

	//检查更新一个平台还是所有平台
	PlatformList := map[uint32]*db.Platform{}
	if getPlatform == 0 { //所有平台
		PlatformList = config.Cfg.Platform
	} else { //一个平台
		if _, ok := config.Cfg.Platform[getPlatform]; ok {
			PlatformList[getPlatform] = config.Cfg.Platform[getPlatform]
		}
	}

	//开始重建缓存
	for platform, _ := range PlatformList {

		//清空rom表
		if err := (&db.Rom{Platform: platform}).DeleteByPlatform(); err != nil {
			return err
		}

		//清空menu表
		if err := (&db.Menu{Platform: platform}).DeleteByPlatform(); err != nil {
			return err
		}

	}

	//收缩数据库
	db.Vacuum()

	if _, err := utils.Window.Call("CB_clearDB"); err != nil {
	}

	return nil

}
