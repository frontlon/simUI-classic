package modules

import (
	"encoding/json"
	"fmt"
	"simUI/code/config"
	"simUI/code/db"
	"simUI/code/utils"
	"simUI/code/utils/go-sciter"
	"time"
)

// 更新url
var upgradeUrl string = "http://upgrade.simui.net/check.html"

type Version struct {
	Id      uint64 //版本id
	Version string //版本号
	Date    string //日期
}

// 启动时自动检测更新
func BootCheckUpgrade() {
	//检测是否启动更新
	if config.Cfg.Default.EnableUpgrade == "0" {
		return
	}
	go func() {
		body, err := utils.HttpGet(upgradeUrl, nil)
		if err != nil {
			return
		}
		ver := &Version{}
		if err := json.Unmarshal(body, &ver); err != nil {
			return
		}

		//如果是启动检测，则验证是否需要显示
		if ver.Id > uint64(utils.ToInt(config.Cfg.Default.UpgradeId)) {
			if _, err := utils.Window.Call("upgrade", sciter.NewValue(string(body))); err != nil {
			}
		}
	}()
}

// 检查更新
func CheckUpgrade() string {
	body, err := utils.HttpGet(upgradeUrl, nil)
	if err != nil {
		return "error"
	}

	ver := &Version{}
	if err := json.Unmarshal(body, &ver); err != nil {
	}

	if ver.Id > uint64(utils.ToInt(config.Cfg.Default.UpgradeId)) {
		return string(body)
	}

	return ""
}

// 升级数据库
func UpgradeDB() {

	//当去当前sql更新num
	oldNum := 0
	cfg, err := (&db.Config{}).GetField("sql_update_num")
	if err == nil {
		oldNum = utils.ToInt(cfg.SqlUpdateNum)
	}

	//读取升级sql列表
	sqls := db.DbUpdateSqlList(config.Cfg.UpgradeId)
	currNum := len(sqls)

	//无需升级
	if oldNum == currNum {
		return
	}

	//备份数据库
	dbName, _ := db.GetDbFileName()
	newName := utils.GetFileName(dbName)
	newExt := utils.GetFileExt(dbName)
	t := time.Now().Format("20060102_150405")
	newFile := config.Cfg.CachePath + newName + "_bak_" + t + newExt
	utils.FileCopy(dbName, newFile)

	//升级sql
	newSqls := sqls[oldNum:]
	for _, sql := range newSqls {
		if sql == "" {
			continue
		}
		fmt.Println("升级SQL：", sql)
		if err = db.Exec(sql); err != nil {
			fmt.Println(err)
		}
	}

	//更新sql_num
	(&db.Config{}).UpdateField("sql_update_num", currNum)

	//系统alert提示
	utils.ShowAlertAndExit(config.Cfg.Lang["DbUpgradeTitle"], config.Cfg.Lang["DbUpgradeContent"])
}
