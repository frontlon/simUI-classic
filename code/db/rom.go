package db

import (
	"VirtualNesGUI/code/utils"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"strings"
)

var ROM_PAGE_NUM = 100; //每页加载rom数量

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

func (*Rom) TableName() string {
	return "rom"
}

//插入rom数据
func (m *Rom) Add() error {

	result := getDb().Create(&m)

	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return result.Error
}

//根据条件，查询多条数据
func (*Rom) Get(pages int, platform uint32, menu string, keyword string) ([]*Rom, error) {

	volist := []*Rom{}
	where := map[string]interface{}{}
	if platform != 0 {
		where["platform"] = platform
	}

	if menu != "" {
		if menu == "favorite" {
			where["star"] = 1
		} else {
			where["menu"] = menu
		}
	}
	where["pname"] = ""

	if keyword != "" {
		where["name LIKE"] = "%" + keyword + "%"
	}

	offset := pages * ROM_PAGE_NUM

	result := getDb().Select("id,name,menu,platform,rom_path").Where(where).Order("pinyin ASC").Limit(ROM_PAGE_NUM).Offset(offset).Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return volist, result.Error
}

//读取子rom
func (*Rom) GetSubRom(platform uint32, pname string) ([]*Rom, error) {

	volist := []*Rom{}

	if platform == 0 || pname == "" {
		return volist, nil
	}

	result := getDb().Select("id,name,pname").Where("platform=? AND pname=?", platform, pname).Order("pinyin ASC").Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return volist, result.Error
}

//根据id查询一条数据
func (*Rom) GetById(id uint64) (*Rom, error) {

	vo := &Rom{}

	result := getDb().Where("id=?", id).First(&vo)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return vo, result.Error
}

//根据拼音筛选
func (*Rom) GetByPinyin(pages int, platform uint32, menu string, keyword string) ([]*Rom, error) {
	where := map[string]interface{}{}

	if platform != 0 {
		where["platform"] = platform
	}

	if menu != "" {
		where["menu"] = menu
	}

	where["pname"] = ""
	offset := pages * ROM_PAGE_NUM
	volist := []*Rom{}
	field := "id,name,menu,platform,rom_path"
	result := getDb().Select(field).Order("pinyin ASC").Limit(ROM_PAGE_NUM).Offset(offset)
	if keyword == "#" {

		//查询0-9数字rom
		subWhere := "pinyin LIKE '0%'"
		for i := 1; i <= 9; i++ {
			subWhere += " OR pinyin LIKE '" + utils.ToString(i) + "%'"
		}
		result.Where(where).Where(subWhere).Find(&volist)
	} else {
		result.Where(where).Where("pinyin LIKE ?",keyword+"%").Find(&volist)
	}

	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return volist, result.Error
}

//查询star
/*func (*Rom) GetByStar(platform uint32, star uint8) (*Rom, error) {
	vo := &Rom{}

	where := ""
	if platform != 0 {
		where = " platform=" + utils.ToString(platform) + " AND "
	}

	sql := "SELECT * FROM rom WHERE " + where + " star = " + utils.ToString(star)
	rows := sqlite.QueryRow(sql)
	err := rows.Scan(&vo.Id, &vo.Platform, &vo.Menu, &vo.Name, &vo.Pname, &vo.RomPath, &vo.Star, &vo.Pinyin,&vo.Md5)
	return vo, err
}
*/
//根据满足条件的rom数量
func (m *Rom) Count(platform uint32, menu string, keyword string) (int, error) {
	count := 0
	where := map[string]interface{}{
	}

	if platform != 0 {
		where["platform"] = platform
		where["pname"] = ""
	}
	if menu != "" {
		where["menu"] = menu
	}
	if keyword != "" {
		where["name LIKE"] = "%" + keyword + "%'"
	}
	result := getDb().Table(m.TableName()).Where(where).Count(&count)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return count, result.Error
}

//更新喜爱状态
func (m *Rom) UpdateStar() error {
	result := getDb().Table(m.TableName()).Where("id=?", m.Id).Update("star", m.Star)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

//更新模拟器
func (m *Rom) UpdateSimulator() error {
	result := getDb().Table(m.TableName()).Where("id=?", m.Id).Update("sim_id", m.SimId)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

//删除一个平台下的所有rom数据
func (m *Rom) DeleteByPlatform() (error) {
	result := getDb().Where("platform=? ", m.Platform).Delete(&m)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

//读取id列表
func (sim *Rom) GetIdsByPlatform(platform uint32, menu string) ([]uint64, error) {
	volist := []*Rom{}
	result := getDb().Select("id").Where("platform=? AND menu=?", platform, menu).Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	ids := []uint64{}
	for _, v := range volist {
		ids = append(ids, v.Id)
	}
	return ids, result.Error
}

//根据一组dm5，查询存在的md5，用于取交集
func (sim *Rom) GetMd5ByMd5(platform uint32, uniq []string) ([]string, error) {
	volist := []*Rom{}
	result := getDb().Select("md5").Where("platform=? AND md5 in (?)", platform, uniq).Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	md5List := []string{}
	for _, v := range volist {
		md5List = append(md5List, v.Md5)
	}
	return md5List, result.Error
}

//删除指定平台下，不存在的rom
func (sim *Rom) DeleteNotExists(platform uint32, uniq []string) ([]string, error) {

	sql_query := ""
	sql := ""
	ids := []string{}
	//如果为空，说明目录下没有rom，全部删除
	if len(uniq) == 0 {
		sql_query = "SELECT id FROM rom WHERE platform = " + utils.ToString(platform)
		sql = "DELETE FROM rom WHERE platform = " + utils.ToString(platform)
	} else {
		uniqStr := strings.Join(uniq, "\",\"")
		uniqStr = "\"" + uniqStr + "\""
		sql_query = "SELECT id FROM rom WHERE platform = " + utils.ToString(platform) + " AND md5 not in (" + uniqStr + ")"
		sql = "DELETE FROM rom WHERE platform = " + utils.ToString(platform) + " AND md5 not in (" + uniqStr + ")"
	}

	//先把要删除的id查询出来
	rows, err := sqlite.Query(sql_query)
	defer rows.Close()
	if err != nil {
		return ids, err
	} else {
		for rows.Next() {
			var id = ""
			err = rows.Scan(&id)
			ids = append(ids, id)
		}
	}

	//删除记录
	if _, err := sqlite.Exec(sql); err != nil {
		return ids, err
	}
	return ids, nil
}

//删除不存在的平台下的所有rom
func (m *Rom) ClearByPlatform(platforms []string) (error) {
	result := getDb().Where("platform not in (?) AND menu=?", platforms).Delete(&m)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

//清空表数据
func (m *Rom) Truncate() (error) {
	result := getDb().Delete(&m)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}
