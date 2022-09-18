package controller

import (
	"simUI/code/modules"
	"simUI/code/utils"
	"simUI/code/utils/go-sciter"
)

/**
 * 定义view用function
 **/

func ImageController() {

	//生成压缩图片
	utils.Window.DefineFunction("CreateOptimizedImage", func(args ...*sciter.Value) *sciter.Value {
		platform := uint32(utils.ToInt(args[0].String()))
		opt := args[1].String()

		go func() {
			err := modules.CreateOptimizedImage(platform, opt)
			if err != nil {
				utils.WriteLog(err.Error())
				utils.ErrorMsg(err.Error())
			}
		}()
		return sciter.NullValue()
	})
}
