package modules

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"simUI/code/config"
	"simUI/code/db"
	"simUI/code/utils"
	"strings"
)

/**
 * 创建缓存
 **/
func CreateRomData(platform uint32) ([]*db.Rom, map[string]*db.Menu, error) {

	romlist := []*db.Rom{}

	menuList := map[string]*db.Menu{}                                   //分类目录
	RomPath := config.Cfg.Platform[platform].RomPath                    //rom文件路径
	RomExt := strings.Split(config.Cfg.Platform[platform].RomExts, ",") //rom扩展名
	BaseInfo, err := GetRomBase(platform)

	if err != nil {
		return nil, nil, errors.New(config.Cfg.Lang["CsvFormatError"] + err.Error())
	}

	//进入循环，遍历文件
	if err := filepath.Walk(RomPath,
		func(p string, f os.FileInfo, err error) error {

			if f == nil {
				return err
			}

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
			romExt := strings.ToLower(path.Ext(p)) //获取文件后缀

			//如果该文件是游戏rom
			if f.IsDir() == false && utils.InSliceString(romExt, RomExt) {
				romName := utils.GetFileName(f.Name())
				title := romName

				//先读取基础数据，如果没有基础数据，则读取别名
				baseName := ""
				base := &RomBase{}
				aliasName := ""
				if _, ok := BaseInfo[title]; ok {
					base = BaseInfo[title]
					if BaseInfo[title].Name != "" {
						baseName = BaseInfo[title].Name
					}
				}

				if baseName != "" {
					title = baseName
				} else if aliasName != "" {
					title = aliasName
				}

				fileMd5 := GetRomMd5(utils.ToString(platform), title)
				infoMd5 := GetRomMd5(title, p, base.Type, base.Year, base.Publisher, base.Country, base.Translate)
				//如果游戏名称存在分隔符，说明是子游戏
				menu := ConstMenuRootKey //无目录，读取默认参数
				//定义目录，如果有子目录，则记录子目录名称
				if len(subpath) > 1 {
					menu = subpath[0]
				}

				//如果游戏名称存在分隔符，说明是子游戏
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
						Pinyin:        utils.TextToPinyin(sub[1]),
						InfoMd5:       infoMd5,
						FileMd5:       fileMd5,
						SimConf:       "{}",
						BaseType:      base.Type,
						BaseYear:      base.Year,
						BasePublisher: base.Publisher,
						BaseCountry:   base.Country,
						BaseTranslate: base.Translate,
					}
					romlist = append(romlist, sinfo)
				} else { //不是子游戏
					//去掉扩展名，生成标题
					rinfo := &db.Rom{
						Menu:          menu,
						Name:          title,
						Platform:      platform,
						RomPath:       p,
						Star:          0,
						Hide:          0,
						Pinyin:        utils.TextToPinyin(title),
						InfoMd5:       infoMd5,
						FileMd5:       fileMd5,
						SimConf:       "{}",
						BaseType:      base.Type,
						BaseYear:      base.Year,
						BasePublisher: base.Publisher,
						BaseCountry:   base.Country,
						BaseTranslate: base.Translate,
					}

					romlist = append(romlist, rinfo)

					//分类列表
					if menu != ConstMenuRootKey {
						menuList[menu] = &db.Menu{
							Platform: platform,
							Name:     menu,
							Pinyin:   utils.TextToPinyin(menu),
						}
					}
				}

			}
			return nil
		}); err != nil {
		fmt.Println(err)
	}

	return romlist, menuList, nil
}

/**
 * 删除不存在平台的缓存数据
 **/
func ClearPlatform() error {
	pfs := []string{}
	for k, _ := range config.Cfg.Platform {
		pfs = append(pfs, utils.ToString(k))
	}

	//清空不存在的平台（rom表）
	if err := (&db.Rom{}).ClearByPlatform(pfs); err != nil {
		return err
	}

	//清空不存在的平台（menu表）
	if err := (&db.Menu{}).ClearByPlatform(pfs); err != nil {
		return err
	}
	return nil
}

/**
 * 更新rom cache
 **/
func UpdateRomDB(platform uint32, romlist []*db.Rom) error {

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
	addAndSubFileUniq := append(addFileUniq, subFileUniq...)  //增加的和删除的

	//整理出要添加的数据体
	addData := []*db.Rom{}
	updateData := []*db.Rom{}
	for _, v := range romlist {
		if utils.InSliceString(v.FileMd5, addFileUniq) {
			addData = append(addData, v) //添加的数据
		}
		if !utils.InSliceString(v.FileMd5, addAndSubFileUniq) {
			updateData = append(updateData, v) //添加的数据
		}
	}

	//在已有数据中查找info_md5不一致的数据，就是修改的数据
	updateIssetData := []*db.Rom{}
	for _, v := range updateData {
		if !utils.InSliceString(v.InfoMd5, DbInfoMd5) {
			updateIssetData = append(updateIssetData, v) //添加的数据
		}
	}

	//删除不存在的rom
	err := (&db.Rom{}).DeleteByMd5(platform, subFileUniq)
	if err != nil {
		return err
	}

	//添加新rom
	(&db.Rom{}).BatchAdd(addData)

	//更新现有rom
	(&db.Rom{}).BatchUpdate(updateIssetData)

	return nil
}

/**
 * 更新rom cache
 **/
func UpdateMenuDB(platform uint32, menumap map[string]*db.Menu) error {

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
func UpdateFilterDB(platform uint32) {

	baseType, _ := (&db.Rom{}).GetFilter(platform, "base_type")
	baseYear, _ := (&db.Rom{}).GetFilter(platform, "base_year")
	basePublisher, _ := (&db.Rom{}).GetFilter(platform, "base_publisher")
	baseCountry, _ := (&db.Rom{}).GetFilter(platform, "base_country")
	baseTranslate, _ := (&db.Rom{}).GetFilter(platform, "base_translate")

	filters := []*db.Filter{}
	for _, v := range baseType {
		data := &db.Filter{
			Name:     v,
			Type:     "base_type",
			Platform: platform,
		}
		filters = append(filters, data)
	}

	for _, v := range baseYear {
		data := &db.Filter{
			Name:     v,
			Type:     "base_year",
			Platform: platform,
		}
		filters = append(filters, data)
	}

	for _, v := range basePublisher {
		data := &db.Filter{
			Name:     v,
			Type:     "base_publisher",
			Platform: platform,
		}
		filters = append(filters, data)
	}

	for _, v := range baseCountry {
		data := &db.Filter{
			Name:     v,
			Type:     "base_country",
			Platform: platform,
		}
		filters = append(filters, data)
	}

	for _, v := range baseTranslate {
		data := &db.Filter{
			Name:     v,
			Type:     "base_translate",
			Platform: platform,
		}
		filters = append(filters, data)
	}

	(&db.Filter{}).BatchAdd(filters)

}
