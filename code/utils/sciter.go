package utils

import (
	"fmt"
	"simUI/code/utils/go-sciter"
	"simUI/code/utils/go-sciter/window"
)

var Window *window.Window //窗体

//调用alert框
func ErrorMsg(err string) *sciter.Value {

	if _, err := Window.Call("errorBox", sciter.NewValue(err)); err != nil {
	}
	return sciter.NullValue()
}

//调用loading框
func Loading(str string, platform string) *sciter.Value {
	if _, err := Window.Call("startLoading", sciter.NewValue(str), sciter.NewValue(platform)); err != nil {
	}
	return sciter.NullValue()
}

//检查当前窗口激活状态
func CheckWinActive() bool {
	active, err := Window.Call("checkWinActive")
	if err != nil{
		fmt.Println(err)
	}
	return active.Bool()
}

//调用视图中的方向控制【手柄】
func ViewDirection(dir int) bool {
	active, err := Window.Call("joystickDirection", sciter.NewValue(dir))
	if err != nil{
		fmt.Println(err)
	}
	return active.Bool()
}

//调用视图中的按钮控制【手柄】
func ViewButton(btn int) bool {
	active, err := Window.Call("joystickButton", sciter.NewValue(btn))
	if err != nil{
		fmt.Println(err)
	}
	return active.Bool()
}
