package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type Menu struct {
	Name     string
	Platform uint32
	Pinyin   string
	Sort     uint32
}

func (*Menu) TableName() string {
	return "menu"
}

//写入cate数据
func (m *Menu) Add() error {
	result := getDb().Create(&m)

	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return nil
}

//根据条件，查询多条数据
func (*Menu) GetByPlatform(platform uint32) ([]*Menu, error) {
	where := map[string]interface{}{}

	if platform > 0 {
		where["platform"] = platform
	}

	volist := []*Menu{}
	result := getDb().Select("name,platform").Where(where).Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return volist, nil
}

//删除一个平台下不存在的目录
func (*Menu) DeleteNotExists(platform uint32, menus []string) error {

	result := &gorm.DB{}
	m := &Menu{}
	if len(menus) == 0 {
		result = getDb().Where("platform=?", platform).Delete(&m)
	} else {
		result = getDb().Where("platform=?", platform).Not("name", menus).Delete(&m)
	}

	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return result.Error
}

//删除不存在的平台下的所有menu
func (*Menu) ClearByPlatform(platforms []string) error {
	m := &Menu{}
	result := getDb().Not("platform", platforms).Delete(&m)

	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return result.Error
}

//根据一组名称，查询存在的名称，用于取交集
func (*Menu) GetMenuByNames(platform uint32, names []string) ([]string, error) {

	nameList := []string{}
	if len(names) == 0 {
		return nameList, nil
	}

	volist := []*Menu{}
	result := getDb().Select("name").Where("platform = (?) AND name in (?)", platform, names).Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	//仅读取name
	if len(volist) > 0 {
		for _, v := range volist {
			nameList = append(nameList, v.Name)
		}
	}
	return nameList, result.Error
}

//清空表数据
func (m *Menu) Truncate() error {
	result := getDb().Delete(&m)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

//更新排序
func (m *Menu) UpdateSortByName() error {
	result := getDb().Table(m.TableName()).Where("name = (?)", m.Name).Update("sort", m.Sort)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}
