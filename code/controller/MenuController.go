package controller

import (
	"VirtualNesGUI/code/db"
	"VirtualNesGUI/code/utils"
	"encoding/json"
	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
)

var constMenuRootKey = "_7b9"                                                //根子目录游戏的Menu参数

/**
 * 定义view用function
 **/

func MenuController(w *window.Window) {


	//读取目录列表
	w.DefineFunction("GetMenuList", func(args ...*sciter.Value) *sciter.Value {
		platform := uint32(utils.ToInt(args[0].String()))
		//读取数据库
		menu, err := (&db.Menu{}).GetByPlatform(platform) //从数据库中读取当前平台的分类目录
		if err != nil {
			WriteLog(err.Error())
			return ErrorMsg(w, err.Error())
		}
		//读取根目录下是否有rom
		count, err := (&db.Rom{}).Count(platform, constMenuRootKey, "")
		if err != nil {
			WriteLog(err.Error())
			return ErrorMsg(w, err.Error())
		}
		newMenu := []*db.Menu{}

		//读取根目录下有rom，则显示未分类文件夹
		if count > 0 {
			root := &db.Menu{
				Name:     constMenuRootKey,
				Platform: platform,
			}
			newMenu = append(newMenu, root)
			newMenu = append(newMenu, menu...)
		} else {
			newMenu = menu
		}

		if err != nil {
			WriteLog(err.Error())
			return ErrorMsg(w, err.Error())
		}



		getjson, _ := json.Marshal(newMenu)
		return sciter.NewValue(string(getjson))
	})

	//更新菜单排序
	w.DefineFunction("UpdateMenuSort", func(args ...*sciter.Value) *sciter.Value {
		platform := uint32(utils.ToInt(args[0].String())) //平台
		data := args[1].String()

		d := make(map[string]uint32)
		json.Unmarshal([]byte(data), &d)

		if len(d) == 0 {
			return sciter.NullValue()
		}

		for name, val := range d {
			if name == "" {
				continue
			}
			menu := &db.Menu{
				Name:     name,
				Platform: platform,
				Sort:     val,
			}
			if err := menu.UpdateSortByName(); err != nil {
			}
		}
		return sciter.NewValue("1")
	})

}
