package controller

import (
	"VirtualNesGUI/code/db"
	"VirtualNesGUI/code/utils"
	"encoding/json"
	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
)

/**
 * 定义view用function
 **/

func RomCmdController(w *window.Window) {

	//读取rom独立模拟器cmd数据
	w.DefineFunction("GetRomCmd", func(args ...*sciter.Value) *sciter.Value {
		romId := uint64(utils.ToInt(args[0].String()))
		simId := uint32(utils.ToInt(args[1].String()))

		//数据库中读取rom详情
		rom, _ := (&db.Rom{}).GetSimConf(romId, simId)

		romJson, _ := json.Marshal(&rom)
		return sciter.NewValue(string(romJson))
	})

	//更新rom独立模拟器参数
	w.DefineFunction("UpdateRomCmd", func(args ...*sciter.Value) *sciter.Value {
		id := uint64(utils.ToInt(args[0].String()))
		simId := uint32(utils.ToInt(args[1].String()))
		data := args[2].String()
		d := make(map[string]string)
		json.Unmarshal([]byte(data), &d)

		//如果当前配置和模拟器默认配置一样，则删除该记录
		if d["cmd"] == "" && d["unzip"] == "2" {
			if err := (&db.Rom{}).DelSimConf(id, simId); err != nil {
				WriteLog(err.Error())
				return ErrorMsg(w, err.Error())
			}
			return sciter.NullValue()
		}

		//开始更新
		if err := (&db.Rom{}).UpdateSimConf(id, simId, d["cmd"], uint8(utils.ToInt(d["unzip"]))); err != nil {
			WriteLog(err.Error())
			return ErrorMsg(w, err.Error())
		}

		return sciter.NullValue()
	})
}
