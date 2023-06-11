package modules

import (
	"simUI/code/db"
	"simUI/code/utils"
)

//添加模拟器
func AddSimulator(data map[string]interface{}) (*db.Simulator, error) {
	pfId := uint32(utils.ToInt(data["platform"]))

	sim := &db.Simulator{
		Name:     utils.ToString(data["name"]),
		Platform: pfId,
		Path:     utils.ToString(data["path"]),
		Cmd:      utils.ToString(data["cmd"]),
		Lua:      utils.ToString(data["lua"]),
		Unzip:    uint8(utils.ToInt(data["unzip"])),
		Pinyin:   utils.TextToPinyin(utils.ToString(data["name"])),
	}
	id, err := sim.Add()

	//更新默认模拟器
	if utils.ToInt(data["default"]) == 1 {
		err = sim.UpdateDefault(pfId, id)
		if err != nil {
			return sim, err
		}
	}
	sim.Id = id
	return sim, nil
}

//更新模拟器
func UpdateSimulator(data map[string]interface{}) (*db.Simulator, error) {
	id := uint32(utils.ToInt(data["id"]))
	pfId := uint32(utils.ToInt(data["platform"]))
	def := uint8(utils.ToInt(data["default"]))
	sim := &db.Simulator{
		Id:       id,
		Name:     utils.ToString(data["name"]),
		Platform: pfId,
		Path:     utils.ToString(data["path"]),
		Cmd:      utils.ToString(data["cmd"]),
		Lua:      utils.ToString(data["lua"]),
		Pinyin:   utils.TextToPinyin(utils.ToString(data["name"])),
		Unzip:    uint8(utils.ToInt(data["unzip"])),
	}

	//更新模拟器
	if err := sim.UpdateById(); err != nil {
		return sim, err
	}

	//如果设置了默认模拟器，则更新默认模拟器
	if def == 1 {
		if err := sim.UpdateDefault(pfId, id); err != nil {
			return sim, err
		}
	}
	return sim, nil
}

func SetRomSimulator(romIds []uint64, simId uint32) error {

	roms, _ := (&db.Rom{}).GetByIds(romIds)

	fileMd5List := []string{}
	for _, v := range roms {
		fileMd5List = append(fileMd5List, v.FileMd5)
	}

	if err := (&db.Rom{}).UpdateSimIdByIds(romIds, simId); err != nil {
		utils.WriteLog(err.Error())
		return err
	}

	if err := (&db.RomSetting{}).UpdateSimIds(roms[0].Platform, fileMd5List, simId); err != nil {
		utils.WriteLog(err.Error())
		return err
	}

	return nil
}
