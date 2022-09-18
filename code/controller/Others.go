package controller

import (
	"encoding/json"
	"os"
	"simUI/code/config"
	"simUI/code/db"
	"simUI/code/utils"
	"simUI/code/utils/go-sciter"
)

/**
 * 定义view用function
 **/

func OthersController() {

	//读取关于信息
	utils.Window.DefineFunction("GetConstInfo", func(args ...*sciter.Value) *sciter.Value {

		result := map[string]string{}
		fileDate := utils.GetFileUpdateDate(os.Args[0])

		//打包时间
		result["buildTime"] = fileDate.Format("2006-01-02 15:04")
		//rom列表每页加载数量
		result["romListPageSize"] = utils.ToString(db.ROM_PAGE_NUM)
		//读取路径分隔符
		result["separator"] = config.Cfg.Separator
		//读取视图路径
		result["viewPath"] = config.Cfg.ViewPath

		getjson, _ := json.Marshal(result)
		return sciter.NewValue(string(getjson))
	})

}
