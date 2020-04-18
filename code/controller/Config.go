package controller

import (
	"VirtualNesGUI/code/config"
	"VirtualNesGUI/code/db"
	"VirtualNesGUI/code/modules"
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
				if err := config.InitConf(); err != nil {
					return ErrorMsg(w, err.Error())
				}
			}
			getjson, _ := json.Marshal(config.Cfg)
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

	//备份配置文件
	w.DefineFunction("BackupConfig", func(args ...*sciter.Value) *sciter.Value {
		p := args[0].String()

		go func() {
			err := modules.BackupConfig(p)
			if err != nil {
				WriteLog(err.Error())
				ErrorMsg(w, err.Error())
			}
		}()
		return sciter.NullValue()
	})

	//还原配置文件
	w.DefineFunction("RestoreConfig", func(args ...*sciter.Value) *sciter.Value {
		p := args[0].String()

		go func() {
			err := modules.RestoreConfig(p)
			if err != nil {
				WriteLog(err.Error())
				ErrorMsg(w, err.Error())
			}
		}()
		return sciter.NullValue()
	})

	//读取还原数据的统计信息
	w.DefineFunction("GetRestoreInfo", func(args ...*sciter.Value) *sciter.Value {
		p := args[0].String()
		info, err := modules.GetRestoreInfo(p)
		if err != nil {
			WriteLog(err.Error())
			ErrorMsg(w, err.Error())
		}
		jsonInfo, _ := json.Marshal(&info)
		return sciter.NewValue(string(jsonInfo))
	})
}
