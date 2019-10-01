package db

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
)


type Menu struct {
	Name     string
	Platform int64
	Pinyin string
}

//写入cate数据
func (*Menu) Add(menulist *map[string]*Menu) error {

	//关闭同步
	sqlite.Exec("PRAGMA synchronous = OFF;")

	stmt, err := sqlite.Prepare("INSERT INTO menu (`name`,platform,pinyin) values(?,?,?)")

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	//开始写入父rom
	for _, v := range *menulist {
		_, err := stmt.Exec(v.Name,v.Platform,v.Pinyin);
		if err != nil {
		}
	}

	return nil
}

//根据条件，查询多条数据
func (*Menu) GetByPlatform(platform int64) ([]*Menu, error) {

	volist := []*Menu{}

	where := ""

	if platform != 0 {
		where += " WHERE platform = " + strconv.FormatInt(platform,10)
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
