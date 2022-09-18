package controller

import (
	"encoding/json"
	"simUI/code/modules"
	"simUI/code/utils"
	"simUI/code/utils/go-sciter"
)

/**
 * 定义view用function
 **/

func StrategyFilesController() {

	//读取攻略文件
	utils.Window.DefineFunction("GetStrategyFile", func(args ...*sciter.Value) *sciter.Value {
		id := uint64(utils.ToInt(args[0].String()))
		volist,err := modules.GetStrategyFile(id)
		if err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}
		jsonInfo, _ := json.Marshal(volist)
		return sciter.NewValue(string(jsonInfo))
	})

	//上传文件
	utils.Window.DefineFunction("UploadStrategyFile", func(args ...*sciter.Value) *sciter.Value {
		id := uint64(utils.ToInt(args[0].String()))
		name := args[1].String()
		p := args[2].String()

		relPath, err := modules.UploadStrategyFile(id, name, p)
		if err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}

		return sciter.NewValue(relPath)
	})

	//更新配置
	utils.Window.DefineFunction("UpdateStrategyFiles", func(args ...*sciter.Value) *sciter.Value {
		id := uint64(utils.ToInt(args[0].String()))
		data := args[1].String()
		if err := modules.UpdateStrategyFiles(id, data); err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}
		return sciter.NullValue()
	})

	//打开攻略文件
	utils.Window.DefineFunction("OpenStrategyFiles", func(args ...*sciter.Value) *sciter.Value {
		p := args[0].String()
		if err := modules.OpenStrategyFiles(p); err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}
		return sciter.NullValue()
	})

}
