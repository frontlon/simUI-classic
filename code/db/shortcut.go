package db

import (
	"VirtualNesGUI/code/utils"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type Shortcut struct {
	Id   uint32
	Name string
	Path string
	Sort uint32
}

//写入数据
func (v *Shortcut) Add() (int64,error) {
	stmt, err := sqlite.Prepare("INSERT INTO shortcut (sort) values(?)")
	if err != nil {
		fmt.Println(err.Error())
		return 0,err
	}
	result, err := stmt.Exec(v.Sort);
	if err != nil {
		fmt.Println(err.Error())
		return 0,err
	}
	return result.LastInsertId()
}

//读取所有数据
func (sim *Shortcut) GetAll() ([]*Shortcut, error) {
	volist := []*Shortcut{}
	sql := "SELECT * FROM rom WHERE ORDER BY sort ASC"
	rows, err := sqlite.Query(sql)
	if err != nil {
		return volist, err
	}
	for rows.Next() {
		vo := &Shortcut{}
		err = rows.Scan(&vo.Id, &vo.Name, &vo.Path, &vo.Sort)
		if err != nil {
			return volist, err
		}
		volist = append(volist, vo)
	}
	return volist, nil
}

//更新一条记录
func (r *Shortcut) UpdateById() error {
	sql := `UPDATE shortcut SET `
	sql += `name = '` + utils.ToString(r.Name) + `'`
	sql += `, path = '` + utils.ToString(r.Path) + `'`
	sql += `, sort = ` + utils.ToString(r.Sort)
	sql += ` WHERE id = ` + utils.ToString(r.Id)
	_, err := sqlite.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

//删除一条记录
func (r *Shortcut) DeleteById() (error) {
	sql := "DELETE FROM shortcut WHERE id = " + utils.ToString(r.Id)
	_, err := sqlite.Exec(sql)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}
