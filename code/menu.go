package main

import (
	"io/ioutil"
	"strings"
)

var constMenuList = &Menu{}

type Menu struct{
	Fc []*MenuInfo
	Sfc []*MenuInfo
	Md []*MenuInfo
	Pce []*MenuInfo
	Gb []*MenuInfo
	Arcade []*MenuInfo
}

//菜单信息
type MenuInfo struct {
	Name  string //菜单名称
}

/**
 * 根据游戏平台，读取相对应的菜单列表
 **/

func GetMenuData(platform string) error {

	platform = strings.ToLower(platform)
	romPath := ""
	switch platform {
	case "fc":
		romPath = Config.Fc.RomPath
	case "sfc":
		romPath = Config.Sfc.RomPath
	case "md":
		romPath = Config.Md.RomPath
	case "pce":
		romPath = Config.Pce.RomPath
	case "gb":
		romPath = Config.Gb.RomPath
	case "arcade":
		romPath = Config.Arcade.RomPath
	}

	dir_list, e := ioutil.ReadDir(romPath)
	if e != nil {
		return e
	}

	menulist := []*MenuInfo{}

	for _, v := range dir_list {
		if v.IsDir() == true {
			des := &MenuInfo{
				Name:   v.Name(),
			}
			menulist = append(menulist, des)
		}
	}
	switch platform {
	case "fc":
		constMenuList.Fc = menulist
	case "sfc":
		constMenuList.Sfc = menulist
	case "md":
		constMenuList.Md = menulist
	case "pce":
		constMenuList.Pce = menulist
	case "gb":
		constMenuList.Gb = menulist
	case "arcade":
		constMenuList.Arcade = menulist
	}

	return e
}