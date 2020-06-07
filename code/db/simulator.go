package db

import (
	"simUI/code/utils"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type Simulator struct {
	Id       uint32
	Name     string
	Platform uint32
	Path     string
	Cmd      string
	Unzip    uint8
	Default  uint8
	Pinyin   string
	Sort     uint32
}

func (*Simulator) TableName() string {
	return "simulator"
}

//写入数据
func (m *Simulator) Add() (uint32, error) {
	result := getDb().Create(&m)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return m.Id, result.Error
}

func (m *Simulator) BatchAdd(simulators []*Simulator) {

	if len(simulators) == 0 {
		return
	}

	tx := getDb().Begin()
	for _, v := range simulators {
		tx.Create(&v)
	}
	tx.Commit()
}

//根据ID查询一个模拟器参数
func (*Simulator) GetById(id uint32) (*Simulator, error) {
	vo := &Simulator{}
	result := getDb().Select("id, platform, name, path, cmd, unzip,`default`").Where("id=?", id).First(&vo)
	return vo, result.Error
}

//根据条件，查询多条数据
func (*Simulator) GetByPlatform(platform uint32) ([]*Simulator, error) {

	volist := []*Simulator{}
	where := ""

	if platform != 0 {
		where += "platform = '" + utils.ToString(platform) + "'"
	}

	result := getDb().Select("id, platform, name, path, cmd, unzip,`default`").Where(where).Order("sort ASC,`default` DESC,pinyin ASC").Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return volist, nil
}

//读取所有模拟器
func (*Simulator) GetAll() ([]*Simulator, error) {
	volist := []*Simulator{}
	result := getDb().Order("sort ASC,`default` DESC,pinyin ASC").Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return volist, nil
}

//更新
func (m *Simulator) UpdateById() error {

	create := map[string]interface{}{
		"name":   m.Name,
		"path":   m.Path,
		"cmd":    m.Cmd,
		"unzip":  m.Unzip,
		"pinyin": m.Pinyin,
	}
	result := getDb().Table(m.TableName()).Where("id=(?)", m.Id).Updates(create)

	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

//更新模拟器为默认模拟器
func (m *Simulator) UpdateDefault(platform uint32, id uint32) error {

	//先将平台下的所有参数都设为0
	result := getDb().Table(m.TableName()).Where("platform=? AND `default`=?", platform, 1).Update("default", 0)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	//将指定的模拟器更换为默认模拟器
	result = getDb().Table(m.TableName()).Where("id=?", id).Update("default", 1)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return result.Error
}

//删除一个模拟器
func (m *Simulator) DeleteById() error {
	result := getDb().Delete(&m)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

//删除一个平台下的所有模拟器
func (m *Simulator) DeleteByPlatform() error {
	result := getDb().Where("platform=?", m.Platform).Delete(&m)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

//更新排序
func (m *Simulator) UpdateSortById() error {
	result := getDb().Table(m.TableName()).Where("id=?", m.Id).Update("sort", m.Sort)
	if result.Error != nil {
		fmt.Println(result.Error.Error())
	}
	return result.Error
}
