package db

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	Lang                  string // 当前语言
	Theme                 string // 当前主题
	Platform              string // 当前平台
	Menu                  string // 当前菜单
	Thumb                 string // 当前缩略图显示哪个模块
	RomlistStyle          string // 当前列表样式
	RomlistSize           string // 当前缩放等级
	RomlistFontBackground string // 是否显示字体背景
	RomlistMargin         string // 是否显示模块间距
	RomlistDirection      uint8  // 模块显示方向（自动、横向、竖向）
	RomlistColumn         string // 列表列
	RomlistFontSize       string // 字体大小
	RomlistNameType       string // 显示名称类型(0:别名;1:文件名)
	SearchEngines         string // 搜索引擎地址
	RootPath              string // 当前根目录
	WindowWidth           string // 当前窗口宽度
	WindowHeight          string // 当前窗口高度
	WindowState           string // 当前窗口显示状态
	UpgradeId             string // 当前版本id
	EnableUpgrade         string // 启用更新
	PanelPlatform         string // 当前是否显示平台面板
	PanelMenu             string // 当前是否显示菜单面板
	PanelSidebar          string // 当前是否显示侧边栏
	PanelPlatformWidth    string // 平台面板宽度
	PanelMenuWidth        string // 菜单面板宽度
	PanelSidebarWidth     string // 侧边栏宽度
	SoftName              string // 软件名称
	BackgroundImage       string // 背景图片
	BackgroundRepeat      string // 背景循环方式
	BackgroundOpacity     string // 背景透明度
	BackgroundFuzzy       string // 背景模糊
	BackgroundMask        string // 背景遮罩图
	WallpaperImage        string // 侧边栏图图片
	Cursor                string // 鼠标指针
	VideoVolume           string // 视频默认音量状态
	MusicPlayer           string //音乐播放器路径
	ThumbOrders           string //图集排序列
	RomlistOrders         string //rom列表排序方式
	SqlUpdateNum          uint32 //sql升级进度id
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

//根据id查询一条数据
func (*Config) GetField(field string) (*Config, error) {

	vo := &Config{}
	result := getDb().Select(field).Where("id=1").First(&vo)

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
