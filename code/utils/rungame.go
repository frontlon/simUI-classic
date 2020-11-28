package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var LAST_PROCESS int = 0

/**
 * 运行游戏
 **/
func RunGame(exeFile string, cmd []string) error {

	switch runtime.GOOS {
	case "darwin":
		runGameDarwin(exeFile, cmd)
	case "windows":
		runGameWindows(exeFile, cmd)
	case "linux":
		//ROOTPATH = ""
	}

	return nil
}

//windows
func runGameWindows(exeFile string, cmd []string) error {
	if exeFile == "" {
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

//macos
func runGameDarwin(exeFile string, cmd []string) error {
	if exeFile == "" {
		exeFile = "open"
	} else {
		if IsDir(exeFile) {
			exeFile += "/Contents/MacOS/" + getDarwinAppName(exeFile)
		}
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

//在窗口中打开文件夹
func OpenFolderByWindow(fileName string) error {

	isDir := IsDir(fileName)
	switch runtime.GOOS {

	case "darwin":
		if isDir == true {
			if err := exec.Command(`open`, fileName).Start(); err != nil {
				return err
			}
		} else {
			if err := exec.Command(`open`, `-R`, fileName).Start(); err != nil {
				return err
			}
		}
	case "windows":
		if isDir == true {
			if err := exec.Command(`explorer`, fileName).Start(); err != nil {
				return err
			}
		} else {
			if err := exec.Command(`explorer`, `/select,`, `/n,`, fileName).Start(); err != nil {
				return err
			}
		}

	case "linux":
	}

	return nil


}

/**
 * 关闭游戏
 **/
func KillGame() error {

	if LAST_PROCESS == 0 {
		return nil
	}

	switch runtime.GOOS {

	case "darwin":
		//
	case "windows":
		c := exec.Command("taskkill.exe", "/T", "/PID", ToString(LAST_PROCESS))
		c.Start()
	case "linux":
	}

	LAST_PROCESS = 0
	return nil
}

//从info.plist中读取应用程序名称
func getDarwinAppName(p string) string {

	fi, err := os.Open(p + "/Contents/Info.plist")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return ""
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	isset := false
	for {
		a, _, c := br.ReadLine()
		str := string(a)
		if isset == true {
			str = strings.Replace(str, " ", "", -1)
			str = strings.Replace(str, "\t", "", -1)
			str = strings.Replace(str, "<string>", "", -1)
			str = strings.Replace(str, "</string>", "", -1)
			return str
		}
		key := strings.Index(str, "CFBundleExecutable")
		if key > -1 {
			isset = true
		}
		if c == io.EOF {
			break
		}
	}
	return ""
}
