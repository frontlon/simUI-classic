package db

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"simUI/code/utils"
)

type Platform struct {
	Id             uint32
	Name           string
	Icon           string
	RomExts        string
	RomPath        string
	ThumbPath      string
	SnapPath       string
	PosterPath     string
	PackingPath    string
	TitlePath      string
	BackgroundPath string
	DocPath        string
	StrategyPath   string
	VideoPath      string
	Romlist        string
	Rombase        string
	Pinyin         string
	Sort           uint32
	SimList        map[uint32]*Simulator `gorm:"-"` //模拟器列表
	UseSim         *Simulator            `gorm:"-"` //当前使用的模拟器
}

func (*Platform) TableName() string {
	return "platform"
}

//添加平台
func (m *Platform) Add() (uint32, error) {
	result := getDb().Create(&m)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return uint32(m.Id), result.Error
}

//批量插入数据
func (m *Platform) BatchAdd(platforms []*Platform) {

	if len(platforms) == 0 {
		return
	}

	tx := getDb().Begin()
	for _, v := range platforms {
		tx.Create(&v)
	}
	tx.Commit()
}

//根据条件，查询多条数据
func (*Platform) GetAll() ([]*Platform, error) {
	volist := []*Platform{}
	result := getDb().Order("sort ASC,pinyin ASC").Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return volist, nil
}

//根据ID查询一个平台参数
func (*Platform) GetById(id uint32) (*Platform, error) {
	vo := &Platform{}
	result := getDb().Where("id=?", id).Order("sort ASC,pinyin ASC").First(&vo)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return vo, result.Error
}

//更新平台信息
func (m *Platform) UpdateById() error {

	create := map[string]interface{}{
		"name":            m.Name,
		"icon":            m.Icon,
		"rom_exts":        m.RomExts,
		"rom_path":        m.RomPath,
		"thumb_path":      m.ThumbPath,
		"snap_path":       m.SnapPath,
		"poster_path":     m.PosterPath,
		"packing_path":    m.PackingPath,
		"title_path":      m.TitlePath,
		"background_path": m.BackgroundPath,
		"strategy_path":   m.StrategyPath,
		"video_path":      m.VideoPath,
		"doc_path":        m.DocPath,
		"romlist":         m.Romlist,
		"rombase":         m.Rombase,
		"pinyin":          m.Pinyin,
	}

	result := getDb().Table(m.TableName()).Where("id=?", m.Id).Updates(create)
	if result.Error != nil {
		fmt.Println(result.Error.Error())
	}
	return result.Error
}

//更新平台的一个字段
func (m *Platform) UpdateFieldById(field string, value interface{}) error {

	switch field {
	case "id", "sort":
		value = utils.ToInt(value)
	default:
		value = utils.ToString(value)
	}

	result := getDb().Table(m.TableName()).Where("id=?", m.Id).Update(field, value)
	if result.Error != nil {
		fmt.Println(result.Error.Error())
	}
	return result.Error
}

//删除一个平台
func (m *Platform) DeleteById() error {
	result := getDb().Delete(&m)
	if result.Error != nil {
		fmt.Println(result.Error.Error())
	}
	return result.Error
}

//更新排序
func (m *Platform) UpdateSortById() error {
	result := getDb().Table(m.TableName()).Where("id=?", m.Id).Update("sort", m.Sort)
	if result.Error != nil {
		fmt.Println(result.Error.Error())
	}
	return result.Error
}
