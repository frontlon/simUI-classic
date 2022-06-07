package modules

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"simUI/code/compoments"
	"simUI/code/config"
	"simUI/code/db"
	"simUI/code/utils"
	"strings"
	"time"
)

//运行游戏
func RunGame(romId uint64, simId uint32) error {

	//数据库中读取rom详情
	rom, err := (&db.Rom{}).GetById(romId)
	if err != nil {
		return err
	}

	//如果是相对路径，转换成绝对路径
	if !utils.IsAbsPath(rom.RomPath) {
		rom.RomPath = config.Cfg.Platform[rom.Platform].RomPath + config.Cfg.Separator + rom.RomPath
	}

	//检测rom文件是否存在
	if utils.FileExists(rom.RomPath) == false {
		return errors.New(config.Cfg.Lang["RomNotFound"] + rom.RomPath)
	}

	//读取父游戏信息
	var parent = &db.Rom{}
	if rom.Pname != "" {
		parent, _ = (&db.Rom{}).GetByFileMd5(rom.Pname)
	}

	if !utils.IsAbsPath(parent.RomPath) {
		parent.RomPath = config.Cfg.Platform[rom.Platform].RomPath + config.Cfg.Separator + parent.RomPath
	}

	ext := strings.ToLower(utils.GetFileExt(rom.RomPath))

	//记录运行信息
	rid := romId
	fileMd5 := rom.FileMd5
	if parent.Id != 0 {
		rid = parent.Id
		fileMd5 = parent.FileMd5
	}
	_ = (&db.Rom{}).UpdateRunNumAndTime(rid)
	_ = (&db.RomSetting{
		Platform:    rom.Platform,
		FileMd5:     fileMd5,
		RunLasttime: time.Now().Unix(),
	}).UpdateRunNumAndTime()

	//运行游戏
	if utils.InSliceString(ext, config.RUN_EXTS) {
		//直接运行exe
		err = runGameExe(rom)
	} else if utils.InSliceString(ext, config.EXPLORER_EXTS) {
		//通过explorer运行
		err = runGameExplorer(rom)
	} else {
		//依赖模拟器
		err = runGameSimulator(rom, parent, simId)
	}

	if err != nil {
		return err
	}
	return nil
}

/**
 * 直接运行exe
 **/
func runGameExe(rom *db.Rom) error {
	cmd := []string{}
	return compoments.RunGame(rom.RomPath, cmd)
}

/**
 * 通过explorer运行
 **/
func runGameExplorer(rom *db.Rom) error {
	cmd := []string{rom.RomPath}
	return compoments.RunGame("", cmd)
}

/**
 * 通过模拟器运行游戏
 **/
func runGameSimulator(rom *db.Rom, parent *db.Rom, simId uint32) error {

	romId := rom.Id
	romCmd := &db.SimConf{}
	if parent.Id != 0 {
		romCmd, _ = (&db.Rom{}).GetSimConf(parent.Id, simId)
	} else {
		romCmd, _ = (&db.Rom{}).GetSimConf(romId, simId)
	}

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

	//解压zip包
	err := errors.New("")
	if (sim.Unzip == 1 && romCmd.Unzip == 0) || romCmd.Unzip == 1 {
		RomExts := strings.Split(config.Cfg.Platform[rom.Platform].RomExts, ",")
		rom.RomPath, err = compoments.UnzipRom(rom.RomPath, RomExts)

		if err != nil {
			return err
		}
		if rom.RomPath == "" {
			return errors.New(config.Cfg.Lang["UnzipExeNotFound"])
		}

		//如果指定了执行文件
		if romCmd.File != "" {
			rom.RomPath = utils.GetFileAbsPath(rom.RomPath) + "/" + romCmd.File
		}

	}

	//加载运行参数
	cmd := []string{}

	//运行游戏前，先杀掉之前运行的程序
	if err = compoments.KillGame(); err != nil {
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
	} else if sim.Lua != "" {
		simLua = sim.Lua
	}

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
			cmd[k] = strings.ReplaceAll(cmd[k], `{RomMainFullpath}`, parent.RomPath)
			cmd[k] = strings.ReplaceAll(cmd[k], `{RomMainName}`, utils.GetFileName(parent.RomPath))
		}
	}

	//运行lua脚本
	if simLua != "" {
		cmdStr := utils.SlicetoString(" ", cmd)
		compoments.CallLua(sim.Lua, sim.Path, cmdStr)
	}

	//模拟器运行游戏
	if err := compoments.RunGame(sim.Path, cmd); err != nil {
		return err
	}

	//数据更新完成后，页面回调，更新页面DOM
	if _, err := utils.Window.Call("CB_runGame"); err != nil {
		fmt.Println(err)
	}

	return nil
}
