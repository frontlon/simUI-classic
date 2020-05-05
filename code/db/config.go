package db

import (
	"simUI/code/utils"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	Lang          string // 当前语言
	Theme         string // 当前主题
	Platform      uint32 // 当前平台
	Menu          string // 当前菜单
	Thumb         string // 当前缩略图显示哪个模块
	RomlistStyle  uint8  // 当前列表样式
	RomlistZoom   uint8  // 当前缩放等级
	SearchEngines string // 搜索引擎地址
	RootPath      string // 当前根目录
	WindowWidth   uint16 // 当前窗口宽度
	WindowHeight  uint16 // 当前窗口高度
	WindowState   uint8  // 当前窗口显示状态
	RenameType    uint8  // 当前rom重命名类型
	FontSize      uint8  // 当前rom列表字体大小
	UpgradeId     uint64 // 当前版本id
	EnableUpgrade uint8  // 启用更新
	PanelPlatform uint8  // 当前是否显示平台面板
	PanelMenu     uint8  // 当前是否显示菜单面板
	PanelSidebar  uint8  // 当前是否显示侧边栏
	SoftName      string // 软件名称
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
	case "platform",
		"romlist_style",
		"romlist_zoom",
		"window_width",
		"window_height",
		"window_state",
		"font_size",
		"rename_type",
		"upgrade_id",
		"enable_upgrade",
		"panel_platform",
		"panel_menu",
		"panel_sidebar":
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

func (m *Config) Add() {
	getDb().Create(&m)
}

//清空表数据
func (m *Config) Truncate() error {
	result := getDb().Delete(&m)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}
