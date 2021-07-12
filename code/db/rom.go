package db

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"simUI/code/utils"
	"strings"
)

var ROM_PAGE_NUM = 80 //每页加载rom数量
var TIP_NUM = 100     //每tip个更新提示一次

type Rom struct {
	Id            uint64
	Pname         string // 所属主游戏
	Menu          string // 菜单名称
	Name          string // 游戏名称
	Platform      uint32 // 平台
	RomPath       string // rom路径
	SimId         uint32 // 正在使用的模拟器id
	SimConf       string // 模拟器参数独立配置
	Star          uint8  // 喜好，星级
	Hide          uint8  // 是否隐藏
	BaseType      string // 游戏类型，如RPG
	BaseYear      string // 游戏年份
	BasePublisher string // 游戏出品公司
	BaseCountry   string // 游戏国家
	BaseTranslate string // 汉化组
	BaseVersion   string // 版本
	Pinyin        string // 拼音索引
	InfoMd5       string // 信息md5，包含资料信息
	FileMd5       string // 文件md5，仅包含平台和文件名
}

func (*Rom) TableName() string {
	return "rom"
}

//插入rom数据

func (m *Rom) BatchUpdate(romlist []*Rom) {
	tx := getDb().Begin()
	create := map[string]string{}
	count := len(romlist)
	utils.Loading("3/3开始更新缓存(1/"+utils.ToString(count)+")", "")
	for k, v := range romlist {
		create = map[string]string{
			"menu":           v.Menu,
			"name":           v.Name,
			"pname":          v.Pname,
			"rom_path":       v.RomPath,
			"base_type":      v.BaseType,
			"base_year":      v.BaseYear,
			"base_publisher": v.BasePublisher,
			"base_country":   v.BaseCountry,
			"base_translate": v.BaseTranslate,
			"base_version":   v.BaseVersion,
			"pinyin":         v.Pinyin,
			"info_md5":       v.InfoMd5,
		}
		getDb().Table(m.TableName()).Where("file_md5 = ?", v.FileMd5).Update(create)
		if k%TIP_NUM == 0 {
			utils.Loading("3/3开始更新缓存("+utils.ToString(k+1)+"/"+utils.ToString(count)+")", "")
		}
	}
	tx.Commit()
}

func (m *Rom) BatchAdd(romlist []*Rom) {

	if len(romlist) == 0 {
		return
	}

	repeatList := map[string]bool{}
	count := len(romlist)
	utils.Loading("[2/3]开始写入缓存(1/"+utils.ToString(count)+")", "")
	tx := getDb().Begin()
	for k, v := range romlist {

		//如果存在则不写入
		if _, ok := repeatList[v.InfoMd5]; ok {
			continue
		}

		tx.Create(&v)
		//记录md5数据
		repeatList[v.InfoMd5] = true;

		if k%TIP_NUM == 0 {
			utils.Loading("[2/3]开始写入缓存("+utils.ToString(k+1)+"/"+utils.ToString(count)+")", "")
		}
	}
	tx.Commit()
}

//根据条件，查询多条数据
func (*Rom) Get(showHide uint8, pages int, platform uint32, menu string, keyword string, baseType string, basePublisher string, baseYear string, baseCountry string, baseTranslate string, baseVersion string) ([]*Rom, error) {

	volist := []*Rom{}
	where := map[string]interface{}{}
	if showHide == 0 {
		where["hide"] = 0
	}
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
	if baseTranslate != "" {
		where["base_translate"] = baseTranslate
	}
	if baseVersion != "" {
		where["base_version"] = baseVersion
	}
	likeWhere := ""
	if keyword != "" {
		likeWhere = `name LIKE "%` + keyword + `%"`
	}

	offset := pages * ROM_PAGE_NUM
	result := getDb().Select("*").Where(where).Where(likeWhere).Order("pinyin ASC").Limit(ROM_PAGE_NUM).Offset(offset).Find(&volist)
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

	result := getDb().Select("id,name,pname,star,rom_path").Where("platform=? AND pname=?", platform, pname).Order("pinyin ASC").Find(&volist)
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

//根据id查询一条数据
func (*Rom) GetByIds(ids []uint64) ([]*Rom, error) {

	vo := []*Rom{}

	result := getDb().Where("id in(?)", ids).Find(&vo)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return vo, result.Error
}

//根据拼音筛选
func (*Rom) GetByPinyin(showHide uint8, pages int, platform uint32, menu string, keyword string) ([]*Rom, error) {
	where := map[string]interface{}{}

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
	if showHide == 0 {
		where["hide"] = 0
	}

	where["pname"] = ""
	offset := pages * ROM_PAGE_NUM
	volist := []*Rom{}
	result := getDb().Select("*").Order("pinyin ASC").Limit(ROM_PAGE_NUM).Offset(offset)
	if keyword == "#" {

		//查询0-9数字rom
		subWhere := "pinyin LIKE '0%'"
		for i := 1; i <= 9; i++ {
			subWhere += " OR pinyin LIKE '" + utils.ToString(i) + "%'"
		}
		result.Where(where).Where(subWhere).Find(&volist)
	} else {
		keyword = strings.ToLower(keyword)
		result.Where(where).Where("pinyin LIKE ?", keyword+"%").Find(&volist)
	}

	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return volist, result.Error
}

//根据满足条件的rom数量
func (m *Rom) Count(showHide uint8, platform uint32, menu string, keyword string, baseType string, basePublisher string, baseYear string, baseCountry string, baseTranslate string, baseVersion string) (int, error) {
	count := 0
	where := map[string]interface{}{
	}

	if platform != 0 {
		where["platform"] = platform
		where["pname"] = ""
	}

	if showHide == 0 {
		where["hide"] = 0
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
	if baseTranslate != "" {
		where["base_translate"] = baseTranslate
	}
	if baseTranslate != "" {
		where["base_version"] = baseVersion
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

//根据平台id，查询数据
func (*Rom) GetByPlatform(platform uint32) ([]*Rom, error) {

	volist := []*Rom{}
	result := getDb().Select("*").Where("platform = ?", platform).Order("pinyin ASC").Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return volist, result.Error
}

//根据平台id，查询主游戏
func (*Rom) GetMasterRomByPlatform(platform uint32) ([]*Rom, error) {

	volist := []*Rom{}
	result := getDb().Select("*").Where("platform = ? AND pname = ''", platform).Order("pinyin ASC").Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return volist, result.Error
}

//读取相关游戏
func (*Rom) GetRelatedGames(id uint64) ([]*Rom, error) {

	vo, _ := (&Rom{}).GetById(id)
	platform := vo.Platform
	baseType := vo.BaseType

	volist := []*Rom{}
	//先读取同类型游戏
	result := getDb().Select("*").Where("id != ? AND platform = ? AND pname='' AND  base_type = ?", id, platform, baseType).Order("random()").Limit(6).Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	//如果找不到同类型游戏，则读取平台下的随机游戏
	if len(volist) == 0 {
		result = getDb().Select("*").Where("id != ? AND platform = ? AND  pname=''", id, platform).Order("random()").Limit(6).Find(&volist)
		if result.Error != nil {
			fmt.Println(result.Error)
		}
	}

	return volist, result.Error

}

//读取相关游戏
func (*Rom) GetIdsByNames(platform uint32, names []string) (map[string]uint64, error) {

	volist := []*Rom{}
	//先读取同类型游戏
	result := getDb().Select("id,name").Where("platform = ? AND name in (?)", platform, names).Find(&volist)
	if result.Error != nil {
		return nil, result.Error
	}
	data := map[string]uint64{}
	for _, v := range volist {
		data[v.Name] = v.Id
	}

	return data, result.Error

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

//批量更新隐藏状态
func (m *Rom) UpdateHideByIds(ids []uint64, hide uint8) error {
	if len(ids) == 0 {
		return nil
	}
	result := getDb().Table(m.TableName()).Where("id in (?)", ids).Update("hide", hide)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

//批量更新喜好状态
func (m *Rom) UpdateStarByIds(ids []uint64, star uint8) error {
	if len(ids) == 0 {
		return nil
	}
	result := getDb().Table(m.TableName()).Where("id in (?)", ids).Update("star", star)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

//更新模拟器
func (m *Rom) UpdateSimIdByIds(romIds []string, simId uint32) error {
	result := getDb().Table(m.TableName()).Where("id in (?)", romIds).Update("sim_id", simId)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

//更新rom地址
func (m *Rom) UpdateRomPath() error {
	result := getDb().Table(m.TableName()).Where("id=?", m.Id).Update("rom_path", m.RomPath)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

//更新子游戏pname
func (m *Rom) UpdateSubRomPname(oldName string, newName string) error {
	result := getDb().Table(m.TableName()).Where("pname=?", oldName).Update("pname", newName)
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

//删除一个rom
func (m *Rom) DeleteSubRom(pname string) error {
	result := getDb().Where("pname=? ", pname).Delete(&m)
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
	utils.Loading("[1/3]开始清理缓存(1/"+utils.ToString(count)+")", "")
	for k, uniq := range uniqs {
		subsql += uniq + "','"
		if k%990 == 0 {
			sql = "DELETE FROM rom where file_md5 in ('" + subsql + "')"
			tx := getDb().Begin()
			tx.Exec(sql)
			result := tx.Commit()
			if result.Error != nil {
				fmt.Println(result.Error)
			}
			subsql = ""
		}
		if k%TIP_NUM == 0 {
			utils.Loading("[1/3]开始清理缓存("+utils.ToString(k+1)+"/"+utils.ToString(count)+")", "")
		}
	}

	//删除剩余数据
	if subsql != "" {
		sql = "DELETE FROM rom where file_md5 in ('" + subsql + "')"
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
func (sim *Rom) GetMd5ByPlatform(platform uint32) ([]string, []string, error) {
	volist := []*Rom{}
	result := getDb().Select("file_md5,info_md5").Where("platform=?", platform).Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	infoMd5List := []string{}
	fileMd5List := []string{}
	for _, v := range volist {
		infoMd5List = append(infoMd5List, v.InfoMd5)
		fileMd5List = append(fileMd5List, v.FileMd5)
	}
	return fileMd5List, infoMd5List, result.Error
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
	case "base_translate":
		for _, v := range volist {
			create = append(create, v.BaseTranslate)
		}
		break
	case "base_version":
		for _, v := range volist {
			create = append(create, v.BaseVersion)
		}
		break
	}

	return create, result.Error
}

//更新游戏资料
func (m *Rom) UpdateRomBase(id uint64) error {

	create := map[string]string{
		"base_type":      m.BaseType,
		"base_year":      m.BaseYear,
		"base_publisher": m.BasePublisher,
		"base_country":   m.BaseCountry,
		"base_translate": m.BaseTranslate,
		"base_version":   m.BaseVersion,
		"name":           m.Name,
		"pinyin":         utils.TextToPinyin(m.Name),
	}

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
