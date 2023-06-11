package components

import (
	"fmt"
	"simUI/code/db"
	"simUI/code/utils"
)

// 开始写入rom_subgame数据
func SaveRomConfigSubRom(subRom map[uint32]map[string]string) {
	for platformId, roms := range subRom {

		//读取rom数据
		romList, _ := (&db.Rom{}).GetByPlatform(platformId)
		romMap := map[string]string{}
		for _, v := range romList {
			romMap[utils.GetFileName(v.RomPath)] = v.FileMd5
		}

		//读取现有子游戏配置
		subGameListMap, _ := (&db.RomSubGame{}).GetByPlatformToMap(platformId)

		data := []*db.RomSubGame{}
		for subName, masterName := range roms {

			if _, ok := romMap[subName]; !ok {
				continue
			}
			if _, ok := romMap[masterName]; !ok {
				continue
			}

			//如果subgame表里已经有数据，则跳过
			if _, ok := subGameListMap[romMap[subName]]; ok {
				continue
			}

			d := &db.RomSubGame{
				Platform: platformId,
				FileMd5:  romMap[subName],
				Pname:    romMap[masterName],
			}
			data = append(data, d)
		}

		if len(data) > 0 {
			fmt.Println("开始写入数据库", len(data))
			(&db.RomSubGame{}).BatchAdd(data)
		}
	}
}

// 开始写入rom_config数据
func SaveRomConfigSetting(romSetting map[uint32][]map[string]string) {
	for platformId, roms := range romSetting {
		//读取rom数据
		romList, _ := (&db.Rom{}).GetByPlatform(platformId)
		romMap := map[string]string{}
		for _, v := range romList {
			romMap[utils.GetFileName(v.RomPath)] = v.FileMd5
		}

		settingList, _ := (&db.RomSetting{}).GetByPlatformToMap(platformId)

		data := []*db.RomSetting{}
		for _, rom := range roms {

			if _, ok := romMap[rom["name"]]; !ok {
				continue
			}

			if _, ok := settingList[romMap[rom["name"]]]; ok {
				continue
			}

			d := &db.RomSetting{
				Platform:    platformId,
				FileMd5:     romMap[rom["name"]],
				Star:        uint8(utils.ToInt(rom["Star"])),
				Hide:        uint8(utils.ToInt(rom["Hide"])),
				RunNum:      uint64(utils.ToInt(rom["RunNum"])),
				RunLasttime: int64(utils.ToInt(rom["RunLasttime"])),
				SimId:       uint32(utils.ToInt(rom["SimId"])),
				SimConf:     rom["SimConf"],
				Menu:        rom["Menu"],
				Complete:    uint8(utils.ToInt(rom["Complete"])),
			}
			data = append(data, d)
		}
		if len(data) > 0 {
			(&db.RomSetting{}).BatchAdd(data)
		}
	}
}
