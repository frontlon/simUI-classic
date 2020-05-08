package utils

import (
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
