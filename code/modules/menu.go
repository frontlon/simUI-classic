package modules

import (
	"simUI/code/db"
)

//读取菜单列表
func GetMenuList(platform uint32) ([]*db.Menu, error) {
	newMenu := []*db.Menu{}

	menu, err := (&db.Menu{}).GetByPlatform(platform) //从数据库中读取当前平台的分类目录
	if err != nil {
		return newMenu, err
	}
	//读取根目录下是否有rom
	count, err := (&db.Rom{}).Count(platform, ConstMenuRootKey, "")
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

//更新菜单排序
func UpdateMenuSort() {

}
