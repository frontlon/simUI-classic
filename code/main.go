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
	"strings"
)

//var constMainFile = "D:\\work\\go\\src\\VirtualNesGUI\\code\\res\\main.html" //主文件路径（测试用）
var constMainFile = "this://app/main.html" //主文件路径（正式）
var constTitle = "Tiny Game UI" //软件名称
var constMenuRootKey = "_7b9" //根子目录游戏的Menu参数

func main() {

	//初始化配置
	InitConf()

	//定义标题
	if Config.General.Title != ""{
		constTitle = Config.General.Title
	}

	//创建window窗口
	w, err := window.New(
		sciter.SW_TITLEBAR|
		sciter.SW_RESIZEABLE|
		sciter.SW_CONTROLS|
		sciter.SW_MAIN|
		sciter.SW_ENABLE_DEBUG,
		&sciter.Rect{Left: 0, Top: 0, Right: 1280, Bottom: 768});

	if err != nil {
		log.Fatal(err);
	}

	//设置view权限
	w.SetOption(sciter.SCITER_SET_SCRIPT_RUNTIME_FEATURES,sciter.ALLOW_SYSINFO);

	//设置回调
	w.SetCallback(newHandler(w.Sciter))

	//解析资源
	w.OpenArchive(res)

	//加载文件
	err = w.LoadFile(constMainFile);
	if err != nil {
		w.Call("errorBox",sciter.NewValue("加载文件失败"))
		return
	}

	//缩略图thumbs
	_, err = os.Stat(Config.Fc.ThumbPath)
	if err != nil {
		os.Mkdir(Config.Fc.ThumbPath, os.ModePerm) //创建目录
	}

	//设置标题
	w.SetTitle(constTitle);

	//初始化平台列表
	GetPlatformData()

	//初始化菜单列表
	err = GetMenuData(Config.General.Platform)

	//初始化rom列表
	getRomList(Config.General.Platform)

	//定义view函数
	defineViewFunction(w)

	//显示窗口
	w.Show();

	//初始化view数据
	menu := []*MenuInfo{}
	rom := []*Rominfo{}
	switch(Config.General.Platform){
	case "fc":
		menu = constMenuList.Fc
		rom = constRomList.Fc
	case "sfc":
		menu = constMenuList.Sfc
		rom = constRomList.Sfc
	case "md":
		menu = constMenuList.Md
		rom = constRomList.Md
	case "pce":
		menu = constMenuList.Pce
		rom = constRomList.Pce
	case "gb":
		menu = constMenuList.Gb
		rom = constRomList.Gb
	case "arcade":
		menu = constMenuList.Arcade
		rom = constRomList.Arcade
	default:
		menu = constMenuList.Fc
		rom = constRomList.Fc
	}

	jsonPlatform, err := json.Marshal(&constPlatformList)
	jsonMenu, err := json.Marshal(&menu)
	jsonRom, err := json.Marshal(&rom)

	w.Call("initView",
		sciter.NewValue(constTitle),
		sciter.NewValue(Config.General.Platform),
		sciter.NewValue(string(jsonPlatform)),
		sciter.NewValue(string(jsonMenu)),
		sciter.NewValue(string(jsonRom)),
	)

	//运行窗口，进入消息循环
	w.Run();
}

/**
 * 定义view用function
 **/
func defineViewFunction(w *window.Window) {
	//运行游戏
	w.DefineFunction("RunGame", func(args ...*sciter.Value) *sciter.Value {
		platform := args[0].String()
		filepath := args[1].String()
		err := runGame(platform,filepath);
		if err != ""{
			w.Call("errorBox",sciter.NewValue(err))
		}
		return sciter.NullValue()
	})

	//更新配置文件
	w.DefineFunction("UpdateConfig", func(args ...*sciter.Value) *sciter.Value {
		section := args[0].String()
		field := args[1].String()
		value := args[2].String()
		err := updateConfig(section,field,value);
		if err != nil{
			w.Call("errorBox",sciter.NewValue("更新配置文件错误:",err.Error()))
		}
		return sciter.NullValue()
	})

	//读取目录列表
	w.DefineFunction("GetMenuList", func(args ...*sciter.Value) *sciter.Value {
		platform := strings.Trim(args[0].String()," ")

		//开始读取游戏列表，如果没有读取，重新读取
		menuList := []*MenuInfo{}
		switch(platform){
		case "fc":
			if len(constMenuList.Fc) == 0 {
				GetMenuData(platform)
			}
			menuList = constMenuList.Fc
		case "sfc":
			if len(constMenuList.Sfc) == 0 {
				GetMenuData(platform)
			}
			menuList = constMenuList.Sfc
		case "md":
			if len(constMenuList.Md) == 0 {
				GetMenuData(platform)
			}
			menuList = constMenuList.Md
		case "pce":
			if len(constMenuList.Pce) == 0 {
				GetMenuData(platform)
			}
			menuList = constMenuList.Pce
		case "gb":
			if len(constMenuList.Gb) == 0 {
				GetMenuData(platform)
			}
			menuList = constMenuList.Gb
		case "arcade":
			if len(constMenuList.Arcade) == 0 {
				GetMenuData(platform)
			}
			menuList = constMenuList.Arcade
		}
		jsonMenu, _ := json.Marshal(&menuList)

		return sciter.NewValue(string(jsonMenu))
	})

	//读取游戏列表
	w.DefineFunction("GetGameList", func(args ...*sciter.Value) *sciter.Value {
		platform := strings.Trim(args[0].String()," ")
		catname := strings.Trim(args[1].String()," ")
		keyword := strings.Trim(args[2].String()," ")

		//开始读取游戏列表，如果没有读取，重新读取
		romlist :=[]*Rominfo{}
		switch(platform){
		case "fc":
			if len(constRomList.Fc) == 0 {
				getRomList(platform)
			}
			romlist = constRomList.Fc
		case "sfc":
			if len(constRomList.Sfc) == 0 {
				getRomList(platform)
			}
			romlist = constRomList.Sfc
		case "md":
			if len(constRomList.Md) == 0 {
				getRomList(platform)
			}
			romlist = constRomList.Md
		case "pce":
			if len(constRomList.Pce) == 0 {
				getRomList(platform)
			}
			romlist = constRomList.Pce
		case "gb":
			if len(constRomList.Gb) == 0 {
				getRomList(platform)
			}
			romlist = constRomList.Gb
		case "arcade":
			if len(constRomList.Arcade) == 0 {
				getRomList(platform)
			}
			romlist = constRomList.Arcade
		}

		if catname == "全部游戏" && keyword == ""{
			jsonRom, _ := json.Marshal(&romlist)
			return sciter.NewValue(string(jsonRom))
		}else if catname == "未分类"{
			catname = constMenuRootKey
		}

		newlist := []*Rominfo{}
		for _,v := range romlist{

			if catname == "全部游戏"{
				//关键字搜索
				if strings.Contains(v.Title, keyword) {
					newlist = append(newlist,v)
				}
			}else{
				if catname == v.Menu{
					if keyword != "" {
						//关键字搜索
						if strings.Contains(v.Title, keyword) {
							newlist = append(newlist,v)
						}
					}else{
						//非关键字搜索
						newlist = append(newlist,v)
					}
				}
			}
		}
		jsonRom, _ := json.Marshal(&newlist)

		return sciter.NewValue(string(jsonRom))
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