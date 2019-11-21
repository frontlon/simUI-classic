package db

import (
	"VirtualNesGUI/code/utils"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"strings"
)

type Rom struct {
	Id       uint64
	Pname    string // 所属主游戏
	Menu     string // 菜单名称
	Name     string // 游戏名称
	Platform uint32 // 平台
	RomPath  string // rom路径
	Star     uint8  // 喜好，星级
	SimId    uint32 // 使用的模拟器id
	RunNum   uint64 // 运行次数
	RunTime  uint32 // 最后运行时间
	Pinyin   string // 拼音索引
	Md5      string // 文件Md5
}


//插入rom数据
func (r *Rom) Add() error {

	stmt, err := sqlite.Prepare("INSERT INTO rom (`name`,pname,menu,platform,rom_path,pinyin) values(?,?,?,?,?,?)")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	//开始写入父rom
	_, err = stmt.Exec(r.Name, r.Pname, r.Menu, r.Platform, r.RomPath, r.Pinyin);
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

//根据条件，查询多条数据
func (*Rom) Get(pages int, platform string, menu string, keyword string) ([]*Rom, error) {
	num := 50 //每页显示100个

	volist := []*Rom{}
	field := "id,name,menu,platform,rom_path";
	sql := "SELECT " + field + " FROM rom WHERE 1=1"
	if platform != "0" {
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

	sql += "  ORDER BY pinyin ASC LIMIT " + utils.ToString(num)

	if pages > 0 {
		offset := pages * num
		sql += " OFFSET " + utils.ToString(offset)
	}

	rows, err := sqlite.Query(sql)
	if err != nil {
		return volist, err
	}
	for rows.Next() {
		vo := &Rom{}
		err = rows.Scan(&vo.Id, &vo.Name, &vo.Menu, &vo.Platform, &vo.RomPath)
		volist = append(volist, vo)
	}

	return volist, err
}

//根据条件，查询多条数据
func (*Rom) GetSubRom(platform uint32, pname string) ([]*Rom, error) {

	volist := []*Rom{}

	if platform == 0 || pname == "" {
		return volist, nil
	}

	field := "id,name,pname";
	sql := "SELECT " + field + " FROM rom WHERE 1=1"
	sql += " AND platform = '" + utils.ToString(platform) + "' AND pname = '" + pname + "'"
	sql += " ORDER BY pinyin ASC"

	rows, err := sqlite.Query(sql)
	if err != nil {
		return volist, err
	}
	for rows.Next() {
		vo := &Rom{}
		err = rows.Scan(&vo.Id, &vo.Name, &vo.Pname)
		volist = append(volist, vo)
	}

	return volist, err
}

//根据id查询一条数据
func (*Rom) GetById(id string) (*Rom, error) {
	vo := &Rom{}
	sql := "SELECT * FROM rom where id= '" + id + "'"
	rows := sqlite.QueryRow(sql)
	err := rows.Scan(&vo.Id, &vo.Platform, &vo.Menu, &vo.Name, &vo.Pname, &vo.RomPath, &vo.Star, &vo.SimId, &vo.RunNum, &vo.RunTime, &vo.Pinyin)
	return vo, err
}

//根据拼音筛选
func (*Rom) GetByPinyin(pages int, platform string, menu string, keyword string) ([]*Rom, error) {
	num := 50 //每页显示100个
	volist := []*Rom{}
	field := "id,name,menu,platform,rom_path";
	sql := ""
	pf := ""
	if platform != "0" {
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
	sql += " ORDER BY pinyin ASC LIMIT " + utils.ToString(num)
	if pages > 0 {
		offset := pages * num
		sql += " OFFSET " + utils.ToString(offset)
	}

	rows, err := sqlite.Query(sql)
	if err != nil {
		return volist, err
	}

	for rows.Next() {
		vo := &Rom{}
		err = rows.Scan(&vo.Id, &vo.Name, &vo.Menu, &vo.Platform, &vo.RomPath)
		volist = append(volist, vo)
	}

	return volist, err
}

//查询star
func (*Rom) GetByStar(platform string, star uint8) (*Rom, error) {
	vo := &Rom{}

	where := ""
	if platform != "0" {
		where = " platform=" + platform + " AND "
	}

	sql := "SELECT * FROM rom WHERE " + where + " star = " + utils.ToString(star)
	rows := sqlite.QueryRow(sql)
	err := rows.Scan(&vo.Id, &vo.Platform, &vo.Menu, &vo.Name, &vo.Pname, &vo.RomPath, &vo.Star, &vo.Pinyin)
	return vo, err
}

//根据满足条件的rom数量
func (*Rom) Count(platform uint32, menu string, keyword string) (int, error) {
	count := 0
	sql := "SELECT count(*) as count FROM rom WHERE 1=1"
	if platform != 0 {
		sql += " AND platform = '" + utils.ToString(platform) + "' AND pname=''"
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
	sql += `star = ` + utils.ToString(r.Star)
	sql += ` WHERE id = ` + utils.ToString(r.Id)
	_, err := sqlite.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

//更新模拟器
func (r *Rom) UpdateSimulator() error {
	sql := `UPDATE rom SET `
	sql += `sim_id = ` + utils.ToString(r.SimId)
	sql += ` WHERE id = ` + utils.ToString(r.Id)
	_, err := sqlite.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

//删除一个平台下的所有rom数据
func (sim *Rom) DeleteByPlatform() (error) {
	sql := "DELETE FROM rom WHERE platform = " + utils.ToString(sim.Platform)
	_, err := sqlite.Exec(sql)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

//读取id列表
func (sim *Rom) GetIdsByPlatform(platform uint32, menu string) ([]uint64, error) {

	ids := []uint64{}
	sql := "SELECT id FROM rom WHERE platform = " + utils.ToString(platform) + " AND menu = '" + menu + "'"
	rows, err := sqlite.Query(sql)
	if err != nil {
		return ids, err
	}

	for rows.Next() {
		id := 0
		err = rows.Scan(&id)
		ids = append(ids, uint64(id))
	}
	return ids, err
}

//根据一组dm5，查询存在的md5，用于取交集
func (sim *Rom) GetMd5ByMd5(platform uint32, uniq []string) ([]string,error) {

	uniqStr := strings.Join(uniq, "\",\"")
	uniqStr = "\"" + uniqStr + "\""

	md5List := []string{}
	sql := "SELECT id FROM rom WHERE platform = " + utils.ToString(platform) + " AND md5 in (" + uniqStr + ")"
	rows, err := sqlite.Query(sql)
	if err != nil {
		return md5List, err
	}

	for rows.Next() {
		md5 := ""
		err = rows.Scan(&md5)
		md5List = append(md5List, md5)
	}
	return md5List, err
}

//删除指定平台下，不存在的rom
func (sim *Rom) DeleteNotExists(platform uint32, uniq []string) (error) {

	if len(uniq) == 0 {
		return nil
	}

	uniqStr := strings.Join(uniq, "\",\"")
	uniqStr = "\"" + uniqStr + "\""

	sql := "DELETE FROM rom WHERE platform = " + utils.ToString(platform) + " AND md5 not in (" + uniqStr + ")"

	//删除主记录
	_, err := sqlite.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}



//删除不存在的平台下的所有rom
func (sim *Rom) ClearByPlatform(platforms []string) (error) {

	sql := "DELETE FROM rom "

	if len(platforms) > 0 {
		namesStr := strings.Join(platforms, ",")
		sql += " WHERE platform not in (" + namesStr + ")"
	}

	_, err := sqlite.Exec(sql)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}
