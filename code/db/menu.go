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
func (*Menu) GetByPlatform(platform uint32, pages uint32) ([]*Menu, error) {
	where := map[string]interface{}{}

	if platform > 0 {
		where["platform"] = platform
	}

	volist := []*Menu{}

	pageNum := 200
	offset := int(pages) * pageNum
	result := getDb().Select("name,platform").Where(where).Order("sort ASC,pinyin ASC").Limit(pageNum).Offset(offset).Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return volist, nil
}

//读取所有菜单数据
func (*Menu) GetAll() ([]*Menu, error) {
	volist := []*Menu{}
	result := getDb().Select("name,platform").Order("sort ASC,pinyin ASC").Find(&volist)
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
		return nil
	}

	//数据量不会很大，慢慢删。
	tx := getDb().Begin()
	for _, v := range menus {
		tx.Where("platform=(?) AND name=(?)", platform, v).Delete(&m)
	}
	result = tx.Commit()

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

//读取一个平台下的所有menu数据
func (*Menu) GetAllNamesByPlatform(platform uint32) ([]string, error) {

	nameList := []string{}

	volist := []*Menu{}
	result := getDb().Select("name").Where("platform = (?)", platform).Find(&volist)
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
