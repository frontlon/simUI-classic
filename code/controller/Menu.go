package controller

import (
	"simUI/code/db"
	"simUI/code/modules"
	"simUI/code/utils"
	"encoding/json"
	"simUI/code/utils/go-sciter"
)


/**
 * 定义view用function
 **/

func MenuController() {


	//读取目录列表
	utils.Window.DefineFunction("GetMenuList", func(args ...*sciter.Value) *sciter.Value {
		platform := uint32(utils.ToInt(args[0].String()))
		page := uint32(utils.ToInt(args[1].String()))
		//读取数据库

		menu,err := modules.GetMenuList(uint32(utils.ToInt(platform)),page)
		if err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}

		getjson, _ := json.Marshal(menu)
		return sciter.NewValue(string(getjson))
	})

	//更新菜单排序
	utils.Window.DefineFunction("UpdateMenuSort", func(args ...*sciter.Value) *sciter.Value {
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


	//读取所有平台的菜单列表
	utils.Window.DefineFunction("GetAllPlatformMenuList", func(args ...*sciter.Value) *sciter.Value {

		lists , err := modules.GetAllPlatformMenuList()
		if err != nil {
			utils.WriteLog(err.Error())
		}

		jsonStr, _ := json.Marshal(lists)
		return sciter.NewValue(string(jsonStr))
	})

	//添加菜单
	utils.Window.DefineFunction("AddMenu", func(args ...*sciter.Value) *sciter.Value {

		platform := uint32(utils.ToInt(args[0].String())) //平台
		name := args[1].String()

		modules.AddMenu(platform,name)

		return sciter.NullValue()
	})

	//菜单重命名
	utils.Window.DefineFunction("MenuRename", func(args ...*sciter.Value) *sciter.Value {

		platform := uint32(utils.ToInt(args[0].String())) //平台
		oldName := args[1].String()
		newName := args[2].String()

		modules.MenuRename(platform,oldName,newName)

		return sciter.NullValue()
	})

	//删除菜单
	utils.Window.DefineFunction("DeleteMenu", func(args ...*sciter.Value) *sciter.Value {

		platform := uint32(utils.ToInt(args[0].String())) //平台
		name := args[1].String()

		modules.DeleteMenu(platform,name)

		return sciter.NullValue()
	})


}
