package config

import (
	"bufio"
	"errors"
	"github.com/go-ini/ini"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"simUI/code/db"
	"simUI/code/utils"
	"simUI/code/utils/go-sciter/window"
	"strings"
)

//配置文件
var (
	Cfg      *ConfStruct                                                                      //公共配置
	ENV      string                                                                           //环境配置
	DOC_EXTS = []string{".txt", ".md"}                                                        //doc文档支持的扩展名
	PIC_EXTS = []string{".png", ".jpg", ".gif", ".ico", ".jpeg", ".bmp", ".wmv", ".mp4", ".avi"} //支持的图片类型
	RUN_EXTS = []string{
		".html", ".htm", ".mht", ".mhtml", ".url",
		".pdf", ".chm", ".doc", ".docx", ".ppt", ".pptx", "xls", "xlsx", ".rtf",
		".exe", ".com", ".cmd", ".bat", ".lnk",
	}                                                                                         //可直接运行的doc文档支持的扩展名
	Window *window.Window                                                                     //窗体
)

//配置文件
type ConfStruct struct {
	ViewPath     string                  //代码路径
	RootPath     string                  //exe文件的当前路径
	Separator    string                  //分隔符
	CachePath    string                  //缓存路径
	UnzipPath    string                  //rom解压路径
	Default      *db.Config              //默认配置
	Shortcut     []*db.Shortcut          //快捷工具
	LangList     map[string]string       //语言列表
	Theme        map[string]*ThemeStruct //主题列表
	Lang         map[string]string       //语言项
	Platform     map[uint32]*db.Platform //平台及对应的模拟器列表（无序）
	PlatformList []*db.Platform          //平台及对应的模拟器列表（有序）
}

//主题配置
type ThemeStruct struct {
	Name   string            //主题名称
	Path   string            //文件路径
	Params map[string]string //主题各项参数
}

type RomBase struct {
	RomName   string // rom文件名
	Name      string // 游戏名称
	Type      string // 类型
	Year      string // 年份
	Platform  string // 平台
	Publisher string // 出品公司
}

/*
 初始化读取配置
 @author frontLon
*/
func InitConf() error {

	err := errors.New("")

	//更新缓存前，需要将工作目录换成默认目录
	if err := os.Chdir(Cfg.RootPath); err != nil {
		return err
	}

	Cfg.Default, err = getDefault()

	if err != nil {
		return err
	}
	Cfg.LangList, err = getLangList()
	if err != nil {
		return err
	}
	Cfg.Lang, err = getLang(Cfg.Default.Lang)
	if err != nil {
		return err
	}
	Cfg.Theme, err = getTheme()
	if err != nil {
		return err
	}
	Cfg.PlatformList, Cfg.Platform, err = getPlatform()
	if err != nil {
		return err
	}

	Cfg.Shortcut, err = getShortcut()
	if err != nil {
		return err
	}

	return nil
}

//读取平台列表
func getPlatform() ([]*db.Platform, map[uint32]*db.Platform, error) {
	platformList, _ := (&db.Platform{}).GetAll()
	platform := map[uint32]*db.Platform{}
	for k, v := range platformList {
		platform[v.Id] = v

		if v.Icon != "" {
			platformList[k].Icon, _ = filepath.Abs(v.Icon)
			platform[v.Id].Icon = platformList[k].Icon
		}

		if v.DocPath != "" {
			platformList[k].DocPath, _ = filepath.Abs(v.DocPath)
			platform[v.Id].DocPath = platformList[k].DocPath
		}

		if v.StrategyPath != "" {
			platformList[k].StrategyPath, _ = filepath.Abs(v.StrategyPath)
			platform[v.Id].StrategyPath = platformList[k].StrategyPath
		}

		if v.RomPath != "" {
			platformList[k].RomPath, _ = filepath.Abs(v.RomPath)
			platform[v.Id].RomPath = platformList[k].RomPath
		}

		if v.ThumbPath != "" {
			platformList[k].ThumbPath, _ = filepath.Abs(v.ThumbPath)
			platform[v.Id].ThumbPath = platformList[k].ThumbPath
		}

		if v.SnapPath != "" {
			platformList[k].SnapPath, _ = filepath.Abs(v.SnapPath)
			platform[v.Id].SnapPath = platformList[k].SnapPath
		}

		if v.PosterPath != "" {
			platformList[k].PosterPath, _ = filepath.Abs(v.PosterPath)
			platform[v.Id].PosterPath = platformList[k].PosterPath
		}

		if v.PackingPath != "" {
			platformList[k].PackingPath, _ = filepath.Abs(v.PackingPath)
			platform[v.Id].PackingPath = platformList[k].PackingPath
		}

		if v.TitlePath != "" {
			platformList[k].TitlePath, _ = filepath.Abs(v.TitlePath)
			platform[v.Id].TitlePath = platformList[k].TitlePath
		}

		if v.CassettePath != "" {
			platformList[k].CassettePath, _ = filepath.Abs(v.CassettePath)
			platform[v.Id].CassettePath = platformList[k].CassettePath
		}
		if v.IconPath != "" {
			platformList[k].IconPath, _ = filepath.Abs(v.IconPath)
			platform[v.Id].IconPath = platformList[k].IconPath
		}
		if v.GifPath != "" {
			platformList[k].GifPath, _ = filepath.Abs(v.GifPath)
			platform[v.Id].GifPath = platformList[k].GifPath
		}
		if v.BackgroundPath != "" {
			platformList[k].BackgroundPath, _ = filepath.Abs(v.BackgroundPath)
			platform[v.Id].BackgroundPath = platformList[k].BackgroundPath
		}
		if v.VideoPath != "" {
			platformList[k].VideoPath, _ = filepath.Abs(v.VideoPath)
			platform[v.Id].VideoPath = platformList[k].VideoPath
		}

		if v.FilesPath != "" {
			platformList[k].FilesPath, _ = filepath.Abs(v.FilesPath)
			platform[v.Id].FilesPath = platformList[k].FilesPath
		}

		if v.Rombase != "" {
			platformList[k].Rombase, _ = filepath.Abs(v.Rombase)
			platform[v.Id].Rombase = platformList[k].Rombase
		}

		//填充模拟器列表
		DBSim := &db.Simulator{}
		simList, _ := DBSim.GetByPlatform(v.Id)
		vomap := map[uint32]*db.Simulator{}
		for _, v := range simList {
			v.Path, _ = filepath.Abs(v.Path)
			vomap[v.Id] = v
		}
		platform[v.Id].SimList = vomap
		platformList[k].SimList = vomap
		platform[v.Id].UseSim = &db.Simulator{}
		//找到默认模拟器
		for _, sim := range simList {
			//当前正在使用的模拟器
			if sim.Default == 1 { //如果有默认模拟器
				platformList[k].UseSim = sim
				platform[v.Id].UseSim = sim
			} else if platformList[k].UseSim.Id == 0 { //如果没有默认模拟器，读取第一个
				platformList[k].UseSim = sim
				platform[v.Id].UseSim = sim
			}
		}
	}
	return platformList, platform, nil
}

//读取缓存配置
func getDefault() (*db.Config, error) {
	vo, err := (&db.Config{}).Get()

	//如果软件名称是图片，则转换为绝对路径

	if strings.ToLower(vo.SoftName) != "simui" && utils.FileExists(vo.SoftName) == true {
		vo.SoftName, _ = filepath.Abs(vo.SoftName)
	}

	//如果背景图片文件存在，则转换为绝对路径
	if utils.FileExists(vo.BackgroundOpacity) == true {
		vo.BackgroundOpacity, _ = filepath.Abs(vo.BackgroundOpacity)
	}

	//如果鼠标指针文件存在，则转换为绝对路径
	if utils.FileExists(vo.Cursor) == true {
		vo.Cursor, _ = filepath.Abs(vo.Cursor)
	}

	if err != nil {
		return vo, err
	}
	//查看当前选定平台值是否是正常的
	isset := false
	for _, v := range Cfg.Platform {
		if uint32(utils.ToInt(vo.Platform)) == v.Id {
			isset = true
			break
		}
	}

	//如果没有匹配上platform，则读取config中的第一项
	if vo.Platform != "0" {
		if isset == false {
			for _, v := range Cfg.Platform {
				vo.Platform = utils.ToString(v.Id)
				//修复配置文件
				if err := (&db.Config{}).UpdateField("platform", utils.ToString(vo.Platform)); err != nil {
				}
				break
			}
		}
	}

	return vo, nil
}

//读取主题列表
func getTheme() (map[string]*ThemeStruct, error) {
	dirPth := Cfg.RootPath + "theme" + Cfg.Separator
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
				if len(strarr) == 2 {
					//标题
					if id == "" {
						first := strings.Index(strarr[1], "(")
						last := strings.Index(strarr[1], ")")
						id = strarr[1][first+1 : last]
						continue
					}
					//内容
					first := strings.Index(strarr[0], "(")
					last := strings.Index(strarr[0], ")")
					key := strings.Trim(strarr[0][first+1:last], " ")
					value := strings.Trim(strings.Replace(strarr[1], ";", "", 1), " ")
					if key != "" && value != "" {
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
		err := errors.New(Cfg.Lang["ThemeFileNotFound"])
		return themelist, err
	}

	//如果当前的主题不存在，则将第一个主题更新到数据库
	if _, ok := themelist[Cfg.Default.Theme]; !ok {
		themeId := ""
		for k, _ := range themelist {
			themeId = k
			break
		}
		if err := (&db.Config{}).UpdateField("theme", themeId); err != nil {
			return themelist, err
		}
		Cfg.Default.Theme = themeId
	}

	return themelist, nil
}

//读取语言参数配置
func getLang(lang string) (map[string]string, error) {
	langpath := Cfg.RootPath + "lang" + Cfg.Separator
	fpath := langpath + lang + ".ini"
	section := make(map[string]string)

	//如果默认语言不存在，则读取列表中的其他语言
	if !utils.FileExists(fpath) {
		if len(Cfg.LangList) > 0 {
			for langName, langFile := range Cfg.LangList {
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
	dirPth := Cfg.RootPath + "lang" + Cfg.Separator
	lists, _ := ioutil.ReadDir(dirPth)
	for _, fi := range lists {
		if !fi.IsDir() { // 忽略目录
			name := strings.TrimSuffix(fi.Name(), path.Ext(fi.Name()))
			lang[name] = fi.Name()
		}
	}
	return lang, nil
}

//读取快捷工具列表
func getShortcut() ([]*db.Shortcut, error) {
	shortcutList, _ := (&db.Shortcut{}).GetAll()
	for k, v := range shortcutList {
		shortcutList[k].Path, _ = filepath.Abs(v.Path)
	}
	return shortcutList, nil
}

//读取全部资源目录
func GetResPath(platformId uint32) map[string]string {

	platform := Cfg.Platform[platformId]
	res := map[string]string{}
	res["thumb"] = platform.ThumbPath
	res["snap"] = platform.SnapPath
	res["poster"] = platform.PosterPath
	res["packing"] = platform.PackingPath
	res["title"] = platform.TitlePath
	res["cassette"] = platform.CassettePath
	res["icon"] = platform.IconPath
	res["gif"] = platform.GifPath
	res["background"] = platform.BackgroundPath
	res["video"] = platform.VideoPath
	res["doc"] = platform.DocPath
	res["strategy"] = platform.StrategyPath
	res["files"] = platform.FilesPath
	return res
}

//读取资源类型名
func GetResExts() map[string][]string {
	res := map[string][]string{}
	res["thumb"] = PIC_EXTS
	res["snap"] = PIC_EXTS
	res["poster"] = PIC_EXTS
	res["packing"] = PIC_EXTS
	res["title"] = PIC_EXTS
	res["cassette"] = PIC_EXTS
	res["icon"] = PIC_EXTS
	res["gif"] = PIC_EXTS
	res["background"] = PIC_EXTS
	res["video"] = PIC_EXTS
	res["doc"] = DOC_EXTS
	res["strategy"] = DOC_EXTS
	res["files"] = RUN_EXTS
	return res
}
