package db

import (
	"VirtualNesGUI/code/utils"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type Shortcut struct {
	Id     uint32
	Name   string
	Path   uint32
	Type   uint8
	Pinyin string
}

//写入数据
func (v *Shortcut) Add() error {
	stmt, err := sqlite.Prepare("INSERT INTO shortcut (`name`,`type`,path,pinyin) values(?,?,?,?)")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	_, err = stmt.Exec(v.Name, v.Type, v.Path, v.Pinyin);
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

//根据id查询一条记录
func (*Shortcut) GetById(id uint32) (*Shortcut, error) {
	vo := &Shortcut{}
	sql := "SELECT * FROM shortcut where id= " + utils.ToString(id)
	rows := sqlite.QueryRow(sql)
	err := rows.Scan(&vo.Id, &vo.Name, &vo.Path, &vo.Type, &vo.Pinyin)
	return vo, err
}

//读取一个类型下的所有数据
func (sim *Shortcut) GetByType(t uint8) ([]*Shortcut,error) {
	volist := []*Shortcut{}
	sql := "SELECT * FROM rom WHERE type = " + utils.ToString(t) + " ORDER BY pinyin ASC"
	rows, err := sqlite.Query(sql)
	if err != nil {
		return volist, err
	}
	for rows.Next() {
		vo := &Shortcut{}
		err = rows.Scan(&vo.Id, &vo.Name, &vo.Path, &vo.Type, &vo.Pinyin)
		if err != nil {
			return volist, err
		}
		volist = append(volist, vo)
	}
	return volist,nil
}
