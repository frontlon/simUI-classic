package controller

import (
	"simUI/code/modules"
	"simUI/code/utils"
	"simUI/code/utils/go-sciter"
)

/**
 * 定义view用function
 **/

//检查外设
func JoystickController() {
	utils.Window.DefineFunction("checkJoystick", func(args ...*sciter.Value) *sciter.Value {
		status := modules.CheckJoystick()
		return sciter.NewValue(utils.ToString(status))
	})
}
