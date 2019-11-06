package db

import (
	"VirtualNesGUI/code/utils"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"strings"
)

type Platform struct {
	Id           uint32
	Name         string
	RomExts      []string
	RomPath      string
	ThumbPath    string
	SnapPath     string
	DocPath      string
	StrategyPath string
	Romlist      string
	Pinyin       string
	SimList      map[uint32]*Simulator
	UseSim       *Simulator //当前使用的模拟器
}

//添加平台
func (v *Platform) Add() (uint32, error) {

	stmt, err := sqlite.Prepare("INSERT INTO platform (`name`, rom_exts, rom_path, thumb_path, snap_path, doc_path, strategy_path, romlist, pinyin) values(?,?,?,?,?,?,?,?,?)")

	if err != nil {
		fmt.Println(err.Error())
		return 0, err
	}

	//开始写入父rom
	exts := ""
	res, err := stmt.Exec(v.Name, exts, v.RomPath, v.ThumbPath, v.SnapPath, v.DocPath, v.StrategyPath, v.Romlist, v.Pinyin);
	if err != nil {
	}
	id, _ := res.LastInsertId()
	return uint32(id), err
}

//根据条件，查询多条数据
func (*Platform) GetAll() (map[uint32]*Platform, error) {

	volist := map[uint32]*Platform{}
	exts := ""
	sql := "SELECT id,`name`, rom_exts, rom_path, thumb_path, snap_path, doc_path,strategy_path, romlist FROM platform  ORDER BY pinyin ASC"

	rows, err := sqlite.Query(sql)
	if err != nil {
		return volist, err
	}
	for rows.Next() {
		v := &Platform{}
		err = rows.Scan(&v.Id, &v.Name, &exts, &v.RomPath, &v.ThumbPath, &v.SnapPath, &v.DocPath, &v.StrategyPath, &v.Romlist)
		if err != nil {
			return volist, err
		}
		v.RomExts = strings.Split(exts, ",") //拆分rom扩展名
		volist[v.Id] = v
	}
	return volist, nil
}

//根据ID查询一个平台参数
func (*Platform) GetById(id uint32) (*Platform, error) {
	v := &Platform{}
	exts := ""
	field := "id,`name`, rom_exts, rom_path, thumb_path, snap_path, doc_path, strategy_path,romlist"
	sql := "SELECT " + field + " FROM platform WHERE id = " + utils.ToString(id)
	rows := sqlite.QueryRow(sql)
	err := rows.Scan(&v.Id, &v.Name, &exts, &v.RomPath, &v.ThumbPath, &v.SnapPath, &v.DocPath, &v.StrategyPath, &v.Romlist)
	v.RomExts = strings.Split(exts, ",") //拆分rom扩展名
	return v, err
}

//更新一个字段
func (pf *Platform) UpdateById() error {
	sql := `UPDATE platform SET `
	sql += `name = '` + pf.Name + `'`
	sql += `,rom_exts = '` + strings.Join(pf.RomExts, ",") + `'`
	sql += `,rom_path = '` + pf.RomPath + `'`
	sql += `,thumb_path = '` + pf.ThumbPath + `'`
	sql += `,snap_path = '` + pf.SnapPath + `'`
	sql += `,strategy_path = '` + pf.StrategyPath + `'`
	sql += `,doc_path = '` + pf.DocPath + `'`
	sql += `,romlist = '` + pf.Romlist + `'`
	sql += `,pinyin = '` + pf.Pinyin + `'`
	sql += ` WHERE id = ` + utils.ToString(pf.Id)

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

//删除一个平台
func (pf *Platform) DeleteById() error {
	sql := "DELETE FROM platform WHERE id = " + utils.ToString(pf.Id)
	_, err := sqlite.Exec(sql)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}
