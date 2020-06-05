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
func CreateRomData(platform uint32) (map[string]*db.Rom, map[string]*db.Menu, error) {

	romlist := map[string]*db.Rom{}
	md5list := []string{}

	menuList := map[string]*db.Menu{}                                   //分类目录
	RomPath := config.Cfg.Platform[platform].RomPath                    //rom文件路径
	RomExt := strings.Split(config.Cfg.Platform[platform].RomExts, ",") //rom扩展名
	RomAlias, _ := config.GetRomAlias(platform)                         //别名配置
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
			newpath := strings.Replace(RomPath, "/", "\\", -1)
			newpath = strings.Replace(p, newpath, "", -1)
			if len(newpath) > 0 && newpath[0:1] == "\\" {
				newpath = strings.Replace(newpath, "\\", "", 1)
			}
			subpath := strings.Split(newpath, "\\")
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

				if _, ok := RomAlias[title]; ok {
					if RomAlias[title] != "" {
						aliasName = RomAlias[title]
					}
				}
				if baseName != "" {
					title = baseName
				} else if aliasName != "" {
					title = aliasName
				}

				pathMd5 := GetPathMd5(title, p, base.Type, base.Year, base.Developer, base.Publisher) //路径md5，可变
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
						Pinyin:        utils.TextToPinyin(sub[1]),
						PathMd5:       pathMd5,
						SimConf:       "{}",
						BaseType:      base.Type,
						BaseYear:      base.Year,
						BasePublisher: base.Publisher,
					}

					romlist[pathMd5] = sinfo
					md5list = append(md5list, sinfo.PathMd5)
				} else { //不是子游戏
					//去掉扩展名，生成标题
					rinfo := &db.Rom{
						Menu:          menu,
						Name:          title,
						Platform:      platform,
						RomPath:       p,
						Star:          0,
						Pinyin:        utils.TextToPinyin(title),
						PathMd5:       pathMd5,
						SimConf:       "{}",
						BaseType:      base.Type,
						BaseYear:      base.Year,
						BasePublisher: base.Publisher,
					}

					romlist[pathMd5] = rinfo
					md5list = append(md5list, rinfo.PathMd5)

					fmt.Println("rinfo", rinfo)
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
func UpdateRomDB(platform uint32, romlist map[string]*db.Rom) error {

	md5s := []string{} //磁盘文件
	err := errors.New("")
	for k, _ := range romlist {
		md5s = append(md5s, k) //文件的md5
	}

	//数据库中读取md5和fileid
	DbMd5s, _ := (&db.Rom{}).GetMd5ByPlatform(platform)
	addUniq := utils.SliceDiff(md5s, DbMd5s) //新增的
	subUniq := utils.SliceDiff(DbMd5s, md5s) //删除的

	//2.删除不存在的rom
	err = (&db.Rom{}).DeleteByMd5(platform, subUniq)
	if err != nil {
		return err
	}

	//3.添加新rom
	(&db.Rom{}).BatchAdd(addUniq, romlist)

	return nil
}

/**
 * 更新rom cache
 **/
func UpdateMenuDB(platform uint32, menumap map[string]*db.Menu) error {

	menus := []string{}
	if len(menumap) > 0 {
		for k, _ := range menumap {
			if k == ConstMenuRootKey {
				continue
			}
			menus = append(menus, k)
		}
	}

	//删除当前平台下不存在的菜单
	if err := (&db.Menu{}).DeleteNotExists(platform, menus); err != nil {
	}

	//查询已存在的记录
	issetName, err := (&db.Menu{}).GetMenuByNames(platform, menus)
	if err != nil {
		return err
	}

	//取出需要写入数据库的rom数据。
	saveMenulist := []*db.Menu{}
	for _, v := range menumap {
		if utils.InSliceString(v.Name, issetName) == false {
			saveMenulist = append(saveMenulist, v)
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
func GetPathMd5(par ...string) string {
	str := strings.Join(par, ",")
	return utils.Md5(str)
}

/**
 * 更新filter cache
 **/
func UpdateFilterDB() {
	_ = (&db.Filter{}).Truncate()

	baseType, _ := (&db.Rom{}).GetFilter("BaseType")
	baseYear, _ := (&db.Rom{}).GetFilter("BaseYear")
	basePlatform, _ := (&db.Rom{}).GetFilter("BasePlatform")
	basePublisher, _ := (&db.Rom{}).GetFilter("BasePublisher")

	filters := []*db.Filter{}
	for _, v := range baseType {
		data := &db.Filter{
			Name: v.BaseType,
		}
		filters = append(filters, data)
	}

	for _, v := range baseYear {
		data := &db.Filter{
			Name: v.BaseYear,
		}
		filters = append(filters, data)
	}

	for _, v := range basePlatform {
		data := &db.Filter{
			Name: v.BasePlatform,
		}
		filters = append(filters, data)
	}

	for _, v := range basePublisher {
		data := &db.Filter{
			Name: v.BasePlatform,
		}
		filters = append(filters, data)
	}

	(&db.Filter{}).BatchAdd(filters)

}
