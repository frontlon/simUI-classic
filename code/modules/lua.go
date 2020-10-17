package modules

import "github.com/Shopify/go-lua"

//调用Lua代码
func callLua(f string, cmd string) {
	go func() {
		var luaState *lua.State
		luaState = lua.NewState()
		lua.OpenLibraries(luaState)
		if err := lua.DoFile(luaState, f); err != nil {
			return;
		}

		// 调用lua函数
		luaState.Global("main")
		// 传递参数给lua函数
		luaState.PushString(cmd)
		luaState.ProtectedCall(1, 0, 0)
	}()
}
