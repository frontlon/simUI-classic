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
					utils.WriteLog(err.Error())
					return utils.ErrorMsg(err.Error())
				}
			}

			//开始重建缓存
			for platform, _ := range PlatformList {

				//第一步：创建rom缓存
				romlist, err := modules.CreateRomData(platform)
				if err != nil {
					utils.WriteLog(err.Error())
					return utils.ErrorMsg(err.Error())
				}

				//第二步：更新rom数据
				if err := modules.UpdateRomDB(platform, romlist); err != nil {
					utils.WriteLog(err.Error())
					return utils.ErrorMsg(err.Error())
				}

				//第三步：读取rom目录
				menu ,err := modules.CreateMenuList(platform)
				if err != nil {
					utils.WriteLog(err.Error())
					return utils.ErrorMsg(err.Error())
				}

				//第四步：更新menu数据
				if err := modules.UpdateMenuDB(platform, menu); err != nil {
					utils.WriteLog(err.Error())
					return utils.ErrorMsg(err.Error())
				}

				//第五步：更新filter数据
				modules.UpdateFilterDB(platform);

			}

			//收缩数据库
			db.Vacuum()

			//数据更新完成后，页面回调，更新页面DOM
			if _, err := utils.Window.Call("CB_createCache", sciter.NewValue("")); err != nil {
				fmt.Print(err)
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
				romBase, _ := modules.GetRomBase(platform)

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
