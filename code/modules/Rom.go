package modules

import (
	"VirtualNesGUI/code/config"
	"VirtualNesGUI/code/db"
	"VirtualNesGUI/code/utils"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var ConstSeparator = "__"     //rom子分隔符
var ConstMenuRootKey = "_7b9" //根子目录游戏的Menu参数

type RomDetail struct {
	Info            *db.Rom   //基础信息
	DocContent      string    //简介内容
	StrategyContent string    //攻略内容
	StrategyFile    string    //攻略文件
	Sublist         []*db.Rom //子游戏
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
		sim = config.C.Platform[rom.Platform].UseSim
		if sim == nil {
			return errors.New(config.C.Lang["SimulatorNotFound"])
		}
	} else {
		if config.C.Platform[rom.Platform].SimList == nil {
			return errors.New(config.C.Lang["SimulatorNotFound"])
		}
		sim = config.C.Platform[rom.Platform].SimList[simId];
	}

	//检测执行文件是否存在
	_, err = os.Stat(sim.Path)
	if err != nil {
		return err
	}

	//如果是相对路径，转换成绝对路径
	if !strings.Contains(rom.RomPath, ":") {
		rom.RomPath = config.C.Platform[rom.Platform].RomPath + config.C.Separator + rom.RomPath;
	}

	//解压zip包
	if (sim.Unzip == 1 && romCmd.Unzip == 2) || romCmd.Unzip == 1 {
		RomExts := strings.Split(config.C.Platform[rom.Platform].RomExts, ",")
		rom.RomPath, err = UnzipRom(rom.RomPath, RomExts)
		if err != nil {
			return err
		}
		if rom.RomPath == "" {
			return errors.New(config.C.Lang["UnzipExeNotFound"])
		}
	}

	//检测rom文件是否存在
	if utils.FileExists(rom.RomPath) == false {
		return errors.New(config.C.Lang["RomNotFound"] + rom.RomPath)
	}

	//加载运行参数
	cmd := []string{}

	ext := utils.GetFileExt(rom.RomPath)

	//运行游戏前，先杀掉之前运行的程序
	if err = utils.KillGame(); err != nil {
		return err
	}
	//如果是可执行程序，则不依赖模拟器直接运行
	if utils.InSliceString(ext, config.RUN_EXTS) {
		cmd = append(cmd, rom.RomPath)
		err = utils.RunGame("explorer", cmd)
	} else {
		//如果依赖模拟器

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
				cmd[k] = strings.ReplaceAll(cmd[k], `{RomExt}`, utils.GetFileExt(filename))
				cmd[k] = strings.ReplaceAll(cmd[k], `{RomFullPath}`, rom.RomPath)
			}
		}
		err = utils.RunGame(sim.Path, cmd)
	}
	return nil
}

//右键打开文件夹
func OpenFolder(id uint64, opt string, simId uint32) error {

	info, err := (&db.Rom{}).GetById(id)
	platform := config.C.Platform[info.Platform] //读取当前平台信息
	if err != nil {
		return err
	}
	romName := utils.GetFileName(filepath.Base(info.RomPath)) //读取文件名 }
	fileName := ""
	isDir := false
	switch opt {
	case "rom":
		fileName = platform.RomPath + config.C.Separator + info.RomPath
	case "thumb":
		if platform.ThumbPath != "" {
			for _, v := range config.PIC_EXTS {
				fileName = platform.ThumbPath + config.C.Separator + romName + v

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
				fileName = platform.SnapPath + config.C.Separator + romName + v
				if utils.FileExists(fileName) {
					break
				} else {
					fileName = ""
				}
			}
			if fileName == "" {
				isDir = true
				fileName = platform.SnapPath + config.C.Separator
			}
		}

	case "poster":
		if platform.PosterPath != "" {
			for _, v := range config.PIC_EXTS {
				fileName = platform.PosterPath + config.C.Separator + romName + v
				if utils.FileExists(fileName) {
					break
				} else {
					fileName = ""
				}
			}
			if fileName == "" {
				isDir = true
				fileName = platform.PosterPath + config.C.Separator
			}
		}
	case "packing":
		if platform.PackingPath != "" {
			for _, v := range config.PIC_EXTS {
				fileName = platform.PackingPath + config.C.Separator + romName + v
				if utils.FileExists(fileName) {
					break
				} else {
					fileName = ""
				}
			}
			if fileName == "" {
				isDir = true
				fileName = platform.PackingPath + config.C.Separator
			}
		}
	case "doc":
		if platform.DocPath != "" {
			for _, v := range config.DOC_EXTS {
				fileName = platform.DocPath + config.C.Separator + romName + v
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
			for _, v := range config.DOC_EXTS {
				fileName = platform.StrategyPath + config.C.Separator + romName + v
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
		return errors.New(config.C.Lang["PathNotFound"])
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

	//读取文档内容
	romName := utils.GetFileName(filepath.Base(info.RomPath)) //生成新文件的完整绝路路径地址
	if config.C.Platform[info.Platform].DocPath != "" {
		docFileName := "";
		for _, v := range config.DOC_EXTS {
			docFileName = config.C.Platform[info.Platform].DocPath + config.C.Separator + romName + v
			res.DocContent = GetDocContent(docFileName)
			if res.DocContent != "" {
				break
			}
		}
	}

	if config.C.Platform[info.Platform].StrategyPath != "" {
		//检测攻略可执行文件是否存在
		strategyFileName := "";
		for _, v := range config.RUN_EXTS {
			strategyFileName = config.C.Platform[info.Platform].StrategyPath + config.C.Separator + romName + v
			if utils.FileExists(strategyFileName) {
				res.StrategyFile = strategyFileName
				break
			}
		}

		//如果没有执行运行的文件，则读取文档内容
		if strategyFileName != "" {
			for _, v := range config.DOC_EXTS {
				strategyFileName = config.C.Platform[info.Platform].StrategyPath + config.C.Separator + romName + v
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
	return content
}

//更新展示图片
func UpdateRomThumbs(typeId int, id uint64, newPath string) (string, error) {

	rom := &db.Rom{
		Id: id,
	}

	//设定新的文件名
	vo, err := rom.GetById(id)
	if err != nil {
		return "", err
	}

	//下载文件
	res, err := http.Get(newPath)
	if err != nil {
		return "", err
	}

	//下载成功后，备份原文件
	platformPath := ""
	//原图存在，则备份
	if typeId == 1 {
		platformPath = config.C.Platform[vo.Platform].ThumbPath
	} else {
		platformPath = config.C.Platform[vo.Platform].SnapPath
	}

	if platformPath == "" {
		return "", errors.New(config.C.Lang["NoSetThumbDir"])
	}

	//开始备份原图
	bakFolder := config.C.CachePath + "thumb_bak/"
	RomFileName := utils.GetFileName(vo.RomPath)

	//检测bak文件夹是否存在，不存在则创建bak目录
	folder := utils.FolderExists(bakFolder)
	if folder == false {

		if err := utils.CreateDir(bakFolder); err != nil {
			return "", err
		}
	}
	for _, ext := range config.PIC_EXTS {
		oldFileName := platformPath + config.C.Separator + RomFileName + ext //老图片文件名
		if utils.FileExists(oldFileName) {
			bakFileName := RomFileName + "_" + utils.ToString(time.Now().Unix()) + ext //生成备份文件名
			err := os.Rename(oldFileName, bakFolder+bakFileName)                       //移动文件
			if err != nil {
				return "", err
			}
		}
	}

	//生成新文件
	platformPathAbs, err := filepath.Abs(platformPath) //读取平台图片路径

	newFileName := platformPathAbs + config.C.Separator + RomFileName + utils.GetFileExt(newPath) //生成新文件的完整绝路路径地址
	f, err := os.Create(newFileName)
	if err != nil {
		return "", err
	}
	io.Copy(f, res.Body)
	return newFileName, nil
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
		if err := (&db.Rom{}).UpdateSimConf(id, simId, data["cmd"], uint8(utils.ToInt(data["unzip"]))); err != nil {
			return err
		}
	}
	return nil
}
