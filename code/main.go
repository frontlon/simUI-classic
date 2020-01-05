package main

import (
	"VirtualNesGUI/code/db"
	"VirtualNesGUI/code/utils"
	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
	"os"
	"path/filepath"
	"strings"
)

var separator = string(os.PathSeparator) //系统路径分隔符
//路径分隔符
var constMenuRootKey = "_7b9"                                                //根子目录游戏的Menu参数
var constMainFile = "D:\\work\\go\\src\\VirtualNesGUI\\code\\view\\main.html" //主文件路径（测试用）
//var constMainFile = "this://app/main.html" //主文件路径（正式）

func main() {

	defer func() {
		if r := recover(); r != nil {
			WriteLog(utils.ToString(r))
		}
	}()

	//连接数据库
	db.Conn()

	//初始化配置
	Config = &ConfStruct{}
	var rootpath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	Config.RootPath = rootpath + separator //当前软件的绝对路径
	Config.Separator = separator //系统的目录分隔符
	Config.CachePath = rootpath + separator + "cache" + separator//缓存路径
	Config.UnzipPath = Config.CachePath + "unzip" + separator//rom解压路径

	errConf := InitConf()

	width := Config.Default.WindowWidth
	height := Config.Default.WindowHeight

	//创建window窗口
	w, err := window.New(
		sciter.SW_MAIN|
			//sciter.SW_RESIZEABLE|
			//sciter.SW_CONTROLS|
			sciter.SW_ENABLE_DEBUG,
		&sciter.Rect{Left: 0, Top: 0, Right: int32(width), Bottom: int32(height)});
	if err != nil {
		WriteLog(err.Error())
	}

	//设置view权限
	w.SetOption(sciter.SCITER_SET_SCRIPT_RUNTIME_FEATURES, sciter.ALLOW_SYSINFO | sciter.ALLOW_FILE_IO |sciter.ALLOW_SOCKET_IO);

	//设置回调
	w.SetCallback(newHandler(w.Sciter))

	//解析资源
	w.OpenArchive(res)

	//加载文件
	err = w.LoadFile(constMainFile);
	if err != nil {
		errorMsg(w, err.Error())
		return
	}

	//配置出先错误
	if errConf != nil{
		errorMsg(w, errConf.Error())
		os.Exit(1)
		return
	}

	if len(Config.Lang) == 0{
		WriteLog("没有找到语言文件或语言文件为空\nNo language files or language files found empty")
		errorMsg(w, "没有找到语言文件或语言文件为空\nNo language files or language files found empty")
		os.Exit(1)
		return
	}

	//设置标题
	w.SetTitle(Config.Lang["SoftName"]);
	//定义view函数
	defineViewFunction(w)
	//显示窗口
	w.Show();
	//运行窗口，进入消息循环
	w.Run();
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

//调用alert框
func errorMsg(w *window.Window,err string) *sciter.Value {

	if _, err := w.Call("errorBox", sciter.NewValue(err)); err != nil {
	}
	return sciter.NullValue();
}
