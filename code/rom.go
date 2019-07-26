package main

import (
	"github.com/axgle/mahonia"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

var constRomList = &Rom{}
var constRomAlias = map[string]map[string]string{} //rom别名列表，用于启动游戏
var constRomPath = map[string]map[string]string{}  //路径列表，用于启动游戏
var constPageLimit = 100                           //rom读取分页大小
var constCurrentRomCount int                       //当前选项的rom总数
var constSeparator = "__"                          //rom子分隔符
var constSubRomList = map[string][]*Rominfo{}      //子游戏缓存列表
type Rom struct {
	Platform map[string][]*Rominfo
}

//游戏信息
type Rominfo struct {
	Title string //标题
	Menu  string //目录
	Thumb string //缩略图
	Video string //视频地址
}

//读取fc rom列表
func getRomList(platform string) {
	platform = strings.ToLower(platform)
	constRomPath[platform] = make(map[string]string)

	//读取平台配置
	conf := Config.Platform[platform]

	//路径最后加入反斜杠
	conf.RomPath = conf.RomPath + separator

	romlist := []*Rominfo{}

	//载入别名
	if _, ok := constRomAlias[platform]; !ok {
		constRomAlias[platform] = getRomAlias(platform)

	}

	//进入循环，遍历文件
	if err := filepath.Walk(conf.RomPath,
		func(p string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}

			//整理目录格式，并转换为数组
			newpath := strings.Replace(conf.RomPath, "/", "\\", -1)
			if newpath[0:2] == ".\\" {
				p = ".\\" + p
			}
			newpath = strings.Replace(p, newpath, "", -1)
			subpath := strings.Split(newpath, "\\")
			romExt := strings.ToLower(path.Ext(p)) //获取文件后缀

			//如果该文件是游戏rom
			if f.IsDir() == false && CheckPlatformExt(conf.RomExt, romExt) {
				title := GetRomName(f.Name())
				//如果有别名配置，则读取别名
				if _, ok := constRomAlias[platform][title]; ok {
					title = constRomAlias[platform][title]
				}

				//如果游戏名称存在分隔符，说明是子游戏
				if strings.Contains(title, constSeparator) {

					//拆分文件名
					sub := strings.Split(title, constSeparator)

					//去掉扩展名，生成标题
					sinfo := &Rominfo{
						Title: sub[1],
						Menu:  "",
						Thumb: "",
						Video: "", //video在详情时在加载
					}

					//子列表
					subKey := sub[0]
					if _, ok := constRomAlias[platform][sub[0]]; ok {
						subKey = constRomAlias[platform][sub[0]]
					}

					constSubRomList[subKey] = append(constSubRomList[subKey], sinfo)
					//rom路径列表
					constRomPath[platform][title] = p
				} else {
					menu := constMenuRootKey //无目录，读取默认参数
					//定义目录，如果有子目录，则记录子目录名称
					if len(subpath) > 1 {
						menu = subpath[0]
					}

					//读取缩略图
					_, thumb := getThumb(conf.ThumbPath, GetRomName(f.Name()))
					//去掉扩展名，生成标题
					rinfo := &Rominfo{
						Title: title,
						Menu:  menu,
						Thumb: thumb,
						Video: "",
					}
					//rom列表
					romlist = append(romlist, rinfo)
					//rom path路径列表
					constRomPath[platform][title] = p
				}
			}
			return nil
		}); err != nil {
	}

	//赋值给全局变量
	constRomList.Platform[platform] = romlist
}

/**
 * 读取游戏缩略图
 **/
func getThumb(thumbPath string, title string) (bool, string) {

	img := thumbPath + separator + title + ".png"
	isset := true
	if !exists(img) {
		img = thumbPath + separator + ".." + separator + "_DEF_.png"
		isset = false
	}
	return isset, img
}

/**
 * 读取游戏视频
 **/
func getVideo(platform string, title string) string {
	video := ""
	file := Config.Platform[platform].VideoPath + separator + title + ".gif"
	if exists(file) {
		video = file
	}
	return video
}

/**
 * 读取游戏介绍
 **/
func getDoc(platform string, title string) string {
	file := Config.Platform[platform].DescPath + separator + title + ".txt"
	text, err := ioutil.ReadFile(file)
	content := ""
	if err != nil {
		return content
	}
	enc := mahonia.NewDecoder("gbk")
	content = enc.ConvertString(string(text))
	return content
}

/**
 * 运行游戏
 **/
func runGame(platform string, name string) string {

	exeFile := Config.Platform[platform].FileExe

	//检测执行文件是否存在
	_, err := os.Stat(exeFile)
	if err != nil {
		return err.Error()
	}

	//检测rom文件是否存在
	if exists(constRomPath[platform][name]) == false {
		return Config.Lang["RomNotFound"] + constRomPath[platform][name]
	}

	cmd := exec.Command(exeFile,constRomPath[platform][name])
	if err := cmd.Start(); err != nil {
		return err.Error()
	}
	return ""
}

/**
 * 检测文件是否存在（文件夹也返回false）
 **/
func exists(path string) bool {

	if path == "" {
		return false
	}

	finfo, err := os.Stat(path)
	isset := false
	if err != nil || finfo.IsDir() == true {
		isset = false
	} else {
		isset = true
	}
	return isset
}

/**
 * 去掉rom扩展名，从文件名中读取Rom名称
 **/
func GetRomName(filename string) string {
	return strings.TrimSuffix(filename, path.Ext(filename))
}

//检查文件扩展名是否存在于该平台中
func CheckPlatformExt(exts []string, file string) bool {
	for _, v := range exts {
		if v == file {
			return true
		}
	}
	return false
}
