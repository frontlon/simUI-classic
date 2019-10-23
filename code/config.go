package main

import (
	"VirtualNesGUI/code/db"
	"VirtualNesGUI/code/utils"
	"bufio"
	"fmt"
	"github.com/go-ini/ini"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

//配置文件
var Config *ConfStruct
//读取配置文件

//配置文件
type ConfStruct struct {
	RootPath  string                  //exe文件的当前路径
	Separator string                  //exe文件的当前路径
	Default   *db.Config              //默认配置
	LangList  map[string]string       //语言列表
	Theme     map[string]*ThemeStruct //主题列表
	Lang      map[string]string       //语言项
	Platform  map[uint32]*db.Platform //平台及对应的模拟器列表
}

type ThemeStruct struct {
	Name   string
	Path   string
	Params map[string]string
}

/*
 初始化读取配置
 @author frontLon
 @return strucctfilepath.Abs
*/
func InitConf() {
	var rootpath, _ = filepath.Abs(filepath.Dir(os.Args[0])) //exe运行文件路径
	Config = &ConfStruct{}
	//配置全局参数
	Config.Platform = getPlatform()
	Config.Default = getDefault()
	Config.LangList = getLangList()
	Config.Lang = getLang(Config.Default.Lang)
	Config.Theme = getTheme()
	Config.RootPath = rootpath + separator //exe文件的绝对路径
	Config.Separator = separator           //系统的目录分隔符
}

//读取平台列表
func getPlatform() map[uint32]*db.Platform {
	DBSim := &db.Simulator{}
	platform, _ := (&db.Platform{}).GetAll()
	for _, v := range platform {
		platform[v.Id].SimList, _ = DBSim.GetByPlatform(v.Id) //填充模拟器
		platform[v.Id].DocPath, _ = filepath.Abs(platform[v.Id].DocPath)
		platform[v.Id].Romlist, _ = filepath.Abs(platform[v.Id].Romlist)
		platform[v.Id].StrategyPath, _ = filepath.Abs(platform[v.Id].StrategyPath)
		platform[v.Id].RomPath, _ = filepath.Abs(platform[v.Id].RomPath)
		platform[v.Id].ThumbPath, _ = filepath.Abs(platform[v.Id].ThumbPath)
		platform[v.Id].SnapPath, _ = filepath.Abs(platform[v.Id].SnapPath)
		platform[v.Id].UseSim = &db.Simulator{}
		for sk, sim := range platform[v.Id].SimList {
			//当前正在使用的模拟器
			if sim.Default == 1 {
				platform[v.Id].UseSim = sim
			}
			//模拟器路径转换为绝对路径
			platform[v.Id].SimList[sk].Path, _ = filepath.Abs(platform[v.Id].SimList[sk].Path)
		}
	}
	return platform
}

//读取缓存配置
func getDefault() *db.Config {
	vo, _ := (&db.Config{}).Get()
	//查看当前选定平台值是否是正常的
	isset := false
	for _, v := range (Config.Platform) {
		if vo.Platform == v.Id {
			isset = true
			break
		}
	}

	//如果没有匹配上platform，则读取config中的第一项
	if vo.Platform != 0 {
		if isset == false {
			for _, v := range (Config.Platform) {
				vo.Platform = v.Id
				//修复配置文件
				if err := (&db.Config{}).UpdateField("platform", utils.ToString(vo.Platform)); err != nil {
				}
				break
			}
		}
	}
	//相对路径转换为绝对路径
	vo.Book, _ = filepath.Abs(vo.Book)

	return vo
}

//读取主题列表
func getTheme() map[string]*ThemeStruct {
	dirPth, _ := filepath.Abs("theme")
	lists, _ := ioutil.ReadDir(dirPth)
	themelist := map[string]*ThemeStruct{}
	for _, fi := range lists {
		ext := strings.ToLower(path.Ext(fi.Name())) //获取文件后缀
		if !fi.IsDir() && ext == ".css" { // 忽略目录

			filename := dirPth + separator + fi.Name()
			file, err := os.Open(filename) //打开文件

			if err != nil {
				fmt.Println(err.Error())
			}
			scanner := bufio.NewScanner(file) //扫描文件
			lineText := ""
			//只读取第一行
			id := ""
			params := make(map[string]string)
			isnode := false
			for scanner.Scan() {
				lineText = scanner.Text()
				//过滤掉注释部分
				if strings.Index(lineText, `*/`) != -1 {
					isnode = false
					continue
				}
				if isnode == true {
					continue
				}
				if strings.Index(lineText, `/*`) != -1 {
					isnode = true
					if strings.Index(lineText, `*/`) != -1 {
						isnode = false
					}
					continue
				}
				strarr := strings.Split(lineText, ":")
				if (len(strarr) == 2) {
					//标题
					if id == "" {
						first := strings.Index(strarr[1], "(");
						last := strings.Index(strarr[1], ")");
						id = strarr[1][first+1 : last]
						continue
					}
					//内容
					first := strings.Index(strarr[0], "(");
					last := strings.Index(strarr[0], ")");
					key := strings.Trim(strarr[0][first+1:last], " ")
					value := strings.Trim(strings.Replace(strarr[1], ";", "", 1), " ")
					if key != "" && value != "" {
						if (key == "window-background-image" || key == "desc-background-image") {
							value = dirPth + separator + value
						}

						params[key] = value
					}
				}
			}
			themelist[id] = &ThemeStruct{
				Name:   GetFileName(fi.Name()),
				Path:   filename,
				Params: params,
			}
			file.Close()
		}
	}
	return themelist
}

//读取ROM别名配置参数
func getRomAlias(platform uint32) map[string]string {
	section := make(map[string]string)
	file, err := ini.LoadSources(ini.LoadOptions{IgnoreInlineComment: true}, Config.Platform[platform].Romlist)
	if err == nil {
		section = file.Section("Alias").KeysHash()
	}
	return section
}

//读取语言参数配置
func getLang(lang string) map[string]string {
	langpath := Config.RootPath + "lang" + separator
	fpath := langpath + lang + ".ini"
	section := make(map[string]string)

	//如果默认语言不存在，则读取列表中的其他语言
	if !Exists(fpath) {
		if len(Config.LangList) > 0 {
			for langName, langFile := range Config.LangList {
				fpath = langpath + langFile
				//如果找到其他语言，则将第一项更新到数据库配置中
				if err := (&db.Config{}).UpdateField("lang", langName); err != nil {
				}
				break
			}
		}
	}

	file, err := ini.LoadSources(ini.LoadOptions{IgnoreInlineComment: true}, fpath)
	if err == nil {
		section = file.Section("").KeysHash()
	}
	return section
}

//读取语言文件列表
func getLangList() map[string]string {
	lang := make(map[string]string)
	dirPth, _ := filepath.Abs("lang")
	lists, _ := ioutil.ReadDir(dirPth)
	for _, fi := range lists {
		if !fi.IsDir() { // 忽略目录
			name := strings.TrimSuffix(fi.Name(), path.Ext(fi.Name()))
			lang[name] = fi.Name()
		}
	}
	return lang
}
