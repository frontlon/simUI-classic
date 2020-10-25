package modules

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"simUI/code/config"
	"simUI/code/db"
	"simUI/code/utils"
	"strings"
)

var ConstSeparator = "__"     //rom子分隔符
var ConstMenuRootKey = "_7b9" //根子目录游戏的Menu参数

type RomDetail struct {
	Info            *db.Rom         //rom信息
	DocContent      string          //简介内容
	StrategyContent string          //攻略内容
	StrategyFile    string          //攻略文件
	Sublist         []*db.Rom       //子游戏
	Simlist         []*db.Simulator //模拟器
	RomFileSize     string          //rom文件大小
}

//运行游戏
func RunGame(romId uint64, simId uint32) error {

	//数据库中读取rom详情
	rom, err := (&db.Rom{}).GetById(romId)
	if err != nil {
		return err
	}

	romCmd, _ := (&db.Rom{}).GetSimConf(romId, simId)

	sim := &db.Simulator{}
	if simId == 0 {
		sim = config.Cfg.Platform[rom.Platform].UseSim
		if sim == nil {
			return errors.New(config.Cfg.Lang["SimulatorNotFound"])
		}
	} else {
		if config.Cfg.Platform[rom.Platform].SimList == nil {
			return errors.New(config.Cfg.Lang["SimulatorNotFound"])
		}
		sim = config.Cfg.Platform[rom.Platform].SimList[simId]
	}

	//如果是相对路径，转换成绝对路径
	if !strings.Contains(rom.RomPath, ":") {
		rom.RomPath = config.Cfg.Platform[rom.Platform].RomPath + config.Cfg.Separator + rom.RomPath
	}

	//解压zip包
	if (sim.Unzip == 1 && romCmd.Unzip == 2) || romCmd.Unzip == 1 {
		RomExts := strings.Split(config.Cfg.Platform[rom.Platform].RomExts, ",")
		rom.RomPath, err = UnzipRom(rom.RomPath, RomExts)
		if err != nil {
			return err
		}
		if rom.RomPath == "" {
			return errors.New(config.Cfg.Lang["UnzipExeNotFound"])
		}

		//如果指定了执行文件
		if romCmd.File != "" {
			rom.RomPath = utils.GetFilePath(rom.RomPath) + "/" + romCmd.File
		}

	}

	//检测rom文件是否存在
	if utils.FileExists(rom.RomPath) == false {
		return errors.New(config.Cfg.Lang["RomNotFound"] + rom.RomPath)
	}

	//加载运行参数
	cmd := []string{}

	ext := utils.GetFileExt(rom.RomPath)

	//运行游戏前，先杀掉之前运行的程序
	if err = utils.KillGame(); err != nil {
		return err
	}

	simCmd := ""
	simLua := ""

	if romCmd.Cmd != "" {
		simCmd = romCmd.Cmd
	} else if sim.Cmd != "" {
		simCmd = sim.Cmd
	}

	if romCmd.Lua != "" {
		simLua = romCmd.Lua
	} else if sim.Cmd != "" {
		simLua = sim.Lua
	}


	//如果是可执行程序，则不依赖模拟器直接运行
	if utils.InSliceString(ext, config.RUN_EXTS) {
		//如果lua脚本存在，则运行lua脚本
		if simLua != "" {
			callLua("", rom.RomPath)
		} else {
			//运行游戏
			cmd = append(cmd, rom.RomPath)
			err = utils.RunGame("", cmd)
		}
	} else {
		//如果依赖模拟器

		//检测模拟器文件是否存在
		_, err = os.Stat(sim.Path)
		if err != nil {
			return errors.New(config.Cfg.Lang["SimulatorNotFound"])
		}



		if simCmd == "" {
			cmd = append(cmd, rom.RomPath)
		} else {
			//如果rom运行参数存在，则使用rom的参数
			cmd = strings.Split(simCmd, " ")
			filename := filepath.Base(rom.RomPath) //exe运行文件路径
			//替换变量
			for k, _ := range cmd {
				cmd[k] = strings.ReplaceAll(cmd[k], `{RomName}`, utils.GetFileName(filename))
				cmd[k] = strings.ReplaceAll(cmd[k], `{RomExt}`, utils.GetFileExt(filename))
				cmd[k] = strings.ReplaceAll(cmd[k], `{RomFullPath}`, rom.RomPath)
			}
		}

		//运行lua脚本
		if simLua != "" {
			cmdStr := utils.SlicetoString(" ", cmd)
			callLua(sim.Path, cmdStr)
		} else {
			//运行游戏
			err = utils.RunGame(sim.Path, cmd)
		}

	}
	return nil
}

//右键打开文件夹
func OpenFolder(id uint64, opt string, simId uint32) error {

	info, err := (&db.Rom{}).GetById(id)
	platform := config.Cfg.Platform[info.Platform] //读取当前平台信息
	if err != nil {
		return err
	}
	romName := utils.GetFileName(filepath.Base(info.RomPath)) //读取文件名 }
	fileName := ""
	isDir := false
	switch opt {
	case "rom":
		fileName = platform.RomPath + config.Cfg.Separator + info.RomPath
	case "thumb":
		if platform.ThumbPath != "" {
			for _, v := range config.PIC_EXTS {
				fileName = platform.ThumbPath + config.Cfg.Separator + romName + v

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
			for _, v := range config.PIC_EXTS {
				fileName = platform.SnapPath + config.Cfg.Separator + romName + v
				if utils.FileExists(fileName) {
					break
				} else {
					fileName = ""
				}
			}
			if fileName == "" {
				isDir = true
				fileName = platform.SnapPath + config.Cfg.Separator
			}
		}

	case "poster":
		if platform.PosterPath != "" {
			for _, v := range config.PIC_EXTS {
				fileName = platform.PosterPath + config.Cfg.Separator + romName + v
				if utils.FileExists(fileName) {
					break
				} else {
					fileName = ""
				}
			}
			if fileName == "" {
				isDir = true
				fileName = platform.PosterPath + config.Cfg.Separator
			}
		}
	case "packing":
		if platform.PackingPath != "" {
			for _, v := range config.PIC_EXTS {
				fileName = platform.PackingPath + config.Cfg.Separator + romName + v
				if utils.FileExists(fileName) {
					break
				} else {
					fileName = ""
				}
			}
			if fileName == "" {
				isDir = true
				fileName = platform.PackingPath + config.Cfg.Separator
			}
		}
	case "doc":
		if platform.DocPath != "" {
			for _, v := range config.DOC_EXTS {
				fileName = platform.DocPath + config.Cfg.Separator + romName + v
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
	case "title":
		if platform.DocPath != "" {
			for _, v := range config.PIC_EXTS {
				fileName = platform.TitlePath + config.Cfg.Separator + romName + v
				if utils.FileExists(fileName) {
					break
				} else {
					fileName = ""
				}
			}
			if fileName == "" {
				isDir = true
				fileName = platform.TitlePath
			}
		}
	case "background":
		if platform.BackgroundPath != "" {
			for _, v := range config.PIC_EXTS {
				fileName = platform.BackgroundPath + config.Cfg.Separator + romName + v
				if utils.FileExists(fileName) {
					break
				} else {
					fileName = ""
				}
			}
			if fileName == "" {
				isDir = true
				fileName = platform.BackgroundPath
			}
		}
	case "video":
		if platform.VideoPath != "" {
			for _, v := range config.PIC_EXTS {
				fileName = platform.VideoPath + config.Cfg.Separator + romName + v
				if utils.FileExists(fileName) {
					break
				} else {
					fileName = ""
				}
			}
			if fileName == "" {
				isDir = true
				fileName = platform.VideoPath
			}
		}
	case "strategy":
		if platform.StrategyPath != "" {
			for _, v := range config.DOC_EXTS {
				fileName = platform.StrategyPath + config.Cfg.Separator + romName + v
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
				return err
			}
		} else {
			if err := exec.Command(`explorer`, `/select,`, `/n,`, fileName).Start(); err != nil {
				return err
			}
		}
	} else {
		return errors.New(config.Cfg.Lang["PathNotFound"])
	}
	return nil
}

//读取rom详情
func GetGameDetail(id uint64) (*RomDetail, error) {

	res := &RomDetail{}
	//游戏游戏详细数据
	info, err := (&db.Rom{}).GetById(id)

	if err != nil {
		return res, err
	}
	//子游戏列表
	sub, _ := (&db.Rom{}).GetSubRom(info.Platform, info.Name)

	res.Info = info
	res.StrategyFile = ""
	res.Sublist = sub
	res.Simlist, _ = (&db.Simulator{}).GetByPlatform(info.Platform)

	//获取rom文件大小
	if res.Info.RomPath != "" {
		fi := config.Cfg.Platform[info.Platform].RomPath + config.Cfg.Separator + res.Info.RomPath
		f, err := os.Stat(fi)
		if err == nil {
			if f.Size() < 1024 {
				res.RomFileSize = fmt.Sprintf("%.2fB", float64(f.Size())/float64(1))
			} else if f.Size() < (1024 * 1024) {
				res.RomFileSize = fmt.Sprintf("%.2fKB", float64(f.Size())/float64(1024))
			} else if f.Size() < (1024 * 1024 * 1024) {
				res.RomFileSize = fmt.Sprintf("%.2fMB", float64(f.Size())/float64(1024*1024))
			} else if f.Size() < (1024 * 1024 * 1024 * 1024) {
				res.RomFileSize = fmt.Sprintf("%.2fGB", float64(f.Size())/float64(1024*1024*1024))
			}
		}
	}

	for k, v := range res.Simlist {
		if res.Simlist[k].Path != "" {
			res.Simlist[k].Path, _ = filepath.Abs(v.Path)
		}
	}

	//读取文档内容
	romName := utils.GetFileName(filepath.Base(info.RomPath)) //生成新文件的完整绝路路径地址
	if config.Cfg.Platform[info.Platform].DocPath != "" {
		docFileName := ""
		for _, v := range config.DOC_EXTS {
			docFileName = config.Cfg.Platform[info.Platform].DocPath + config.Cfg.Separator + romName + v
			res.DocContent = GetDocContent(docFileName)
			if res.DocContent != "" {
				break
			}
		}
	}

	if config.Cfg.Platform[info.Platform].StrategyPath != "" {
		//检测攻略可执行文件是否存在
		strategyFileName := ""
		for _, v := range config.RUN_EXTS {
			strategyFileName = config.Cfg.Platform[info.Platform].StrategyPath + config.Cfg.Separator + romName + v
			if utils.FileExists(strategyFileName) {
				res.StrategyFile = strategyFileName
				break
			}
		}

		//如果没有执行运行的文件，则读取文档内容
		if strategyFileName != "" {
			for _, v := range config.DOC_EXTS {
				strategyFileName = config.Cfg.Platform[info.Platform].StrategyPath + config.Cfg.Separator + romName + v
				res.StrategyContent = GetDocContent(strategyFileName)
				if res.StrategyContent != "" {
					break
				}
			}
		}

	}
	return res, nil
}

/**
 * 读取游戏介绍文本
 **/
func GetDocContent(f string) string {
	if f == "" {
		return ""
	}
	text, err := ioutil.ReadFile(f)
	content := ""
	if err != nil {
		return content
	}
	content = string(text)

	if !utils.IsUTF8(content) {
		content = utils.ToUTF8(content)
	}

	return content
}

//更新模拟器独立参数
func UpdateRomCmd(id uint64, simId uint32, data map[string]string) error {
	if data["cmd"] == "" && data["unzip"] == "2" {
		//如果当前配置和模拟器默认配置一样，则删除该记录
		if err := (&db.Rom{}).DelSimConf(id, simId); err != nil {
			return err
		}
	} else {
		//开始更新
		if err := (&db.Rom{}).UpdateSimConf(id, simId, data["cmd"], uint8(utils.ToInt(data["unzip"])), data["file"]); err != nil {
			return err
		}
	}
	return nil
}

//读取rom以及相关资源
func DeleteRomAndRes(id uint64) error {

	//游戏游戏详细数据
	info, err := (&db.Rom{}).GetById(id)
	if err != nil {
		return err
	}

	fname := utils.GetFileName(info.RomPath)
	platform := config.Cfg.Platform[info.Platform]

	go func() {
		romFiles, _ := utils.ScanDirByKeyword(platform.RomPath, fname)
		for _, f := range romFiles {
			utils.FileDelete(f)
		}
	}()
	go func() {
		thumbFiles, _ := utils.ScanDirByKeyword(platform.ThumbPath, fname)
		for _, f := range thumbFiles {
			utils.FileDelete(f)
		}
	}()
	go func() {
		backgroundFiles, _ := utils.ScanDirByKeyword(platform.BackgroundPath, fname)
		for _, f := range backgroundFiles {
			utils.FileDelete(f)
		}
	}()
	go func() {
		packingFiles, _ := utils.ScanDirByKeyword(platform.PackingPath, fname)
		for _, f := range packingFiles {
			utils.FileDelete(f)
		}
	}()
	go func() {
		posterFiles, _ := utils.ScanDirByKeyword(platform.PosterPath, fname)
		for _, f := range posterFiles {
			utils.FileDelete(f)
		}
	}()
	go func() {
		snapFiles, _ := utils.ScanDirByKeyword(platform.SnapPath, fname)
		for _, f := range snapFiles {
			utils.FileDelete(f)
		}
	}()
	go func() {
		titleFiles, _ := utils.ScanDirByKeyword(platform.TitlePath, fname)
		for _, f := range titleFiles {
			utils.FileDelete(f)
		}
	}()
	go func() {
		videoFiles, _ := utils.ScanDirByKeyword(platform.VideoPath, fname)
		for _, f := range videoFiles {
			utils.FileDelete(f)
		}
	}()
	go func() {
		docFiles, _ := utils.ScanDirByKeyword(platform.DocPath, fname)
		for _, f := range docFiles {
			utils.FileDelete(f)
		}
	}()
	go func() {
		strategyFiles, _ := utils.ScanDirByKeyword(platform.StrategyPath, fname)
		for _, f := range strategyFiles {
			utils.FileDelete(f)
		}
	}()
	return nil
}
