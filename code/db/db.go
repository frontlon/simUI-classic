package db

import (
	"database/sql"
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var sqlite = &sql.DB{}
var engine *gorm.DB

//连接数据库
func Conn() {
	//连接数据库
	err := errors.New("")
	engine, err = gorm.Open("sqlite3", "data.dll")
	if err != nil {
		panic("连接数据库失败")
	}
	//调试模式下 打印日志
	engine.LogMode(false)

	//禁用同步模式
	engine.Exec("PRAGMA synchronous = OFF;")
}

func getDb() *gorm.DB {
	return engine
}
