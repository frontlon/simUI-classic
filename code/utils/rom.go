package utils

import (
	"os"
	"os/exec"
	"path/filepath"
)

var LAST_PROCESS int = 0

/**
 * 运行游戏
 **/
func RunGame(exeFile string, cmd []string) error {

	if exeFile == ""{
		exeFile = "explorer"
	}

	//更改程序运行目录
	if err := os.Chdir(filepath.Dir(exeFile)); err != nil {
		return err
	}

	result := exec.Command(exeFile, cmd...)

	if err := result.Start(); err != nil {
		return err
	}

	//保存进程id
	LAST_PROCESS = result.Process.Pid

	return nil
}

/**
 * 关闭游戏
 **/
func KillGame() error {

	if LAST_PROCESS == 0 {
		return nil
	}

	c := exec.Command("taskkill.exe", "/T", "/PID", ToString(LAST_PROCESS))
	err := c.Start()

	LAST_PROCESS = 0
	return err
}
