package main

import (
	"github.com/go-ini/ini"
	"os"
	"path/filepath"
	"strings"
)

// 配置文件
//var confFile = "D:\\work\\go\\src\\VirtualNesGUI\\app\\config.ini"
var confFile = "./config.ini"

//读取配置文件
var confSource, _ = ini.LoadSources(ini.LoadOptions{IgnoreInlineComment: true}, confFile)

var Config *ConfStruct

type ConfStruct struct {
	General *GeneralStruct
	Fc *PfStruct
	Sfc *PfStruct
	Md *PfStruct
	Pce *PfStruct
	Gb *PfStruct
	Arcade *PfStruct
}

//通用信息结构体
type GeneralStruct struct {
	Title  string
	Platform  string
}

//平台信息结构体
type PfStruct struct {
	Enable    string
	Title     string
	RomPath   string
	ThumbPath string
	FileExe   string
	FileExt []string
}

/*
 初始化读取配置
 @author frontLon
 @return strucct
*/
func InitConf() {

	//读取通用配置
	general := &GeneralStruct{}
	if err := confSource.Section("General").MapTo(general);err != nil{}

	//读取Fc配置
	fc := &PfStruct{}
	if err := confSource.Section("Fc").MapTo(fc);err != nil{}
	sfc := &PfStruct{}
	if err := confSource.Section("Sfc").MapTo(sfc);err != nil{}
	md := &PfStruct{}
	if err := confSource.Section("Md").MapTo(md);err != nil{}
	pce := &PfStruct{}
	if err := confSource.Section("Pce").MapTo(pce);err != nil{}
	gb := &PfStruct{}
	if err := confSource.Section("Gb").MapTo(gb);err != nil{}
	arcade := &PfStruct{}
	if err := confSource.Section("Arcade").MapTo(arcade);err != nil{}


	//缩略图不支持相对路径，转换为绝对路径
	if fc.ThumbPath != "" && !strings.Contains(fc.ThumbPath, ":"){
		dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		fc.ThumbPath = dir + "\\"+fc.ThumbPath
	}

	if sfc.ThumbPath != "" && !strings.Contains(sfc.ThumbPath, ":"){
		dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		sfc.ThumbPath = dir + "\\"+sfc.ThumbPath
	}

	if md.ThumbPath != "" && !strings.Contains(md.ThumbPath, ":"){
		dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		md.ThumbPath = dir + "\\"+md.ThumbPath
	}

	if pce.ThumbPath != "" && !strings.Contains(pce.ThumbPath, ":"){
		dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		pce.ThumbPath = dir + "\\"+pce.ThumbPath
	}

	if gb.ThumbPath != "" && !strings.Contains(gb.ThumbPath, ":"){
		dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		gb.ThumbPath = dir + "\\"+gb.ThumbPath
	}

	if arcade.ThumbPath != "" && !strings.Contains(arcade.ThumbPath, ":"){
		dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		arcade.ThumbPath = dir + "\\"+arcade.ThumbPath
	}

	//配置全局参数
	Config = &ConfStruct{
		General : general,
		Fc:fc,
		Sfc:sfc,
		Md:md,
		Pce:pce,
		Gb:gb,
		Arcade:arcade,
	}
}

func updateConfig(section string,field string,value string) error {
	confSource.Section(section).Key(field).SetValue(value)
	err :=confSource.SaveTo(confFile);
	return err
}
