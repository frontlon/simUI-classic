package utils

import (
	"os"
	"syscall"
	"unsafe"
)

func intPtr(n int) uintptr {
	return uintptr(n)
}
func strPtr(s string) uintptr {
	return uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(s)))
}

//调用windows的alert框
func ShowWindowsAlert(tittle, text string) {
	user32dll, _ := syscall.LoadLibrary("user32.dll")
	user32 := syscall.NewLazyDLL("user32.dll")
	MessageBoxW := user32.NewProc("MessageBoxW")
	MessageBoxW.Call(intPtr(0), strPtr(text), strPtr(tittle), intPtr(0))
	defer syscall.FreeLibrary(user32dll)
}

//弹出系统alert后关闭应用
func ShowAlertAndExit(tittle, text string) {
	ShowWindowsAlert(tittle, text)
	os.Exit(0)
}
