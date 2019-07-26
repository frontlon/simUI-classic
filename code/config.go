package main

import (
	"github.com/go-ini/ini"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//配置文件
var confFile = "./config.ini"
var Config *ConfStruct
var constPlatformList = []*Platform{}
var constThemeList = map[string]string{}
//读取配置文件
var confSource, _ = ini.LoadSources(ini.LoadOptions{IgnoreInlineComment: true}, confFile)

//配置文件
type ConfStruct struct {
	RootPath string
	Default    *DefaultStruct
	Lang     map[string]string
	Platform map[string]*PfStruct
}

//缓存配置结构体
type DefaultStruct struct {
	Platform      string
	Romlist       string
	ZoomValue     string
	ZoomDirection string
	WindowLeft    string
	WindowTop     string
	WindowWidth   string
	WindowHeight  string
	Theme         string
	Lang          string
}

//平台信息结构体
type PfStruct struct {
	Enable    string
	Title     string
	RomPath   string
	ThumbPath string
	VideoPath string
	DescPath  string
	FileExe   string
	RomExt    []string
}

//前端用平台列表
type Platform struct {
	Name  string
	Value string
}

/*
 初始化读取配置
 @author frontLon
 @return strucct
*/
func InitConf() {
	Config = &ConfStruct{}
	//配置全局参数
	Config.Platform = getSectionPlatform()
	Config.Default = getSectionDefault()
	Config.Lang = getLang(Config.Default.Lang)
	Config.RootPath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
}

//读取平台列表
func getSectionPlatform() map[string]*PfStruct {
	platform := make(map[string]*PfStruct)

	//读取section
	var section = confSource.Section("platform")
	plat := &Platform{}
	for _, v := range (section.ChildSections()) {

		//如果该section禁用，则不载入
		if v.Key("Enable").String() != "1" {
			continue
		}

		name := strings.Split(v.Name(), ".")
		//拼装缩略图绝对路径，解决找不到路径的问题
		thumb, _ := filepath.Abs(v.Key("ThumbPath").String())
		video, _ := filepath.Abs(v.Key("VideoPath").String())
		rom, _ := filepath.Abs(v.Key("RomPath").String())

		//定义平台config
		platform[name[1]] = &PfStruct{
			Enable:    v.Key("Enable").String(),
			Title:     v.Key("Title").String(),
			FileExe:   v.Key("FileExe").String(),
			RomPath:   rom,
			ThumbPath: thumb,
			VideoPath: video,
			DescPath:  v.Key("DescPath").String(),
			RomExt:    strings.Split(v.Key("RomExt").String(), ","), //拆分rom扩展名
		}

		//定义平台常量
		plat = &Platform{
			Name:  v.Key("Title").String(),
			Value: name[1],
		}
		constPlatformList = append(constPlatformList, plat)
	}
	return platform
}

//读取缓存配置
func getSectionDefault() *DefaultStruct {
	section := confSource.Section("Default")
	platform := section.Key("Platform").String()
	romlist := section.Key("Romlist").String()
	zoomvalue := section.Key("ZoomValue").String()
	zoomdirection := section.Key("ZoomDirection").String()
	windowleft := section.Key("WindowLeft").String()
	windowtop := section.Key("WindowTop").String()
	windowwidth := section.Key("WindowWidth").String()
	windowheight := section.Key("WindowHeight").String()
	theme := section.Key("Theme").String()
	lang := section.Key("Lang").String()
	isset := false

	if romlist == "" {
		romlist = "1"
	}

	if zoomvalue == "" {
		zoomvalue = "1"
	}

	if zoomdirection == "" {
		zoomdirection = "1"
	}

	if theme == "" {
		theme = "dark"
	}

	//查看当前选定平台值是否是正常的
	for k, _ := range (Config.Platform) {
		if platform == k {
			isset = true
			break
		}
	}
	//如果没有匹配上platform，则读取config中的第一项
	if isset == false {
		for k, _ := range (Config.Platform) {
			platform = k
			//修复配置文件
			if err := updateConfig("Default", "Platform", platform); err != nil {
			}
			break
		}
	}
	return &DefaultStruct{
		Platform:      platform,
		Romlist:       romlist,
		ZoomValue:     zoomvalue,
		ZoomDirection: zoomdirection,
		WindowLeft:    windowleft,
		WindowTop:     windowtop,
		WindowWidth:   windowwidth,
		WindowHeight:  windowheight,
		Theme:         theme,
		Lang:          lang,
	}
}

//更新配置文件
func updateConfig(section string, field string, value string) error {
	confSource.Section(section).Key(field).SetValue(value)
	err := confSource.SaveTo(confFile);
	return err
}

//读取主题列表
func getThemeList() () {
	dirPth, _ := filepath.Abs("theme")
	lists, _ := ioutil.ReadDir(dirPth)
	for _, fi := range lists {
		if fi.IsDir() { // 忽略目录
			constThemeList[fi.Name()] = dirPth + separator + fi.Name() + separator + "theme.ini"
		}
	}
}

//读取主题配置参数
func getThemeParams(title string) (map[string]string, error) {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	dirPth := dir + separator + "theme" + separator + title + separator + "theme.ini"
	file, err := ini.LoadSources(ini.LoadOptions{IgnoreInlineComment: true}, dirPth)

	section := file.Section("").KeysHash()
	return section, err
}

//读取ROM别名配置参数
func getRomAlias(platform string) map[string]string {
	dirPth := Config.Platform[platform].RomPath + separator + "romlist.ini"
	section := make(map[string]string)
	file, err := ini.LoadSources(ini.LoadOptions{IgnoreInlineComment: true}, dirPth)
	if err == nil {
		section = file.Section("Alias").KeysHash()
	}

	return section
}

//读取语言参数配置
func getLang(lang string) map[string]string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	dirPth := dir + separator + "lang" + separator + lang + ".ini"
	section := make(map[string]string)
	file, err := ini.LoadSources(ini.LoadOptions{IgnoreInlineComment: true}, dirPth)
	if err == nil {
		section = file.Section("").KeysHash()
	}
	return section
}
