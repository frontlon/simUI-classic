package db

import (
	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	Lang          string
	Theme         string
	Platform      uint32
	Menu          string
	Thumb         string
	RomlistStyle  uint8
	RomlistZoom   uint8
	SearchEngines string
	RootPath      string
	WindowWidth   uint16
	WindowHeight  uint16
	WindowState    uint8
}

//根据id查询一条数据
func (*Config) Get() (*Config, error) {
	vo := &Config{}
	sql := "SELECT lang, theme, platform, menu, thumb, romlist_style, romlist_zoom, search_engines, root_path, window_width, window_height, window_state FROM config where id= 1"
	rows := sqlite.QueryRow(sql)
	err := rows.Scan(&vo.Lang, &vo.Theme, &vo.Platform, &vo.Menu, &vo.Thumb, &vo.RomlistStyle, &vo.RomlistZoom, &vo.SearchEngines, &vo.RootPath, &vo.WindowWidth, &vo.WindowHeight, &vo.WindowState)
	return vo, err
}

//更新一个字段
func (*Config) UpdateField(field string, value string) error {
	sql := `UPDATE config SET ` + field + ` = '` + value + `' WHERE id = 1`
	stmt, err := sqlite.Prepare(sql)
	defer stmt.Close()
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
