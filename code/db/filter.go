package db

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type Filter struct {
	Id       string
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

//根据条件，查询多条数据
func (*Filter) GetByPlatform(platform uint32, t string) ([]*Filter, error) {
	volist := []*Filter{}
	where := map[string]interface{}{}
	group := ""
	if platform > 0 {
		where["platform"] = platform
	} else {
		group = "name"
	}
	where["type"] = t

	result := getDb().Select("name").Where(where).Group(group).Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return volist, nil
}

//删除一个平台的数据
func (f *Filter) DeleteByPlatform() error {
	result := getDb().Where("platform = ?", f.Platform).Delete(&f)
	if result.Error != nil {
		fmt.Println(result.Error)
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
