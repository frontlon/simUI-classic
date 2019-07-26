package main

import (
	"io/ioutil"
)

var constMenuList = &Menu{}
var constMenuRootKey = "_7b9" //根子目录游戏的Menu参数

type Menu struct {
	Platform map[string][]*MenuInfo
}

//菜单信息
type MenuInfo struct {
	Name string //菜单名称
}

/**
 * 根据游戏平台，读取相对应的菜单列表
 **/

func GetMenuData(platform string) error {

	dir_list, e := ioutil.ReadDir(Config.Platform[platform].RomPath)
	if e != nil {
		return e
	}

	item := &MenuInfo{}
	menulist := []*MenuInfo{}
	for _, v := range dir_list {
		if v.IsDir() == true {
			item = &MenuInfo{
				Name: v.Name(),
			}
			menulist = append(menulist, item)
		}
	}
	constMenuList.Platform[platform] = menulist
	return e
}
