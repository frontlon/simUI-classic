package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"simUI/code/utils"
	"strings"
	"time"
)

var ROM_PAGE_NUM = 100 //首页每页加载rom数量
var TIP_NUM = 200      //每tip个更新提示一次(更新缓存)

type Rom struct {
	Id            uint64
	Pname         string // 所属主游戏
	Menu          string // 菜单名称
	Name          string // 中文名
	Platform      uint32 // 平台
	RomPath       string // rom路径
	SimId         uint32 // 正在使用的模拟器id
	SimConf       string // 模拟器参数独立配置
	Star          uint8  // 喜好（起名不对，初期设计的遗留问题）
	Score         string // 评分
	Hide          uint8  // 是否隐藏
	Size          string // rom文件大小
	BaseNameEn    string // 英文名
	BaseNameJp    string // 日文名
	BaseType      string // 游戏类型，如RPG
	BaseYear      string // 游戏年份
	BaseProducer  string // 游戏出品公司
	BasePublisher string // 游戏出品公司
	BaseCountry   string // 游戏国家
	BaseTranslate string // 汉化组
	BaseVersion   string // 版本
	BaseOtherA    string // 其他信息A
	BaseOtherB    string // 其他信息B
	BaseOtherC    string // 其他信息C
	BaseOtherD    string // 其他信息D
	Pinyin        string // 拼音索引
	InfoMd5       string // 信息md5，包含资料信息
	FileMd5       string // 文件md5，仅包含平台和文件名
	RunNum        uint64 // 运行次数
	RunLasttime   int64  // 最后运行时间
	Complete      uint8  // 通关状态(0未通关;1已通关;2完美通关)
	SubGames      []*Rom //子游戏列表
}

func (*Rom) TableName() string {
	return "rom"
}

//插入rom数据
func (m *Rom) BatchAdd(romlist []*Rom, showLoading int) {

	if len(romlist) == 0 {
		return
	}

	count := len(romlist)
	if showLoading == 1 {
		utils.Loading("[2/3]开始写入缓存(1/"+utils.ToString(count)+")", "")
	}
	tx := getDb().Begin()
	for k, v := range romlist {

		tx.Create(&v)
		//记录md5数据

		if showLoading == 1 && k%TIP_NUM == 0 {
			utils.Loading("[2/3]开始写入缓存("+utils.ToString(k+1)+"/"+utils.ToString(count)+")", "")
		}
	}
	tx.Commit()
}

//用新的file_md5替换旧的pname和file_md5
func (m *Rom) BatchUpdateFileMd5(platform uint32, lists []map[string]string) error {

	if len(lists) == 0 {
		return nil
	}

	tx := getDb().Begin()
	for _, rom := range lists {
		tx.Table(m.TableName()).Where("platform = ? AND rom_path = ?", platform, rom["romPath"]).Update(map[string]interface{}{"file_md5": rom["newMd5"]})
		tx.Table(m.TableName()).Where("platform = ? AND pname = ?", platform, rom["oldMd5"]).Update(map[string]interface{}{"pname": rom["newMd5"]})
	}
	if err := tx.Commit().Error; err != nil {
		fmt.Println("update错误", err)
	}
	return nil
}

//根据条件，查询多条数据

func (m *Rom) Get(showSubGame uint8, showHide uint8, pages int, platform uint32, menu string, keyword string, baseType string, basePublisher string, baseYear string, baseCountry string, baseTranslate string, baseVersion string, baseProducer string, score string, complete string) ([]*Rom, error) {

	volist := []*Rom{}
	where := map[string]interface{}{}
	likeWhere := "1=1 "

	if platform != 0 {
		where["platform"] = platform
	}

	where["pname"] = ""

	if showHide == 0 {
		where["hide"] = showHide
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

	if baseType != "" {
		likeWhere += ` AND base_type LIKE "%` + baseType + `%"`
	}
	if baseProducer != "" {
		likeWhere += ` AND base_producer LIKE "%` + baseProducer + `%"`
	}
	if basePublisher != "" {
		likeWhere += ` AND base_publisher LIKE "%` + basePublisher + `%"`
	}
	if baseYear != "" {
		likeWhere += ` AND base_year LIKE "` + baseYear + `%"`
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

	if score != "" {
		where["score"] = score
	}

	if complete != "" {
		where["complete"] = complete
	}
	if keyword != "" {
		likeWhere += `AND (name LIKE "%` + keyword + `%" or rom_path LIKE "%` + keyword + `%" or pinyin LIKE "%` + keyword + `%")`
	}

	offset := pages * ROM_PAGE_NUM
	conf, _ := (&Config{}).GetField("romlist_orders")
	sort := "pinyin ASC"
	switch conf.RomlistOrders {
	case "1":
		sort = "pinyin ASC"
	case "2":
		sort = "pinyin DESC"
	case "3":
		sort = "score ASC,pinyin ASC"
	case "4":
		sort = "score DESC,pinyin ASC"
	case "5":
		sort = "base_year ASC,pinyin ASC"
	case "6":
		sort = "base_year DESC,pinyin ASC"
	}

	result := getDb().Select("*").Where(where).Where(likeWhere).Order(sort).Limit(ROM_PAGE_NUM).Offset(offset).Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	//加载子游戏
	if showSubGame == 1 {
		md5s := []string{}
		for _, v := range volist {
			md5s = append(md5s, v.FileMd5)
		}
		subRoms, _ := m.GetSubRomByFileMd5s(md5s)
		for k, v := range volist {
			volist[k].SubGames = subRoms[v.FileMd5]
		}
	}

	return volist, result.Error
}

//读取没有子游戏的rom
func (m *Rom) GetNotSubRom(pages int, showHide uint8, platform uint32, menu string, keyword string) ([]*Rom, error) {

	volist := []*Rom{}
	where := map[string]interface{}{}

	if platform != 0 {
		where["platform"] = platform
	}

	if showHide == 0 {
		where["hide"] = showHide
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

	likeWhere := ""
	if keyword != "" {
		likeWhere = `(name LIKE "%` + keyword + `%" or rom_path LIKE "%` + keyword + `%" or pinyin LIKE "%` + keyword + `%")`
	}

	where["pname"] = ""

	offset := pages * ROM_PAGE_NUM

	//查出重复rom
	repeat := []*Rom{}
	repeatFileMd5 := []string{}
	getDb().Select("pname").Where("platform = ? AND pname != ''", platform).Where(likeWhere).Group("pname").Find(&repeat)
	for _, v := range repeat {
		repeatFileMd5 = append(repeatFileMd5, v.Pname)
	}

	//查询rom数据
	if len(repeatFileMd5) > 0 {
		getDb().Select("*").Where(where).Where(likeWhere).Where("file_md5 not in (?)", repeatFileMd5).Order("pinyin ASC").Offset(offset).Limit(ROM_PAGE_NUM).Find(&volist)
	} else {
		getDb().Select("*").Where(where).Where(likeWhere).Order("pinyin ASC").Offset(offset).Limit(ROM_PAGE_NUM).Find(&volist)
	}

	return volist, nil
}

//读取子rom
func (*Rom) GetSubRom(platform uint32, fileMd5 string) ([]*Rom, error) {

	volist := []*Rom{}

	if platform == 0 || fileMd5 == "" {
		return volist, nil
	}

	result := getDb().Select("id,name,pname,rom_path").Where("platform=? AND pname=?", platform, fileMd5).Order("pinyin ASC").Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return volist, result.Error
}

//批量查询子游戏
func (*Rom) GetSubRomByFileMd5s(fileMd5s []string) (map[string][]*Rom, error) {

	volist := []*Rom{}

	if len(fileMd5s) == 0 {
		return map[string][]*Rom{}, nil
	}

	result := getDb().Select("id,name,pname,rom_path").Where("pname in (?)", fileMd5s).Order("pinyin ASC").Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	//封装成map
	create := map[string][]*Rom{}
	for _, v := range volist {
		create[v.Pname] = append(create[v.Pname], v)
	}

	return create, result.Error
}

//根据id查询一条数据
func (*Rom) GetById(id uint64) (*Rom, error) {

	vo := &Rom{}

	result := getDb().Where("id=?", id).First(&vo)
	if result.Error != nil {
		fmt.Println(result.Error)
		return nil, result.Error
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

//根据file_md5查询一条数据
func (*Rom) GetByFileMd5(fileMd5 string) (*Rom, error) {

	vo := &Rom{}
	result := getDb().Where("file_md5 = ?", fileMd5).Find(&vo)
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

	where["hide"] = showHide

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

//读取一个平台下的所有md5
func (m *Rom) GetMd5ByPlatform(platform uint32) ([]string, []string, error) {
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

//读取一个平台下的file md5
func (m *Rom) GetFileMd5ByPlatform(platform uint32) (map[string]string, error) {
	volist := []*Rom{}
	result := getDb().Select("rom_path,file_md5").Where("platform=?", platform).Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	fileMd5List := map[string]string{}
	for _, v := range volist {
		fileMd5List[utils.GetFileNameAndExt(v.RomPath)] = v.FileMd5
	}
	return fileMd5List, result.Error
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
	volist := []*Rom{}

	result := getDb().Select("*").Where("id != ? AND platform = ? AND pname='' AND hide = 0 ", id, vo.Platform).Order("run_num desc").Limit(6).Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return volist, result.Error
}

//根据name读取id
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

//根据目录名读取rom
func (*Rom) GetByMenu(platform uint32, menu string) ([]*Rom, error) {

	volist := []*Rom{}

	result := getDb().Where("platform = ? AND menu = ?", platform, menu).Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return volist, result.Error
}

//根据目录名列表读取rom
func (*Rom) GetMasterByMenus(platform uint32, menus []string) ([]*Rom, error) {

	volist := []*Rom{}

	result := getDb().Where("platform = ? AND pname = '' and menu in (?)", platform, menus).Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return volist, result.Error
}

//根据目录列表读取rom
func (r *Rom) GetMasterAndSubByMenus(platform uint32, menus []string) ([]*Rom, error) {

	vo := []*Rom{}
	t := r.TableName()
	subSql := `select file_md5 from ` + t + ` where platform = ? and menu in (?) and pname =''`
	sql := `select * from ` + t + ` where (platform = ? and menu in (?)) or pname in(` + subSql + `)`
	err := getDb().Raw(sql, platform, menus, platform, menus).Scan(&vo).Error
	if err != nil {
		return nil, err
	}

	return vo, err
}

//根据父id列表读取rom
func (*Rom) GetMasterAndSubByMasterIds(ids []uint64) ([]*Rom, error) {

	vo := []*Rom{}

	subSql := `select file_md5 from rom where id in (?) and pname =''`
	sql := `select * from rom where id in (?) or pname in(` + subSql + `)`

	err := getDb().Raw(sql, ids, ids).Scan(&vo).Error
	if err != nil {
		return nil, err
	}

	return vo, err
}

//根据满足条件的rom数量
func (m *Rom) Count(showHide uint8, platform uint32, menu string, keyword string, baseType string, basePublisher string, baseYear string, baseCountry string, baseTranslate string, baseVersion string, baseProducer string, score string, complete string) (int, error) {
	count := 0
	where := map[string]interface{}{}
	likeWhere := "1=1 "
	if platform != 0 {
		where["platform"] = platform
	}

	where["pname"] = ""

	if showHide == 0 {
		where["hide"] = showHide
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
		likeWhere += ` AND base_type LIKE "%` + baseType + `%"`
	}
	if baseProducer != "" {
		likeWhere += ` AND base_producer LIKE "%` + baseProducer + `%"`
	}
	if basePublisher != "" {
		likeWhere += ` AND base_publisher LIKE "%` + basePublisher + `%"`
	}
	if baseYear != "" {
		likeWhere += ` AND base_year LIKE "` + baseYear + `%"`
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
	if score != "" {
		where["score"] = score
	}
	if complete != "" {
		where["complete"] = complete
	}
	if keyword != "" {
		likeWhere += ` AND name LIKE "%` + keyword + `%"`
	}

	result := getDb().Table(m.TableName()).Where(where).Where(likeWhere).Count(&count)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return count, result.Error
}

//根据拼音筛选
func (m *Rom) CountByPinyin(showHide uint8, pages int, platform uint32, menu string, keyword string) (int, error) {
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

	where["hide"] = showHide

	where["pname"] = ""
	count := 0

	result := &gorm.DB{}
	if keyword == "#" {

		//查询0-9数字rom
		subWhere := "pinyin LIKE '0%'"
		for i := 1; i <= 9; i++ {
			subWhere += " OR pinyin LIKE '" + utils.ToString(i) + "%'"
		}
		result = getDb().Table(m.TableName()).Select("*").Where(where).Where(subWhere).Count(&count)
	} else {
		keyword = strings.ToLower(keyword)
		result = getDb().Table(m.TableName()).Select("*").Where(where).Where("pinyin LIKE ?", keyword+"%").Count(&count)
	}

	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return count, result.Error
}

//读取没有子游戏的rom
func (m *Rom) CountNotSubGame(showHide uint8, platform uint32, menu string, keyword string) (int, error) {

	count := 0
	where := map[string]interface{}{}

	if platform != 0 {
		where["platform"] = platform
	}

	if showHide == 0 {
		where["hide"] = showHide
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

	likeWhere := ""
	if keyword != "" {
		likeWhere = `name LIKE "%` + keyword + `%" or rom_path LIKE "%` + keyword + `%"`
	}

	where["pname"] = ""

	//查出重复rom
	repeat := []*Rom{}
	repeatFileMd5 := []string{}
	getDb().Select("pname").Where("platform = ? AND pname != ''", platform).Where(likeWhere).Group("pname").Find(&repeat)
	for _, v := range repeat {
		repeatFileMd5 = append(repeatFileMd5, v.Pname)
	}

	//查询rom数据
	if len(repeatFileMd5) > 0 {
		getDb().Table(m.TableName()).Where(where).Where(likeWhere).Where("file_md5 not in (?)", repeatFileMd5).Count(&count)
	} else {
		getDb().Table(m.TableName()).Where(where).Where(likeWhere).Count(&count)
	}

	return count, nil
}

//更新名称
func (m *Rom) UpdateName() error {

	create := map[string]interface{}{
		"name":     m.Name,
		"pinyin":   m.Pinyin,
		"rom_path": m.RomPath,
		"info_md5": m.InfoMd5,
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
func (m *Rom) UpdateSimIdByIds(romIds []uint64, simId uint32) error {
	result := getDb().Table(m.TableName()).Where("id in (?)", romIds).Update("sim_id", simId)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

//更新rom地址
func (m *Rom) UpdateRomPath() error {

	update := map[string]interface{}{
		"rom_path": m.RomPath,
		"menu":     m.Menu,
	}

	result := getDb().Table(m.TableName()).Where("id=?", m.Id).Update(update)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

//更新游戏资料
func (m *Rom) UpdateRomBase(id uint64) error {

	create := map[string]string{
		"base_type":      m.BaseType,
		"base_year":      m.BaseYear,
		"base_producer":  m.BaseProducer,
		"base_publisher": m.BasePublisher,
		"base_country":   m.BaseCountry,
		"base_translate": m.BaseTranslate,
		"base_version":   m.BaseVersion,
		"score":          m.Score,
		"base_name_en":   m.BaseNameEn,
		"base_name_jp":   m.BaseNameJp,
		"base_other_a":   m.BaseOtherA,
		"base_other_b":   m.BaseOtherB,
		"base_other_c":   m.BaseOtherC,
		"base_other_d":   m.BaseOtherD,
		"name":           m.Name,
		"pinyin":         utils.TextToPinyin(m.Name),
		"info_md5":       m.InfoMd5,
	}

	result := getDb().Table(m.TableName()).Where("id=?", id).Updates(create)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

//更新运行次数和最后运行时间
func (m *Rom) UpdateRunNumAndTime(id uint64) error {
	result := getDb().Table(m.TableName()).Where("id=?", id).Update("run_num", gorm.Expr("run_num + 1"))
	getDb().Table(m.TableName()).Where("id=?", id).Update("run_lasttime", time.Now().Unix())
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

//更新子游戏pname
func (m *Rom) ReplaceSubRomPname(platform uint32, oldName string, newName string) error {
	result := getDb().Table(m.TableName()).Where("platform = ? AND pname=?", platform, oldName).Update("pname", newName)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

//更新一个游戏的pane（设为子游戏）
func (m *Rom) UpdatePnameById(id uint64, newName string) error {
	result := getDb().Table(m.TableName()).Where("id = ?", id).Update("pname", newName)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

//更新目录
func (m *Rom) UpdateMenu(platform uint32, oldName string, newName string) error {
	result := getDb().Table(m.TableName()).Where("platform = ? AND menu = ? ", platform, oldName).Update("menu", newName)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

//更新评分
func (m *Rom) UpdateScore(id uint64, score string) error {

	vo := &Rom{}
	result := getDb().Where("id=?", id).First(&vo)
	if result.Error != nil {
		fmt.Println(result.Error)
		return result.Error
	}

	infoMd5 := utils.GetRomMd5(vo.Name, vo.RomPath, vo.BaseType, vo.BaseYear, vo.BaseProducer, vo.BasePublisher, vo.BaseCountry, vo.BaseTranslate, vo.BaseVersion, vo.BaseNameEn, vo.BaseNameJp, vo.BaseOtherA, vo.BaseOtherB, vo.BaseOtherC, vo.BaseOtherD, score, vo.Size)
	create := map[string]string{
		"score":    score,
		"info_md5": infoMd5,
	}

	update := getDb().Table(m.TableName()).Where("id = ?", id).Updates(create)
	if update.Error != nil {
		fmt.Println(update.Error)
	}
	return update.Error
}

//更新通关状态
func (m *Rom) UpdateComplete(id uint64, status uint8) error {
	result := getDb().Table(m.TableName()).Where("id = ?", id).Update("complete", status)
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
func (m *Rom) DeleteSubRom(platform uint32, pname string) error {
	result := getDb().Where("platform = ? AND pname=? ", platform, pname).Delete(&m)
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

func (m *Rom) DeleteByMd5(platform uint32, uniqs []string) {

	if len(uniqs) == 0 {
		return
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

	return
}

//删除重复rom
func (m *Rom) DeleteRepeat(platform uint32) {
	sql := "DELETE FROM rom WHERE platform = " + utils.ToString(platform) + " AND id NOT IN (SELECT max(id) FROM rom WHERE platform = " + utils.ToString(platform) + " GROUP BY info_md5)"
	sql2 := "DELETE FROM rom WHERE platform = " + utils.ToString(platform) + " AND id NOT IN (SELECT max(id) FROM rom WHERE platform = " + utils.ToString(platform) + " GROUP BY file_md5)"
	tx := getDb().Begin()
	tx.Exec(sql)
	tx.Exec(sql2)
	tx.Commit()
}

//删除不存在的平台下的所有rom
func (m *Rom) ClearByNotPlatform(platforms []string) error {
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

//有空游戏统计信息
func (m *Rom) TruncateGameStat() error {

	create := map[string]string{
		"run_lasttime": "0",
		"run_num":      "0",
	}

	result := getDb().Table(m.TableName()).Updates(create)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}
