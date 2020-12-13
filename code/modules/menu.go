package modules

import (
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

func AddMenu(platform uint32,name string){

}

func MenuRename(platform uint32,oldName string,newName string){

}

func DeleteMenu(platform uint32,name string){

}
