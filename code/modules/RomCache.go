package modules

import (
	"VirtualNesGUI/code/config"
	"VirtualNesGUI/code/db"
	"VirtualNesGUI/code/utils"
	"fmt"
	"os"
	"path"
	"path/filepath"
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

	//进入循环，遍历文件
	if err := filepath.Walk(RomPath,
		func(p string, f os.FileInfo, err error) (error) {

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

				//如果有别名配置，则读取别名
				if _, ok := RomAlias[title]; ok {
					if RomAlias[title] != "" {
						if RomAlias[title] == "-"{ //如果是-，则忽略这个rom
							return nil
						}
						title = RomAlias[title]
					}
				}



				//py := TextToPinyin(title)
				md5 := GetFileUniqId(title, p, f)
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
						Name:     sub[1],
						Pname:    sub[0],
						RomPath:  p,
						Menu:     menu,
						Platform: platform,
						Star:     0,
						Pinyin:   utils.TextToPinyin(sub[1]),
						Md5:      md5,
						SimConf:  "{}",
					}

					romlist[md5] = sinfo
					md5list = append(md5list, sinfo.Md5)
				} else { //不是子游戏
					//去掉扩展名，生成标题
					rinfo := &db.Rom{
						Menu:     menu,
						Name:     title,
						Platform: platform,
						RomPath:  p,
						Star:     0,
						Pinyin:   utils.TextToPinyin(title),
						Md5:      md5,
						SimConf:  "{}",
					}

					romlist[md5] = rinfo
					md5list = append(md5list, rinfo.Md5)

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

	fileUniqs := []string{} //磁盘文件

	for k, _ := range romlist {
		fileUniqs = append(fileUniqs, k)
	}

	DbUniqs, _ := (&db.Rom{}).GetMd5ByPlatform(platform) //数据库的md5

	addUniq := utils.SliceDiff(fileUniqs, DbUniqs)
	subUniq := utils.SliceDiff(DbUniqs, fileUniqs)

	//删除rom数据
	err := (&db.Rom{}).DeleteByMd5(platform, subUniq)
	if err != nil {
		return err
	}

	//保存新数据到数据库rom表
	(&db.Rom{}).BatchAdd(addUniq, romlist)

	//这些变量可能过大，清空变量
	romlist = map[string]*db.Rom{}
	addUniq = []string{}
	subUniq = []string{}
	DbUniqs = []string{}

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

/**
 * 读取文件唯一标识
 **/
func GetFileUniqId(title string, p string, f os.FileInfo) string {
	str := title + p + utils.ToString(f.Size()) + f.ModTime().String()
	return utils.Md5(str)
}
