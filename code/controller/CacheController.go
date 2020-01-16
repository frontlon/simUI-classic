package controller

import (
	"VirtualNesGUI/code/db"
	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
)

/**
 * 定义view用function
 **/

func CacheController(w *window.Window) {

	//删除所有缓存
	w.DefineFunction("TruncateRomCache", func(args ...*sciter.Value) *sciter.Value {

		//清空rom表
		if err := (&db.Rom{}).Truncate(); err != nil {
			WriteLog(err.Error())
			return ErrorMsg(w, err.Error())
		}

		//清空menu表
		if err := (&db.Menu{}).Truncate(); err != nil {
			WriteLog(err.Error())
			return ErrorMsg(w, err.Error())
		}

		return sciter.NullValue()
	})

	//生成所有缓存
	w.DefineFunction("CreateRomCache", func(args ...*sciter.Value) *sciter.Value {

		go func() *sciter.Value {

			//先检查平台，将不存在的平台数据先干掉
			if err := ClearPlatform(); err != nil {
				WriteLog(err.Error())
				return ErrorMsg(w, err.Error())
			}

			//开始重建缓存
			for platform, _ := range Config.Platform {

				//创建rom数据
				romlist, menu, err := CreateRomCache(platform)
				if err != nil {
					WriteLog(err.Error())
					return ErrorMsg(w, err.Error())
				}

				//更新rom数据
				if err := UpdateRomDB(platform, romlist); err != nil {
					WriteLog(err.Error())
					return ErrorMsg(w, err.Error())
				}

				//更新menu数据
				if err := UpdateMenuDB(platform, menu); err != nil {
					WriteLog(err.Error())
					return ErrorMsg(w, err.Error())
				}

			}

			//数据更新完成后，页面回调，更新页面DOM
			if _, err := w.Call("createAllCache"); err != nil {
			}
			return sciter.NullValue()
		}()

		return sciter.NullValue()
	})
}
