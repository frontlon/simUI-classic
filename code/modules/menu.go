package modules

import (
	"errors"
	"fmt"
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
	count, err := (&db.Rom{}).Count(platform, ConstMenuRootKey, "", "", "", "", "", "")
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

func AddMenu(platform uint32,name string) error {
	folder := config.Cfg.Platform[platform].RomPath + config.Cfg.Separator + name

	if utils.FolderExists(folder){
		return errors.New(config.Cfg.Lang["MenuExists"])
	}
	fmt.Println("folder:",folder)
	if err := utils.CreateDir(folder);err != nil{
		return err
	}
	return nil
}

func MenuRename(platform uint32,oldName string,newName string) error {
	oldMenu := config.Cfg.Platform[platform].RomPath + config.Cfg.Separator + oldName
	newMenu := config.Cfg.Platform[platform].RomPath + config.Cfg.Separator + newName

	if !utils.FolderExists(oldMenu){
		return errors.New(config.Cfg.Lang["MenuIsNotExists"])
	}
	if err := utils.FolderMove(oldMenu,newMenu);err != nil{
		return err
	}
	return nil

}

func DeleteMenu(platform uint32,name string) error {
	folder := config.Cfg.Platform[platform].RomPath + config.Cfg.Separator + name

	if !utils.FolderExists(folder){
		return errors.New(config.Cfg.Lang["MenuIsNotExists"])
	}

	dir, _ := ioutil.ReadDir(folder)
	if len(dir) > 0{
		return errors.New(config.Cfg.Lang["MenuIsNotEmpty"])
	}


	if err := utils.DeleteDir(folder);err != nil{
		return err
	}
	return nil
}
