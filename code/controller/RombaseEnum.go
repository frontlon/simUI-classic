package controller

import (
	"encoding/json"
	"fmt"
	"simUI/code/db"
	"simUI/code/modules"
	"simUI/code/utils"
	"simUI/code/utils/go-sciter"
)

/**
 * 定义view用function
 **/

func RomBaseEnumController() {

	//读取rom独立模拟器cmd数据
	utils.Window.DefineFunction("GetByType", func(args ...*sciter.Value) *sciter.Value {
		t := args[0].String()
		volist, _ := (&db.RombaseEnum{}).GetByType(t)
		romJson, _ := json.Marshal(&volist)
		return sciter.NewValue(string(romJson))
	})

	//更新rom独立模拟器参数
	utils.Window.DefineFunction("UpdateByType", func(args ...*sciter.Value) *sciter.Value {
		t := args[0].String()
		data := args[0].String()

		create := []string{}

		if err := json.Unmarshal([]byte(data), &create); err != nil {
			fmt.Println(err.Error())
		}

		err := modules.UpdateRomBaseEnum(t,create)
		if err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}
		return sciter.NullValue()
	})
}
