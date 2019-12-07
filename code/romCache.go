package main

import (
	"VirtualNesGUI/code/db"
	"VirtualNesGUI/code/utils"
	"os"
	"path"
	"path/filepath"
	"strings"
)


/**
 * 创建缓存
 **/
func CreateRomCache(platform uint32) ([]*db.Rom, map[string]*db.Menu, error) {
	romlist := []*db.Rom{}
	menuList := map[string]*db.Menu{}            //分类目录
	RomPath := Config.Platform[platform].RomPath //rom文件路径
	RomExt := Config.Platform[platform].RomExts  //rom扩展名
	RomAlias, _ := getRomAlias(platform)         //别名配置

	//进入循环，遍历文件
	if err := filepath.Walk(RomPath,
		func(p string, f os.FileInfo, err error) (error) {
			if f == nil {
				return err
			}

			//转换为相对路径
			p = strings.Replace(p,RomPath + separator,"",-1)

			//整理目录格式，并转换为数组
			newpath := strings.Replace(RomPath, "/", "\\", -1)
			newpath = strings.Replace(p, newpath, "", -1)
			if len(newpath) > 0 && newpath[0:1] == "\\" {
				newpath = strings.Replace(newpath, "\\", "", 1)
			}
			subpath := strings.Split(newpath, "\\")
			romExt := strings.ToLower(path.Ext(p)) //获取文件后缀

			//如果该文件是游戏rom
			if f.IsDir() == false && utils.InSliceString(romExt,RomExt) {
				romName := utils.GetFileName(f.Name())
				title := romName
				//如果有别名配置，则读取别名
				if _, ok := RomAlias[title]; ok {
					title = RomAlias[title]
				}

				py := TextToPinyin(title)
				md5 := GetFileUniqId(f)
				//如果游戏名称存在分隔符，说明是子游戏
				menu := constMenuRootKey //无目录，读取默认参数
				//定义目录，如果有子目录，则记录子目录名称
				if len(subpath) > 1 {
					menu = subpath[0]
				}

				//如果游戏名称存在分隔符，说明是子游戏
				if strings.Contains(title, constSeparator) {

					//拆分文件名
					sub := strings.Split(title, constSeparator)

					//去掉扩展名，生成标题
					sinfo := &db.Rom{
						Name:     sub[1],
						Pname:    sub[0],
						RomPath:  p,
						Menu:     menu,
						Platform: platform,
						Star:     0,
						Pinyin:   py,
						Md5:      md5,
					}
					romlist = append(romlist, sinfo)
				} else { //不是子游戏
					//去掉扩展名，生成标题
					rinfo := &db.Rom{
						Menu:     menu,
						Name:     title,
						Platform: platform,
						RomPath:  p,
						Star:     0,
						Pinyin:   py,
						Md5:      md5,
					}

					//rom列表
					romlist = append(romlist, rinfo)
					//分类列表
					if menu != constMenuRootKey {
						menuList[menu] = &db.Menu{
							Platform: platform,
							Name:     menu,
							Pinyin:   TextToPinyin(menu),
						}
					}
				}
			}
			return nil
		}); err != nil {
	}

	return romlist, menuList,nil
}


/**
 * 删除不存在平台的缓存数据
 **/
func ClearPlatform() error{
	pfs := []string{}
	for k, _ := range Config.Platform {
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
func UpdateRomDB(platform uint32,romlist []*db.Rom) error{

	uniqs := []string{}

	for _,v := range romlist{
		uniqs = append(uniqs,v.Md5)
	}

	//删除当前平台下，不存在的rom
	if err := (&db.Rom{}).DeleteNotExists(platform, uniqs); err != nil {
		return err
	}

	//查询已存在的记录
	issetMd5, err := (&db.Rom{}).GetMd5ByMd5(platform, uniqs)
	if err != nil {
		return err
	}

	//取出需要写入数据库的rom数据。
	saveRomlist := []*db.Rom{}

	for _, v := range romlist {
		if utils.InSliceString(v.Md5, issetMd5) == false {
			saveRomlist = append(saveRomlist, v)
		}
	}

	//保存新数据到数据库rom表
	if len(saveRomlist) > 0 {
		for _, v := range saveRomlist {
			if err := v.Add(); err != nil {
				return err
			}
		}
	}

	//这些变量可能过大，清空变量
	romlist = []*db.Rom{}
	saveRomlist = []*db.Rom{}

	return nil
}




/**
 * 更新rom cache
 **/
func UpdateMenuDB(platform uint32,menumap map[string]*db.Menu) error{

	menus := []string{}
	if len(menumap) > 0{
		for k,_ := range menumap{
			if k == constMenuRootKey{
				continue
			}
			menus = append(menus,k)
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

