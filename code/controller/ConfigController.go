package controller

import (
	"VirtualNesGUI/code/db"
	"encoding/json"
	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
)



/**
 * 定义view用function
 **/

func ConfigController(w *window.Window) {
	w.DefineFunction("InitData", func(args ...*sciter.Value) *sciter.Value {

		ctype := args[0].String()
		isfresh := args[1].String()

		data := ""
		switch (ctype) {
		case "config": //读取配置
			//初始化配置
			if (isfresh == "1") {
				//如果是刷新，则重新生成配置项
				if err := InitConf(); err != nil {
					return ErrorMsg(w, err.Error())
				}
			}
			getjson, _ := json.Marshal(Config)
			data = string(getjson)
		}
		return sciter.NewValue(data)
	})

	//更新配置文件
	w.DefineFunction("UpdateConfig", func(args ...*sciter.Value) *sciter.Value {
		field := args[0].String()
		value := args[1].String()

		err := (&db.Config{}).UpdateField(field, value)



		if err != nil {
			WriteLog(err.Error())
			return ErrorMsg(w, err.Error())
		}
		return sciter.NullValue()
	})
}
