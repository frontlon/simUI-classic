package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

/**
 * 软件启动器
 * 模拟双击软件效果，调用系统关联来启动应用
 * 实现类似cmd的 start 1.exe 效果
 */
func main() {

	if len(os.Args) < 2 {
		return
	}

	boot := "open"
	if runtime.GOOS == "windows" {
		boot = "explorer"
	}

	if err := os.Chdir(filepath.Dir(os.Args[1])); err != nil {
		return
	}
	result := exec.Command(boot, os.Args[1])
	if err := result.Start(); err != nil {
		return
	}
}
