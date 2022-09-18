package controller

import (
	"simUI/code/modules"
	"simUI/code/utils"
	"simUI/code/utils/go-sciter"
)

/**
 * 定义view用function
 **/

func BackupController() {

	//备份rom配置
	utils.Window.DefineFunction("BackupRomConfig", func(args ...*sciter.Value) *sciter.Value {
		p := args[0].String()

		if err := modules.BackupRomConfig(p); err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}

		return sciter.NullValue()
	})

	//恢复rom配置
	utils.Window.DefineFunction("RestoreRomConfig", func(args ...*sciter.Value) *sciter.Value {
		p := args[0].String()

		if err := modules.RestoreRomConfig(p); err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}
		return sciter.NullValue()
	})

}
