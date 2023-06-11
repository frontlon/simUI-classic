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

func BackupController() {

	//备份rom配置
	utils.Window.DefineFunction("BackupRomConfig", func(args ...*sciter.Value) *sciter.Value {
		p := args[0].String()
		platform := uint32(utils.ToInt(args[1].String()))
		if err := modules.BackupRomConfig(p, platform); err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}

		return sciter.NullValue()
	})

	//恢复rom配置
	utils.Window.DefineFunction("RestoreRomConfig", func(args ...*sciter.Value) *sciter.Value {
		p := args[0].String()
		platform := uint32(utils.ToInt(args[1].String()))
		if err := modules.RestoreRomConfig(p, platform); err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}
		return sciter.NullValue()
	})

	//合并DB - 读取检测数据
	utils.Window.DefineFunction("GetMergeDbData", func(args ...*sciter.Value) *sciter.Value {
		p := args[0].String()
		volist, err := modules.GetMergeDbData(p)
		if err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}
		jsonInfo, _ := json.Marshal(&volist)
		return sciter.NewValue(string(jsonInfo))
	})

	//合并DB - 合并数据
	utils.Window.DefineFunction("MergeDB", func(args ...*sciter.Value) *sciter.Value {
		p := args[0].String()
		platformIds := []uint32{}
		platformJson := []string{}
		simulatorIds := []string{}
		json.Unmarshal([]byte(args[1].String()), &platformJson)
		json.Unmarshal([]byte(args[2].String()), &simulatorIds)
		for _, v := range platformJson {
			platformIds = append(platformIds, uint32(utils.ToInt(v)))
		}
		go func() {
			if err := modules.MergeDB(p, platformIds, simulatorIds); err != nil {
				utils.WriteLog(err.Error())
				utils.ErrorMsg(err.Error())
			}
		}()
		return sciter.NullValue()
	})

}
