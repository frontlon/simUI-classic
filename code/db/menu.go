package db

import (
	"VirtualNesGUI/code/utils"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"strings"
)

type Menu struct {
	Name     string
	Platform uint32
	Pinyin   string
	Sort     uint32
}

//写入cate数据
func (v *Menu) Add() error {
	stmt, err := sqlite.Prepare("INSERT INTO menu (`name`,platform,pinyin) values(?,?,?)")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	_, err = stmt.Exec(v.Name, v.Platform, v.Pinyin);
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

//根据条件，查询多条数据
func (*Menu) GetByPlatform(platform uint32) ([]*Menu, error) {

	volist := []*Menu{}

	where := ""

	if platform != 0 {
		where += " WHERE platform = " + utils.ToString(platform)
	}
	sql := "SELECT name,platform FROM menu " + where + " ORDER BY sort ASC,pinyin ASC"

	rows, err := sqlite.Query(sql)
	if err != nil {
		return volist, err
	}
	for rows.Next() {
		v := &Menu{}
		err = rows.Scan(&v.Name, &v.Platform)
		if err != nil {
			return volist, err
		}
		volist = append(volist, v)
	}
	return volist, nil
}

//删除一个平台下不存在的目录
func (*Menu) DeleteNotExists(platform uint32, menus []string) (error) {

	sql := ""
	//如果没有菜单，则全部删除
	if len(menus) == 0 {
		sql = "DELETE FROM menu WHERE platform = " + utils.ToString(platform)
	} else {
		menuStr := strings.Join(menus, "\",\"")
		menuStr = "\"" + menuStr + "\""
		sql = "DELETE FROM menu WHERE platform = " + utils.ToString(platform) + " AND name not in (" + menuStr + ")"
	}
	_, err := sqlite.Exec(sql)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

//删除不存在的平台下的所有menu
func (*Menu) ClearByPlatform(platforms []string) (error) {

	sql := "DELETE FROM menu "

	if len(platforms) > 0 {
		namesStr := strings.Join(platforms, ",")
		sql += " WHERE platform not in (" + namesStr + ")"
	}

	_, err := sqlite.Exec(sql)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

//根据一组名称，查询存在的名称，用于取交集
func (*Menu) GetMenuByNames(platform uint32, names []string) ([]string, error) {

	nameStr := strings.Join(names, "\",\"")
	nameStr = "\"" + nameStr + "\""

	nameList := []string{}
	sql := "SELECT name FROM menu WHERE platform = " + utils.ToString(platform) + " AND name in (" + nameStr + ")"
	rows, err := sqlite.Query(sql)
	if err != nil {
		return nameList, err
	}

	for rows.Next() {
		n := ""
		err = rows.Scan(&n)
		nameList = append(nameList, n)
	}
	return nameList, err
}

//清空表数据
func (*Menu) Truncate() (error) {
	_, err := sqlite.Exec("DELETE FROM menu")
	if err != nil {
		return err
	}
	return nil
}

//更新排序
func (menu *Menu) UpdateSortByName() error {
	sql := `UPDATE menu SET `
	sql += `sort = ` + utils.ToString(menu.Sort)
	sql += ` WHERE name = '` + utils.ToString(menu.Name) + `'`

	stmt, err := sqlite.Prepare(sql)
	if err != nil {
		return err
	}
	_, err2 := stmt.Exec()
	if err2 != nil {
		return err2
	} else {
		return nil
	}
}