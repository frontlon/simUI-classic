package modules

import (
	"errors"
	"os"
	"path/filepath"
	"simUI/code/components"
	"simUI/code/config"
	"simUI/code/db"
	"simUI/code/utils"
	"strings"
	"time"
)

// 运行游戏
func RunGame(romId uint64, simId uint32) error {
	//数据库中读取rom详情
	romInfo, err := (&db.Rom{}).GetById(romId)
	if err != nil {
		return err
	}

	//如果是相对路径，转换成绝对路径
	if !utils.IsAbsPath(romInfo.RomPath) {
		romInfo.RomPath = config.Cfg.Platform[romInfo.Platform].RomPath + config.Cfg.Separator + romInfo.RomPath
	}

	//检测rom文件是否存在
	if utils.FileExists(romInfo.RomPath) == false {
		return errors.New(config.Cfg.Lang["RomNotFound"] + romInfo.RomPath)
	}

	ext := strings.ToLower(utils.GetFileExt(romInfo.RomPath))

	//检查path文件
	cmd := []string{}
	pfile := ""
	if ext == ".path" {
		pfile, cmd = components.GetPathFile(romInfo.RomPath)
		ext = strings.ToLower(utils.GetFileExt(pfile))
	} else if ext == ".slnk" {
		pfile, cmd = components.GetSlnkFile(romInfo.RomPath)
		ext = strings.ToLower(utils.GetFileExt(pfile))
	}

	//运行游戏
	if utils.InSliceString(ext, config.RUN_EXTS) {
		//直接运行exe
		if pfile != "" {
			romInfo.RomPath = pfile
		}
		err = runGameExe(romInfo.RomPath, cmd)
	} else if utils.InSliceString(ext, config.EXPLORER_EXTS) {
		//依赖 explorer 启动
		if pfile != "" {
			romInfo.RomPath = pfile
		}
		err = runGameExplorer(romInfo.RomPath)
	} else {
		//依赖模拟器
		if pfile != "" {
			cmd = append([]string{pfile}, cmd...)
		}
		err = runGameSimulator(romInfo, cmd, simId)
	}

	if err != nil {
		return err
	}

	//记录运行次数和时间
	rid := romId
	fileMd5 := romInfo.FileMd5
	if romInfo.Pname != "" {
		parent, _ := (&db.Rom{}).GetByFileMd5(romInfo.Pname)
		if parent.Id != 0 {
			rid = parent.Id
			fileMd5 = parent.FileMd5
		}
	}
	_ = (&db.Rom{}).UpdateRunNumAndTime(rid)
	_ = (&db.RomSetting{
		Platform:    romInfo.Platform,
		FileMd5:     fileMd5,
		RunLasttime: time.Now().Unix(),
	}).UpdateRunNumAndTime()

	return nil
}

/**
 * 直接运行exe
 **/
func runGameExe(romPath string, cmd []string) error {
	return components.RunGame(romPath, cmd)
}

/**
 * 通过explorer运行
 **/
func runGameExplorer(romPath string) error {
	cmd := []string{romPath}
	return components.RunGame("", cmd)
}

/**
 * 通过模拟器运行游戏
 **/
func runGameSimulator(rom *db.Rom, cmd []string, simId uint32) error {

	//读取父游戏信息
	var parent = &db.Rom{}
	if rom.Pname != "" {
		parent, _ = (&db.Rom{}).GetByFileMd5(rom.Pname)
	}

	if !utils.IsAbsPath(parent.RomPath) {
		parent.RomPath = config.Cfg.Platform[rom.Platform].RomPath + config.Cfg.Separator + parent.RomPath
	}

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
		if config.Cfg.Platform[rom.Platform].SimList == nil || len(config.Cfg.Platform[rom.Platform].SimList) == 0 {
			return errors.New(config.Cfg.Lang["SimulatorNotFound"])
		}
		sim = config.Cfg.Platform[rom.Platform].SimList[simId]

		if sim == nil {
			sim = config.Cfg.Platform[rom.Platform].UseSim
		}
	}

	//检测模拟器文件是否存在
	_, err := os.Stat(sim.Path)
	if err != nil {
		return errors.New(config.Cfg.Lang["SimulatorNotFound"])
	}

	//解压后运行 - 解压zip包
	if (sim.Unzip == 1 && romCmd.Unzip == 0) || romCmd.Unzip == 1 {
		RomExts := strings.Split(config.Cfg.Platform[rom.Platform].RomExts, ",")
		rom.RomPath, err = components.UnzipRom(rom.RomPath, RomExts)

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

	if len(cmd) == 0 {
		simCmd := ""
		if romCmd.Cmd != "" {
			simCmd = romCmd.Cmd
		} else if sim.Cmd != "" {
			simCmd = sim.Cmd
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
				cmd[k] = strings.ReplaceAll(cmd[k], `{RomAlias}`, rom.Name)
				cmd[k] = strings.ReplaceAll(cmd[k], `{RomExt}`, utils.GetFileExt(filename))
				cmd[k] = strings.ReplaceAll(cmd[k], `{RomFullPath}`, rom.RomPath)
				cmd[k] = strings.ReplaceAll(cmd[k], `{RomMainFullpath}`, parent.RomPath)
				cmd[k] = strings.ReplaceAll(cmd[k], `{RomMainName}`, utils.GetFileName(parent.RomPath))

			}
		}
	}

	//运行游戏前，先杀掉之前运行的程序
	if err = components.KillGame(); err != nil {
		return err
	}

	//模拟器运行游戏
	if err = components.RunGame(sim.Path, cmd); err != nil {
		return err
	}

	return nil
}
