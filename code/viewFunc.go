package main

import (
	"VirtualNesGUI/code/db"
	"VirtualNesGUI/code/utils"
	"encoding/json"
	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

/**
 * 定义view用function
 **/

func defineViewFunction(w *window.Window) {

	w.DefineFunction("InitData", func(args ...*sciter.Value) *sciter.Value {

		ctype := args[0].String()
		isfresh := args[1].String()

		data := ""
		switch (ctype) {
		case "config": //读取配置
			//初始化配置
			if (isfresh == "1") {
				//如果是刷新，则重新生成配置项
				InitConf()
			}
			getjson, _ := json.Marshal(Config)
			data = string(getjson)
		}
		return sciter.NewValue(data)
	})

	//运行游戏
	w.DefineFunction("RunGame", func(args ...*sciter.Value) *sciter.Value {

		romId := uint64(utils.ToInt(args[0].String()))
		simId := uint32(utils.ToInt(args[1].String()))

		//数据库中读取rom详情
		rom, err := (&db.Rom{}).GetById(romId)
		romCmd,_ := (&db.RomCmd{RomId:romId, SimId:simId,}).Get()

		if err != nil {
			return errorMsg(w, err.Error())
		}

		sim := &db.Simulator{}
		if simId == 0 {
			sim = Config.Platform[rom.Platform].UseSim
			if sim == nil {
				return errorMsg(w, Config.Lang["SimulatorNotFound"])
			}
		} else {
			if Config.Platform[rom.Platform].SimList == nil {
				return errorMsg(w, Config.Lang["SimulatorNotFound"])
			}
			sim = Config.Platform[rom.Platform].SimList[simId];
		}

		//检测执行文件是否存在
		_, err = os.Stat(sim.Path)
		if err != nil {
			return errorMsg(w, err.Error())
		}

		//解压zip包
		if sim.Unzip == 1 || romCmd.Unzip == 1{
			rom.RomPath,err = UnzipRom(rom.RomPath, Config.Platform[rom.Platform].RomExts)
			if err != nil{
				return errorMsg(w, err.Error())
			}
		}else{
			rom.RomPath = Config.Platform[rom.Platform].RomPath + separator + rom.RomPath;
		}

		//检测rom文件是否存在
		if utils.FileExists(rom.RomPath) == false {
			return errorMsg(w, Config.Lang["RomNotFound"]+rom.RomPath)
		}

		//加载运行参数
		cmd := []string{}

		ext := utils.GetFileExt(rom.RomPath)

		//如果rom运行参数存在，则使用rom的参数
		if romCmd.Cmd != ""{
			sim.Cmd = romCmd.Cmd
		}

		//运行游戏前，先杀掉之前运行的程序
		if err = killGame();err != nil {
			return errorMsg(w, err.Error())
		}
		//如果是可执行程序，则不依赖模拟器直接运行
		if utils.InSliceString(ext, RUN_EXTS) {
			cmd = append(cmd, rom.RomPath)
			err = runGame("explorer", cmd)
		} else {
			//如果依赖模拟器
			if sim.Cmd == "" {
				cmd = append(cmd, rom.RomPath)
			} else {
				cmd = strings.Split(sim.Cmd, " ")
				filename := filepath.Base(rom.RomPath) //exe运行文件路径
				//替换变量
				for k, _ := range cmd {
					cmd[k] = strings.ReplaceAll(cmd[k], `{RomName}`, utils.GetFileName(filename))
					cmd[k] = strings.ReplaceAll(cmd[k], `{RomExt}`, utils.GetFileExt(filename))
					cmd[k] = strings.ReplaceAll(cmd[k], `{RomFullPath}`, rom.RomPath)
				}
			}
			err = runGame(sim.Path, cmd)
		}

		if err != nil {
			return errorMsg(w, err.Error())
		}
		return sciter.NullValue()
	})

	//运行电子书
	w.DefineFunction("RunBook", func(args ...*sciter.Value) *sciter.Value {

		//检测执行文件是否存在
		_, err := os.Stat(Config.Default.Book)
		if err != nil {
			return errorMsg(w, Config.Lang["BookNotFound"])
		}

		err = runGame(Config.Default.Book, []string{})
		if err != nil {
			return errorMsg(w, err.Error())
		}
		return sciter.NullValue()
	})

	//运行攻略文件
	w.DefineFunction("RunStrategy", func(args ...*sciter.Value) *sciter.Value {
		f := args[0].String()
		if (f == ""){
			return sciter.NullValue()
		}
		if err := runGame("explorer", []string{f});err != nil{
			return errorMsg(w, err.Error())
		}
		return sciter.NullValue()

	})

	//打开rom目录
	w.DefineFunction("OpenFolder", func(args ...*sciter.Value) *sciter.Value {

		id := uint64(utils.ToInt(args[0].String()))
		opt := args[1].String()
		simId := uint32(utils.ToInt(args[2].String()))
		info, err := (&db.Rom{}).GetById(id)
		platform := Config.Platform[info.Platform] //读取当前平台信息
		if err != nil {
			return errorMsg(w, err.Error())
		}
		romName := utils.GetFileName(filepath.Base(info.RomPath)) //读取文件名 }
		fileName := ""
		isDir := false
		switch opt {
		case "rom":
			fileName = platform.RomPath + separator + info.RomPath
		case "thumb":
			if platform.ThumbPath != "" {
				for _, v := range PIC_EXTS {
					fileName = platform.ThumbPath + separator + romName + v

					if utils.FileExists(fileName) {
						break
					} else {
						fileName = ""
					}
				}

				if fileName == "" {
					isDir = true
					fileName = platform.ThumbPath
				}

			}
		case "snap":
			if platform.SnapPath != "" {
				for _, v := range PIC_EXTS {
					fileName = platform.SnapPath + separator + romName + v
					if utils.FileExists(fileName) {
						break
					} else {
						fileName = ""
					}
				}
				if fileName == "" {
					isDir = true
					fileName = platform.SnapPath + separator
				}
			}
		case "doc":
			if platform.DocPath != "" {
				for _, v := range DOC_EXTS {
					fileName = platform.DocPath + separator + romName + v
					if utils.FileExists(fileName) {
						break
					} else {
						fileName = ""
					}
				}
				if fileName == "" {
					isDir = true
					fileName = platform.DocPath
				}
			}
		case "strategy":
			if platform.StrategyPath != "" {
				for _, v := range DOC_EXTS {
					fileName = platform.StrategyPath + separator + romName + v
					if utils.FileExists(fileName) {
						break
					} else {
						fileName = ""
					}
				}
				if fileName == "" {
					isDir = true
					fileName = platform.StrategyPath
				}
			}
		case "sim":
			if _, ok := platform.SimList[simId]; ok {
				fileName = platform.SimList[simId].Path
			}
		}

		if fileName != "" {
			if isDir == true {
				if err := exec.Command(`explorer`, fileName).Start(); err != nil {
					return errorMsg(w, err.Error())
				}
			} else {
				if err := exec.Command(`explorer`, `/select,`, `/n,`, fileName).Start(); err != nil {
					return errorMsg(w, err.Error())
				}
			}
		} else {
			return errorMsg(w, Config.Lang["PathNotFound"])
		}

		return sciter.NullValue()
	})

	//更新配置文件
	w.DefineFunction("UpdateConfig", func(args ...*sciter.Value) *sciter.Value {
		field := args[0].String()
		value := args[1].String()

		err := (&db.Config{}).UpdateField(field, value)

		if err != nil {
			return errorMsg(w, err.Error())
		}
		return sciter.NullValue()
	})

	//删除所有缓存
	w.DefineFunction("TruncateRomCache", func(args ...*sciter.Value) *sciter.Value {

		//清空rom表
		if err := (&db.Rom{}).Truncate(); err != nil {
			return errorMsg(w, err.Error())
		}

		//清空menu表
		if err := (&db.Menu{}).Truncate(); err != nil {
			return errorMsg(w, err.Error())
		}

		return sciter.NullValue()
	})

	//生成所有缓存
	w.DefineFunction("CreateRomCache", func(args ...*sciter.Value) *sciter.Value {
		//先检查平台，将不存在的平台数据先干掉
		if err := ClearPlatform(); err != nil {
			return errorMsg(w, err.Error())
		}

		//开始重建缓存
		for platform, _ := range Config.Platform {

			//创建rom数据
			romlist, menu, err := CreateRomCache(platform)
			if err != nil {
				return errorMsg(w, err.Error())
			}

			//更新rom数据
			if err := UpdateRomDB(platform, romlist); err != nil {
				return errorMsg(w, err.Error())
			}

			//更新menu数据
			if err := UpdateMenuDB(platform, menu); err != nil {
				return errorMsg(w, err.Error())
			}

		}
		return sciter.NullValue()
	})

	//读取目录列表
	w.DefineFunction("GetMenuList", func(args ...*sciter.Value) *sciter.Value {
		platform := uint32(utils.ToInt(args[0].String()))

		//读取数据库
		menu, err := (&db.Menu{}).GetByPlatform(platform) //从数据库中读取当前平台的分类目录

		//读取根目录下是否有rom
		count, err := (&db.Rom{}).Count(platform, constMenuRootKey, "")
		newMenu := []*db.Menu{}

		//读取根目录下有rom，则显示未分类文件夹
		if count > 0 {
			root := &db.Menu{
				Name:     constMenuRootKey,
				Platform: platform,
			}
			newMenu = append(newMenu, root)
			newMenu = append(newMenu, menu...)
		} else {
			newMenu = menu
		}

		if err != nil {
			return errorMsg(w, err.Error())
		}
		getjson, _ := json.Marshal(newMenu)
		return sciter.NewValue(string(getjson))
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

	//更新菜单排序
	w.DefineFunction("UpdateMenuSort", func(args ...*sciter.Value) *sciter.Value {
		platform := uint32(utils.ToInt(args[0].String())) //平台
		data := args[1].String()

		d := make(map[string]uint32)
		json.Unmarshal([]byte(data), &d)

		if len(d) == 0 {
			return sciter.NullValue()
		}

		for name, val := range d {
			if name == "" {
				continue
			}
			menu := &db.Menu{
				Name:     name,
				Platform: platform,
				Sort:     val,
			}
			if err := menu.UpdateSortByName(); err != nil {
			}
		}
		return sciter.NewValue("1")
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
		res := &RomDetail{}
		//游戏游戏详细数据
		info, err := (&db.Rom{}).GetById(id)
		if err != nil {
			return errorMsg(w, err.Error())
		}
		//子游戏列表
		sub, _ := (&db.Rom{}).GetSubRom(info.Platform, info.Name)

		res.Info = info
		res.StrategyFile = ""
		res.Sublist = sub

		//读取文档内容
		romName := utils.GetFileName(filepath.Base(info.RomPath)) //生成新文件的完整绝路路径地址
		if Config.Platform[info.Platform].DocPath != "" {
			docFileName := "";
			for _, v := range DOC_EXTS {
				docFileName = Config.Platform[info.Platform].DocPath + separator + romName + v
				res.DocContent = getDocContent(docFileName)
				if res.DocContent != "" {
					break
				}
			}
		}

		if Config.Platform[info.Platform].StrategyPath != "" {
			//检测攻略可执行文件是否存在
			strategyFileName := "";
			for _, v := range RUN_EXTS {
				strategyFileName = Config.Platform[info.Platform].StrategyPath + separator + romName + v
				if utils.FileExists(strategyFileName){
					res.StrategyFile = strategyFileName
					break
				}
			}

			//如果没有执行运行的文件，则读取文档内容
			if strategyFileName != ""{
				for _, v := range DOC_EXTS {
					strategyFileName = Config.Platform[info.Platform].StrategyPath + separator + romName + v
					res.StrategyContent = getDocContent(strategyFileName)
					if res.StrategyContent != "" {
						break
					}
				}
			}

		}

		jsonMenu, _ := json.Marshal(&res)
		return sciter.NewValue(string(jsonMenu))
	})

	//添加一个平台
	w.DefineFunction("AddPlatform", func(args ...*sciter.Value) *sciter.Value {
		name := args[0].String()
		platform := &db.Platform{
			Name:   name,
			Pinyin: TextToPinyin(name),
		}
		id, err := platform.Add()
		if err != nil {
			return errorMsg(w, err.Error())
		}
		return sciter.NewValue(utils.ToString(id))
	})

	//删除一个平台
	w.DefineFunction("DelPlatform", func(args ...*sciter.Value) *sciter.Value {
		id := uint32(utils.ToInt(args[0].String()))
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
		err := rom.DeleteByPlatform()
		//删除模拟器
		err = sim.DeleteByPlatform()
		//删除平台
		err = platform.DeleteById()

		if err != nil {
			return errorMsg(w, err.Error())
		}
		return sciter.NewValue("1")
	})

	//更新平台信息
	w.DefineFunction("UpdatePlatform", func(args ...*sciter.Value) *sciter.Value {
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
		d["strategy"] = strings.TrimRight(d["strategy"], `\`)
		d["strategy"] = strings.TrimRight(d["strategy"], `/`)
		d["doc"] = strings.TrimRight(d["doc"], `\`)
		d["doc"] = strings.TrimRight(d["doc"], `/`)

		exts := strings.Split(d["exts"], ",")

		platform := &db.Platform{
			Id:           id,
			Name:         d["name"],
			RomExts:      exts,
			RomPath:      d["rom"],
			ThumbPath:    d["thumb"],
			SnapPath:     d["snap"],
			StrategyPath: d["strategy"],
			DocPath:      d["doc"],
			Romlist:      d["romlist"],
			Pinyin:       TextToPinyin(d["name"]),
		}

		err := platform.UpdateById()
		if err != nil {
			return errorMsg(w, err.Error())
		}
		return sciter.NewValue("1")
	})

	//更新平台排序
	w.DefineFunction("UpdatePlatformSort", func(args ...*sciter.Value) *sciter.Value {
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
				return errorMsg(w, err.Error())
			}
		}
		return sciter.NewValue("1")
	})

	//添加模拟器
	w.DefineFunction("AddSimulator", func(args ...*sciter.Value) *sciter.Value {
		data := args[0].String()
		d := make(map[string]string)
		json.Unmarshal([]byte(data), &d)
		pfId := uint32(utils.ToInt(d["platform"]))

		sim := &db.Simulator{
			Name:     d["name"],
			Platform: pfId,
			Path:     d["path"],
			Cmd:      d["cmd"],
			Pinyin:   TextToPinyin(d["name"]),
		}
		id, err := sim.Add()

		//更新默认模拟器
		if utils.ToInt(d["default"]) == 1 {
			err = sim.UpdateDefault(pfId, id)
			if err != nil {
				return errorMsg(w, err.Error())
			}
		}
		sim.Id = id
		jsonData, _ := json.Marshal(&sim)
		return sciter.NewValue(string(jsonData))
	})

	//更新模拟器
	w.DefineFunction("UpdateSimulator", func(args ...*sciter.Value) *sciter.Value {
		data := args[0].String()
		d := make(map[string]string)
		json.Unmarshal([]byte(data), &d)
		id := uint32(utils.ToInt(d["id"]))
		pfId := uint32(utils.ToInt(d["platform"]))
		def := uint8(utils.ToInt(d["default"]))
		sim := &db.Simulator{
			Id:       id,
			Name:     d["name"],
			Platform: pfId,
			Path:     d["path"],
			Cmd:      d["cmd"],
			Pinyin:   TextToPinyin(d["name"]),
		}

		//更新模拟器
		if err := sim.UpdateById(); err != nil {
			return errorMsg(w, err.Error())
		}

		//如果设置了默认模拟器，则更新默认模拟器
		if def == 1 {
			if err := sim.UpdateDefault(pfId, id); err != nil {
				return errorMsg(w, err.Error())
			}
		}

		jsonData, _ := json.Marshal(&sim)
		return sciter.NewValue(string(jsonData))
	})

	//删除一个模拟器
	w.DefineFunction("DelSimulator", func(args ...*sciter.Value) *sciter.Value {
		id := uint32(utils.ToInt(args[0].String()))

		sim := &db.Simulator{
			Id: id,
		}

		//删除模拟器
		err := sim.DeleteById()
		if err != nil {
			return errorMsg(w, err.Error())
		}

		//删除rom独立模拟器cmd配置
		if err = (&db.RomCmd{SimId:id}).ClearBySimId();err != nil{
			return errorMsg(w, err.Error())
		}
		return sciter.NewValue(args[0].String())
	})

	//读取平台详情
	w.DefineFunction("GetPlatformById", func(args ...*sciter.Value) *sciter.Value {
		id := uint32(utils.ToInt(args[0].String()))

		//游戏游戏详细数据
		info, err := (&db.Platform{}).GetById(id)
		if err != nil {
			return errorMsg(w, err.Error())
		}
		jsonInfo, _ := json.Marshal(&info)
		return sciter.NewValue(string(jsonInfo))
	})

	//读取平台列表
	w.DefineFunction("GetPlatform", func(args ...*sciter.Value) *sciter.Value {
		//游戏游戏详细数据
		info, err := (&db.Platform{}).GetAll()
		if err != nil {
			return errorMsg(w, err.Error())
		}
		jsonInfo, _ := json.Marshal(&info)
		return sciter.NewValue(string(jsonInfo))
	})

	//读取模拟器详情
	w.DefineFunction("GetSimulatorById", func(args ...*sciter.Value) *sciter.Value {
		id := uint32(utils.ToInt(args[0].String()))

		//游戏游戏详细数据
		info, err := (&db.Simulator{}).GetById(id)
		if err != nil {
			return errorMsg(w, err.Error())
		}
		jsonInfo, _ := json.Marshal(&info)
		return sciter.NewValue(string(jsonInfo))
	})

	//读取一个平台下的所有模拟器
	w.DefineFunction("GetSimulatorByPlatform", func(args ...*sciter.Value) *sciter.Value {
		id := uint32(utils.ToInt(args[0].String()))
		//游戏游戏详细数据
		info, err := (&db.Simulator{}).GetByPlatform(id)
		if err != nil {
			return errorMsg(w, err.Error())
		}
		jsonInfo, _ := json.Marshal(&info)
		return sciter.NewValue(string(jsonInfo))
	})

	//读取rom独立模拟器cmd数据
	w.DefineFunction("GetRomCmd", func(args ...*sciter.Value) *sciter.Value {
		romId := uint64(utils.ToInt(args[0].String()))
		simId := uint32(utils.ToInt(args[1].String()))

		//数据库中读取rom详情
		rom, _ := (&db.RomCmd{RomId:romId, SimId:simId,}).Get()

		romJson, _ := json.Marshal(&rom)
		return sciter.NewValue(string(romJson))
	})

	//添加rom独立模拟器参数
	w.DefineFunction("AddRomCmd", func(args ...*sciter.Value) *sciter.Value {
		romId := uint64(utils.ToInt(args[0].String()))
		simId := uint32(utils.ToInt(args[1].String()))
		data := args[2].String()
		d := make(map[string]string)
		json.Unmarshal([]byte(data), &d)

		romCmd := &db.RomCmd{
			RomId:romId,
			SimId:simId,
			Cmd:d["cmd"],
			Unzip: uint8(utils.ToInt(d["unzip"])),
		}


		//如果当前配置和模拟器默认配置一样，则不用添加
		sim,_ := (&db.Simulator{}).GetById(simId)
		if romCmd.Cmd == sim.Cmd && romCmd.Unzip == sim.Unzip{
			return sciter.NullValue()
		}

		if err := romCmd.Add();err != nil{
			return errorMsg(w, err.Error())
		}

		return sciter.NullValue()
	})

	//更新rom独立模拟器参数
	w.DefineFunction("UpdateRomCmd", func(args ...*sciter.Value) *sciter.Value {
		id := uint32(utils.ToInt(args[0].String()))
		simId := uint32(utils.ToInt(args[1].String()))
		data := args[2].String()
		d := make(map[string]string)
		json.Unmarshal([]byte(data), &d)

		romCmd := &db.RomCmd{
			Id:id,
			Cmd:d["cmd"],
			Unzip: uint8(utils.ToInt(d["unzip"])),
		}

		//如果当前配置和模拟器默认配置一样，则删除该记录
		sim,_ := (&db.Simulator{}).GetById(simId)
		if romCmd.Cmd == sim.Cmd && romCmd.Unzip == sim.Unzip{
			if err := romCmd.DeleteById();err != nil{
				return errorMsg(w, err.Error())
			}
			return sciter.NullValue()
		}

		//开始更新
		if err := romCmd.UpdateCmd();err != nil{
			return errorMsg(w, err.Error())
		}

		return sciter.NullValue()
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

		err := rom.UpdateStar()
		//更新数据
		if err != nil {
			return errorMsg(w, err.Error())
		}
		return sciter.NewValue("1")
	})

	//设置rom的模拟器
	w.DefineFunction("SetRomSimulator", func(args ...*sciter.Value) *sciter.Value {
		romId := uint64(utils.ToInt(args[0].String()))
		simId := uint32(utils.ToInt(args[1].String()))
		//更新rom表
		rom := &db.Rom{
			Id:    romId,
			SimId: simId,
		}
		err := rom.UpdateSimulator()
		//更新数据
		if err != nil {
			return errorMsg(w, err.Error())
		}
		return sciter.NewValue("1")
	})

	//更新rom图片
	w.DefineFunction("UpdateRomThumbs", func(args ...*sciter.Value) *sciter.Value {
		typeid := utils.ToInt(args[0].String())
		id := uint64(utils.ToInt(args[1].String()))
		newpath := args[2].String()

		rom := &db.Rom{
			Id: id,
		}

		//设定新的文件名
		vo, _ := rom.GetById(id)

		//下载文件
		res, err := http.Get(newpath)
		if err != nil {
			return errorMsg(w, err.Error())
		}

		//下载成功后，备份原文件
		platformPath := ""
		//原图存在，则备份
		if typeid == 1 {
			platformPath = Config.Platform[vo.Platform].ThumbPath
		} else {
			platformPath = Config.Platform[vo.Platform].SnapPath
		}

		if platformPath == "" {
			return errorMsg(w, Config.Lang["NoSetThumbDir"])
		}

		//开始备份原图
		bakFolder := Config.RootPath + "bak" + separator
		RomFileName := utils.GetFileName(vo.RomPath)

		//检测bak文件夹是否存在，不存在这创建bak目录
		folder := utils.PathExists(bakFolder)
		if folder == false {
			_ = os.Mkdir(bakFolder, os.ModePerm);
		}
		for _, ext := range PIC_EXTS {
			oldFileName := RomFileName + ext //老图片文件名
			if utils.FileExists(oldFileName) {
				bakFileName := RomFileName + "_" + utils.ToString(time.Now().Unix()) + ext //生成备份文件名
				err := os.Rename(oldFileName, bakFolder+bakFileName)                       //移动文件
				if err != nil {
					return errorMsg(w, err.Error())
				}
			}
		}

		//生成新文件
		platformPathAbs, err := filepath.Abs(platformPath) //读取平台图片路径

		newFileName := platformPathAbs + separator + RomFileName + utils.GetFileExt(newpath) //生成新文件的完整绝路路径地址
		f, err := os.Create(newFileName)
		if err != nil {
			return errorMsg(w, err.Error())
		}
		io.Copy(f, res.Body)

		return sciter.NewValue(newFileName)
	})

}
