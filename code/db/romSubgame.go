package db

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"math"
)

type RomSubGame struct {
	Id       uint64
	Platform uint32 // 平台
	FileMd5  string // file_md5
	Pname    string // 父file_md5
}

func (*RomSubGame) TableName() string {
	return "rom_subgame"
}

// 根据平台id查询数据
func (*RomSubGame) GetByPlatform(platform uint32) ([]*RomSubGame, error) {

	volist := []*RomSubGame{}

	result := getDb().Where("platform=?", platform).Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return volist, result.Error
}

// 根据平台id查询数据，返回map
func (*RomSubGame) GetByPlatformToMap(platform uint32) (map[string]string, error) {
	vo := []*RomSubGame{}
	result := getDb().Where("platform=?", platform).Find(&vo)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	//转换成map
	volist := map[string]string{}
	for _, v := range vo {
		volist[v.FileMd5] = v.Pname
	}
	return volist, result.Error
}

// 根据平台id查询file_md5
func (*RomSubGame) GetFileMd5ByPlatform(platform uint32) ([]string, error) {

	volist := []*RomSubGame{}

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

func (m *RomSubGame) BatchAdd(romlist []*RomSubGame) {

	if len(romlist) == 0 {
		return
	}
	tx := getDb().Begin()
	for _, v := range romlist {
		tx.Create(&v)
	}
	tx.Commit()
}

// 用新的file_md5替换旧的pname和file_md5
func (m *RomSubGame) BatchUpdateFileMd5(platform uint32, lists []map[string]string) error {

	if len(lists) == 0 {
		return nil
	}

	tx := getDb().Begin()
	for _, rom := range lists {
		tx.Table(m.TableName()).Where("platform = ? AND file_md5 = ?", platform, rom["oldMd5"]).Update(map[string]interface{}{"file_md5": rom["newMd5"]})
		tx.Table(m.TableName()).Where("platform = ? AND pname = ?", platform, rom["oldMd5"]).Update(map[string]interface{}{"pname": rom["newMd5"]})
	}
	if err := tx.Commit().Error; err != nil {
		fmt.Println("update错误", err)
	}
	return nil
}

// 绑定子游戏
func (m *RomSubGame) UpdatePname(platform uint32, SlaveFileMd5 string, pname string) error {

	if platform == 0 || SlaveFileMd5 == "" || pname == "" {
		return nil
	}

	//初始化数据
	_ = m.InitData(platform, SlaveFileMd5)

	//更新数据
	result := getDb().Table(m.TableName()).Where("platform=? AND file_md5 = ?", platform, SlaveFileMd5).Update("pname", pname)
	if result.Error != nil {
		fmt.Println(result.Error)
		return result.Error
	}

	return nil
}

// 解绑子游戏
func (m *RomSubGame) DeleteByFileMd5(platform uint32, fileMd5 string) error {

	if platform == 0 || fileMd5 == "" {
		return nil
	}

	result := getDb().Where("platform = ? AND file_md5 = ?", platform, fileMd5).Delete(&m)
	if result.Error != nil {
		fmt.Println(result.Error)
		return result.Error
	}
	return result.Error
}

// 删除记录
func (m *RomSubGame) DeleteByFileMd5s(platform uint32, md5List []string) error {

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

// 初始化数据，如果没有数据，则生成一条
func (m *RomSubGame) InitData(platform uint32, fileMd5 string) error {
	count := 0
	getDb().Table(m.TableName()).Where("platform=? AND file_md5 = ?", platform, fileMd5).Count(&count)

	if count == 0 {
		create := &RomSubGame{
			Platform: platform,
			FileMd5:  fileMd5,
			Pname:    "",
		}
		result := getDb().Create(&create)
		if result.Error != nil {
			fmt.Println(result.Error)
			return result.Error
		}
	}
	return nil
}

// 删除不存在的平台下的所有数据
func (*RomSubGame) ClearByNotPlatform(platforms []string) error {

	m := &RomSubGame{}
	result := getDb().Not("platform", platforms).Delete(&m)

	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return result.Error
}
