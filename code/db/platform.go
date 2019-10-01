package db

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
	"strings"
)

type Platform struct {
	Id        int64
	Name      string
	RomExts   []string
	RomPath   string
	ThumbPath string
	VideoPath string
	SnapPath  string
	DocPath   string
	Romlist   string
	Status    int64
	Pinyin    string
	SimList   map[int64]*Simulator
	UseSim    *Simulator //当前使用的模拟器
}

//添加平台
func (v *Platform) Add() (int64,error) {

	//关闭同步
	sqlite.Exec("PRAGMA synchronous = OFF;")

	stmt, err := sqlite.Prepare("INSERT INTO platform (`name`, rom_exts, rom_path, thumb_path, snap_path, video_path, doc_path, romlist, status, pinyin) values(?,?,?,?,?,?,?,?,?,?)")

	if err != nil {
		fmt.Println(err.Error())
		return 0,err
	}

	//开始写入父rom
	exts := ""
	res, err := stmt.Exec(v.Name, exts, v.RomPath, v.ThumbPath, v.SnapPath,v.VideoPath, v.DocPath, v.Romlist, v.Status,v.Pinyin);
	if err != nil {
	}
	id, _ := res.LastInsertId()
	return id,err
}

//根据条件，查询多条数据
func (*Platform) GetAll() (map[int64]*Platform, error) {

	volist := map[int64]*Platform{}
	exts := ""
	sql := "SELECT id,`name`, rom_exts, rom_path, thumb_path, snap_path,video_path, doc_path, romlist, status FROM platform  ORDER BY pinyin ASC"

	rows, err := sqlite.Query(sql)
	if err != nil {
		return volist, err
	}
	for rows.Next() {
		v := &Platform{}
		err = rows.Scan(&v.Id, &v.Name, &exts, &v.RomPath, &v.ThumbPath, &v.SnapPath,&v.VideoPath, &v.DocPath, &v.Romlist, &v.Status)
		if err != nil {
			return volist, err
		}
		v.RomExts = strings.Split(exts, ",") //拆分rom扩展名
		volist[v.Id] = v
	}
	return volist, nil
}

//根据ID查询一个平台参数
func (*Platform) GetById(id int64) (*Platform, error) {
	v := &Platform{}
	exts := ""
	field := "id,`name`, rom_exts, rom_path, thumb_path, snap_path,video_path, doc_path, romlist, status"
	sql := "SELECT " + field + " FROM platform WHERE id = " + strconv.Itoa(int(id))
	rows := sqlite.QueryRow(sql)
	err := rows.Scan(&v.Id, &v.Name, &exts, &v.RomPath, &v.ThumbPath,&v.SnapPath, &v.VideoPath, &v.DocPath, &v.Romlist, &v.Status)
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
	sql += `,video_path = '` + pf.VideoPath + `'`
	sql += `,doc_path = '` + pf.DocPath + `'`
	sql += `,romlist = '` + pf.Romlist + `'`
	sql += `,status = '` + strconv.Itoa(int(pf.Status)) + `'`
	sql += `,pinyin = '` + pf.Pinyin + `'`
	sql += ` WHERE id = ` + strconv.Itoa(int(pf.Id))

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

//删除一个喜好
func (pf *Platform) Delete() (error) {
	sql := "DELETE FROM platform WHERE id = "+ strconv.Itoa(int(pf.Id))
	_, err := sqlite.Exec(sql)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}
