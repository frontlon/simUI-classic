package modules

import (
	"fmt"
	"github.com/go-ini/ini"
	"simUI/code/compoments"
	"simUI/code/config"
	"simUI/code/db"
	"simUI/code/utils"
	"strings"
)

/**
 * 导出rom配置
 */
func BackupRomConfig(p string) error {

	go func() {
		subRom := map[uint32]map[string]string{}
		romSetting := map[uint32]map[string]*db.RomSetting{}
		for platformId, _ := range config.Cfg.Platform {
			//读取rom数据
			romList, _ := (&db.Rom{}).GetByPlatform(platformId)
			romMap := map[string]string{}
			for _, v := range romList {
				romMap[v.FileMd5] = utils.GetFileName(v.RomPath)
			}

			//读取子游戏关系
			subList, _ := (&db.RomSubGame{}).GetByPlatform(platformId)
			subMap := map[string]string{}
			for _, v := range subList {

				if _, ok := romMap[v.FileMd5]; !ok {
					continue
				}
				if _, ok := romMap[v.Pname]; !ok {
					continue
				}

				subMap[utils.Base64Encode(romMap[v.FileMd5])] = utils.Base64Encode(romMap[v.Pname])
			}
			subRom[platformId] = subMap

			//读取rom配置
			settings, _ := (&db.RomSetting{}).GetByPlatform(platformId)
			settingsMap := map[string]*db.RomSetting{}
			for _, v := range settings {
				if _, ok := romMap[v.FileMd5]; !ok {
					continue
				}
				settingsMap[utils.Base64Encode(romMap[v.FileMd5])] = v
			}

			romSetting[platformId] = settingsMap
		}
		compoments.WriteRomConfigToIni(p, subRom, romSetting)

		if _, err := utils.Window.Call("CB_romConfigBackup"); err != nil {
		}

		fmt.Println("导出完成")
	}()
	return nil
}

/**
 * 导入rom配置
 */
func RestoreRomConfig(p string) error {

	f, err := ini.Load(p)
	if err != nil {
		return err
	}

	go func() {

		cfg := f.Sections()

		subRom := map[uint32]map[string]string{}
		romSetting := map[uint32][]map[string]string{}
		//解析数据
		for _, section := range cfg {
			sectionName := strings.Split(section.Name(), ".")
			if len(sectionName) < 2 {
				continue
			}
			platformId := uint32(utils.ToInt(sectionName[1]))
			//romSetting[platformId] = []map[string]string{}
			if strings.Contains(section.Name(), "subGame") {
				//子游戏
				subMap := map[string]string{}
				for _, v := range section.Keys() {
					subMap[utils.Base64Decode(v.Name())] = utils.Base64Decode(v.Value())
				}
				subRom[platformId] = subMap
			} else if strings.Contains(section.Name(), "setting") {
				//rom配置
				data := map[string]string{}
				data["name"] = utils.Base64Decode(sectionName[2])
				for _, v := range section.Keys() {
					data[v.Name()] = v.Value()
				}
				fmt.Println(data)
				romSetting[platformId] = append(romSetting[platformId], data)
			}
		}

		//开始写入rom_subgame数据
		compoments.SaveRomConfigSubRom(subRom)

		//开始写入rom_config数据
		compoments.SaveRomConfigSetting(romSetting)

		if _, err := utils.Window.Call("CB_romConfigRestore"); err != nil {
		}

		fmt.Println("导入完成")
	}()
	return nil

}
