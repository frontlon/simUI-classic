package db

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
)

type Favorite struct {
	Name     string
	Platform int64
	Star     int64
}


//插入数据，如果存在则更新
func (fav *Favorite) UpSert() (error) {

	sql := "REPLACE INTO favorite (platform,`name`, star) VALUES (?,?,?)"
	stmt, err := sqlite.Prepare(sql)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	_, err = stmt.Exec(strconv.Itoa(int(fav.Platform)),fav.Name, strconv.Itoa(int(fav.Star)))
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

//删除一个喜好
func (fav *Favorite) Delete() (error) {
	sql := "DELETE FROM favorite WHERE platform = "+ strconv.Itoa(int(fav.Platform)) +" AND name = '"+fav.Name+"'"
	_, err := sqlite.Exec(sql)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

//根据平台，查询数据
func (fav *Favorite) GetByPlatform() (map[string]string, error) {
	volist := map[string]string{}
	sql := "SELECT `name`,star FROM favorite WHERE platform = "+ strconv.Itoa(int(fav.Platform))
	rows, err := sqlite.Query(sql)
	if err != nil {
		return volist, err
	}
	for rows.Next() {
		name := ""
		star := ""
		err = rows.Scan(&name,&star)
		if err != nil {
			return volist, err
		}
		volist[name] = star
	}
	return volist, nil
}


//清理老数据
func (fav *Favorite) Clear() (error) {
	subsql := "SELECT name FROM rom" //读取全部rom的name
	sql := "DELETE FROM favorite WHERE Name not in (" + subsql + ")"
	_, err := sqlite.Exec(sql)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}