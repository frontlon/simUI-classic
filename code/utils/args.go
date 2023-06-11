package utils

import (
	"flag"
)

var args map[string]string

//解析并读取命令行参数
func GetCmdArgs(name string) string {
	if len(args) == 0 {
		var db string
		flag.StringVar(&db, "db", "", "")
		flag.Parse()
		args = map[string]string{
			"db": db,
		}
	}
	if _, ok := args[name]; ok {
		return args[name]
	} else {
		return ""
	}
}
