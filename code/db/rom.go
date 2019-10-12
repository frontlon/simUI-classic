package db

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
	"strings"
)

type Rom struct {
	Id           int64
	Pname        string // 所属主游戏
	Menu         string // 菜单名称
	Name         string // 游戏名称
	Platform     int64  // 平台
	RomPath      string // rom路径
	ThumbPath    string // 缩略图路径
	SnapPath     string // 场景截图路径
	DocPath      string // 介绍文档路径
	StrategyPath string // 攻略文档路径
	Star         int64  // 喜好，星级
	SimId        int64  // 使用的模拟器id
	RunNum       int64  // 运行次数
	RunTime      int64  // 最后运行时间
	Pinyin       string // 拼音索引
}

//更新插入，(这里不能使用replace)
func (r *Rom) UpdateSert() error {

	//先查询记录是否存在
	vo := &Rom{}
	sql := "SELECT id FROM rom "
	sql += ` WHERE platform = ` + strconv.Itoa(int(r.Platform)) + ` AND pname = '`+ sql +`' AND name = '` + r.Name + `'`

	rows := sqlite.QueryRow(sql)
	err := rows.Scan(&vo.Id)
	if err != nil{
		return err
	}

	//记录不存在，则新建记录
	if vo == nil{
		r.Add()
		return nil
	}

	//如果记录存在，则更新记录
	sql = `UPDATE rom SET `
	sql += `rom_path = '` + r.RomPath + `'`
	sql += ` ,thumb_path = '` + r.ThumbPath + `'`
	sql += ` ,snap_path = '` + r.SnapPath + `'`
	sql += ` ,doc_path = '` + r.DocPath + `'`
	sql += ` ,strategy_path = '` + r.StrategyPath + `' `
	sql += ` WHERE id = ` + strconv.Itoa(int(vo.Id))
	_, err = sqlite.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}


//插入rom数据
func (r *Rom) Add() error {

	//关闭同步
	sqlite.Exec("PRAGMA synchronous = OFF;")

	stmt, err := sqlite.Prepare("INSERT INTO rom (`name`,pname,menu,platform,rom_path,thumb_path,snap_path,doc_path,strategy_path,pinyin) values(?,?,?,?,?,?,?,?,?,?)")

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	//开始写入父rom
	_, err = stmt.Exec(r.Name, r.Pname, r.Menu, r.Platform, r.RomPath, r.ThumbPath, r.SnapPath, r.DocPath, r.StrategyPath, r.Pinyin);
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}


/*

func (*Rom) Add(romlist *[]*Rom) error {

	lastIds := make(map[string]int64)

	//关闭同步
	sqlite.Exec("PRAGMA synchronous = OFF;")

	stmt, err := sqlite.Prepare("REPLACE INTO rom (`name`,pname,menu,platform,rom_path,thumb_path,snap_path,doc_path,strategy_path,pinyin) values(?,?,?,?,?,?,?,?,?,?)")

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	//开始写入父rom
	for _, r := range *romlist {
		res, err := stmt.Exec(r.Name, r.Pname, r.Menu, r.Platform, r.RomPath, r.ThumbPath, r.SnapPath, r.DocPath, r.StrategyPath, r.Pinyin);
		if err != nil {
		}
		id, _ := res.LastInsertId()
		lastIds[r.Name] = id
	}

	return nil
}

*/

//根据条件，查询多条数据
func (*Rom) Get(pages int, platform string, menu string, keyword string) ([]*Rom, error) {
	num := 50 //每页显示100个

	volist := []*Rom{}
	field := "id,name,menu,thumb_path";
	sql := "SELECT " + field + " FROM rom WHERE 1=1"
	if platform != "" {
		sql += " AND platform = '" + platform + "'"
	}

	if menu != "" {
		if menu == "favorite" {
			sql += " AND star = '1'"
		} else {
			sql += " AND menu = '" + menu + "'"
		}
	}

	sql += " AND pname = ''"

	if keyword != "" {
		sql += " AND name LIKE '%" + keyword + "%'"
	}

	sql += "  ORDER BY pinyin ASC LIMIT " + strconv.Itoa(num)

	if pages > 0 {
		offset := pages * num
		sql += " OFFSET " + strconv.Itoa(offset)
	}

	rows, err := sqlite.Query(sql)
	if err != nil {
		return volist, err
	}
	for rows.Next() {
		vo := &Rom{}
		err = rows.Scan(&vo.Id, &vo.Name, &vo.Menu, &vo.ThumbPath)
		volist = append(volist, vo)
	}

	return volist, err
}

//根据条件，查询多条数据
func (*Rom) GetSubRom(platform int64, pname string) ([]*Rom, error) {

	volist := []*Rom{}

	if platform == 0 || pname == "" {
		return volist, nil
	}

	field := "id,name,pname,rom_path";
	sql := "SELECT " + field + " FROM rom WHERE 1=1"
	sql += " AND platform = '" + string(platform) + "' AND pname = '" + pname + "'"
	sql += " ORDER BY pinyin ASC"

	rows, err := sqlite.Query(sql)
	if err != nil {
		return volist, err
	}
	for rows.Next() {
		vo := &Rom{}
		err = rows.Scan(&vo.Id, &vo.Name, &vo.Pname, &vo.RomPath)
		volist = append(volist, vo)
	}

	return volist, err
}

//根据id查询一条数据
func (*Rom) GetById(id string) (*Rom, error) {
	vo := &Rom{}
	sql := "SELECT * FROM rom where id= '" + id + "'"
	rows := sqlite.QueryRow(sql)
	err := rows.Scan(&vo.Id, &vo.Platform, &vo.Menu, &vo.Name, &vo.Pname, &vo.RomPath, &vo.ThumbPath, &vo.SnapPath, &vo.DocPath, &vo.StrategyPath,&vo.Star,&vo.SimId, &vo.Pinyin)
	return vo, err
}

//根据拼音筛选
func (*Rom) GetByPinyin(pages int, platform string, menu string, keyword string) ([]*Rom, error) {
	num := 50 //每页显示100个
	volist := []*Rom{}
	field := "id,name,menu,thumb_path";
	sql := ""
	pf := ""
	if platform != "" {
		pf = " platform=" + platform + " AND "
	}

	if menu != "" {
		pf = " menu = '" + menu + "' AND "
	}

	if keyword == "#" {
		subsql := "SELECT id FROM rom WHERE " + pf + " (pinyin LIKE 'a%'"
		//查询b-z
		for i := 98; i <= 122; i++ {
			subsql += " OR pinyin LIKE '" + string(i) + "%'"
		}
		subsql += ")"
		sql += "SELECT " + field + " FROM rom WHERE " + pf + " id not in (" + subsql + ")"
	} else {
		sql = "SELECT " + field + " FROM rom WHERE " + pf + " pinyin LIKE '" + keyword + "%'"
	}
	sql += " ORDER BY pinyin ASC LIMIT " + strconv.Itoa(num)
	if pages > 0 {
		offset := pages * num
		sql += " OFFSET " + strconv.Itoa(offset)
	}

	rows, err := sqlite.Query(sql)
	if err != nil {
		return volist, err
	}
	for rows.Next() {
		vo := &Rom{}
		err = rows.Scan(&vo.Id, &vo.Name, &vo.Menu, &vo.ThumbPath, &vo.SnapPath)
		volist = append(volist, vo)
	}
	return volist, err
}

//查询star
func (*Rom) GetByStar(platform string, star int64) (*Rom, error) {
	vo := &Rom{}

	where := ""
	if platform != "" {
		where = " platform=" + platform + " AND "
	}

	sql := "SELECT * FROM rom WHERE " + where + " star = " + strconv.Itoa(int(star))
	rows := sqlite.QueryRow(sql)
	err := rows.Scan(&vo.Id, &vo.Platform, &vo.Menu, &vo.Name, &vo.Pname, &vo.RomPath, &vo.ThumbPath, &vo.SnapPath, &vo.DocPath,&vo.StrategyPath, &vo.Star, &vo.Pinyin)
	return vo, err
}

//根据满足条件的rom数量
func (*Rom) Count(platform string, menu string, keyword string) (int, error) {
	count := 0
	sql := "SELECT count(*) as count FROM rom WHERE 1=1"
	if platform != "" {
		sql += " AND platform = '" + platform + "'"
	}
	if menu != "" {
		sql += " AND menu = '" + menu + "'"
	}
	if keyword != "" {
		sql += " AND name LIKE '%" + keyword + "%'"
	}
	rows := sqlite.QueryRow(sql)
	err := rows.Scan(&count)
	return count, err
}

//更新喜爱状态
func (r *Rom) UpdateStar() error {
	sql := `UPDATE rom SET `
	sql += `star = ` + strconv.Itoa(int(r.Star))
	sql += ` WHERE platform = ` + strconv.Itoa(int(r.Platform)) + ` AND name = '` + r.Name + `'`
	_, err := sqlite.Exec(sql)
	if err != nil {
		return err
	}
	return nil

}

//更新图片地址
func (r *Rom) UpdatePic() error {

	sql := `UPDATE rom SET `

	if r.SnapPath != "" {
		sql += " snap_path = '" + r.SnapPath + "'"
	} else if r.ThumbPath != "" {
		sql += " thumb_path = '" + r.ThumbPath + "'"
	} else {
		return nil
	}

	sql += ` WHERE id = ` + strconv.Itoa(int(r.Id))
	_, err := sqlite.Exec(sql)
	if err != nil {
		return err
	}
	return nil

}

//删除一个平台下的所有rom数据
func (sim *Rom) DeleteByPlatform() (error) {
	sql := "DELETE FROM rom WHERE platform = " + strconv.Itoa(int(sim.Platform))
	_, err := sqlite.Exec(sql)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

//删除不存在的rom
func (sim *Rom) DeleteByNotExists(platform int64,names []string) (error) {
	namesStr:=strings.Join(names,"\",\"")
	namesStr = "\""+namesStr+"\""
	sql := "DELETE FROM rom WHERE platform = " + strconv.Itoa(int(sim.Platform)) + "AND name not in ("+namesStr+")"
	_, err := sqlite.Exec(sql)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}