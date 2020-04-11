package db

import (
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type RomCmd struct {
	Id    uint32
	RomId uint64
	SimId uint32
	Cmd   string
	Unzip uint8
}

type SimConf struct {
	Cmd   string
	Unzip uint8
}

//设置一个Rom的模拟器配置
func (*Rom) GetSimConf(romId uint64, simId uint32) (*SimConf, error) {

	vo := &Rom{}

	result := getDb().Select("sim_conf").Where("id=?", romId).First(&vo)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	sim := map[uint32]*SimConf{}

	json.Unmarshal([]byte(vo.SimConf), &sim)

	return sim[simId], result.Error
}

//设置rom模拟器参数
func (*Rom) UpdateSimConf(romId uint64, simId uint32, cmd string, unzip uint8) error {

	vo := &Rom{}

	result := getDb().Select("sim_conf").Where("id=?", romId).First(&vo)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	sim := map[uint32]*SimConf{}

	json.Unmarshal([]byte(vo.SimConf), &sim)

	sim[simId].Cmd = cmd
	sim[simId].Unzip = unzip

	jsonInfo, _ := json.Marshal(&sim)

	//更新到数据库
	result = getDb().Where("id=? AND sim_id=?", romId, simId).Update("sim_conf", jsonInfo)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return result.Error
}

//删除一个rom模拟器参数
func (*Rom) DelSimConf(romId uint64, simId uint32) error {

	vo := &Rom{}

	result := getDb().Select("sim_conf").Where("id=?", romId).First(&vo)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	sim := map[uint32]*SimConf{}

	json.Unmarshal([]byte(vo.SimConf), &sim)

	delete(sim, simId)

	jsonInfo, _ := json.Marshal(&sim)

	//更新到数据库
	result = getDb().Where("id=? AND sim_id=?", romId, simId).Update("sim_conf", jsonInfo)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return result.Error
}

func (*RomCmd) TableName() string {
	return "rom_cmd"
}

//插入rom数据
func (m *RomCmd) Add() error {

	result := getDb().Create(&m)

	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return result.Error
}

//查询
/*func (m *RomCmd) Get() (*RomCmd, error) {
	vo := &RomCmd{}
	where := map[string]interface{}{
		"rom_id": m.RomId,
		"sim_id": m.SimId,
	}
	result := getDb().Select("id,rom_id,sim_id,cmd,unzip").Where(where).First(&vo)
	return vo, result.Error
}*/

//更新cmd参数
/*func (m *RomCmd) UpdateCmd() error {
	result := getDb().Where("id=?", m.Id).Updates(m)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}*/

//删除模拟器时，删除所有rom参数
/*func (m *RomCmd) ClearBySimId() (error) {
	result := getDb().Where("sim_id=?", m.SimId).Delete(&m)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}*/

//删除一个游戏的模拟器参数记录
/*func (m *RomCmd) DeleteById() (error) {
	result := getDb().Where("id=?", m.Id).Delete(&m)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}*/

//根据id列表删除数据
/*func (m *RomCmd) DeleteByRomIds(ids []string) (error) {
	if len(ids) == 0 {
		return nil
	}

	result := getDb().Where("id in (?) ", ids).Delete(&m)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}
*/