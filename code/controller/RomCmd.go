package controller

import (
	"VirtualNesGUI/code/db"
	"VirtualNesGUI/code/modules"
	"VirtualNesGUI/code/utils"
	"encoding/json"
	"github.com/sciter-sdk/go-sciter"
)

/**
 * 定义view用function
 **/

func RomCmdController() {

	//读取rom独立模拟器cmd数据
	utils.Window.DefineFunction("GetRomCmd", func(args ...*sciter.Value) *sciter.Value {
		romId := uint64(utils.ToInt(args[0].String()))
		simId := uint32(utils.ToInt(args[1].String()))

		//数据库中读取rom详情
		rom, _ := (&db.Rom{}).GetSimConf(romId, simId)

		romJson, _ := json.Marshal(&rom)
		return sciter.NewValue(string(romJson))
	})

	//更新rom独立模拟器参数
	utils.Window.DefineFunction("UpdateRomCmd", func(args ...*sciter.Value) *sciter.Value {
		id := uint64(utils.ToInt(args[0].String()))
		simId := uint32(utils.ToInt(args[1].String()))
		data := args[2].String()
		d := make(map[string]string)
		json.Unmarshal([]byte(data), &d)

		err := modules.UpdateRomCmd(id, simId, d)
		if err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}

		return sciter.NullValue()
	})
}
