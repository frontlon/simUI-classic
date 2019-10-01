package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var sqlite, _ = sql.Open("sqlite3", `./db.db`)

//清空rom数据
func DbClear() {


	//关闭同步
	sqlite.Exec("PRAGMA synchronous = OFF;")
	//清空数据
	sqlite.Exec(`DELETE FROM sqlite_sequence WHERE name = "rom"`)
	sqlite.Exec(`DELETE FROM sqlite_sequence WHERE name = "menu"`)
	sqlite.Exec(`DELETE FROM rom`)
	sqlite.Exec(`DELETE FROM menu`)
}