/*
package main
var res = []byte{
*/
package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)








func TestInsert(T *testing.T) {
	db, err := sql.Open("sqlite3", "./db.db")
	fmt.Println(err)
	stmt, err := db.Prepare("INSERT INTO user (name, age) values(?,?)")
	res, err := stmt.Exec("zhangsan", "13")
	fmt.Println(res)
}



