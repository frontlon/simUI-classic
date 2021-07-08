package controller

import (
	"fmt"
	"simUI/code/config"
	"simUI/code/db"
	"simUI/code/modules"
	"simUI/code/utils"
	"simUI/code/utils/go-sciter"
)

/**
 * 定义view用function
 **/

func CacheController() {

	//删除所有缓存
	utils.Window.DefineFunction("TruncateRomCache", func(args ...*sciter.Value) *sciter.Value {

		//清空rom表
		if err := (&db.Rom{}).Truncate(); err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}

		//清空menu表
		if err := (&db.Menu{}).Truncate(); err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}

		return sciter.NullValue()
	})

	//生成所有缓存
	utils.Window.DefineFunction("CreateRomCache", func(args ...*sciter.Value) *sciter.Value {
		var getPlatform uint32 = 0
		if len(args) > 0 {
			getPlatform = uint32(utils.ToInt(args[0].String()))
		}

		//如果没有任何平台，则不用更新
		if len(config.Cfg.Platform) == 0 {
			if _, err := utils.Window.Call("CB_createCache", sciter.NewValue("")); err != nil {
			}
			return sciter.NullValue()
		}
		go func() *sciter.Value {
			if err := modules.CreateRomCache(getPlatform);err != nil{
				utils.WriteLog(err.Error())
				return utils.ErrorMsg(err.Error())
			}
			return sciter.NullValue()
		}()
		return sciter.NullValue()
	})

	//清理csv文件
	utils.Window.DefineFunction("ClearRombase", func(args ...*sciter.Value) *sciter.Value {

		go func() {
			PlatformList := config.Cfg.Platform
			total := 0 //总清理数
			for platform, _ := range PlatformList {
				nameList := []string{}
				count := 0
				romList, _ := (&db.Rom{}).GetByPlatform(platform)
				for _, v := range romList {
					nameList = append(nameList,utils.GetFileName(v.RomPath))
				}
				//读取csv文件数据
				romBase, _ := modules.GetRomBaseList(platform)

				if romBase == nil{
					continue
				}

				for a, _ := range romBase {
					//如果rom列表中无此游戏，则清理
					if !utils.InSliceString(a, nameList) {
						delete(romBase, a)
						count ++
					}

				}
				if count > 0 {
					_ = modules.CoverRomBaseFile(platform, romBase)
					total += count
				}
			}

			//数据更新完成后，页面回调，更新页面DOM
			if _, err := utils.Window.Call("CB_clearRombase", sciter.NewValue(utils.ToString(total))); err != nil {
				fmt.Print(err)
			}
		}();
		return sciter.NullValue()

	})

}
