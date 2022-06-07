package db

import (
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type SimConf struct {
	Cmd   string
	Unzip uint8
	File  string
	Lua   string
}

//设置一个Rom的模拟器配置
func (*Rom) GetSimConf(romId uint64, simId uint32) (*SimConf, error) {

	vo := &Rom{}

	result := getDb().Select("sim_conf").Where("id=?", romId).First(&vo)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	sim := map[uint32]*SimConf{}
	if vo.SimConf != "" {
		if err := json.Unmarshal([]byte(vo.SimConf), &sim); err != nil {
			fmt.Println(err.Error())
			return &SimConf{Unzip: 2}, nil
		}
		if _, ok := sim[simId]; ok {
			return sim[simId], nil
		} else {
			return &SimConf{Unzip: 2}, nil
		}
	}

	return &SimConf{}, nil
}

//设置rom模拟器参数
func (m *Rom) UpdateSimConf(romId uint64, simId uint32, cmd string, unzip uint8, file string, lua string) error {

	vo := &Rom{}
	result := getDb().Select("sim_conf,file_md5,platform").Where("id=?", romId).First(&vo)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	sim := map[uint32]*SimConf{}

	if err := json.Unmarshal([]byte(vo.SimConf), &sim); err != nil {
		fmt.Println(err.Error())
	}
	sim[simId] = &SimConf{
		Cmd:   cmd,
		Unzip: unzip,
		File:  file,
		Lua:   lua,
	}
	jsonInfo, _ := json.Marshal(&sim)

	//更新到rom数据库
	result = getDb().Table(m.TableName()).Where("id=?", romId).Update("sim_conf", jsonInfo)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	//更新到rom_conf数据库
	result = getDb().Table((&RomSetting{}).TableName()).Where("file_md5=?", vo.FileMd5).Update("sim_conf", jsonInfo)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return result.Error
}

//更新模拟器配置
func (m *Rom) UpdateSimConfById(id uint64, conf string) error {
	result := getDb().Table(m.TableName()).Where("id = ?", id).Update("sim_conf", conf)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}
