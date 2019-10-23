package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var sqlite = &sql.DB{}

//连接数据库
func Conn() {
	//连接数据库
	sqlite, _ = sql.Open("sqlite3", `./db.db`)
	//关闭同步
	sqlite.Exec("PRAGMA synchronous = OFF;")
}