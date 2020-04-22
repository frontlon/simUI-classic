package controller

import (
	"VirtualNesGUI/code/config"
	"VirtualNesGUI/code/db"
	"VirtualNesGUI/code/modules"
	"VirtualNesGUI/code/utils"
	"fmt"
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

		var getPlatform uint32 = 0
		if len(args) > 0{
			getPlatform = uint32(utils.ToInt(args[0].String()))
		}

		if len(config.Cfg.Platform) == 0{
			if _, err := w.Call("CB_createCache",sciter.NewValue("")); err != nil {
			}
			return sciter.NullValue()
		}

		go func() *sciter.Value {

			//检查更新一个平台还是所有平台
			PlatformList := map[uint32]*db.Platform{}
			if getPlatform == 0 { //所有平台
				PlatformList = config.Cfg.Platform
			} else { //一个平台
				if _, ok := config.Cfg.Platform[getPlatform]; ok {
					PlatformList[getPlatform] = config.Cfg.Platform[getPlatform]
				}
			}

			//先检查平台，将不存在的平台数据先干掉
			if getPlatform == 0 {
				if err := modules.ClearPlatform(); err != nil {
					WriteLog(err.Error())
					return ErrorMsg(w, err.Error())
				}
			}
			//开始重建缓存
			for platform, _ := range PlatformList {

				romlist, menu, err := modules.CreateRomData(platform)

				if err != nil {
					WriteLog(err.Error())
					return ErrorMsg(w, err.Error())
				}

				//更新rom数据
				if err := modules.UpdateRomDB(platform, romlist); err != nil {
					WriteLog(err.Error())
					return ErrorMsg(w, err.Error())
				}
				//更新menu数据
				if err := modules.UpdateMenuDB(platform, menu); err != nil {
					WriteLog(err.Error())
					return ErrorMsg(w, err.Error())
				}

			}

			//数据更新完成后，页面回调，更新页面DOM
			if _, err := w.Call("CB_createCache",sciter.NewValue("")); err != nil {
				fmt.Print(err)
			}
			return sciter.NullValue()
		}()

		return sciter.NullValue()
	})
}