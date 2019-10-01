package main

import (
	"VirtualNesGUI/code/db"
	"github.com/axgle/mahonia"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

var constSeparator = "__" //rom子分隔符

type RomDetail struct {
	Info       *db.Rom
	DocContent string
	Sublist    []*db.Rom
}

/**
 * 读取游戏介绍文本
 **/
func getDocContent(f string) string {
	if f == "" {
		return ""
	}
	text, err := ioutil.ReadFile(f)
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
func runGame(platform int64, romfile string, sim int64) string {
	exeFile := Config.Platform[platform].UseSim.Path
	if sim != 0 {
		exeFile = Config.Platform[platform].SimList[sim].Path;
	}
	//检测执行文件是否存在
	_, err := os.Stat(exeFile)
	if err != nil {
		return err.Error()
	}

	//检测rom文件是否存在
	if Exists(romfile) == false {
		return Config.Lang["RomNotFound"] + romfile
	}

	//验证桥接程序是否存在
	bridge := filepath.Dir(exeFile) + separator + "tplugin.exe"
	cmd := &exec.Cmd{}
	_, ok := os.Stat(bridge)
	if ok == nil {
		cmd = exec.Command(bridge, exeFile, romfile)
	} else {
		cmd = exec.Command(exeFile, romfile)
	}

	if err := cmd.Start(); err != nil {
		return err.Error()
	}

	return ""
}



/**
 * 创建缓存
 **/
func CreateRomCache(platform int64) error {
	romlist := []*db.Rom{}
	menuList := map[string]*db.Menu{}              //分类目录
	RomPath := Config.Platform[platform].RomPath   //rom文件路径
	RomExt := Config.Platform[platform].RomExts    //rom扩展名
	ThumbList := GetMaterialUrl("thumb", platform) //缩略图
	VideoList := GetMaterialUrl("video", platform) //视频
	SnapList := GetMaterialUrl("snap", platform)   //截图
	DocList := GetMaterialUrl("doc", platform)     //文档
	RomAlias := getRomAlias(platform)              //别名配置

	//进入循环，遍历文件
	if err := filepath.Walk(RomPath,
		func(p string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}

			//整理目录格式，并转换为数组
			newpath := strings.Replace(RomPath, "/", "\\", -1)
			newpath = strings.Replace(p, newpath, "", -1)
			if len(newpath) > 0 && newpath[0:1] == "\\" {
				newpath = strings.Replace(newpath, "\\", "", 1)
			}
			subpath := strings.Split(newpath, "\\")
			romExt := strings.ToLower(path.Ext(p)) //获取文件后缀

			//如果该文件是游戏rom
			if f.IsDir() == false && CheckPlatformExt(RomExt, romExt) {
				romName := GetFileName(f.Name())
				title := romName
				//如果有别名配置，则读取别名
				if _, ok := RomAlias[title]; ok {
					title = RomAlias[title]
				}

				py := TextToPinyin(title)

				//如果游戏名称存在分隔符，说明是子游戏
				menu := constMenuRootKey //无目录，读取默认参数
				//定义目录，如果有子目录，则记录子目录名称
				if len(subpath) > 1 {
					menu = subpath[0]
				}

				thumb := ""
				snap := ""
				video := ""
				doc := ""

				if _, ok := ThumbList[romName]; ok {
					thumb = ThumbList[romName]
				}

				if _, ok := SnapList[romName]; ok {
					snap = SnapList[romName]
				}

				if _, ok := VideoList[romName]; ok {
					video = VideoList[romName]
				}

				if _, ok := DocList[romName]; ok {
					doc = DocList[romName]
				}

				//如果游戏名称存在分隔符，说明是子游戏
				if strings.Contains(title, constSeparator) {

					//拆分文件名
					sub := strings.Split(title, constSeparator)

					//去掉扩展名，生成标题
					sinfo := &db.Rom{
						Name:      sub[1],
						Pname:     sub[0],
						RomPath:   p,
						Menu:      menu,
						Platform:  platform,
						ThumbPath: thumb,
						SnapPath:  snap,
						VideoPath: video,
						DocPath:   doc,
						Star:      0,
						Pinyin:    py,
					}
					romlist = append(romlist, sinfo)
				} else {

					//去掉扩展名，生成标题
					rinfo := &db.Rom{
						Menu:      menu,
						Name:      title,
						Platform:  platform,
						RomPath:   p,
						ThumbPath: thumb,
						SnapPath:  snap,
						VideoPath: video,
						DocPath:   doc,
						Star:      0,
						Pinyin:    py,
					}

					//rom列表
					romlist = append(romlist, rinfo)
					//分类列表
					if menu != constMenuRootKey {
						menuList[menu] = &db.Menu{
							Platform: platform,
							Name:     menu,
							Pinyin:   TextToPinyin(menu),
						}
					}

				}
			}
			return nil
		}); err != nil {
	}

	//保存数据到数据库rom表
	if len(romlist) > 0 {
		if err := (&db.Rom{}).Add(&romlist); err != nil {
		}
	}
	//保存数据到数据库cate表
	if len(menuList) > 0 {

		if err := (&db.Menu{}).Add(&menuList); err != nil {
		}
	}

	//写入完成后清理变量
	romlist = []*db.Rom{}
	menuList = make(map[string]*db.Menu)

	return nil
}

//读取资源文件url
func GetMaterialUrl(stype string, platform int64) map[string]string {
	getpath := ""
	exts := []string{}
	list := make(map[string]string)
	switch stype {
	case "video":
		getpath = Config.Platform[platform].VideoPath;
		exts = []string{".gif"}
	case "thumb":
		getpath = Config.Platform[platform].ThumbPath;
		exts = []string{".jpg", ".bmp", ".png", ".jpeg"}
	case "snap":
		getpath = Config.Platform[platform].SnapPath;
		exts = []string{".jpg", ".bmp", ".png", ".jpeg"}
	case "doc":
		getpath = Config.Platform[platform].DocPath;
		exts = []string{".txt"}
	}

	//如果参数为空，不向下执行
	if getpath == "" || len(exts) == 0 {
		return list
	}

	//进入循环，遍历文件
	if err := filepath.Walk(getpath,
		func(p string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}
			romExt := strings.ToLower(path.Ext(p)) //获取文件后缀
			//如果是规定的扩展名，则记录数据
			if f.IsDir() == false && CheckPlatformExt(exts, romExt) {
				list[GetFileName(f.Name())] = p
			}
			return nil
		}); err != nil {
	}
	return list
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