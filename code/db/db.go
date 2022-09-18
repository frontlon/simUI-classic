package db

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"os"
	"path/filepath"
	"simUI/code/utils"
)

var engine *gorm.DB
var LogMode bool = true
var maxVar = 990 //sqlite最大参数个数

//连接数据库
func Conn() error {
	//连接数据库
	dbFile, err := GetDbFileName()
	if err != nil {
		panic(err.Error())
	}

	engine, err = gorm.Open("sqlite3", dbFile)
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

//读取数据库文件名称
func GetDbFileName() (string, error) {

	dbPath, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	//先读文件同名数据库文件
	dbFile := dbPath + "/" + utils.GetFileName(os.Args[0]) + ".dll"

	//再读默认数据库文件
	if !utils.FileExists(dbFile) {
		dbFile = dbPath + "/data.dll"
	}

	if !utils.FileExists(dbFile) {
		fmt.Println("数据库文件data.dll不存在")
		return "", errors.New("数据库文件data.dll不存在\n Database does not exist")
	}

	return dbFile, nil

}
