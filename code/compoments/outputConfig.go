package compoments

import (
	"fmt"
	"github.com/go-ini/ini"
	"simUI/code/db"
	"simUI/code/utils"
	"strings"
)

func WriteRomConfigToIni(p string, subRom map[uint32]map[string]string, romSetting map[uint32]map[string]*db.RomSetting) {

	var outputCfg = ini.Empty()

	//子游戏
	for platformId, roms := range subRom {
		section := "subGame." + utils.ToString(platformId)
		for k, v := range roms {
			outputCfg.Section(section).Key(k).SetValue(v)
		}
	}

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
		}
	}

	if err := outputCfg.SaveTo(p); err != nil {
		fmt.Println(err)
	}

}
