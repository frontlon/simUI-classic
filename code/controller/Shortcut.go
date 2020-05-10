package controller

import (
	"simUI/code/config"
	"simUI/code/db"
	"simUI/code/utils"
	"encoding/json"
	"simUI/code/utils/go-sciter"
	"os"
)

/**
 * 定义view用function
 **/

func ShortcutController() {

	//读取快捷工具
	utils.Window.DefineFunction("GetShortcut", func(args ...*sciter.Value) *sciter.Value {
		volist, err := (&db.Shortcut{}).GetAll()
		if err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}
		romJson, _ := json.Marshal(&volist)
		return sciter.NewValue(string(romJson))
	})

	//添加快捷工具
	utils.Window.DefineFunction("AddShortcut", func(args ...*sciter.Value) *sciter.Value {

		count,err := (&db.Shortcut{}).Count()
		if err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}
		count++
		shortcut := &db.Shortcut{
			Sort: uint32(count),
		}

		id, err := shortcut.Add();
		if err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}
		return sciter.NewValue(utils.ToString(id))
	})

	//更新快捷工具
	utils.Window.DefineFunction("UpdateShortcut", func(args ...*sciter.Value) *sciter.Value {
		data := args[0].String()
		d := make(map[string]string)
		json.Unmarshal([]byte(data), &d)
		shortcut := &db.Shortcut{
			Id:   uint32(utils.ToInt(d["id"])),
			Name: d["name"],
			Path: d["path"],
		}
		if err := shortcut.UpdateById(); err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}
		return sciter.NullValue()
	})

	//删除快捷工具
	utils.Window.DefineFunction("DelShortcut", func(args ...*sciter.Value) *sciter.Value {
		id := uint32(utils.ToInt(args[0].String()))

		shortcut := &db.Shortcut{
			Id: id,
		}

		if err := shortcut.DeleteById(); err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}
		return sciter.NewValue("1")
	})

	//更新快捷工具排序
	utils.Window.DefineFunction("UpdateShortcutSort", func(args ...*sciter.Value) *sciter.Value {
		data := args[0].String()

		d := make(map[uint32]uint32)
		json.Unmarshal([]byte(data), &d)

		if len(d) == 0 {
			return sciter.NullValue()
		}

		for id, sort := range d {
			shortcut := &db.Shortcut{
				Id:   id,
				Sort: sort,
			}
			if err := shortcut.UpdateSortById(); err != nil {
			}
		}
		return sciter.NewValue("1")
	})


	//运行快捷工具
	utils.Window.DefineFunction("RunShortcut", func(args ...*sciter.Value) *sciter.Value {

		shortcut := args[0].String()

		//检测执行文件是否存在
		_, err := os.Stat(shortcut)
		if err != nil {
			utils.WriteLog(config.Cfg.Lang["ShortcutNotExists"])
			return utils.ErrorMsg( config.Cfg.Lang["ShortcutNotExists"])
		}

		err = utils.RunGame("explorer", []string{shortcut})
		if err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}
		return sciter.NullValue()
	})
}
