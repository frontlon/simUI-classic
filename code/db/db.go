package db

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"io/ioutil"
	"os"
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

	//禁用同步模式
	engine.Exec("PRAGMA synchronous = OFF;")
	return nil
}

func getDb() *gorm.DB {
	return engine
}

//收缩数据库
func Vacuum() {
	engine.Exec("VACUUM;")
}

//升级数据库
func UpgradeDB() {

	filename := "upgrade.dat"

	if !utils.FileExists(filename){
		return
	}


	/*f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()
	defer os.Remove(filename)

	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行
		if err != nil || io.EOF == err {
			break
		}
		//alter table mydownload add column 'IsFree' varchar(100) default '1'
		getDb().Exec(line)
	}*/

	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("read file fail", err)
		return
	}
	defer os.Remove(filename)
	defer f.Close()

	sql, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("read to fd fail", err)
		return
	}
	getDb().Exec(string(sql))

	return


}
