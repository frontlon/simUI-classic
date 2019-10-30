package db

import (
	"VirtualNesGUI/code/utils"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)


type Menu struct {
	Name     string
	Platform uint32
	Pinyin string
}

//写入cate数据
func (v *Menu) Add() error {
	stmt, err := sqlite.Prepare("INSERT INTO menu (`name`,platform,pinyin) values(?,?,?)")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	_, err = stmt.Exec(v.Name,v.Platform,v.Pinyin);
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
	sql := "SELECT name,platform FROM menu " + where + " ORDER BY pinyin ASC"

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

//清理菜单数据
func (*Menu) ClearMenu(platform uint32) error {
	where := ""
	if platform > 0{
		where = " WHERE platform = "+ utils.ToString(platform)
	}
	if _,err :=sqlite.Exec(`DELETE FROM menu` + where);err != nil{
		return err
	}
	return nil
}
