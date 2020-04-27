package controller

import (
	"VirtualNesGUI/code/db"
	"VirtualNesGUI/code/modules"
	"VirtualNesGUI/code/utils"
	"encoding/json"
	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
	"strings"
)

/**
 * 定义view用function
 **/

func RomController(w *window.Window) {

	//运行游戏
	w.DefineFunction("RunGame", func(args ...*sciter.Value) *sciter.Value {

		romId := uint64(utils.ToInt(args[0].String()))
		simId := uint32(utils.ToInt(args[1].String()))

		err := modules.RunGame(romId, simId);
		if err != nil {
			WriteLog(err.Error())
			return ErrorMsg(w, err.Error())
		}
		return sciter.NullValue()
	})

	//运行攻略文件
	w.DefineFunction("RunStrategy", func(args ...*sciter.Value) *sciter.Value {
		f := args[0].String()
		if (f == "") {
			return sciter.NullValue()
		}
		if err := utils.RunGame("explorer", []string{f}); err != nil {
			WriteLog(err.Error())
			return ErrorMsg(w, err.Error())
		}
		return sciter.NullValue()

	})

	//打开rom目录
	w.DefineFunction("OpenFolder", func(args ...*sciter.Value) *sciter.Value {

		id := uint64(utils.ToInt(args[0].String()))
		opt := args[1].String()
		simId := uint32(utils.ToInt(args[2].String()))

		err := modules.OpenFolder(id, opt, simId)
		if err != nil {
			WriteLog(err.Error())
			return ErrorMsg(w, err.Error())
		}

		return sciter.NullValue()
	})

	//读取游戏列表
	w.DefineFunction("GetGameList", func(args ...*sciter.Value) *sciter.Value {
		platform := uint32(utils.ToInt(args[0].String()))        //平台
		catname := strings.Trim(args[1].String(), " ")           //分类
		keyword := strings.Trim(args[2].String(), " ")           //关键字
		num := strings.Trim(args[3].String(), " ")               //字母索引
		page := utils.ToInt(strings.Trim(args[4].String(), " ")) //分页数

		newlist := []*db.Rom{}
		if num == "" {
			newlist, _ = (&db.Rom{}).Get(page, platform, catname, keyword)
		} else {
			//按拼音查询
			newlist, _ = (&db.Rom{}).GetByPinyin(page, platform, catname, num)
		}

		jsonRom, _ := json.Marshal(newlist)
		return sciter.NewValue(string(jsonRom))
	})

	//读取游戏数量
	w.DefineFunction("GetGameCount", func(args ...*sciter.Value) *sciter.Value {
		platform := uint32(utils.ToInt(args[0].String()))
		catname := strings.Trim(args[1].String(), " ")
		keyword := strings.Trim(args[2].String(), " ")
		count, _ := (&db.Rom{}).Count(platform, catname, keyword)
		return sciter.NewValue(utils.ToString(count))
	})

	//读取rom详情
	w.DefineFunction("GetGameDetail", func(args ...*sciter.Value) *sciter.Value {
		id := uint64(utils.ToInt(args[0].String()))
		res, err := modules.GetGameDetail(id)
		if err != nil {
			WriteLog(err.Error())
			return ErrorMsg(w, err.Error())
		}
		jsonMenu, _ := json.Marshal(&res)
		return sciter.NewValue(string(jsonMenu))
	})

	//设为我的最爱
	w.DefineFunction("SetFavorite", func(args ...*sciter.Value) *sciter.Value {
		id := uint64(utils.ToInt(args[0].String()))
		star := uint8(utils.ToInt(args[1].String()))

		//更新rom表
		rom := &db.Rom{
			Id:   id,
			Star: star,
		}

		//更新数据
		if err := rom.UpdateStar(); err != nil {
			WriteLog(err.Error())
			return ErrorMsg(w, err.Error())
		}

		return sciter.NewValue("1")
	})

	//下载rom图片
	w.DefineFunction("DownloadRomThumbs", func(args ...*sciter.Value) *sciter.Value {
		typeName := args[0].String()
		id := uint64(utils.ToInt(args[1].String()))
		newPath := args[2].String()
		newFileName, err := modules.DownloadRomThumbs(typeName, id, newPath)
		if err != nil {
			WriteLog(err.Error())
			return ErrorMsg(w, err.Error())
		}
		return sciter.NewValue(newFileName)
	})

	//编辑图片
	w.DefineFunction("EditRomThumbs", func(args ...*sciter.Value) *sciter.Value {
		typeName := args[0].String()
		id := uint64(utils.ToInt(args[1].String()))
		newPath := args[2].String()
		newFileName, err := modules.EditRomThumbs(typeName, id, newPath)
		if err != nil {
			WriteLog(err.Error())
			return ErrorMsg(w, err.Error())
		}
		return sciter.NewValue(newFileName)
	})

	//删除图片
	w.DefineFunction("DeleteThumbs", func(args ...*sciter.Value) *sciter.Value {
		typeName := args[0].String()
		id := uint64(utils.ToInt(args[1].String()))
		err := modules.DeleteThumbs(typeName, id)
		if err != nil {
			WriteLog(err.Error())
			return ErrorMsg(w, err.Error())
		}
		return sciter.NullValue()
	})

	//重命名
	w.DefineFunction("RomRename", func(args ...*sciter.Value) *sciter.Value {
		setType := uint8(utils.ToInt(args[0].String())) //1:alias,2:filename
		id := uint64(utils.ToInt(args[1].String()))
		name := args[2].String()

		err := modules.RomRename(setType, id, name)
		if err != nil {
			WriteLog(err.Error())
			return ErrorMsg(w, err.Error())
		}
		return sciter.NullValue()

	})

}

