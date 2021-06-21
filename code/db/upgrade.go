package db

import (
	"bufio"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"os"
	"simUI/code/utils"
)

//升级数据库
func UpgradeDB() {

	db, _ := sql.Open("sqlite3", "./data.dll")

	filename := "upgrade.sql"
	f, err := os.Open(filename)
	if err != nil {
		return
	}

	defer os.Remove(filename)
	defer f.Close()
	defer db.Close()

	br := bufio.NewReader(f)
	for {
		a, _, c := br.ReadLine()

		if c == io.EOF {
			break
		}

		if len(a) == 0 {
			continue
		}
		if _,err := db.Exec(string(a));err != nil{
			utils.WriteLog(err.Error())
		}
	}
}
