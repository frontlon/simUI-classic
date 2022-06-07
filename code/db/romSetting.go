package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"math"
	"simUI/code/utils"
)

type RomSetting struct {
	Id          uint64
	Platform    uint32 // 平台
	FileMd5     string // file_md5
	Star        uint8  // 喜好，星级
	Hide        uint8  // 是否隐藏
	RunNum      uint64 //运行次数
	RunLasttime int64  //最后运行时间
	SimId       uint32 // 正在使用的模拟器id
	SimConf     string // 模拟器参数独立配置
}

func (*RomSetting) TableName() string {
	return "rom_setting"
}

//批量添加
func (m *RomSetting) BatchAdd(romlist []*RomSetting) {

	if len(romlist) == 0 {
		return
	}
	tx := getDb().Begin()
	for _, v := range romlist {
		tx.Create(&v)
	}
	tx.Commit()
}

//用新的file_md5替换旧的pname和file_md5
func (m *RomSetting) BatchUpdateFileMd5(platform uint32,lists []map[string]string) error {

	if len(lists) == 0 {
		return nil
	}

	tx := getDb().Begin()
	for _, rom := range lists {
		tx.Table(m.TableName()).Where("platform = ? AND file_md5 = ?", platform,rom["oldMd5"]).Update(map[string]interface{}{"file_md5": rom["newMd5"]})
	}
	if err := tx.Commit().Error; err != nil {
		fmt.Println("update错误", err)
	}
	return nil
}

//根据平台id查询数据
func (*RomSetting) GetByPlatform(platform uint32) ([]*RomSetting, error) {

	volist := []*RomSetting{}

	result := getDb().Where("platform=?", platform).Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return volist, result.Error
}

//根据平台id查询数据，返回map
func (*RomSetting) GetByPlatformToMap(platform uint32) (map[string]*RomSetting, error) {

	vo := []*RomSetting{}

	result := getDb().Where("platform=?", platform).Find(&vo)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	//转换成map
	volist := map[string]*RomSetting{}
	for _, v := range vo {
		volist[v.FileMd5] = v
	}

	return volist, result.Error
}

//根据平台id查询file_md5
func (*RomSetting) GetFileMd5ByPlatform(platform uint32) ([]string, error) {

	volist := []*RomSetting{}

	result := getDb().Select("file_md5").Where("platform=?", platform).Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	md5List := []string{}
	for _, v := range volist {
		md5List = append(md5List, v.FileMd5)
	}

	return md5List, result.Error
}

//更新喜爱状态
func (m *RomSetting) UpdateStar() error {
	//初始化数据
	_ = m.InitData(m.Platform, m.FileMd5)

	//更新数据
	result := getDb().Table(m.TableName()).Where("platform=? AND file_md5 = ?", m.Platform, m.FileMd5).Update("star", m.Star)
	if result.Error != nil {
		fmt.Println(result.Error)
		return result.Error
	}

	return nil
}

//更新隐藏状态
func (m *RomSetting) UpdateHide() error {

	//初始化数据
	_ = m.InitData(m.Platform, m.FileMd5)

	//更新数据
	result := getDb().Table(m.TableName()).Where("platform=? AND file_md5 = ?", m.Platform, m.FileMd5).Update("hide", m.Hide)
	if result.Error != nil {
		fmt.Println(result.Error)
		return result.Error
	}

	return nil
}

//更新运行次数和最后运行时间
func (m *RomSetting) UpdateRunNumAndTime() error {

	//初始化数据
	_ = m.InitData(m.Platform, m.FileMd5)

	//更新数据
	create := map[string]interface{}{
		"run_num":      gorm.Expr("run_num + 1"),
		"run_lasttime": m.RunLasttime,
	}

	result := getDb().Table(m.TableName()).Where("platform=? AND file_md5 = ?", m.Platform, m.FileMd5).Updates(create)
	if result.Error != nil {
		fmt.Println(result.Error)
		return result.Error
	}

	return nil
}

//rom更换模拟器
func (m *RomSetting) UpdateSimIds(platform uint32, fileMd5List []string, simId uint32) error {

	//初始化数据
	for _, fileMd5 := range fileMd5List {
		_ = m.InitData(platform, fileMd5)
	}

	//更新数据
	result := getDb().Table(m.TableName()).Where("platform=? AND file_md5 in (?)", platform, fileMd5List).Update("sim_id", simId)
	if result.Error != nil {
		fmt.Println(result.Error)
		return result.Error
	}

	return nil
}

//更新rom 模拟器配置
func (m *RomSetting) UpdateSimConf() error {

	//初始化数据
	_ = m.InitData(m.Platform, m.FileMd5)

	//更新数据
	result := getDb().Table(m.TableName()).Where("platform=? AND file_md5 = ?", m.Platform, m.FileMd5).Update("sim_conf", m.SimConf)
	if result.Error != nil {
		fmt.Println(result.Error)
		return result.Error
	}

	return nil
}

//批量更新隐藏状态
func (m *RomSetting) UpdateHideByFileMd5(platform uint32, md5s []string, hide uint8) error {

	if len(md5s) == 0 {
		return nil
	}

	//先读取出有的数据
	volist := []*RomSetting{}
	getDb().Table(m.TableName()).Select("file_md5").Where("platform =? AND file_md5 in (?)", platform, md5s).Find(&volist)

	Md5List := []string{}
	for _, v := range volist {
		Md5List = append(Md5List, utils.GetFileName(v.FileMd5))
	}

	isset := []string{}
	notIsset := []string{}
	for _, name := range md5s {
		if utils.InSliceString(name, Md5List) {
			isset = append(isset, name)
		} else {
			notIsset = append(notIsset, name)
		}
	}

	//更新老数据
	if len(isset) > 0 {
		result := getDb().Table(m.TableName()).Where("platform =? AND file_md5 in (?)", platform, isset).Update("hide", hide)
		if result.Error != nil {
			fmt.Println(result.Error)
			return result.Error
		}
	}
	tx := getDb().Begin()

	//写入新数据
	if len(notIsset) > 0 {
		for _, md5 := range notIsset {
			v := &RomSetting{
				FileMd5:  md5,
				Platform: platform,
				Hide:     hide,
			}
			tx.Create(&v)
		}
	}
	tx.Commit()

	return nil
}

//删除记录
func (m *RomSetting) DeleteByFileMd5s(platform uint32, md5List []string) error {

	if len(md5List) == 0 {
		return nil
	}

	listLen := len(md5List)

	ceil := int(math.Ceil(float64(listLen) / float64(maxVar)))

	for i := 0; i < ceil; i++ {
		start := i * maxVar
		end := (i + 1) * maxVar
		if end > listLen {
			end = listLen
		}
		list := md5List[start:end]
		getDb().Where("platform = ? AND file_md5 in (?)", platform, list).Delete(&m)
	}

	return nil
}

//初始化数据，如果没有数据，则生成一条
func (m *RomSetting) InitData(platform uint32, fileMd5 string) error {
	count := 0
	getDb().Table(m.TableName()).Where("platform=? AND file_md5 = ?", platform, fileMd5).Count(&count)

	if count == 0 {
		create := &RomSetting{
			Platform: platform,
			FileMd5:  fileMd5,
		}
		result := getDb().Create(&create)
		if result.Error != nil {
			fmt.Println(result.Error)
			return result.Error
		}
	}
	return nil
}

//删除不存在的平台下的所有数据
func (*RomSetting) ClearByNotPlatform(platforms []string) error {
	m := &RomSetting{}
	result := getDb().Not("platform", platforms).Delete(&m)

	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return result.Error
}
