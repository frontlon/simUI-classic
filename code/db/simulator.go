package db

import (
	"VirtualNesGUI/code/utils"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type Simulator struct {
	Id       uint32
	Name     string
	Platform uint32
	Path     string
	Cmd      string
	Default  uint8
	Pinyin   string
}

//写入数据
func (sim *Simulator) Add() (uint32,error) {
	stmt, err := sqlite.Prepare("INSERT INTO simulator (`name`, platform, path, cmd, `default`,pinyin) values(?,?,?,?,?,?)")
	if err != nil {
		fmt.Println(err.Error())
		return 0,err
	}
	res, err := stmt.Exec(sim.Name, sim.Platform, sim.Path, sim.Cmd, sim.Default, sim.Pinyin);
	if err != nil {
		fmt.Println(err.Error())
		return 0,err
	}
	//返回插入的id
	id, err := res.LastInsertId()
	return uint32(id),nil
}


//根据ID查询一个模拟器参数
func (*Simulator) GetById(id uint32) (*Simulator, error) {
	v := &Simulator{}
	field := "id, platform, name, path, cmd, `default`"
	sql := "SELECT "+ field +" FROM simulator WHERE id = " + utils.ToString(id)

	rows := sqlite.QueryRow(sql)
	err := rows.Scan(&v.Id, &v.Platform, &v.Name, &v.Path, &v.Cmd, &v.Default)
	return v, err
}

//根据条件，查询多条数据
func (*Simulator) GetByPlatform(platform uint32) (map[uint32]*Simulator, error) {

	volist:= map[uint32]*Simulator{}
	where := ""

	if platform != 0 {
		where += " where platform = '" + utils.ToString(platform) + "'"
	}
	sql := "SELECT id, platform, name, path, cmd, `default` FROM simulator " + where + " ORDER BY `default` DESC,pinyin ASC"

	rows, err := sqlite.Query(sql)
	if err != nil {
		fmt.Println(err.Error())
		return volist, err
	}
	for rows.Next() {
		v := &Simulator{}
		err = rows.Scan(&v.Id, &v.Platform, &v.Name, &v.Path, &v.Cmd, &v.Default)
		if err != nil {
			fmt.Println(err.Error())
			return volist, err
		}
		volist[v.Id] = v
	}
	return volist, nil
}

//更新
func (sim *Simulator) UpdateById() error {
	sql := `UPDATE simulator SET `
	sql += `name = '` + sim.Name + `'`
	sql += `,path = '` + sim.Path + `'`
	sql += `,cmd = '` + sim.Cmd + `'`
	sql += `,pinyin = '` + sim.Pinyin + `'`
	sql += ` WHERE id = ` + utils.ToString(sim.Id)
	stmt, err := sqlite.Prepare(sql)

	if err != nil {
		return err
	}
	_, err2 := stmt.Exec()
	if err2 != nil {
		return err2
	}else{
		return nil
	}
}

//更新模拟器为默认模拟器
func (*Simulator) UpdateDefault(platform uint32,id uint32) error {

	//先将平台下的所有参数都设为0
	sql := "UPDATE simulator SET `default` = 0 WHERE `platform` = '" + utils.ToString(platform) + "' AND `default` = 1"
	stmt, err := sqlite.Prepare(sql)
	if err != nil {
		return err
	}
	_, err2 := stmt.Exec()
	if err2 != nil {
		return err2
	}

	//将指定的模拟器更换为默认模拟器
	sql = "UPDATE simulator SET `default` = 1 WHERE id = " + utils.ToString(id)
	stmt, err = sqlite.Prepare(sql)
	if err != nil {
		return err
	}
	_, err2 = stmt.Exec()
	if err2 != nil {
		return err2
	}
	return nil
}

//删除一个模拟器
func (sim *Simulator) DeleteById() (error) {
	sql := "DELETE FROM simulator WHERE id = "+ utils.ToString(sim.Id)
	_, err := sqlite.Exec(sql)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

//删除一个平台下的所有模拟器
func (sim *Simulator) DeleteByPlatform() (error) {
	sql := "DELETE FROM simulator WHERE platform = "+ utils.ToString(sim.Platform)
	_, err := sqlite.Exec(sql)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}