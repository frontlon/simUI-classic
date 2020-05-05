package controller

import (
	"simUI/code/db"
	"simUI/code/utils"
	"encoding/json"
	"github.com/sciter-sdk/go-sciter"
	"strings"
)

/**
 * 定义view用function
 **/

func PlatformController() {

	//读取平台详情
	utils.Window.DefineFunction("GetPlatformById", func(args ...*sciter.Value) *sciter.Value {
		id := uint32(utils.ToInt(args[0].String()))

		//游戏游戏详细数据
		info, err := (&db.Platform{}).GetById(id)
		if err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}
		jsonInfo, _ := json.Marshal(&info)

		return sciter.NewValue(string(jsonInfo))
	})

	//读取平台列表
	utils.Window.DefineFunction("GetPlatform", func(args ...*sciter.Value) *sciter.Value {
		//游戏游戏详细数据
		info, err := (&db.Platform{}).GetAll()
		if err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}
		jsonInfo, _ := json.Marshal(&info)
		return sciter.NewValue(string(jsonInfo))
	})

	//添加一个平台
	utils.Window.DefineFunction("AddPlatform", func(args ...*sciter.Value) *sciter.Value {
		name := args[0].String()
		platform := &db.Platform{
			Name:   name,
			Pinyin: utils.TextToPinyin(name),
		}
		id, err := platform.Add()
		if err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}
		return sciter.NewValue(utils.ToString(id))
	})

	//删除一个平台
	utils.Window.DefineFunction("DelPlatform", func(args ...*sciter.Value) *sciter.Value {
		id := uint32(utils.ToInt(args[0].String()))

		go func(id uint32) *sciter.Value {

			platform := &db.Platform{
				Id: id,
			}
			sim := &db.Simulator{
				Platform: id,
			}
			rom := &db.Rom{
				Platform: id,
			}

			//删除rom数据
			if err := rom.DeleteByPlatform(); err != nil {
				utils.WriteLog(err.Error())
				return utils.ErrorMsg(err.Error())
			}
			//删除模拟器
			if err := sim.DeleteByPlatform(); err != nil {
				utils.WriteLog(err.Error())
				return utils.ErrorMsg(err.Error())
			}
			//删除平台
			if err := platform.DeleteById(); err != nil {
				utils.WriteLog(err.Error())
				return utils.ErrorMsg(err.Error())
			}
			return sciter.NullValue()
		}(id)
		return sciter.NewValue("1")
	})

	//更新平台信息
	utils.Window.DefineFunction("UpdatePlatform", func(args ...*sciter.Value) *sciter.Value {
		data := args[0].String()
		d := make(map[string]string)
		json.Unmarshal([]byte(data), &d)
		id := uint32(utils.ToInt(d["id"]))

		//取掉路径结尾路径分隔符
		d["rom"] = strings.TrimRight(d["rom"], `\`)
		d["rom"] = strings.TrimRight(d["rom"], `/`)
		d["thumb"] = strings.TrimRight(d["thumb"], `\`)
		d["thumb"] = strings.TrimRight(d["thumb"], `/`)
		d["snap"] = strings.TrimRight(d["snap"], `\`)
		d["snap"] = strings.TrimRight(d["snap"], `/`)
		d["poster"] = strings.TrimRight(d["poster"], `\`)
		d["poster"] = strings.TrimRight(d["poster"], `/`)
		d["packing"] = strings.TrimRight(d["packing"], `\`)
		d["packing"] = strings.TrimRight(d["packing"], `/`)
		d["video"] = strings.TrimRight(d["video"], `\`)
		d["video"] = strings.TrimRight(d["video"], `/`)
		d["strategy"] = strings.TrimRight(d["strategy"], `\`)
		d["strategy"] = strings.TrimRight(d["strategy"], `/`)
		d["doc"] = strings.TrimRight(d["doc"], `\`)
		d["doc"] = strings.TrimRight(d["doc"], `/`)

		platform := &db.Platform{
			Id:           id,
			Name:         d["name"],
			Icon:         d["icon"],
			RomExts:      d["exts"],
			RomPath:      d["rom"],
			ThumbPath:    d["thumb"],
			SnapPath:     d["snap"],
			PosterPath:   d["poster"],
			PackingPath:  d["packing"],
			StrategyPath: d["strategy"],
			VideoPath:    d["video"],
			DocPath:      d["doc"],
			Romlist:      d["romlist"],
			Pinyin:       utils.TextToPinyin(d["name"]),
		}

		err := platform.UpdateById()
		if err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}
		return sciter.NewValue("1")
	})

	//更新平台排序
	utils.Window.DefineFunction("UpdatePlatformSort", func(args ...*sciter.Value) *sciter.Value {
		data := args[0].String()
		d := make(map[uint32]uint32)
		json.Unmarshal([]byte(data), &d)

		if len(d) == 0 {
			return sciter.NullValue()
		}

		for id, val := range d {
			platform := &db.Platform{
				Id:   id,
				Sort: val,
			}
			err := platform.UpdateSortById()
			if err != nil {
				utils.WriteLog(err.Error())
				return utils.ErrorMsg(err.Error())
			}
		}
		return sciter.NewValue("1")
	})
}
