package main

import (
	"VirtualNesGUI/code/db"
	"VirtualNesGUI/code/utils"
	"bufio"
	"errors"
	"github.com/go-ini/ini"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

//配置文件
var Config *ConfStruct

//配置文件
type ConfStruct struct {
	RootPath  string                  //exe文件的当前路径
	Separator string                  //exe文件的当前路径
	Default   *db.Config              //默认配置
	LangList  map[string]string       //语言列表
	Theme     map[string]*ThemeStruct //主题列表
	Lang      map[string]string       //语言项
	Platform  map[uint32]*db.Platform //平台及对应的模拟器列表（无序）
	PlatformList  []*db.Platform //平台及对应的模拟器列表（有序）
}

//主题配置
type ThemeStruct struct {
	Name   string            //主题名称
	Path   string            //文件路径
	Params map[string]string //主题各项参数
}

/*
 初始化读取配置
 @author frontLon
*/
func InitConf() error {

	err := errors.New("")

	//更新缓存前，需要将工作目录换成默认目录
	if err := os.Chdir(Config.RootPath); err != nil {
		return err
	}
	Config.Default, err = getDefault()
	if err != nil {
		return err
	}
	Config.LangList, err = getLangList()
	if err != nil {
		return err
	}
	Config.Lang, err = getLang(Config.Default.Lang)
	if err != nil {
		return err
	}
	Config.PlatformList,Config.Platform, err = getPlatform()
	if err != nil {
		return err
	}
	Config.Theme, err = getTheme()
	if err != nil {
		return err
	}
	return nil
}

//读取平台列表
func getPlatform() ([]*db.Platform,map[uint32]*db.Platform, error) {
	DBSim := &db.Simulator{}
	platformList, _ := (&db.Platform{}).GetAll()
	platform := map[uint32]*db.Platform{}
	for k, v := range platformList {
		platform[v.Id] = v

		if v.DocPath != "" {
			platformList[k].DocPath,_ = filepath.Abs(v.DocPath)
			platform[v.Id].DocPath = platformList[k].DocPath
		}

		if v.StrategyPath != "" {
			platformList[k].StrategyPath,_ = filepath.Abs(v.StrategyPath)
			platform[v.Id].StrategyPath = platformList[k].StrategyPath
		}

		if v.RomPath != "" {
			platformList[k].RomPath,_ = filepath.Abs(v.RomPath)
			platform[v.Id].RomPath = platformList[k].RomPath
		}

		if v.ThumbPath != "" {
			platformList[k].ThumbPath,_ = filepath.Abs(v.ThumbPath)
			platform[v.Id].ThumbPath = platformList[k].ThumbPath
		}

		if v.SnapPath != "" {
			platformList[k].SnapPath,_ = filepath.Abs(v.SnapPath)
			platform[v.Id].SnapPath = platformList[k].SnapPath
		}

		if v.Romlist != "" {
			platformList[k].Romlist,_ = filepath.Abs(v.Romlist)
			platform[v.Id].Romlist = platformList[k].Romlist
		}

		//填充模拟器列表
		simList,_ := DBSim.GetByPlatform(v.Id)
		platform[v.Id].SimList = simList
		platformList[k].SimList = simList

		platform[v.Id].UseSim = &db.Simulator{}
		//找到默认模拟器
		for sk, sim := range simList {
			//当前正在使用的模拟器
			if sim.Default == 1 {
				platformList[k].UseSim = sim
				platform[v.Id].UseSim = sim
			}
			//模拟器路径转换为绝对路径
			if sim.Path != "" {
				sim.Path,_ = filepath.Abs(sim.Path)
				platformList[k].SimList[sk].Path = sim.Path
				platform[v.Id].SimList[sk].Path = sim.Path
			}
		}
	}
	return platformList,platform, nil
}

//读取缓存配置
func getDefault() (*db.Config, error) {
	vo, err := (&db.Config{}).Get()
	if err != nil{
		return vo,err
	}
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
	if vo.Book != "" {
		vo.Book, _ = filepath.Abs(vo.Book)
	}
	return vo, nil
}

//读取主题列表
func getTheme() (map[string]*ThemeStruct, error) {
	dirPth := Config.RootPath + "theme" + separator
	lists, _ := ioutil.ReadDir(dirPth)

	themelist := map[string]*ThemeStruct{}
	for _, fi := range lists {
		ext := strings.ToLower(path.Ext(fi.Name())) //获取文件后缀
		if !fi.IsDir() && ext == ".css" { // 忽略目录

			filename := dirPth + fi.Name()
			file, err := os.Open(filename) //打开文件

			if err != nil {
				return themelist, err
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
						if (key == "window-background-image" ||
							key == "desc-background-image" ||
							key == "default-thumb-image") {
							value = dirPth + value
						}
						params[key] = value
					}
				}
			}
			themelist[id] = &ThemeStruct{
				Name:   utils.GetFileName(fi.Name()),
				Path:   filename,
				Params: params,
			}
			file.Close()
		}
	}

	if len(themelist) == 0 {
		err := errors.New(Config.Lang["ThemeFileNotFound"])
		return themelist, err
	}

	//如果当前的主题不存在，则将第一个主题更新到数据库
	if _, ok := themelist[Config.Default.Theme]; !ok {
		themeId := ""
		for k,_ := range themelist{
			themeId = k
			break
		}
		if err := (&db.Config{}).UpdateField("theme", themeId); err != nil {
			return themelist, err
		}
		Config.Default.Theme = themeId
	}

	return themelist, nil
}

//读取ROM别名配置参数
func getRomAlias(platform uint32) (map[string]string, error) {
	section := make(map[string]string)
	file, err := ini.LoadSources(ini.LoadOptions{IgnoreInlineComment: true}, Config.Platform[platform].Romlist)
	if err != nil {
		return section, err
	}
	section = file.Section("Alias").KeysHash()
	return section, nil
}

//读取语言参数配置
func getLang(lang string) (map[string]string, error) {
	langpath := Config.RootPath + "lang" + separator
	fpath := langpath + lang + ".ini"
	section := make(map[string]string)

	//如果默认语言不存在，则读取列表中的其他语言
	if !utils.FileExists(fpath) {
		if len(Config.LangList) > 0 {
			for langName, langFile := range Config.LangList {
				fpath = langpath + langFile
				//如果找到其他语言，则将第一项更新到数据库配置中
				if err := (&db.Config{}).UpdateField("lang", langName); err != nil {
					return section, err
				}
				break
			}
		}
	}

	file, err := ini.LoadSources(ini.LoadOptions{IgnoreInlineComment: true}, fpath)

	if err != nil {
		return section, err
	}

	section = file.Section("").KeysHash()
	return section, nil
}

//读取语言文件列表
func getLangList() (map[string]string, error) {
	lang := make(map[string]string)
	dirPth := Config.RootPath + "lang" + separator
	lists, _ := ioutil.ReadDir(dirPth)
	for _, fi := range lists {
		if !fi.IsDir() { // 忽略目录
			name := strings.TrimSuffix(fi.Name(), path.Ext(fi.Name()))
			lang[name] = fi.Name()
		}
	}
	return lang, nil
}
