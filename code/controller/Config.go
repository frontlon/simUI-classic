package controller

import (
	"encoding/json"
	"simUI/code/config"
	"simUI/code/db"
	"simUI/code/modules"
	"simUI/code/utils"
	"simUI/code/utils/go-sciter"
)

/**
 * 定义view用function
 **/

func ConfigController() {

	//初始化数据
	utils.Window.DefineFunction("InitData", func(args ...*sciter.Value) *sciter.Value {

		ctype := args[0].String()
		isfresh := args[1].String()

		data := ""
		switch (ctype) {
		case "config": //读取配置
			//初始化配置
			if (isfresh == "1") {
				//如果是刷新，则重新生成配置项
				if err := config.InitConf(); err != nil {
					return utils.ErrorMsg(err.Error())
				}
			}
			getjson, _ := json.Marshal(config.Cfg)
			data = string(getjson)
		}
		return sciter.NewValue(data)
	})

	//更新配置文件
	utils.Window.DefineFunction("UpdateConfig", func(args ...*sciter.Value) *sciter.Value {
		field := args[0].String()
		value := args[1].String()

		err := (&db.Config{}).UpdateField(field, value)

		if err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}
		return sciter.NullValue()
	})

	//检查更新
	utils.Window.DefineFunction("CheckUpgrade", func(args ...*sciter.Value) *sciter.Value {
		body := modules.CheckUpgrade()
		return sciter.NewValue(body)
	})

	//导出平台设置
	utils.Window.DefineFunction("InputPlatform", func(args ...*sciter.Value) *sciter.Value {
		p := args[0].String()
		modules.InputPlatform(p)
		return sciter.NewValue(1)
	})

	//导出平台设置
	utils.Window.DefineFunction("OutputPlatform", func(args ...*sciter.Value) *sciter.Value {
		platform := uint32(utils.ToInt(args[0].String()))
		p := args[1].String()
		compress := utils.ToInt(args[2].String())
		packRom := utils.ToInt(args[3].String())



		go func() {
			modules.OutputPlatform(platform, p, compress,packRom)
		}()
		return sciter.NewValue(1)
	})

}
