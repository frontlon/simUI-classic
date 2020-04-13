package db

import (
	"VirtualNesGUI/code/utils"
	"fmt"
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
	WindowState   uint8
	RenameType    uint8
	FontSize      uint8
}

func (*Config) TableName() string {
	return "config"
}

//根据id查询一条数据
func (*Config) Get() (*Config, error) {

	vo := &Config{}
	result := getDb().Where("id=1").First(&vo)

	if result.Error != nil {
		fmt.Println(result.Error.Error())
	}

	return vo, result.Error
}

//更新一个字段
func (m *Config) UpdateField(field string, value interface{}) error {

	switch field {
	case "platform", "romlist_style", "romlist_zoom", "window_width", "window_height", "window_state","font_size","RenameType":
		value = utils.ToInt(value)
	default:
		value = utils.ToString(value)
	}
	result := getDb().Table(m.TableName()).Where("id=1").Update(field, value)

	if result.Error != nil {
		fmt.Println(result.Error.Error())
	}

	return result.Error
}
