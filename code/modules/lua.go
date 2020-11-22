package modules

import (
	"github.com/Shopify/go-lua"
	"simUI/code/config"
	"simUI/code/utils"
)

//调用Lua代码
func callLua(luaFile string, simPath string,romPath string) {
	go func() {
		var luaState *lua.State
		luaState = lua.NewState()
		lua.OpenLibraries(luaState)

		if !utils.IsAbsPath(luaFile){
			luaFile = config.Cfg.RootPath + luaFile
		}

		if err := lua.DoFile(luaState, luaFile); err != nil {
			utils.WriteLog("Lua Run Error:" + err.Error())
			return
		}

		// 调用lua函数
		luaState.Global("main")

		// 传递参数给lua函数
		luaState.PushString(config.Cfg.RootPath) //simui根目录
		luaState.PushString(simPath) //模拟器文件
		luaState.PushString(romPath) //rom文件
		if err := luaState.ProtectedCall(3, 0, 0);err != nil{
			utils.WriteLog("Lua Error:" + err.Error())
			return
		}
	}()
}
