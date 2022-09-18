package db

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"math"
)

type Filter struct {
	Platform uint32
	Type     string
	Name     string
}

func (*Filter) TableName() string {
	return "filter"
}

//写入数据
func (m *Filter) BatchAdd(data []*Filter) {

	if len(data) == 0 {
		return
	}

	tx := getDb().Begin()
	for _, v := range data {
		tx.Create(&v)
	}
	tx.Commit()
}

func (*Filter) GetAll() ([]*Filter, error) {
	volist := []*Filter{}

	result := getDb().Select("name,type").Order("name ASC").Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return volist, nil
}

func (*Filter) GetByPlatform(platform uint32) ([]*Filter, error) {
	volist := []*Filter{}

	result := getDb().Select("name,type").Where("platform = ?", platform).Order("name ASC").Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return volist, nil
}

//删除记录
func (m *Filter) DeleteByFileNames(platform uint32, t string, nameList []string) error {

	if len(nameList) == 0 {
		return nil
	}

	listLen := len(nameList)

	ceil := int(math.Ceil(float64(listLen) / float64(maxVar)))

	for i := 0; i < ceil; i++ {
		start := i * maxVar
		end := (i + 1) * maxVar
		if end > listLen {
			end = listLen
		}
		list := nameList[start:end]
		getDb().Where("platform = ? AND type = ? AND name in (?)", platform, t, list).Delete(&m)
	}

	return nil
}

//清空表数据
func (m *Filter) Truncate() error {
	result := getDb().Delete(&m)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}
