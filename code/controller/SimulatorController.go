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

func SimulatorController(w *window.Window) {

	//添加模拟器
	w.DefineFunction("AddSimulator", func(args ...*sciter.Value) *sciter.Value {
		data := args[0].String()
		d := make(map[string]string)
		json.Unmarshal([]byte(data), &d)
		pfId := uint32(utils.ToInt(d["platform"]))

		sim := &db.Simulator{
			Name:     d["name"],
			Platform: pfId,
			Path:     d["path"],
			Cmd:      d["cmd"],
			Pinyin:   TextToPinyin(d["name"]),
		}
		id, err := sim.Add()

		//更新默认模拟器
		if utils.ToInt(d["default"]) == 1 {
			err = sim.UpdateDefault(pfId, id)
			if err != nil {
				WriteLog(err.Error())
				return ErrorMsg(w, err.Error())
			}
		}
		sim.Id = id
		jsonData, _ := json.Marshal(&sim)
		return sciter.NewValue(string(jsonData))
	})

	//更新模拟器
	w.DefineFunction("UpdateSimulator", func(args ...*sciter.Value) *sciter.Value {
		data := args[0].String()
		d := make(map[string]string)
		json.Unmarshal([]byte(data), &d)

		id := uint32(utils.ToInt(d["id"]))
		pfId := uint32(utils.ToInt(d["platform"]))
		def := uint8(utils.ToInt(d["default"]))
		sim := &db.Simulator{
			Id:       id,
			Name:     d["name"],
			Platform: pfId,
			Path:     d["path"],
			Cmd:      d["cmd"],
			Pinyin:   TextToPinyin(d["name"]),
		}

		//更新模拟器
		if err := sim.UpdateById(); err != nil {
			WriteLog(err.Error())
			return ErrorMsg(w, err.Error())
		}

		//如果设置了默认模拟器，则更新默认模拟器
		if def == 1 {
			if err := sim.UpdateDefault(pfId, id); err != nil {
				WriteLog(err.Error())
				return ErrorMsg(w, err.Error())
			}
		}

		jsonData, _ := json.Marshal(&sim)
		return sciter.NewValue(string(jsonData))
	})

	//删除一个模拟器
	w.DefineFunction("DelSimulator", func(args ...*sciter.Value) *sciter.Value {
		id := uint32(utils.ToInt(args[0].String()))

		sim := &db.Simulator{
			Id: id,
		}

		//删除模拟器
		err := sim.DeleteById()
		if err != nil {
			WriteLog(err.Error())
			return ErrorMsg(w, err.Error())
		}

		//删除rom独立模拟器cmd配置
		if err = (&db.RomCmd{SimId: id}).ClearBySimId(); err != nil {
			WriteLog(err.Error())
			return ErrorMsg(w, err.Error())
		}
		return sciter.NewValue(args[0].String())
	})



	//读取模拟器详情
	w.DefineFunction("GetSimulatorById", func(args ...*sciter.Value) *sciter.Value {
		id := uint32(utils.ToInt(args[0].String()))

		//游戏游戏详细数据
		info, err := (&db.Simulator{}).GetById(id)
		if err != nil {
			WriteLog(err.Error())
			return ErrorMsg(w, err.Error())
		}
		jsonInfo, _ := json.Marshal(&info)
		return sciter.NewValue(string(jsonInfo))
	})

	//读取一个平台下的所有模拟器
	w.DefineFunction("GetSimulatorByPlatform", func(args ...*sciter.Value) *sciter.Value {
		id := uint32(utils.ToInt(args[0].String()))

		//游戏游戏详细数据
		info, err := (&db.Simulator{}).GetByPlatform(id)
		if err != nil {
			WriteLog(err.Error())
			return ErrorMsg(w, err.Error())
		}
		jsonInfo, _ := json.Marshal(&info)
		return sciter.NewValue(string(jsonInfo))
	})


	//设置rom的模拟器
	w.DefineFunction("SetRomSimulator", func(args ...*sciter.Value) *sciter.Value {
		romId := uint64(utils.ToInt(args[0].String()))
		simId := uint32(utils.ToInt(args[1].String()))
		//更新rom表
		rom := &db.Rom{
			Id:    romId,
			SimId: simId,
		}

		//更新数据
		if err := rom.UpdateSimulator(); err != nil {
			WriteLog(err.Error())
			return ErrorMsg(w, err.Error())
		}
		return sciter.NewValue("1")
	})
}
