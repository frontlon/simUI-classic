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
	SimId    uint32 // 正在使用的模拟器id
	SimConf  string // 模拟器参数独立配置
	RunNum   uint64 // 运行次数
	RunTime  uint32 // 最后运行时间
	Pinyin   string // 拼音索引
	PathMd5  string // 文件Md5
	FileId   string //唯一标识
}

func (*Rom) TableName() string {
	return "rom"
}

//插入rom数据
func (m *Rom) BatchAdd(uniqs []string, romlist map[string]*Rom) {

	if len(uniqs) == 0 {
		return
	}
	fmt.Println("开始批量写库", len(uniqs))
	tx := getDb().Begin()
	for _, md5 := range uniqs {
		v := romlist[md5]
		fmt.Println("v=", v)
		tx.Create(&v)
	}
	tx.Commit()
}

//根据fileid更新现有的rom
func (m *Rom) BatchUpdateByFileId(fileIds []string, romlist map[string]*Rom) error {
	tx := getDb().Begin()

	for _, v := range romlist {
		if utils.InSliceString(v.FileId, fileIds) {
			result := getDb().Where("file_id=?", v.FileId).Updates(&v)
			if result.Error != nil {
				fmt.Println(result.Error)
			}
		}
	}

	tx.Commit()
	return nil
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
		result.Where(where).Where("pinyin LIKE ?", keyword+"%").Find(&volist)
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

//更新名称
func (m *Rom) UpdateName(setType uint8) error {

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

			createSub := map[string]string{}
			if setType == 1 { //别名文件
				createSub = map[string]string{
					"pname": m.Name,
				}
			} else { //文件名
				createSub = map[string]string{
					"pname":    m.Name,
					"rom_path": newName,
				}
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

func (m *Rom) DeleteByMd5(platform uint32, uniqs []string) error {

	if len(uniqs) == 0 {
		return nil
	}

	sql := "";
	subsql := "";
	for k, uniq := range uniqs {
		subsql += uniq + "','";
		if k%990 == 0 {
			sql = "DELETE FROM rom where path_md5 in ('" + subsql + "')";
			tx := getDb().Begin()
			tx.Exec(sql)
			result := tx.Commit()
			if result.Error != nil {
				fmt.Println(result.Error)
			}
			subsql = ""
		}
	}

	//删除剩余数据
	if subsql != "" {
		sql = "DELETE FROM rom where path_md5 in ('" + subsql + "')";
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
		md5List = append(md5List, v.PathMd5)
	}
	return md5List, result.Error
}

//根据file_id取file_id，交集
func (sim *Rom) GetFileIdByFileId(platform uint32, fileIds []string) ([]string, error) {
	volist := []*Rom{}
	result := getDb().Select("file_id").Where("platform=?", platform).Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	fileIdList := []string{}
	for _, v := range volist {
		if utils.InSliceString(v.FileId, fileIds) {
			fileIdList = append(fileIdList, v.FileId)
		}
	}
	return fileIdList, result.Error
}

//删除不存在的平台下的所有rom
func (m *Rom) ClearByPlatform(platforms []string) (error) {
	result := getDb().Where("platform not in (?)", platforms).Delete(&m)
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
