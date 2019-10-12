package db

import (
	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	Lang          string
	Theme         string
	Platform      int64
	RomlistStyle  int64
	RomlistZoom   int64
	SearchEngines string
	WindowWidth   int64
	WindowHeight  int64
	WindowLeft    int64
	WindowTop     int64
}

//根据id查询一条数据
func (*Config) Get() (*Config, error) {
	vo := &Config{}
	sql := "SELECT lang, theme, platform, romlist_style, romlist_zoom, search_engines,window_width, window_height, window_left, window_top FROM config where id= 1"
	rows := sqlite.QueryRow(sql)
	err := rows.Scan(&vo.Lang, &vo.Theme, &vo.Platform, &vo.RomlistStyle, &vo.RomlistZoom, &vo.SearchEngines,&vo.WindowWidth, &vo.WindowHeight, &vo.WindowLeft, &vo.WindowTop)
	return vo, err
}

//更新一个字段
func (*Config) UpdateField(field string, value string) error {
	sql := `UPDATE config SET ` + field + ` = "` + value + `" WHERE id = 1`
	stmt, err := sqlite.Prepare(sql)
	if err != nil {
		return err
	}
	_, err2 := stmt.Exec()
	if err2 != nil {
		return err2
	} else {
		return nil
	}
}
