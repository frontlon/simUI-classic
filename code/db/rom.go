package db

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"simUI/code/utils"
	"strings"
)

var ROM_PAGE_NUM = 100 //每页加载rom数量

type Rom struct {
	Id            uint64
	Pname         string // 所属主游戏
	Menu          string // 菜单名称
	Name          string // 游戏名称
	Platform      uint32 // 平台
	RomPath       string // rom路径
	Star          uint8  // 喜好，星级
	SimId         uint32 // 正在使用的模拟器id
	SimConf       string // 模拟器参数独立配置
	Hide          uint8  // 是否隐藏
	BaseType      string // 游戏类型，如RPG
	BaseYear      string // 游戏年份
	BasePublisher string // 游戏出品公司
	BaseCountry   string // 游戏国家
	Pinyin        string // 拼音索引
	InfoMd5       string // 信息md5，包含资料信息
	FileMd5       string // 文件md5，仅包含平台和文件名
}

func (*Rom) TableName() string {
	return "rom"
}

//插入rom数据
func (m *Rom) BatchAdd(uniqs []string, romlist map[string]*Rom) {

	if len(uniqs) == 0 {
		return
	}
	tx := getDb().Begin()
	count := len(uniqs)
	for k, md5 := range uniqs {
		v := romlist[md5]
		tx.Create(&v)
		if k%500 == 0 {
			utils.Loading("开始写入缓存("+utils.ToString(k+1)+"/"+utils.ToString(count)+")", "")
		}
	}
	tx.Commit()
}

//根据条件，查询多条数据
func (*Rom) Get(pages int, platform uint32, menu string, keyword string, baseType string, basePublisher string, baseYear string, baseCountry string) ([]*Rom, error) {

	volist := []*Rom{}
	where := map[string]interface{}{}
	where["hide"] = 0
	if platform != 0 {
		where["platform"] = platform
	}

	if menu != "" {
		if menu == "favorite" {
			where["star"] = 1
		} else if menu == "hide" {
			where["hide"] = 1
		} else {
			where["menu"] = menu
		}

	}
	where["pname"] = ""

	if baseType != "" {
		where["base_type"] = baseType
	}

	if basePublisher != "" {
		where["base_publisher"] = basePublisher
	}
	if baseYear != "" {
		where["base_year"] = baseYear
	}
	if baseCountry != "" {
		where["base_country"] = baseCountry
	}
	likeWhere := ""
	if keyword != "" {
		likeWhere = `name LIKE "%` + keyword + `%"`
	}

	offset := pages * ROM_PAGE_NUM
	field := "id,name,menu,platform,rom_path,base_type,base_year,base_publisher,base_country"
	result := getDb().Select(field).Where(where).Where(likeWhere).Order("pinyin ASC").Limit(ROM_PAGE_NUM).Offset(offset).Find(&volist)
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

	result := getDb().Select("id,name,pname,rom_path").Where("platform=? AND pname=?", platform, pname).Order("pinyin ASC").Find(&volist)
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
	field := "id,name,menu,platform,rom_path,base_type,base_year,base_publisher,base_country"
	result := getDb().Select(field).Order("pinyin ASC").Limit(ROM_PAGE_NUM).Offset(offset)
	if keyword == "#" {

		//查询0-9数字rom
		subWhere := "pinyin LIKE '0%'"
		for i := 1; i <= 9; i++ {
			subWhere += " OR pinyin LIKE '" + utils.ToString(i) + "%'"
		}
		result.Where(where).Where(subWhere).Find(&volist)
	} else {
		result.Where(where).Where("pinyin LIKE ?", keyword+"%").Find(&volist)
	}

	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return volist, result.Error
}

//根据满足条件的rom数量
func (m *Rom) Count(platform uint32, menu string, keyword string, baseType string, basePublisher string, baseYear string, baseCountry string) (int, error) {
	count := 0
	where := map[string]interface{}{
	}

	if platform != 0 {
		where["platform"] = platform
		where["pname"] = ""
	}

	if menu != "" {
		if menu == "hide" {
			where["hide"] = 1
		} else if menu == "favorite" {
			where["star"] = 1
		} else {
			where["menu"] = menu
		}
	}

	if baseType != "" {
		where["base_type"] = baseType
	}

	if basePublisher != "" {
		where["base_publisher"] = basePublisher
	}
	if baseYear != "" {
		where["base_year"] = baseYear
	}
	if baseCountry != "" {
		where["base_country"] = baseCountry
	}
	likeWhere := ""
	if keyword != "" {
		likeWhere = `name LIKE "%` + keyword + `%"`
	}

	result := getDb().Table(m.TableName()).Where(where).Where(likeWhere).Count(&count)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return count, result.Error
}

//更新名称
func (m *Rom) UpdateName() error {

	create := map[string]interface{}{
		"name":     m.Name,
		"pinyin":   m.Pinyin,
		"rom_path": m.RomPath,
	}

	vo := &Rom{}
	result := getDb().Select("platform,name").Table(m.TableName()).Where("id=?", m.Id).First(&vo)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	result = getDb().Table(m.TableName()).Where("id=?", m.Id).Updates(create)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	roms := []*Rom{}
	result = getDb().Table(m.TableName()).Where("platform=? AND pname=?", vo.Platform, vo.Name).Find(&roms)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	if len(roms) > 0 {
		for _, v := range roms {
			newName := strings.Replace(v.RomPath, vo.Name+"__", m.Name+"__", 1)
			createSub := map[string]string{
				"pname":    m.Name,
				"rom_path": newName,
			}

			result = getDb().Table(m.TableName()).Where("id=?", v.Id).Updates(createSub)
		}
	}
	return result.Error
}

//更新喜爱状态
func (m *Rom) UpdateStar() error {
	result := getDb().Table(m.TableName()).Where("id=?", m.Id).Update("star", m.Star)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

//更新隐藏状态
func (m *Rom) UpdateHide() error {
	result := getDb().Table(m.TableName()).Where("id=?", m.Id).Update("hide", m.Hide)
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

//删除一个rom
func (m *Rom) DeleteById(id uint64) error {
	result := getDb().Where("id=? ", id).Delete(&m)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

//删除一个平台下的所有rom数据
func (m *Rom) DeleteByPlatform() error {
	result := getDb().Where("platform=? ", m.Platform).Delete(&m)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

func (m *Rom) DeleteByMd5(platform uint32, uniqs []string) error {

	if len(uniqs) == 0 {
		return nil
	}

	sql := ""
	subsql := ""
	count := len(uniqs)
	for k, uniq := range uniqs {
		subsql += uniq + "','"
		if k%990 == 0 {
			sql = "DELETE FROM rom where path_md5 in ('" + subsql + "')"
			tx := getDb().Begin()
			tx.Exec(sql)
			result := tx.Commit()
			if result.Error != nil {
				fmt.Println(result.Error)
			}
			subsql = ""
		}
		if k%500 == 0 {
			utils.Loading("开始清理缓存("+utils.ToString(k+1)+"/"+utils.ToString(count)+")", "")
		}
	}

	//删除剩余数据
	if subsql != "" {
		sql = "DELETE FROM rom where path_md5 in ('" + subsql + "')"
		tx := getDb().Begin()
		tx.Exec(sql)
		result := tx.Commit()
		if result.Error != nil {
			fmt.Println(result.Error)
		}
	}

	return nil
}

//读取一个平台下的所有md5
func (sim *Rom) GetMd5ByPlatform(platform uint32) ([]string, error) {
	volist := []*Rom{}
	result := getDb().Select("path_md5").Where("platform=?", platform).Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	md5List := []string{}
	for _, v := range volist {
		md5List = append(md5List, v.InfoMd5)
	}
	return md5List, result.Error
}

//读取一个过滤器分类
func (sim *Rom) GetFilter(platform uint32, t string) ([]string, error) {
	volist := []*Rom{}

	result := getDb().Select(t).Where("platform = " + utils.ToString(platform) + " AND " + t + " != ''").Group(t).Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	create := []string{}

	switch t {
	case "base_type":
		for _, v := range volist {
			create = append(create, v.BaseType)
		}
		break
	case "base_year":
		for _, v := range volist {
			create = append(create, v.BaseYear)
		}
		break
	case "base_publisher":
		for _, v := range volist {
			create = append(create, v.BasePublisher)
		}
		break
	case "base_country":
		for _, v := range volist {
			create = append(create, v.BaseCountry)
		}
		break
	}

	return create, result.Error
}

//更新喜爱状态
func (m *Rom) UpdateRomBase(id uint64) error {

	create := map[string]string{
		"base_type":      m.BaseType,
		"base_year":      m.BaseYear,
		"base_publisher": m.BasePublisher,
		"base_country":   m.BaseCountry,
		"name":           m.Name,
		"pinyin":         utils.TextToPinyin(m.Name),
	}

	/*if m.Name != "" {
		create["name"] = m.Name
		create["pinyin"] = utils.TextToPinyin(m.Name)
	}*/

	result := getDb().Table(m.TableName()).Where("id=?", id).Updates(create)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

//删除不存在的平台下的所有rom
func (m *Rom) ClearByPlatform(platforms []string) error {
	result := getDb().Where("platform not in (?)", platforms).Delete(&m)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

//清空表数据
func (m *Rom) Truncate() error {
	result := getDb().Delete(&m)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}
