package db

import (
	"VirtualNesGUI/code/utils"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type RomCmd struct {
	Id    uint32
	RomId uint64
	SimId uint32
	Cmd   string
}

//插入数据
func (r *RomCmd) Add() error {

	//写入数据
	stmt, err := sqlite.Prepare("INSERT INTO rom_cmd (rom_id,sim_id,cmd) values(?,?,?)")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	if _, err = stmt.Exec(r.RomId, r.SimId, r.Cmd);err != nil{
		fmt.Println(err.Error())
		return err
	}

	return nil
}

//根据id查询一条数据
func (*RomCmd) Get(romId uint64,simId uint32) (*RomCmd, error) {
	vo := &RomCmd{}
	sql := "SELECT * FROM rom_cmd where rom_id = " + utils.ToString(romId) + " AND sim_id = " + utils.ToString(simId)
	rows := sqlite.QueryRow(sql)
	err := rows.Scan(&vo.Id, &vo.RomId, &vo.SimId, &vo.Cmd)
	return vo, err
}

//更新模拟器命令行参数
func (r *RomCmd) UpdateSimCmd() error {
	sql := `UPDATE rom_cmd SET cmd = ` + utils.ToString(r.Cmd) + ` WHERE id = ` + utils.ToString(r.Id)
	if _, err := sqlite.Exec(sql);err !=nil{
		return err
	}
	return nil
}

//删除一个模拟器的所有cmd数据
func (sim *RomCmd) ClearBySimId() (error) {
	sql := "DELETE FROM rom_cmd WHERE sim_id = " + utils.ToString(sim.SimId)
	if _, err := sqlite.Exec(sql); err != nil{
		return err
	}
	return nil
}
