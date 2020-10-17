package db

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	Lang               string // 当前语言
	Theme              string // 当前主题
	Platform           string // 当前平台
	Menu               string // 当前菜单
	Thumb              string // 当前缩略图显示哪个模块
	RomlistStyle       string // 当前列表样式
	RomlistZoom        string // 当前缩放等级
	SearchEngines      string // 搜索引擎地址
	RootPath           string // 当前根目录
	WindowWidth        string // 当前窗口宽度
	WindowHeight       string // 当前窗口高度
	WindowState        string // 当前窗口显示状态
	FontSize           string // 当前rom列表字体大小
	UpgradeId          string // 当前版本id
	EnableUpgrade      string // 启用更新
	PanelPlatform      string // 当前是否显示平台面板
	PanelMenu          string // 当前是否显示菜单面板
	PanelSidebar       string // 当前是否显示侧边栏
	PanelPlatformWidth string // 平台面板宽度
	PanelMenuWidth     string // 菜单面板宽度
	PanelSidebarWidth  string // 侧边栏宽度
	SoftName           string // 软件名称
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
	result := getDb().Table(m.TableName()).Where("id=1").Update(field, value)

	if result.Error != nil {
		fmt.Println(result.Error.Error())
	}

	return result.Error
}

func (m *Config) Add() {
	getDb().Create(&m)
}
