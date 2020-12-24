package controller

import (
	"encoding/json"
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
		if err := utils.RunGame("", []string{f}); err != nil {
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

		baseType := args[5].String()
		basePublisher := args[6].String()
		baseYear := args[7].String()
		baseCountry := args[8].String()
		baseTranslate := args[9].String()

		if baseType == config.Cfg.Lang["BaseType"] {
			baseType = ""
		}

		if basePublisher == config.Cfg.Lang["BasePublisher"] {
			basePublisher = ""
		}
		if baseYear == config.Cfg.Lang["BaseYear"] {
			baseYear = ""
		}
		if baseCountry == config.Cfg.Lang["BaseCountry"] {
			baseCountry = ""
		}
		if baseTranslate == config.Cfg.Lang["BaseTranslate"] {
			baseTranslate = ""
		}

		newlist := []*db.Rom{}
		if num == "" {
			newlist, _ = (&db.Rom{}).Get(page, platform, catname, keyword, baseType, basePublisher, baseYear, baseCountry, baseTranslate)
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
		keyword := strings.Trim(args[2].String(), " ") //关键字
		baseType := args[3].String()
		basePublisher := args[4].String()
		baseYear := args[5].String()
		baseCountry := args[6].String()
		baseTranslate := args[7].String()
		count, _ := (&db.Rom{}).Count(platform, catname, keyword, baseType, basePublisher, baseYear, baseCountry, baseTranslate)
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

	//读取游戏攻略内容
	utils.Window.DefineFunction("GetGameDoc", func(args ...*sciter.Value) *sciter.Value {
		t := args[0].String()
		id := uint64(utils.ToInt(args[1].String()))
		res, err := modules.GetGameDoc(t, id)
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

	//上传攻略图片
	utils.Window.DefineFunction("UploadStrategyImages", func(args ...*sciter.Value) *sciter.Value {
		id := uint64(utils.ToInt(args[0].String()))
		p := args[1].String()
		relative, err := modules.UploadStrategyImages(id, p)
		if err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}
		return sciter.NewValue(relative)
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
		id := uint64(utils.ToInt(args[0].String()))
		name := args[1].String()

		err := modules.RomRename(id, name)
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

		rom, _ := (&db.Rom{}).GetById(uint64(utils.ToInt(d["id"])))
		romName := utils.GetFileName(rom.RomPath)
		romBase := &modules.RomBase{
			RomName:   romName,
			Name:      d["name"],
			Type:      d["type"],
			Year:      d["year"],
			Publisher: d["publisher"],
			Country:   d["country"],
			Translate: d["translate"],
		}

		//写入配置文件
		if err := modules.WriteRomBaseFile(rom.Platform, romBase); err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}
		name := ""
		if d["name"] == "" {
			name = romName
		} else {
			name = d["name"]
		}

		//更新到数据库
		dbRom := &db.Rom{
			Name:          name,
			BaseType:      d["type"],
			BaseYear:      d["year"],
			BasePublisher: d["publisher"],
			BaseCountry:   d["country"],
			BaseTranslate: d["translate"],
		}
		if err := dbRom.UpdateRomBase(uint64(utils.ToInt(d["id"]))); err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}

		jsonstr, _ := json.Marshal(&dbRom)
		return sciter.NewValue(string(jsonstr))

	})

	//读取rom基础信息
	utils.Window.DefineFunction("GetRomBase", func(args ...*sciter.Value) *sciter.Value {

		id := uint64(utils.ToInt(args[0].String()))
		rom, _ := (&db.Rom{}).GetById(id)

		romName := utils.GetFileName(rom.RomPath)
		baseinfo, _ := modules.GetRomBase(rom.Platform)

		if _, ok := baseinfo[romName]; ok {
			jsonMenu, _ := json.Marshal(baseinfo[romName])
			return sciter.NewValue(string(jsonMenu))
		}

		return sciter.NullValue()
	})

	//读取过滤器列表
	utils.Window.DefineFunction("GetFilter", func(args ...*sciter.Value) *sciter.Value {
		platform := uint32(utils.ToInt(args[0].String()))
		t := args[1].String()
		lists, err := (&db.Filter{}).GetByPlatform(platform, t)
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

		//删除文件
		err := modules.DeleteRomAndRes(id)
		if err != nil {
			utils.WriteLog(err.Error())
		}

		//删除数据库缓存
		err = (&db.Rom{}).DeleteById(id)
		if err != nil {
			utils.WriteLog(err.Error())
		}

		return sciter.NewValue("1")
	})

	//移动rom及相关资源文件
	utils.Window.DefineFunction("MoveRom", func(args ...*sciter.Value) *sciter.Value {

		id := uint64(utils.ToInt(args[0].String()))
		newPlatform := uint32(utils.ToInt(args[1].String()))
		newFolder := args[2].String()

		if err := modules.MoveRom(id, newPlatform, newFolder); err != nil {
			utils.WriteLog(err.Error())
			return utils.ErrorMsg(err.Error())
		}

		return sciter.NullValue()
	})

}
