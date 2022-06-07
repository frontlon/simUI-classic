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

func RomManagerController() {

	//查询重复rom
	utils.Window.DefineFunction("CheckRomRepeat", func(args ...*sciter.Value) *sciter.Value {
		platformId := uint32(utils.ToInt(args[0].String()))

		//数据库中读取rom详情
		romlist, _ := modules.CheckRomRepeat(platformId)

		romJson, _ := json.Marshal(&romlist)
		return sciter.NewValue(string(romJson))
	})

	//移动rom
	utils.Window.DefineFunction("MoveRomByFileManager", func(args ...*sciter.Value) *sciter.Value {
		romPaths := args[0].String()
		newFolder := args[1].String()
		d := []string{}
		_ = json.Unmarshal([]byte(romPaths), &d)

		for _, p := range d {
			if err := modules.MoveRomByFile(p, newFolder); err != nil {
				utils.WriteLog(err.Error())
				return utils.ErrorMsg(err.Error())
			}
		}

		return sciter.NullValue()
	})

	//查询僵尸资源
	utils.Window.DefineFunction("CheckRomZombie", func(args ...*sciter.Value) *sciter.Value {
		platformId := uint32(utils.ToInt(args[0].String()))

		//数据库中读取rom详情
		romlist, _ := modules.CheckRomZombie(platformId)

		romJson, _ := json.Marshal(&romlist)
		return sciter.NewValue(string(romJson))
	})

	//移动僵尸文件
	utils.Window.DefineFunction("MoveZombieByFileManager", func(args ...*sciter.Value) *sciter.Value {
		romPaths := args[0].String()
		newFolder := args[1].String()
		d := []string{}
		_ = json.Unmarshal([]byte(romPaths), &d)

		for _, p := range d {
			if err := modules.MoveZombieByFile(p, newFolder); err != nil {
				utils.WriteLog(err.Error())
				return utils.ErrorMsg(err.Error())
			}
		}

		return sciter.NullValue()
	})

	//删除僵尸文件
	utils.Window.DefineFunction("DeleteZombieByFileManager", func(args ...*sciter.Value) *sciter.Value {
		path := args[0].String()

		if err := utils.FileDelete(path); err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}

		return sciter.NullValue()
	})

	//绑定子游戏
	utils.Window.DefineFunction("BindSubGame", func(args ...*sciter.Value) *sciter.Value {
		pid := uint64(utils.ToInt(args[0].String()))
		sid := uint64(utils.ToInt(args[1].String()))

		vo, err := modules.BindSubGame(pid, sid)
		if err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}
		romJson, _ := json.Marshal(&vo)
		return sciter.NewValue(string(romJson))
	})

	//解绑子游戏
	utils.Window.DefineFunction("UnBindSubGame", func(args ...*sciter.Value) *sciter.Value {
		id := uint64(utils.ToInt(args[0].String()))

		vo, err := modules.UnBindSubGame(id)
		if err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}
		romJson, _ := json.Marshal(&vo)

		return sciter.NewValue(string(romJson))
	})

}
