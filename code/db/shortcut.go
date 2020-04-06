package db

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type Shortcut struct {
	Id   uint32
	Name string
	Path string
	Sort uint32
}

func (*Shortcut) TableName() string {
	return "shortcut"
}

//写入数据
func (m *Shortcut) Add() (int64, error) {

	result := getDb().Create(&m)

	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return int64(m.Id), result.Error
}

//读取所有数据
func (sim *Shortcut) GetAll() ([]*Shortcut, error) {
	volist := []*Shortcut{}
	result := getDb().Order("sort ASC").Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return volist, nil
}

//查询所有记录数
func (*Shortcut) Count() (int, error) {
	count := 0
	result := getDb().Count(&count)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return count, result.Error
}

//更新一条记录
func (m *Shortcut) UpdateById() error {
	create := map[string]interface{}{
		"name": m.Name,
		"path": m.Path,
	}
	result := getDb().Where("id=?", m.Id).Updates(create)

	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

//更新排序
func (m *Shortcut) UpdateSortById() error {
	result := getDb().Where("id=?", m.Id).Update("sort", m.Sort)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

//删除一条记录
func (m *Shortcut) DeleteById() (error) {
	result := getDb().Where("id=?", m.Id).Delete(&m)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}
