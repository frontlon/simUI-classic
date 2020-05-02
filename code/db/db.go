package db

import (
	"VirtualNesGUI/code/utils"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var engine *gorm.DB
var LogMode bool = true
//连接数据库
func Conn() error {
	//连接数据库
	err := errors.New("")

	if !utils.FileExists("data.dll"){
		fmt.Println("数据库文件data.dll不存在")
		return errors.New("数据库文件data.dll不存在\n Database does not exist")
	}

	engine, err = gorm.Open("sqlite3", "data.dll")
	if err != nil {
		panic("连接数据库失败")
	}
	//调试模式下 打印日志
	engine.LogMode(LogMode)

	//禁用同步模式
	engine.Exec("PRAGMA synchronous = OFF;")
	return nil
}

func getDb() *gorm.DB {
	return engine
}

//关闭数据库
func CloseDb(){
	engine.Close()
}
