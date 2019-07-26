/*
package main
var res = []byte{
*/
package main

import (
	"encoding/json"
	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

var separator = string(os.PathSeparator) //路径分隔符

//var constMainFile = "D:\\work\\go\\src\\VirtualNesGUI\\code\\res\\main.html" //主文件路径（测试用）
var constMainFile = "this://app/main.html" //主文件路径（正式）

//rom详情数据
type DetailStruct struct {
	Sublist []*Rominfo
	Doc     string
	Video   string
}

func main() {

	//初始化配置
	InitConf()

	//初始化菜单列表
	constMenuList.Platform = make(map[string][]*MenuInfo)
	if err := GetMenuData(Config.Default.Platform); err != nil {
	}
	//初始化rom列表
	constRomList.Platform = make(map[string][]*Rominfo)
	getRomList(Config.Default.Platform)

	//初始化主题列表
	getThemeList();
	left, _ := strconv.Atoi(Config.Default.WindowLeft)
	top, _ := strconv.Atoi(Config.Default.WindowTop)
	width, _ := strconv.Atoi(Config.Default.WindowWidth)
	height, _ := strconv.Atoi(Config.Default.WindowHeight)

	//创建window窗口
	w, err := window.New(
		sciter.SW_MAIN|
			sciter.SW_ENABLE_DEBUG,
		&sciter.Rect{Left: int32(left), Top: int32(top), Right: int32(width), Bottom: int32(height)});
	if err != nil {
		log.Fatal(err);
	}

	//设置view权限
	w.SetOption(sciter.SCITER_SET_SCRIPT_RUNTIME_FEATURES, sciter.ALLOW_SYSINFO);

	//设置回调
	w.SetCallback(newHandler(w.Sciter))

	//解析资源
	w.OpenArchive(res)

	//加载文件
	err = w.LoadFile(constMainFile);
	if err != nil {
		if _, err := w.Call("errorBox", sciter.NewValue(err.Error())); err != nil {
		}
		return
	}

	//设置标题
	w.SetTitle(Config.Lang["Title"]);

	//定义view函数
	defineViewFunction(w)

	//显示窗口
	w.Show();

	//运行窗口，进入消息循环
	w.Run();
}

/**
 * 定义view用function
 **/
func defineViewFunction(w *window.Window) {

	w.DefineFunction("InitData", func(args ...*sciter.Value) *sciter.Value {
		ctype := args[0].String()
		data := ""
		switch (ctype) {
		case "config": //读取配置
			getjson, _ := json.Marshal(Config)
			data = string(getjson)
		case "platform": //当前平台列表
			getjson, _ := json.Marshal(constPlatformList)
			data = string(getjson)
		case "menu": //当前菜单列表
			getjson, _ := json.Marshal(constMenuList.Platform[Config.Default.Platform])
			data = string(getjson)
		case "themeList": //当前菜单列表
			getjson, _ := json.Marshal(constThemeList)
			data = string(getjson)
		case "rom": //当前rom列表
			page := 0
			constCurrentRomCount = len(constRomList.Platform[Config.Default.Platform]) //当前分类的rom总数
			getjson := []byte{}
			//如果现有rom小于分页数量，则全部显示
			if len(constRomList.Platform[Config.Default.Platform][page*constPageLimit:]) > constPageLimit {
				if len(constRomList.Platform[Config.Default.Platform]) <= constPageLimit {
					getjson, _ = json.Marshal(constRomList.Platform[Config.Default.Platform])
				} else { //如果rom数量太多，则显示少量
					getjson, _ = json.Marshal(constRomList.Platform[Config.Default.Platform][page*constPageLimit : constPageLimit*page+constPageLimit])
				}
			} else {
				getjson, _ = json.Marshal(constRomList.Platform[Config.Default.Platform][page*constPageLimit:])
			}
			data = string(getjson)
		case "romCount": //当前rom列表
			data = strconv.Itoa(constCurrentRomCount)
		}
		return sciter.NewValue(data)
	})

	//运行游戏
	w.DefineFunction("RunGame", func(args ...*sciter.Value) *sciter.Value {
		platform := args[0].String()
		filename := args[1].String()

		err := runGame(platform, filename);
		if err != "" {
			if _, err := w.Call("errorBox", sciter.NewValue(err)); err != nil {
			}
		}
		return sciter.NullValue()
	})

	//打开rom目录
	w.DefineFunction("OpenFolder", func(args ...*sciter.Value) *sciter.Value {
		gtype := args[0].String() //目录类型
		platform := args[1].String()
		p := ""
		switch gtype {
		case "rom":
			p = Config.Platform[platform].RomPath
		case "thumb":
			p = Config.Platform[platform].ThumbPath
		case "video":
			p = Config.Platform[platform].VideoPath
		case "sim":
			exe := Config.Platform[platform].FileExe
			p = filepath.Dir(exe)
		}
		if err := exec.Command(`cmd`, `/c`, `explorer`, p).Start(); err != nil {
			if _, err := w.Call("errorBox", sciter.NewValue(err)); err != nil {
			}
		}
		return sciter.NullValue()
	})

	//更新配置文件
	w.DefineFunction("UpdateConfig", func(args ...*sciter.Value) *sciter.Value {
		section := args[0].String()
		field := args[1].String()
		value := args[2].String()
		err := updateConfig(section, field, value);
		if err != nil {
			if _, err := w.Call("errorBox", sciter.NewValue(err.Error())); err != nil {
			}
		}
		return sciter.NullValue()
	})

	//读取目录列表
	w.DefineFunction("GetMenuList", func(args ...*sciter.Value) *sciter.Value {
		platform := strings.Trim(args[0].String(), " ")

		//开始读取游戏列表，如果没有读取，重新读取
		if _, ok := constMenuList.Platform[platform]; !ok {
			if err := GetMenuData(platform); err != nil {
			}
		}
		menuList := constMenuList.Platform[platform]
		jsonMenu, _ := json.Marshal(&menuList)
		return sciter.NewValue(string(jsonMenu))
	})

	//读取游戏列表
	w.DefineFunction("GetGameList", func(args ...*sciter.Value) *sciter.Value {
		platform := strings.Trim(args[0].String(), " ")
		catname := strings.Trim(args[1].String(), " ")
		keyword := strings.Trim(args[2].String(), " ")
		page, _ := strconv.Atoi(strings.Trim(args[3].String(), " ")) //分页数

		//开始读取游戏列表，如果没有读取，重新读取
		if _, ok := constRomList.Platform[platform]; !ok {
			getRomList(platform);
		}
		romlist := constRomList.Platform[platform]
		newlist := []*Rominfo{}

		if catname == Config.Lang["Uncate"] { //未分类
			catname = constMenuRootKey
		}

		if catname == Config.Lang["AllGames"] && keyword == "" {
			newlist = romlist
		} else {
			for _, v := range romlist {

				if catname == Config.Lang["AllGames"] {
					//关键字搜索
					if strings.Contains(v.Title, keyword) {
						newlist = append(newlist, v)
					}
				} else {
					if catname == v.Menu {
						if keyword != "" {
							//关键字搜索
							if strings.Contains(v.Title, keyword) {
								newlist = append(newlist, v)
							}
						} else {
							//非关键字搜索
							newlist = append(newlist, v)
						}
					}
				}
			}
		}

		constCurrentRomCount = len(newlist) //记录当前分类的rom总数

		jsonRom, _ := json.Marshal([]*Rominfo{})
		if len(newlist) > 0 && page*constPageLimit < len(newlist) {

			//数据大于1页
			if len(newlist[page*constPageLimit:]) > constPageLimit {
				if len(newlist) <= constPageLimit {
					jsonRom, _ = json.Marshal(newlist)
				} else { //如果rom数量太多，则显示少量
					jsonRom, _ = json.Marshal(newlist[page*constPageLimit : constPageLimit*page+constPageLimit])
				}
			} else {
				//数据不足一页，全部显示出去
				jsonRom, _ = json.Marshal(newlist[page*constPageLimit:])
			}
		}
		return sciter.NewValue(string(jsonRom))
	})

	//读取rom详情
	w.DefineFunction("GetRomDetail", func(args ...*sciter.Value) *sciter.Value {
		platform := strings.Trim(args[0].String(), " ")
		title := strings.Trim(args[1].String(), " ")
		detail := &DetailStruct{}

		//加载子rom列表
		if _, ok := constSubRomList[title]; ok {
			detail.Sublist = constSubRomList[title]
		}

		detail.Doc = getDoc(platform, title)
		detail.Video = getVideo(platform, title)

		jsonMenu, _ := json.Marshal(&detail)
		return sciter.NewValue(string(jsonMenu))
	})

	//读取主题参数
	w.DefineFunction("GetThemeParams", func(args ...*sciter.Value) *sciter.Value {
		title := strings.Trim(args[0].String(), " ")
		theme, err := getThemeParams(title);
		if err != nil {
			if _, err := w.Call("errorBox", sciter.NewValue(err)); err != nil {
			}
			return sciter.NewValue(string(""))
		}
		jsonTheme, _ := json.Marshal(&theme)
		return sciter.NewValue(string(jsonTheme))
	})

}

//资源加载
func OnLoadData(s *sciter.Sciter) func(ld *sciter.ScnLoadData) int {
	return func(ld *sciter.ScnLoadData) int {
		uri := ld.Uri()
		if strings.HasPrefix(uri, "this://app/") {
			path := uri[11:]
			data := s.GetArchiveItem(path)
			if data == nil {
				return sciter.LOAD_OK
			}
			s.DataReady(uri, data)
		}
		return sciter.LOAD_OK
	}
}

func newHandler(s *sciter.Sciter) *sciter.CallbackHandler {
	return &sciter.CallbackHandler{
		OnLoadData: OnLoadData(s),
	}
}
