package controller

import (
	"VirtualNesGUI/code/db"
	"VirtualNesGUI/code/utils"
	"encoding/json"
	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

/**
 * 定义view用function
 **/

func RomController(w *window.Window) {

	//运行游戏
	w.DefineFunction("RunGame", func(args ...*sciter.Value) *sciter.Value {

		romId := uint64(utils.ToInt(args[0].String()))
		simId := uint32(utils.ToInt(args[1].String()))

		//数据库中读取rom详情
		rom, err := (&db.Rom{}).GetById(romId)
		if err != nil {
			WriteLog(Config.Lang["GameNotFound"])
			return ErrorMsg(w, Config.Lang["GameNotFound"])
		}

		romCmd, _ := (&db.RomCmd{RomId: romId, SimId: simId,}).Get()

		sim := &db.Simulator{}
		if simId == 0 {
			sim = Config.Platform[rom.Platform].UseSim
			if sim == nil {
				WriteLog(Config.Lang["SimulatorNotFound"])
				return ErrorMsg(w, Config.Lang["SimulatorNotFound"])
			}
		} else {
			if Config.Platform[rom.Platform].SimList == nil {
				WriteLog(Config.Lang["SimulatorNotFound"])
				return ErrorMsg(w, Config.Lang["SimulatorNotFound"])
			}
			sim = Config.Platform[rom.Platform].SimList[simId];
		}

		//检测执行文件是否存在
		_, err = os.Stat(sim.Path)
		if err != nil {
			WriteLog(err.Error())
			return ErrorMsg(w, err.Error())
		}

		//解压zip包
		rom.RomPath = Config.Platform[rom.Platform].RomPath + Config.Separator + rom.RomPath;
		if sim.Unzip == 1 || romCmd.Unzip == 1 {
			rom.RomPath, err = UnzipRom(rom.RomPath, Config.Platform[rom.Platform].RomExts)
			if err != nil {
				WriteLog(err.Error())
				return ErrorMsg(w, err.Error())
			}
			if rom.RomPath == "" {
				return ErrorMsg(w, Config.Lang["UnzipExeNotFound"])

			}
		}

		//检测rom文件是否存在
		if utils.FileExists(rom.RomPath) == false {
			WriteLog(Config.Lang["RomNotFound"] + rom.RomPath)
			return ErrorMsg(w, Config.Lang["RomNotFound"]+rom.RomPath)
		}

		//加载运行参数
		cmd := []string{}

		ext := utils.GetFileExt(rom.RomPath)

		//如果rom运行参数存在，则使用rom的参数
		if romCmd.Cmd != "" {
			sim.Cmd = romCmd.Cmd
		}

		//运行游戏前，先杀掉之前运行的程序
		if err = killGame(); err != nil {
			WriteLog(err.Error())
			return ErrorMsg(w, err.Error())
		}
		//如果是可执行程序，则不依赖模拟器直接运行
		if utils.InSliceString(ext, RUN_EXTS) {
			cmd = append(cmd, rom.RomPath)
			err = runGame("explorer", cmd)
		} else {
			//如果依赖模拟器
			if sim.Cmd == "" {
				cmd = append(cmd, rom.RomPath)
			} else {
				cmd = strings.Split(sim.Cmd, " ")
				filename := filepath.Base(rom.RomPath) //exe运行文件路径
				//替换变量
				for k, _ := range cmd {
					cmd[k] = strings.ReplaceAll(cmd[k], `{RomName}`, utils.GetFileName(filename))
					cmd[k] = strings.ReplaceAll(cmd[k], `{RomExt}`, utils.GetFileExt(filename))
					cmd[k] = strings.ReplaceAll(cmd[k], `{RomFullPath}`, rom.RomPath)
				}
			}
			err = runGame(sim.Path, cmd)
		}

		if err != nil {
			WriteLog(err.Error())
			return ErrorMsg(w, err.Error())
		}
		return sciter.NullValue()
	})

	//运行攻略文件
	w.DefineFunction("RunStrategy", func(args ...*sciter.Value) *sciter.Value {
		f := args[0].String()
		if (f == "") {
			return sciter.NullValue()
		}
		if err := runGame("explorer", []string{f}); err != nil {
			WriteLog(err.Error())
			return ErrorMsg(w, err.Error())
		}
		return sciter.NullValue()

	})

	//打开rom目录
	w.DefineFunction("OpenFolder", func(args ...*sciter.Value) *sciter.Value {

		id := uint64(utils.ToInt(args[0].String()))
		opt := args[1].String()
		simId := uint32(utils.ToInt(args[2].String()))
		info, err := (&db.Rom{}).GetById(id)
		platform := Config.Platform[info.Platform] //读取当前平台信息
		if err != nil {
			WriteLog(err.Error())
			return ErrorMsg(w, err.Error())
		}
		romName := utils.GetFileName(filepath.Base(info.RomPath)) //读取文件名 }
		fileName := ""
		isDir := false
		switch opt {
		case "rom":
			fileName = platform.RomPath + Config.Separator + info.RomPath
		case "thumb":
			if platform.ThumbPath != "" {
				for _, v := range PIC_EXTS {
					fileName = platform.ThumbPath + Config.Separator + romName + v

					if utils.FileExists(fileName) {
						break
					} else {
						fileName = ""
					}
				}

				if fileName == "" {
					isDir = true
					fileName = platform.ThumbPath
				}

			}
		case "snap":
			if platform.SnapPath != "" {
				for _, v := range PIC_EXTS {
					fileName = platform.SnapPath + Config.Separator + romName + v
					if utils.FileExists(fileName) {
						break
					} else {
						fileName = ""
					}
				}
				if fileName == "" {
					isDir = true
					fileName = platform.SnapPath + Config.Separator
				}
			}


		case "poster":
			if platform.PosterPath != "" {
				for _, v := range PIC_EXTS {
					fileName = platform.PosterPath + Config.Separator + romName + v
					if utils.FileExists(fileName) {
						break
					} else {
						fileName = ""
					}
				}
				if fileName == "" {
					isDir = true
					fileName = platform.PosterPath + Config.Separator
				}
			}
		case "packing":
			if platform.PackingPath != "" {
				for _, v := range PIC_EXTS {
					fileName = platform.PackingPath + Config.Separator + romName + v
					if utils.FileExists(fileName) {
						break
					} else {
						fileName = ""
					}
				}
				if fileName == "" {
					isDir = true
					fileName = platform.PackingPath + Config.Separator
				}
			}
		case "doc":
			if platform.DocPath != "" {
				for _, v := range DOC_EXTS {
					fileName = platform.DocPath + Config.Separator + romName + v
					if utils.FileExists(fileName) {
						break
					} else {
						fileName = ""
					}
				}
				if fileName == "" {
					isDir = true
					fileName = platform.DocPath
				}
			}
		case "strategy":
			if platform.StrategyPath != "" {
				for _, v := range DOC_EXTS {
					fileName = platform.StrategyPath + Config.Separator + romName + v
					if utils.FileExists(fileName) {
						break
					} else {
						fileName = ""
					}
				}
				if fileName == "" {
					isDir = true
					fileName = platform.StrategyPath
				}
			}
		case "sim":
			if _, ok := platform.SimList[simId]; ok {
				fileName = platform.SimList[simId].Path
			}
		}

		if fileName != "" {
			if isDir == true {
				if err := exec.Command(`explorer`, fileName).Start(); err != nil {
					WriteLog(err.Error())
					return ErrorMsg(w, err.Error())
				}
			} else {
				if err := exec.Command(`explorer`, `/select,`, `/n,`, fileName).Start(); err != nil {
					WriteLog(err.Error())
					return ErrorMsg(w, err.Error())
				}
			}
		} else {
			WriteLog(Config.Lang["PathNotFound"])
			return ErrorMsg(w, Config.Lang["PathNotFound"])
		}

		return sciter.NullValue()
	})

	//读取游戏列表
	w.DefineFunction("GetGameList", func(args ...*sciter.Value) *sciter.Value {
		platform := uint32(utils.ToInt(args[0].String()))        //平台
		catname := strings.Trim(args[1].String(), " ")           //分类
		keyword := strings.Trim(args[2].String(), " ")           //关键字
		num := strings.Trim(args[3].String(), " ")               //字母索引
		page := utils.ToInt(strings.Trim(args[4].String(), " ")) //分页数

		newlist := []*db.Rom{}
		if num == "" {
			newlist, _ = (&db.Rom{}).Get(page, platform, catname, keyword)
		} else {
			//按拼音查询
			newlist, _ = (&db.Rom{}).GetByPinyin(page, platform, catname, num)
		}

		jsonRom, _ := json.Marshal(newlist)
		return sciter.NewValue(string(jsonRom))
	})

	//读取游戏数量
	w.DefineFunction("GetGameCount", func(args ...*sciter.Value) *sciter.Value {
		platform := uint32(utils.ToInt(args[0].String()))
		catname := strings.Trim(args[1].String(), " ")
		keyword := strings.Trim(args[2].String(), " ")
		count, _ := (&db.Rom{}).Count(platform, catname, keyword)
		return sciter.NewValue(utils.ToString(count))
	})

	//读取rom详情
	w.DefineFunction("GetGameDetail", func(args ...*sciter.Value) *sciter.Value {
		id := uint64(utils.ToInt(args[0].String()))
		res := &RomDetail{}
		//游戏游戏详细数据
		info, err := (&db.Rom{}).GetById(id)
		if err != nil {
			WriteLog(err.Error())
			return ErrorMsg(w, err.Error())
		}
		//子游戏列表
		sub, _ := (&db.Rom{}).GetSubRom(info.Platform, info.Name)

		res.Info = info
		res.StrategyFile = ""
		res.Sublist = sub

		//读取文档内容
		romName := utils.GetFileName(filepath.Base(info.RomPath)) //生成新文件的完整绝路路径地址
		if Config.Platform[info.Platform].DocPath != "" {
			docFileName := "";
			for _, v := range DOC_EXTS {
				docFileName = Config.Platform[info.Platform].DocPath + Config.Separator + romName + v
				res.DocContent = getDocContent(docFileName)
				if res.DocContent != "" {
					break
				}
			}
		}

		if Config.Platform[info.Platform].StrategyPath != "" {
			//检测攻略可执行文件是否存在
			strategyFileName := "";
			for _, v := range RUN_EXTS {
				strategyFileName = Config.Platform[info.Platform].StrategyPath + Config.Separator + romName + v
				if utils.FileExists(strategyFileName) {
					res.StrategyFile = strategyFileName
					break
				}
			}

			//如果没有执行运行的文件，则读取文档内容
			if strategyFileName != "" {
				for _, v := range DOC_EXTS {
					strategyFileName = Config.Platform[info.Platform].StrategyPath + Config.Separator + romName + v
					res.StrategyContent = getDocContent(strategyFileName)
					if res.StrategyContent != "" {
						break
					}
				}
			}

		}

		jsonMenu, _ := json.Marshal(&res)
		return sciter.NewValue(string(jsonMenu))
	})

	//设为我的最爱
	w.DefineFunction("SetFavorite", func(args ...*sciter.Value) *sciter.Value {
		id := uint64(utils.ToInt(args[0].String()))
		star := uint8(utils.ToInt(args[1].String()))

		//更新rom表
		rom := &db.Rom{
			Id:   id,
			Star: star,
		}

		//更新数据
		if err := rom.UpdateStar(); err != nil {
			WriteLog(err.Error())
			return ErrorMsg(w, err.Error())
		}

		return sciter.NewValue("1")
	})

	//更新rom图片
	w.DefineFunction("UpdateRomThumbs", func(args ...*sciter.Value) *sciter.Value {
		typeid := utils.ToInt(args[0].String())
		id := uint64(utils.ToInt(args[1].String()))
		newpath := args[2].String()

		rom := &db.Rom{
			Id: id,
		}

		//设定新的文件名
		vo, err := rom.GetById(id)
		if err != nil {
			WriteLog(err.Error())
			return ErrorMsg(w, err.Error())
		}

		//下载文件
		res, err := http.Get(newpath)
		if err != nil {
			WriteLog(err.Error())
			return ErrorMsg(w, err.Error())
		}

		//下载成功后，备份原文件
		platformPath := ""
		//原图存在，则备份
		if typeid == 1 {
			platformPath = Config.Platform[vo.Platform].ThumbPath
		} else {
			platformPath = Config.Platform[vo.Platform].SnapPath
		}

		if platformPath == "" {
			WriteLog(Config.Lang["NoSetThumbDir"])
			return ErrorMsg(w, Config.Lang["NoSetThumbDir"])
		}

		//开始备份原图
		bakFolder := Config.CachePath + "thumb_bak/"
		RomFileName := utils.GetFileName(vo.RomPath)

		//检测bak文件夹是否存在，不存在这创建bak目录
		folder := utils.FolderExists(bakFolder)
		if folder == false {
			_ = os.Mkdir(bakFolder, os.ModePerm);
		}
		for _, ext := range PIC_EXTS {
			oldFileName := RomFileName + ext //老图片文件名
			if utils.FileExists(oldFileName) {
				bakFileName := RomFileName + "_" + utils.ToString(time.Now().Unix()) + ext //生成备份文件名
				err := os.Rename(oldFileName, bakFolder+bakFileName)                       //移动文件
				if err != nil {
					WriteLog(err.Error())
					return ErrorMsg(w, err.Error())
				}
			}
		}

		//生成新文件
		platformPathAbs, err := filepath.Abs(platformPath) //读取平台图片路径

		newFileName := platformPathAbs + Config.Separator + RomFileName + utils.GetFileExt(newpath) //生成新文件的完整绝路路径地址
		f, err := os.Create(newFileName)
		if err != nil {
			WriteLog(err.Error())
			return ErrorMsg(w, err.Error())
		}
		io.Copy(f, res.Body)

		return sciter.NewValue(newFileName)
	})


	//rom翻页
	w.DefineFunction("scrollLoadRom", func(args ...*sciter.Value) *sciter.Value {
		scrollPos := args[0].String()
		go func(scrollPos string) *sciter.Value {
			//数据更新完成后，页面回调，更新页面DOM
			if _, err := w.Call("scrollLoadRom",sciter.NewValue(scrollPos)); err != nil {
			}
			return sciter.NullValue()
		}(scrollPos)

		return sciter.NullValue()
	})

}
