package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type CustomDB struct {
	Engine *gorm.DB
}

type RomCount struct {
	Platform uint32
	Count    int64
}

func CustomDBConn(p string) (*CustomDB, error) {
	//连接数据库
	conn, err := gorm.Open("sqlite3", p)
	if err != nil {
		return nil, err
	}
	//调试模式下 打印日志
	conn.LogMode(LogMode)

	//配置参数
	conn.Exec("PRAGMA synchronous = OFF;")
	conn.Exec("PRAGMA journal_mode = OFF;")
	conn.Exec("PRAGMA auto_vacuum = 0;")
	conn.Exec("PRAGMA cache_size = 8000;")
	conn.Exec("PRAGMA temp_store = 2;")

	create := &CustomDB{Engine: conn}
	return create, nil
}

func (e *CustomDB) Close() error {
	err := e.Engine.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (e *CustomDB) GetAllPlatform() ([]*Platform, error) {
	volist := []*Platform{}
	result := e.Engine.Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return volist, nil
}

// 根据ID查询一个平台参数
func (e *CustomDB) GetPlatformByIds(ids []uint32) ([]*Platform, error) {
	volist := []*Platform{}
	result := e.Engine.Where("id in (?)", ids).Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return volist, result.Error
}

// 统计rom数据
func (e *CustomDB) GetRomCount() (map[uint32]int64, error) {
	volist := []*RomCount{}

	result := e.Engine.Table((&Rom{}).TableName()).Select("count(*) as count,platform").Group("platform").Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
		return nil, result.Error
	}
	data := map[uint32]int64{}
	if len(volist) > 0 {
		for _, v := range volist {
			data[v.Platform] = v.Count
		}
	}
	return data, nil
}

// 读取rom 配置
func (e *CustomDB) GetRomSettingByPlatform(platform uint32) ([]*RomSetting, error) {
	volist := []*RomSetting{}
	result := e.Engine.Where("platform=?", platform).Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return volist, result.Error
}

// 根据平台id查询数据
func (e *CustomDB) GetSubGameByPlatform(platform uint32) ([]*RomSubGame, error) {
	volist := []*RomSubGame{}
	result := e.Engine.Where("platform=?", platform).Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return volist, result.Error
}

// 读取模拟器
func (e *CustomDB) GetAllSimulator() (map[uint32][]*Simulator, error) {
	volist := []*Simulator{}
	result := e.Engine.Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
		return nil, result.Error
	}
	data := map[uint32][]*Simulator{}
	if len(volist) > 0 {
		for _, v := range volist {
			if _, ok := data[v.Platform]; ok {
				data[v.Platform] = append(data[v.Platform], v)
			} else {
				data[v.Platform] = []*Simulator{v}
			}
		}
	}
	return data, nil
}
