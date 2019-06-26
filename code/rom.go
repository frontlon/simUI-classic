package main

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

var constRomList = &Rom{}

type Rom struct{
	Fc []*Rominfo
	Sfc []*Rominfo
	Md []*Rominfo
	Pce []*Rominfo
	Gb []*Rominfo
	Arcade []*Rominfo
}

//游戏信息
type Rominfo struct {
	Title string //标题
	Menu  string //目录
	Path  string //完整路径
	Thumb string //缩略图
}



//读取fc rom列表
func getRomList(platform string)  {
	platform = strings.ToLower(platform)

	//读取平台配置
	conf := &PfStruct{}
	switch platform {
	case "fc":
		conf = Config.Fc
	case "sfc":
		conf = Config.Sfc
	case "md":
		conf = Config.Md
	case "pce":
		conf = Config.Pce
	case "gb":
		conf = Config.Gb
	case "arcade":
		conf = Config.Arcade
	}

	//路径最后加入反斜杠
	conf.RomPath = conf.RomPath + "\\"

	romlist := []*Rominfo{}

	//进入循环，遍历文件
	if err:= filepath.Walk(conf.RomPath,
		func(p string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}

			//整理目录格式，并转换为数组
			newpath := strings.Replace(conf.RomPath, "/", "\\", -1)
			if newpath[0 : 2] == ".\\"{
				p = ".\\" + p
			}
			newpath = strings.Replace(p, newpath, "", -1)
			subpath := strings.Split(newpath, "\\")
			fileExt := strings.ToLower(path.Ext(p)) //获取文件后缀

			//如果该文件是游戏rom
			if f.IsDir() == false && CheckPlatformExt(conf.FileExt,fileExt){
				menu := constMenuRootKey //无目录，读取默认参数
				//定义目录，如果有子目录，则记录子目录名称
				if len(subpath) > 1 {
					menu = subpath[0]
				}

				//去掉扩展名，生成标题
				rinfo := &Rominfo{
					Title: GetRomName(f.Name()),
					Path:  p,
					Menu:  menu,
					Thumb: getThumb(conf.ThumbPath,GetRomName(f.Name())),
				}
				romlist = append(romlist, rinfo)
			}
			return nil
		});err != nil{}

	//赋值给全局变量
	switch platform {
	case "fc":
		constRomList.Fc = romlist
	case "sfc":
		constRomList.Sfc = romlist
	case "md":
		constRomList.Md = romlist
	case "pce":
		constRomList.Pce = romlist
	case "gb":
		constRomList.Gb = romlist
	case "arcade":
		constRomList.Arcade = romlist
	}

}

/**
 * 读取游戏缩略图
 **/
func  getThumb(thumbPath string,title string) string {
	img := thumbPath + "\\" +title+".png"
	if !exists(img){
		img = thumbPath + "\\..\\_DEF_.png"
	}
	return img
}

/**
 * 运行游戏
 **/
func runGame(platform string,path string) string {

	exeFile := ""
	switch(platform){
	case "fc":
		exeFile = Config.Fc.FileExe
	case "sfc":
		exeFile = Config.Sfc.FileExe
	case "md":
		exeFile = Config.Md.FileExe
	case "pce":
		exeFile = Config.Pce.FileExe
	case "gb":
		exeFile = Config.Gb.FileExe
	case "arcade":
		exeFile = Config.Arcade.FileExe
	}

	//检测执行文件是否存在
	_, err := os.Stat(exeFile)
	if err != nil {
		return "执行程序" + exeFile + "不存在"
	}

	//检测rom文件是否存在
	if exists(path) == false {
		return "rom文件:" + path + "不存在"
	}

	cmd := exec.Command(exeFile,path)
	if err := cmd.Start(); err != nil {
		return "程序启动失败:" + err.Error()
	}
	return ""
}

/**
 * 检测文件是否存在（文件夹也返回false）
 **/
func exists(path string) bool {
	finfo, err := os.Stat(path)
	isset := false
	if err != nil || finfo.IsDir() == true{
		isset =  false
	}else{
		isset = true
	}
	return isset
}

/**
 * 去掉rom扩展名，从文件名中读取Rom名称
 **/
func GetRomName(filename string) string{
	return strings.TrimSuffix(filename, path.Ext(filename))
}

//检查文件扩展名是否存在于该平台中
func CheckPlatformExt(exts []string,file string) bool{
	for _,v := range exts{
		if v == file{
			return true
		}
	}
	return false
}