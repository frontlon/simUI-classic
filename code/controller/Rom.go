package controller

import (
	"simUI/code/db"
	"simUI/code/modules"
	"simUI/code/utils"
	"encoding/json"
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
		if err := utils.RunGame("explorer", []string{f}); err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}
		return sciter.NullValue()

	})

	//打开rom目录
	utils.Window.DefineFunction("OpenFolder", func(args ...*sciter.Value) *sciter.Value {

		id := uint64(utils.ToInt(args[0].String()))
		opt := args[1].String()
		simId := uint32(utils.ToInt(args[2].String()))

		err := modules.OpenFolder(id, opt, simId)
		if err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}

		return sciter.NullValue()
	})

	//读取游戏列表
	utils.Window.DefineFunction("GetGameList", func(args ...*sciter.Value) *sciter.Value {
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
	utils.Window.DefineFunction("GetGameCount", func(args ...*sciter.Value) *sciter.Value {
		platform := uint32(utils.ToInt(args[0].String()))
		catname := strings.Trim(args[1].String(), " ")
		keyword := strings.Trim(args[2].String(), " ")
		count, _ := (&db.Rom{}).Count(platform, catname, keyword)
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

	//设为我的最爱
	utils.Window.DefineFunction("SetFavorite", func(args ...*sciter.Value) *sciter.Value {
		id := uint64(utils.ToInt(args[0].String()))
		star := uint8(utils.ToInt(args[1].String()))

		//更新rom表
		rom := &db.Rom{
			Id:   id,
			Star: star,
		}

		//更新数据
		if err := rom.UpdateStar(); err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}

		return sciter.NewValue("1")
	})

	//设为隐藏
	utils.Window.DefineFunction("SetHide", func(args ...*sciter.Value) *sciter.Value {
		id := uint64(utils.ToInt(args[0].String()))
		ishide := uint8(utils.ToInt(args[1].String()))

		//更新rom表
		rom := &db.Rom{
			Id:   id,
			Hide: ishide,
		}

		//更新数据
		if err := rom.UpdateHide(); err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}

		return sciter.NewValue("1")
	})

	//下载rom图片
	utils.Window.DefineFunction("DownloadRomThumbs", func(args ...*sciter.Value) *sciter.Value {
		typeName := args[0].String()
		id := uint64(utils.ToInt(args[1].String()))
		newPath := args[2].String()
		newFileName, err := modules.DownloadRomThumbs(typeName, id, newPath)
		if err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}
		return sciter.NewValue(newFileName)
	})

	//编辑图片
	utils.Window.DefineFunction("EditRomThumbs", func(args ...*sciter.Value) *sciter.Value {
		typeName := args[0].String()
		id := uint64(utils.ToInt(args[1].String()))
		newPath := args[2].String()
		newFileName, err := modules.EditRomThumbs(typeName, id, newPath)
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
		err := modules.DeleteThumbs(typeName, id)
		if err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}
		return sciter.NullValue()
	})

	//重命名
	utils.Window.DefineFunction("RomRename", func(args ...*sciter.Value) *sciter.Value {
		setType := args[0].String() //1:alias,2:filename
		id := uint64(utils.ToInt(args[1].String()))
		name := args[2].String()

		err := modules.RomRename(setType, id, name)
		if err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}
		return sciter.NullValue()
	})

	//编辑rom基础信息
	utils.Window.DefineFunction("SetRomBase", func(args ...*sciter.Value) *sciter.Value {

		data := args[0].String()
		d := make(map[string]string)
		json.Unmarshal([]byte(data), &d)
		platform := uint32(utils.ToInt(args[0].String()))

		romBase := &modules.RomBase{
			RomName:   d["rom_name"],
			Name:      d["name"],
			Type:      d["type"],
			Year:      d["year"],
			Developer: d["developer"],
			Publisher: d["publisher"],
		}

		//写入配置文件
		if err := modules.WriteRomBaseFile(platform, romBase); err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}

		return sciter.NullValue()
	})

}
