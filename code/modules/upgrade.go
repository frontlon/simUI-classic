package modules

import (
	"VirtualNesGUI/code/config"
	"VirtualNesGUI/code/utils"
	"encoding/json"
	"github.com/sciter-sdk/go-sciter"
)

//更新url
var upgradeUrl string = "http://upgrade.simui.net/check.html"

type Version struct {
	Id      uint64 //版本id
	Version string //版本号
	Date    string //日期
}

//启动时自动检测更新
func BootCheckUpgrade() {
	//检测是否启动更新
	if config.Cfg.Default.EnableUpgrade == 0 {
		return
	}
	go func() {
		body,_ := utils.GetHttp(upgradeUrl)
		ver := &Version{}
		if err := json.Unmarshal([]byte(body), &ver); err != nil {
		}

		//如果是启动检测，则验证是否需要显示
		if ver.Id > config.Cfg.Default.UpgradeId {
			if _, err := config.Cfg.Window.Call("upgrade", sciter.NewValue(string(body))); err != nil {
			}
		}
	}()
}

//检查更新
func CheckUpgrade() string {
	body,err := utils.GetHttp(upgradeUrl)
	if err != nil{
		return "error"
	}

	ver := &Version{}
	if err := json.Unmarshal([]byte(body), &ver); err != nil {
	}

	if ver.Id > config.Cfg.Default.UpgradeId {
		return body
	}

	return ""
}
