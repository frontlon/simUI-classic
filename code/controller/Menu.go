package controller

import (
	"VirtualNesGUI/code/db"
	"VirtualNesGUI/code/modules"
	"VirtualNesGUI/code/utils"
	"encoding/json"
	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
)


/**
 * 定义view用function
 **/

func MenuController(w *window.Window) {


	//读取目录列表
	w.DefineFunction("GetMenuList", func(args ...*sciter.Value) *sciter.Value {
		platform := uint32(utils.ToInt(args[0].String()))
		//读取数据库



		menu,err := modules.GetMenuList(uint32(utils.ToInt(platform)))
		if err != nil {
			WriteLog(err.Error())
			return ErrorMsg(w, err.Error())
		}

		getjson, _ := json.Marshal(menu)
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
