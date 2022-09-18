package controller

import (
	"encoding/json"
	"simUI/code/compoments"
	"simUI/code/modules"
	"simUI/code/utils"
	"simUI/code/utils/go-sciter"
)

/**
 * 定义view用function
 **/

func AudioController() {

	//读取音频文件列表
	utils.Window.DefineFunction("GetAudio", func(args ...*sciter.Value) *sciter.Value {
		id := uint64(utils.ToInt(args[0].String()))
		volist,err := modules.GetAudioList(id)
		if err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}
		jsonInfo, _ := json.Marshal(volist)
		return sciter.NewValue(string(jsonInfo))
	})

	//上传文件
	utils.Window.DefineFunction("UploadAudioFile", func(args ...*sciter.Value) *sciter.Value {
		id := uint64(utils.ToInt(args[0].String()))
		name := args[1].String()
		p := args[2].String()

		relPath, err := modules.UploadAudioFile(id, name, p)
		if err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}

		return sciter.NewValue(relPath)
	})

	//更新配置
	utils.Window.DefineFunction("UpdateAudio", func(args ...*sciter.Value) *sciter.Value {
		id := uint64(utils.ToInt(args[0].String()))
		data := args[1].String()
		if err := modules.UpdateAudio(id, data); err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}
		return sciter.NullValue()
	})

	//播放音频文件
	utils.Window.DefineFunction("PlayAudio", func(args ...*sciter.Value) *sciter.Value {
		urls := []string{}
		json.Unmarshal([]byte(args[0].String()), &urls)

		for k,v:= range urls{
			urls[k] = utils.AbsPath(v)
		}

		if err := compoments.PlayAudio(urls); err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}

		return sciter.NullValue()
	})

}
