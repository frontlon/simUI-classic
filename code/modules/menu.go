package modules

import (
	"errors"
	"io/ioutil"
	"simUI/code/config"
	"simUI/code/db"
	"simUI/code/utils"
)

//读取菜单列表
func GetMenuList(platform uint32, page uint32) ([]*db.Menu, error) {
	newMenu := []*db.Menu{}

	menu, err := (&db.Menu{}).GetByPlatform(platform, page) //从数据库中读取当前平台的分类目录
	if err != nil {
		return newMenu, err
	}
	//读取根目录下是否有rom
	count, err := (&db.Rom{}).Count(0, platform, ConstMenuRootKey, "", "", "", "", "", "", "", "", "", "")
	if err != nil {
		return newMenu, err
	}

	//读取根目录下有rom，则显示未分类文件夹
	if count > 0 {
		root := &db.Menu{
			Name:     ConstMenuRootKey,
			Platform: platform,
		}
		newMenu = append(newMenu, root)
		newMenu = append(newMenu, menu...)
	} else {
		newMenu = menu
	}
	return newMenu, nil
}

//读取所有平台的菜单列表
func GetAllPlatformMenuList() (map[string][]map[string]string, error) {

	platformList, _ := (&db.Platform{}).GetAll()
	menuList, _ := (&db.Menu{}).GetAll()

	create := map[string][]map[string]string{}
	for _, platform := range platformList {

		//平台菜单根目录
		vo := map[string]string{
			"platform": utils.ToString(platform.Id),
			"name":     "/",
		}
		create[platform.Name] = append(create[platform.Name], vo)

		for _, menu := range menuList {
			if platform.Id == menu.Platform {
				vo := map[string]string{
					"platform": utils.ToString(menu.Platform),
					"name":     menu.Name,
				}
				create[platform.Name] = append(create[platform.Name], vo)
			}
		}
	}

	return create, nil

}

func AddMenu(platform uint32, name string, virtual int8) error {

	//读取目录信息
	info := (&db.Menu{}).GetByName(platform, name)

	//检查目录是否存在
	folder := config.Cfg.Platform[platform].RomPath + config.Cfg.Separator + name
	exists := utils.FolderExists(folder)

	//如果是实体目录，才去新建文件夹
	if virtual == 0 {
		if exists {
			return errors.New(config.Cfg.Lang["MenuExists"])
		}
		if err := utils.CreateDir(folder); err != nil {
			return err
		}
	} else {
		if info != nil {
			return errors.New(config.Cfg.Lang["MenuExists"])
		}
		//如果目录真实存在，则转换为实体目录
		if exists {
			virtual = 0
		}
	}

	//更新数据库
	if info == nil {
		if err := (&db.Menu{
			Name:     name,
			Platform: platform,
			Pinyin:   utils.TextToPinyin(name),
			Virtual:  virtual,
			Sort:     0,
		}).Add(); err != nil {
			return err
		}
	}
	return nil
}

func MenuRename(platform uint32, oldName string, newName string) error {
	oldMenu := config.Cfg.Platform[platform].RomPath + config.Cfg.Separator + oldName
	newMenu := config.Cfg.Platform[platform].RomPath + config.Cfg.Separator + newName

	//读取目录信息
	info := (&db.Menu{}).GetByName(platform, oldName)
	if info == nil {
		return errors.New(config.Cfg.Lang["MenuIsNotExists"])
	}

	if info.Virtual == 0 {
		if !utils.FolderExists(oldMenu) {
			return errors.New(config.Cfg.Lang["MenuIsNotExists"])
		}
		if err := utils.FolderMove(oldMenu, newMenu); err != nil {
			return err
		}
	}

	//更新数据库
	(&db.Menu{}).UpdateName(platform, oldName, newName)
	(&db.Rom{}).UpdateMenu(platform, oldName, newName)

	return nil

}

func DeleteMenu(platform uint32, name string) error {

	//读取目录信息
	info := (&db.Menu{}).GetByName(platform, name)

	if info == nil {
		return errors.New(config.Cfg.Lang["MenuIsNotExists"])
	}

	if info.Virtual == 1 {
		volist, _ := (&db.Rom{}).GetByMenu(platform, name)
		if len(volist) > 0 {
			return errors.New(config.Cfg.Lang["MenuIsNotEmpty"])
		}
	} else {
		folder := config.Cfg.Platform[platform].RomPath + config.Cfg.Separator + name
		exists := utils.FolderExists(folder)

		if exists == true {
			dir, _ := ioutil.ReadDir(folder)
			if len(dir) > 0 {
				return errors.New(config.Cfg.Lang["MenuIsNotEmpty"])
			}

			if err := utils.DeleteDir(folder); err != nil {
				return err
			}
		}
	}

	//删除数据库
	if err := (&db.Menu{
		Platform: platform,
		Name:     name,
	}).DeleteByName(); err != nil {
		return err
	}

	return nil
}
