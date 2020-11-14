package db

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"simUI/code/utils"
)

var engine *gorm.DB
var LogMode bool = true
//连接数据库
func Conn() error {
	//连接数据库
	err := errors.New("")

	if !utils.FileExists("data.dll") {
		fmt.Println("数据库文件data.dll不存在")
		return errors.New("数据库文件data.dll不存在\n Database does not exist")
	}

	engine, err = gorm.Open("sqlite3", "data.dll")
	if err != nil {
		panic("连接数据库失败")
	}
	//调试模式下 打印日志
	engine.LogMode(LogMode)

	//配置参数
	engine.Exec("PRAGMA synchronous = OFF;")
	engine.Exec("PRAGMA journal_mode = OFF;")
	engine.Exec("PRAGMA auto_vacuum = 0;")
	engine.Exec("PRAGMA cache_size = 8000;")
	engine.Exec("PRAGMA temp_store = 2;")
	return nil
}

func getDb() *gorm.DB {
	return engine
}

//收缩数据库
func Vacuum() {
	engine.Exec("VACUUM;")
}

