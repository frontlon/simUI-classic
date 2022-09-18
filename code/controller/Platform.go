package controller

import (
	"encoding/json"
	"fmt"
	"simUI/code/config"
	"simUI/code/db"
	"simUI/code/modules"
	"simUI/code/utils"
	"simUI/code/utils/go-sciter"
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

		//中文字符转换
		d["exts"] = strings.ReplaceAll(d["exts"], "，", ",")

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
		d["title"] = strings.TrimRight(d["title"], `\`)
		d["title"] = strings.TrimRight(d["title"], `/`)
		d["background"] = strings.TrimRight(d["background"], `\`)
		d["background"] = strings.TrimRight(d["background"], `/`)
		d["wallpaper"] = strings.TrimRight(d["wallpaper"], `\`)
		d["wallpaper"] = strings.TrimRight(d["wallpaper"], `/`)
		d["cassette"] = strings.TrimRight(d["cassette"], `\`)
		d["cassette"] = strings.TrimRight(d["cassette"], `/`)
		d["icon"] = strings.TrimRight(d["icon"], `\`)
		d["icon"] = strings.TrimRight(d["icon"], `/`)
		d["gif"] = strings.TrimRight(d["gif"], `\`)
		d["gif"] = strings.TrimRight(d["gif"], `/`)
		d["optimized"] = strings.TrimRight(d["optimized"], `\`)
		d["optimized"] = strings.TrimRight(d["optimized"], `/`)
		d["video"] = strings.TrimRight(d["video"], `\`)
		d["video"] = strings.TrimRight(d["video"], `/`)
		d["strategy"] = strings.TrimRight(d["strategy"], `\`)
		d["strategy"] = strings.TrimRight(d["strategy"], `/`)
		d["doc"] = strings.TrimRight(d["doc"], `\`)
		d["doc"] = strings.TrimRight(d["doc"], `/`)
		d["files"] = strings.TrimRight(d["files"], `\`)
		d["files"] = strings.TrimRight(d["files"], `/`)
		d["audio"] = strings.TrimRight(d["audio"], `\`)
		d["audio"] = strings.TrimRight(d["audio"], `/`)
		platform := &db.Platform{
			Id:             id,
			Name:           strings.Trim(d["name"], " "),
			Icon:           d["ico"],
			Tag:            strings.Trim(d["tag"], " "),
			RomExts:        strings.ToLower(d["exts"]),
			RomPath:        d["rom"],
			ThumbPath:      d["thumb"],
			SnapPath:       d["snap"],
			PosterPath:     d["poster"],
			PackingPath:    d["packing"],
			TitlePath:      d["title"],
			BackgroundPath: d["background"],
			WallpaperPath:  d["wallpaper"],
			CassettePath:   d["cassette"],
			IconPath:       d["icon"],
			GifPath:        d["gif"],
			OptimizedPath:  d["optimized"],
			StrategyPath:   d["strategy"],
			VideoPath:      d["video"],
			DocPath:        d["doc"],
			FilesPath:      d["files"],
			AudioPath:      d["audio"],
			Rombase:        d["rombase"],
			Pinyin:         utils.TextToPinyin(d["name"]),
		}

		fmt.Println("d[\"audio\"]", d["audio"])
		err := platform.UpdateById()
		if err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}

		//更新缓存
		_ = config.InitConf()

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

	//快速创建平台
	utils.Window.DefineFunction("CreatePlatform", func(args ...*sciter.Value) *sciter.Value {
		id := uint32(utils.ToInt(args[0].String()))
		err := modules.CreatePlatform(id)
		if err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}

		return sciter.NewValue(config.ENV)
	})

	//更新平台介绍
	utils.Window.DefineFunction("UpdatePlatformDesc", func(args ...*sciter.Value) *sciter.Value {
		platform := uint32(utils.ToInt(args[0].String()))
		desc := args[1].String()
		err := modules.UpdatePlatformDesc(platform, desc)
		if err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}
		return sciter.NullValue()
	})

	//更新平台缩略图
	utils.Window.DefineFunction("UpdatePlatformThumb", func(args ...*sciter.Value) *sciter.Value {
		platform := uint32(utils.ToInt(args[0].String()))
		thumb := args[1].String()
		err := modules.UpdatePlatformThumb(platform, thumb)
		if err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}
		return sciter.NullValue()
	})

	//清空平台缩略图
	utils.Window.DefineFunction("ClearPlatformThumb", func(args ...*sciter.Value) *sciter.Value {
		err := modules.ClearPlatformThumb()
		if err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}
		return sciter.NullValue()
	})

	//更新平台缩略图方向
	utils.Window.DefineFunction("UpdatePlatformThumbDirection", func(args ...*sciter.Value) *sciter.Value {
		platform := uint32(utils.ToInt(args[0].String()))
		dir := args[1].String()
		err := modules.UpdatePlatformThumbDirection(platform, dir)
		if err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}
		return sciter.NullValue()
	})

	//清空平台缩略图方向
	utils.Window.DefineFunction("ClearPlatformThumbDirection", func(args ...*sciter.Value) *sciter.Value {
		err := modules.ClearPlatformThumbDirection()
		if err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}
		return sciter.NullValue()
	})

}
