package components

import (
	"fmt"
	"github.com/go-ini/ini"
	"simUI/code/config"
	"simUI/code/db"
	"simUI/code/utils"
	"strings"
)

func WriteRomConfigToIni(p string, subRom map[uint32]map[string]string, romSetting map[uint32]map[string]*db.RomSetting) {

	var outputCfg = ini.Empty()

	//子游戏关系
	for platformId, roms := range subRom {
		section := "subGame." + utils.ToString(platformId)
		for k, v := range roms {
			outputCfg.Section(section).Key(k).SetValue(v)
		}
	}

	//rom设置
	for platformId, settings := range romSetting {
		for k, v := range settings {
			romName := strings.ReplaceAll(k, ".", "@@*@@")
			section := "setting." + utils.ToString(platformId) + "." + romName
			outputCfg.Section(section).Key("SimId").SetValue(utils.ToString(v.SimId))
			outputCfg.Section(section).Key("SimConf").SetValue(utils.ToString(v.SimConf))
			outputCfg.Section(section).Key("Star").SetValue(utils.ToString(v.Star))
			outputCfg.Section(section).Key("Hide").SetValue(utils.ToString(v.Hide))
			outputCfg.Section(section).Key("RunLasttime").SetValue(utils.ToString(v.RunLasttime))
			outputCfg.Section(section).Key("RunNum").SetValue(utils.ToString(v.RunNum))
			outputCfg.Section(section).Key("Menu").SetValue(utils.ToString(v.Menu))
			outputCfg.Section(section).Key("Complete").SetValue(utils.ToString(v.Complete))
		}
	}

	if err := outputCfg.SaveTo(p); err != nil {
		fmt.Println(err)
	}

}

// 复制文件夹时，生成 src 和 dst 两个路径
func GetSrcAndDstPath(root, p, resType string) (string, string) {
	if root == "" || p == "" {
		return "", ""
	}

	if resType == "simulator" {
		p = utils.GetFilePath(p)
	}

	//绝对路径不拷贝
	src := strings.Replace(p, root, "", 1)
	if utils.IsAbsPath(src) {
		return p, ""
	}
	src = root + src
	dst := config.Cfg.RootPath + p

	//目录不存在，则新建
	fileType, _ := utils.CheckFileOrDir(dst)
	if fileType == 2 {
		if !utils.DirExists(dst) {
			utils.CreateDir(dst)
		}
		//已存在，不复制
		if !utils.IsDirEmpty(dst) {
			return "", ""
		}
	} else {
		//已存在，不复制
		if utils.FileExists(dst) {
			return "", ""
		}
	}
	return src, dst
}
