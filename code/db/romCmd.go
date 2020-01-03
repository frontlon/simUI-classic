package db

import (
	"VirtualNesGUI/code/utils"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"strings"
)

type RomCmd struct {
	Id    uint32
	RomId uint64
	SimId uint32
	Cmd   string
	Unzip uint8
}

//插入rom数据
func (r *RomCmd) Add() error {

	stmt, err := sqlite.Prepare("INSERT INTO rom_cmd (rom_id,sim_id,cmd,unzip) values(?,?,?,?)")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	//开始写入父rom
	_, err = stmt.Exec(r.RomId, r.SimId, r.Cmd,r.Unzip);
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

//查询
func (r *RomCmd) Get() (*RomCmd, error) {
	vo := &RomCmd{}
	sql := "SELECT id,rom_id,sim_id,cmd,unzip FROM rom_cmd WHERE rom_id= " + utils.ToString(r.RomId) + " AND sim_id = " + utils.ToString(r.SimId)
	rows := sqlite.QueryRow(sql)
	err := rows.Scan(&vo.Id, &vo.RomId, &vo.SimId, &vo.Cmd,&vo.Unzip)
	return vo, err
}

//更新cmd参数
func (r *RomCmd) UpdateCmd() error {
	sql := `UPDATE rom_cmd SET `
	sql += `cmd = '` + utils.ToString(r.Cmd) + `'`
	sql += `, unzip = ` + utils.ToString(r.Unzip)
	sql += ` WHERE id = ` + utils.ToString(r.Id)
	_, err := sqlite.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

//删除模拟器时，删除所有rom参数
func (r *RomCmd) ClearBySimId() (error) {
	sql := "DELETE FROM rom_cmd WHERE sim_id = " + utils.ToString(r.SimId)
	_, err := sqlite.Exec(sql)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

//删除一个游戏的模拟器参数记录
func (r *RomCmd) DeleteById() (error) {
	sql := "DELETE FROM rom_cmd WHERE id = " + utils.ToString(r.Id)
	_, err := sqlite.Exec(sql)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

//根据id列表删除数据
func (r *RomCmd) ClearByRomIds(ids []string) (error) {
	idsStr := strings.Join(ids, ",")

	sql := "DELETE FROM rom_cmd WHERE rom_id in (" + idsStr +")"
	_, err := sqlite.Exec(sql)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}