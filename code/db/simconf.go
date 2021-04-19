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
func (m *Rom) UpdateSimConf(romId uint64, simId uint32, cmd string, unzip uint8, file string,lua string) error {

	vo := &Rom{}
	result := getDb().Select("sim_conf,name,platform").Where("id=?", romId).First(&vo)
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
		Lua : lua,
	}
	jsonInfo, _ := json.Marshal(&sim)

	//更新到数据库
	result = getDb().Table(m.TableName()).Where("id=?", romId).Update("sim_conf", jsonInfo)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	//更新所有子游戏
	getDb().Table(m.TableName()).Where("platform=? AND pname=?", vo.Platform, vo.Name).Update("sim_conf", jsonInfo)

	return result.Error
}

//删除一个rom模拟器参数
func (m *Rom) DelSimConf(romId uint64, simId uint32) error {

	vo := &Rom{}

	result := getDb().Select("sim_conf,name,platform").Where("id=?", romId).First(&vo)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	sim := map[uint32]*SimConf{}

	json.Unmarshal([]byte(vo.SimConf), &sim)

	delete(sim, simId)

	jsonInfo, _ := json.Marshal(&sim)

	//更新到数据库
	result = getDb().Table(m.TableName()).Where("id=? AND sim_id=?", romId, simId).Update("sim_conf", jsonInfo)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	//更新所有子游戏
	getDb().Table(m.TableName()).Where("platform=? AND pname=?", vo.Platform, vo.Name).Update("sim_conf", jsonInfo)

	return result.Error
}
