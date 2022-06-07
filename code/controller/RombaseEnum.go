package controller

import (
	"encoding/json"
	"simUI/code/db"
	"simUI/code/modules"
	"simUI/code/utils"
	"simUI/code/utils/go-sciter"
	"strings"
)

/**
 * 定义view用function
 **/

func RomBaseEnumController() {

	//读取rom独立模拟器cmd数据
	utils.Window.DefineFunction("GetRombaseEnumByType", func(args ...*sciter.Value) *sciter.Value {
		t := args[0].String()
		volist, _ := (&db.RombaseEnum{}).GetByType(t)
		result := ""
		for _, v := range volist {
			result += v.Name + "\n"
		}
		return sciter.NewValue(result)
	})

	//根据枚举数据列表
	utils.Window.DefineFunction("GetRombaseEnumList", func(args ...*sciter.Value) *sciter.Value {
		t := args[0].String()
		volist, _ := (&db.RombaseEnum{}).GetByType(t)
		result := []string{}
		for _, v := range volist {
			result  = append(result,v.Name)
		}
		jsonInfo, _ := json.Marshal(result)
		return sciter.NewValue(string(jsonInfo))
	})

	//更新数据
	utils.Window.DefineFunction("UpdateRomBaseEnumByType", func(args ...*sciter.Value) *sciter.Value {
		t := args[0].String()
		data := strings.Trim(args[1].String(), "")

		create := []string{}
		if data != "" {
			create = strings.Split(data, "\n")
		}

		err := modules.UpdateRomBaseEnum(t, create)
		if err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}
		return sciter.NullValue()
	})
}
