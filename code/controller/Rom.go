package controller

import (
	"encoding/json"
	"errors"
	"simUI/code/compoments"
	"simUI/code/config"
	"simUI/code/db"
	"simUI/code/modules"
	"simUI/code/request"
	"simUI/code/utils"
	"simUI/code/utils/go-sciter"
	"strings"
)

/**
 * 定义view用function
 **/

func RomController() {

	//运行游戏
	utils.Window.DefineFunction("RunGame", func(args ...*sciter.Value) *sciter.Value {

		romId := uint64(utils.ToInt(args[0].String()))
		simId := uint32(utils.ToInt(args[1].String()))

		err := modules.RunGame(romId, simId)
		if err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}
		return sciter.NullValue()
	})

	//运行攻略文件
	utils.Window.DefineFunction("RunStrategy", func(args ...*sciter.Value) *sciter.Value {
		f := args[0].String()
		if f == "" {
			return sciter.NullValue()
		}
		if err := compoments.RunGame("", []string{f}); err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}
		return sciter.NullValue()

	})

	//在资源管理器中打开目录
	utils.Window.DefineFunction("OpenFolder", func(args ...*sciter.Value) *sciter.Value {

		id := uint64(utils.ToInt(args[0].String()))
		opt := args[1].String()
		otherId := args[2].String()

		err := modules.OpenFolder(id, opt, otherId)
		if err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}

		return sciter.NullValue()
	})

	//在资源管理器中打开rom目录
	utils.Window.DefineFunction("OpenRomPathFolder", func(args ...*sciter.Value) *sciter.Value {

		platform := uint32(utils.ToInt(args[0].String()))
		err := modules.OpenRomPathFolder(platform)
		if err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}
		return sciter.NullValue()
	})

	//添加游戏bat
	utils.Window.DefineFunction("AddIndieGame", func(args ...*sciter.Value) *sciter.Value {

		platform := uint32(utils.ToInt(args[0].String()))
		menu := args[1].String()
		datastr := args[2].String()

		files := []string{}
		_ = json.Unmarshal([]byte(datastr), &files)

		err := errors.New("")
		err = modules.AddIndieGame(platform, menu, files)

		if err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}
		return sciter.NullValue()
	})

	//读取游戏列表
	utils.Window.DefineFunction("GetGameList", func(args ...*sciter.Value) *sciter.Value {
		data := &request.GetGameList{}

		_ = json.Unmarshal([]byte(args[0].String()), &data)

		if data.BaseYear == config.Cfg.Lang["BaseYear"] {
			data.BaseYear = ""
		}

		newlist := []*db.Rom{}
		if data.Num == "" {
			newlist, _ = (&db.Rom{}).Get(data.ShowSubGame, data.ShowHide, data.Page, data.Platform, data.Catname, data.Keyword, data.BaseType, data.BasePublisher, data.BaseYear, data.BaseCountry, data.BaseTranslate, data.BaseVersion, data.BaseProducer, data.Score, data.Complete)
		} else {
			//按拼音查询
			newlist, _ = (&db.Rom{}).GetByPinyin(data.ShowHide, data.Page, data.Platform, data.Catname, data.Num)
		}

		jsonRom, _ := json.Marshal(newlist)
		return sciter.NewValue(string(jsonRom))
	})

	//读取没有子游戏的主游戏列表
	utils.Window.DefineFunction("GetGameListNotSubGame", func(args ...*sciter.Value) *sciter.Value {
		data := &request.GetGameList{}

		_ = json.Unmarshal([]byte(args[0].String()), &data)
		newlist := []*db.Rom{}
		newlist, _ = (&db.Rom{}).GetNotSubRom(data.Page, data.ShowHide, data.Platform, data.Catname, data.Keyword)

		jsonRom, _ := json.Marshal(newlist)
		return sciter.NewValue(string(jsonRom))
	})

	//根据id列表读取rom
	utils.Window.DefineFunction("GetGameListByIds", func(args ...*sciter.Value) *sciter.Value {

		romIdsStr := strings.Split(args[0].String(), ",")
		ids := []uint64{}
		for _, v := range romIdsStr {
			ids = append(ids, uint64(utils.ToInt(v)))
		}
		newlist, _ := (&db.Rom{}).GetByIds(ids)

		jsonRom, _ := json.Marshal(newlist)
		return sciter.NewValue(string(jsonRom))

	})

	//根据父id夺取子游戏列表
	utils.Window.DefineFunction("GetSubGamesByPid", func(args ...*sciter.Value) *sciter.Value {
		pid := uint64(utils.ToInt(args[0].String())) //romId
		rom, _ := (&db.Rom{}).GetById(pid)
		romlist, _ := (&db.Rom{}).GetSubRom(rom.Platform, rom.FileMd5)
		jsonRom, _ := json.Marshal(romlist)
		return sciter.NewValue(string(jsonRom))
	})

	//读取游戏数量
	utils.Window.DefineFunction("GetGameCount", func(args ...*sciter.Value) *sciter.Value {
		data := &request.GetGameList{}
		_ = json.Unmarshal([]byte(args[0].String()), &data)

		if data.BaseYear == config.Cfg.Lang["BaseYear"] {
			data.BaseYear = ""
		}

		count := 0
		if data.Num == "" {
			count, _ = (&db.Rom{}).Count(data.ShowHide, data.Platform, data.Catname, data.Keyword, data.BaseType, data.BasePublisher, data.BaseYear, data.BaseCountry, data.BaseTranslate, data.BaseVersion, data.BaseProducer, data.Score, data.Complete)
		} else {
			//按拼音查询
			count, _ = (&db.Rom{}).CountByPinyin(data.ShowHide, data.Page, data.Platform, data.Catname, data.Num)
		}

		return sciter.NewValue(utils.ToString(count))
	})
	//读取游戏数量 - 不显示包含子游戏的主rom
	utils.Window.DefineFunction("GetGameCountNotSubGame", func(args ...*sciter.Value) *sciter.Value {
		data := &request.GetGameList{}
		_ = json.Unmarshal([]byte(args[0].String()), &data)
		count, _ := (&db.Rom{}).CountNotSubGame(data.ShowHide, data.Platform, data.Catname, data.Keyword)
		return sciter.NewValue(utils.ToString(count))
	})

	//读取rom详情
	utils.Window.DefineFunction("GetGameDetail", func(args ...*sciter.Value) *sciter.Value {
		id := uint64(utils.ToInt(args[0].String()))
		res, err := modules.GetGameDetail(id)
		if err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}
		jsonMenu, _ := json.Marshal(&res)
		return sciter.NewValue(string(jsonMenu))
	})

	//根据id读取rom基本信息
	utils.Window.DefineFunction("GetGameById", func(args ...*sciter.Value) *sciter.Value {
		id := uint64(utils.ToInt(args[0].String()))
		rom, _ := (&db.Rom{}).GetById(id)
		jsonInfo, _ := json.Marshal(&rom)
		return sciter.NewValue(string(jsonInfo))
	})

	//读取游戏攻略内容
	utils.Window.DefineFunction("GetGameDoc", func(args ...*sciter.Value) *sciter.Value {
		t := args[0].String()
		id := uint64(utils.ToInt(args[1].String()))
		toHtml := uint8(utils.ToInt(args[2].String()))
		res, err := modules.GetGameDoc(t, id, toHtml)
		if err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}
		return sciter.NewValue(res)
	})

	//更新游戏攻略内容
	utils.Window.DefineFunction("SetGameDoc", func(args ...*sciter.Value) *sciter.Value {
		t := args[0].String()
		id := uint64(utils.ToInt(args[1].String()))
		content := args[2].String()
		err := modules.SetGameDoc(t, id, content)
		if err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}
		return sciter.NullValue()
	})

	//删除游戏攻略内容
	utils.Window.DefineFunction("DelGameDoc", func(args ...*sciter.Value) *sciter.Value {
		t := args[0].String()
		id := uint64(utils.ToInt(args[1].String()))
		err := modules.DelGameDoc(t, id)
		if err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}
		return sciter.NullValue()
	})

	//设为我的最爱
	utils.Window.DefineFunction("SetFavorite", func(args ...*sciter.Value) *sciter.Value {
		id := uint64(utils.ToInt(args[0].String()))
		star := uint8(utils.ToInt(args[1].String()))

		//更新数据
		if err := modules.SetFavorite(id, star); err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}

		return sciter.NewValue("1")
	})

	//设为隐藏
	utils.Window.DefineFunction("SetHide", func(args ...*sciter.Value) *sciter.Value {

		id := uint64(utils.ToInt(args[0].String()))
		hide := uint8(utils.ToInt(args[1].String()))

		//更新数据
		if err := modules.SetHide(id, hide); err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}

		return sciter.NewValue("1")
	})

	//批量设为隐藏
	utils.Window.DefineFunction("SetHideBatch", func(args ...*sciter.Value) *sciter.Value {

		romIdsStr := strings.Split(args[0].String(), ",")
		hide := uint8(utils.ToInt(args[1].String()))

		romIds := []uint64{}
		for _, v := range romIdsStr {
			romIds = append(romIds, uint64(utils.ToInt(v)))
		}

		//更新数据
		if err := modules.SetHideBatch(romIds, hide); err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}

		return sciter.NewValue("1")
	})

	//编辑图片
	utils.Window.DefineFunction("EditRomThumbs", func(args ...*sciter.Value) *sciter.Value {
		typeName := args[0].String()
		id := uint64(utils.ToInt(args[1].String()))
		sid := args[2].String()
		newPath := args[3].String()
		ext := args[4].String()
		newFileName, err := modules.EditRomThumbs(typeName, id, sid, newPath, ext)
		if err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}
		return sciter.NewValue(newFileName)
	})

	//删除图片
	utils.Window.DefineFunction("DeleteThumbs", func(args ...*sciter.Value) *sciter.Value {
		typeName := args[0].String()
		id := uint64(utils.ToInt(args[1].String()))
		sid := args[2].String()
		err := modules.DeleteThumbs(typeName, id, sid)
		if err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}
		return sciter.NullValue()
	})

	//设置图片为主图片
	utils.Window.DefineFunction("SetMasterThumbs", func(args ...*sciter.Value) *sciter.Value {
		typeName := args[0].String()
		id := uint64(utils.ToInt(args[1].String()))
		sid := args[2].String()
		err := modules.SetMasterThumbs(typeName, id, sid)
		if err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}
		return sciter.NullValue()
	})

	//下载展示图
	utils.Window.DefineFunction("DownloadThumbs", func(args ...*sciter.Value) *sciter.Value {
		keyword := args[0].String()
		page := utils.ToInt(args[1].String())
		data, err := modules.DownloadThumbs(keyword, page)
		if err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}

		jsonInfo, _ := json.Marshal(&data)
		return sciter.NewValue(string(jsonInfo))
	})

	//重命名
	utils.Window.DefineFunction("RomRename", func(args ...*sciter.Value) *sciter.Value {
		id := uint64(utils.ToInt(args[0].String()))
		name := args[1].String()

		err := modules.RomRename(id, name)
		if err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}
		return sciter.NullValue()
	})

	//批量重命名
	utils.Window.DefineFunction("BatchRomRename", func(args ...*sciter.Value) *sciter.Value {
		data := args[0].String()

		d := []map[string]string{}
		_ = json.Unmarshal([]byte(data), &d)
		err := modules.BatchRomRename(d)
		if err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}
		return sciter.NullValue()
	})

	//编辑rom基础信息中的名称
	utils.Window.DefineFunction("SetRomBaseName", func(args ...*sciter.Value) *sciter.Value {

		id := uint64(utils.ToInt(args[0].String()))
		name := args[1].String()

		if err := modules.SetRomBaseName(id, name); err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}

		return sciter.NullValue()

	})

	//编辑rom基础信息
	utils.Window.DefineFunction("SetRomBase", func(args ...*sciter.Value) *sciter.Value {

		data := args[0].String()

		d := make(map[string]string)
		_ = json.Unmarshal([]byte(data), &d)

		dbRom, err := modules.SetRomBase(d)
		if err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}

		jsonstr, _ := json.Marshal(&dbRom)
		return sciter.NewValue(string(jsonstr))

	})

	//批量编辑rom基础信息
	utils.Window.DefineFunction("BatchSetRomBase", func(args ...*sciter.Value) *sciter.Value {

		data := args[0].String()

		d := []map[string]string{}
		_ = json.Unmarshal([]byte(data), &d)

		err := modules.BatchSetRomBase(d)
		if err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}

		return sciter.NewValue(1)

	})

	//读取rom基础信息
	utils.Window.DefineFunction("GetRomBase", func(args ...*sciter.Value) *sciter.Value {

		id := uint64(utils.ToInt(args[0].String()))
		rom, _ := (&db.Rom{}).GetById(id)

		baseinfo := modules.GetRomBaseById(rom.Platform, utils.GetFileName(rom.RomPath))

		if baseinfo != nil {
			jsonMenu, _ := json.Marshal(baseinfo)
			return sciter.NewValue(string(jsonMenu))
		}

		return sciter.NullValue()
	})

	//读取过滤器列表
	utils.Window.DefineFunction("GetFilter", func(args ...*sciter.Value) *sciter.Value {
		platform := uint32(utils.ToInt(args[0].String()))

		lists, err := modules.GetFilter(platform)
		if err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}
		jsonstr, _ := json.Marshal(&lists)
		return sciter.NewValue(string(jsonstr))
	})

	//删除rom
	utils.Window.DefineFunction("DeleteRom", func(args ...*sciter.Value) *sciter.Value {
		id := uint64(utils.ToInt(args[0].String()))
		deleteRes := utils.ToInt(args[1].String())

		//删除文件
		err := modules.DeleteRomAndRes(id, deleteRes)
		if err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}

		return sciter.NewValue("1")
	})

	//移动rom及相关资源文件
	utils.Window.DefineFunction("MoveRom", func(args ...*sciter.Value) *sciter.Value {

		id := uint64(utils.ToInt(args[0].String()))
		newFolder := args[1].String()

		if err := modules.MoveRom(id, newFolder); err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}

		return sciter.NullValue()
	})

	//读取相关游戏
	utils.Window.DefineFunction("GetRelatedGames", func(args ...*sciter.Value) *sciter.Value {
		id := uint64(utils.ToInt(args[0].String())) //romId
		romlist, _ := (&db.Rom{}).GetRelatedGames(id)
		jsonRom, _ := json.Marshal(romlist)
		return sciter.NewValue(string(jsonRom))
	})

	//设置评分
	utils.Window.DefineFunction("SetScore", func(args ...*sciter.Value) *sciter.Value {

		id := uint64(utils.ToInt(args[0].String()))
		score := args[1].String()

		//更新数据
		if err := modules.SetScore(id, score); err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}

		return sciter.NewValue("1")
	})

	//设置通关状态
	utils.Window.DefineFunction("SetComplete", func(args ...*sciter.Value) *sciter.Value {

		id := uint64(utils.ToInt(args[0].String()))
		status := uint8(utils.ToInt(args[1].String()))

		//更新数据
		if err := modules.SetComplete(id, status); err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}

		return sciter.NewValue("1")
	})

}
