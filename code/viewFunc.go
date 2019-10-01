package main

import (
	"VirtualNesGUI/code/db"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
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
		id := args[0].String()
		simId, _ := strconv.ParseInt(args[1].String(), 10, 64)

		vo, err := (&db.Rom{}).GetById(id)
		if err != nil {
			if _, err := w.Call("errorBox", sciter.NewValue(err.Error())); err != nil {
			}
		}

		errstr := runGame(vo.Platform, vo.RomPath, simId);
		if errstr != "" {
			if _, err := w.Call("errorBox", sciter.NewValue(errstr)); err != nil {
			}
		}
		return sciter.NullValue()
	})

	//打开rom目录
	w.DefineFunction("OpenFolder", func(args ...*sciter.Value) *sciter.Value {
		gtype := args[0].String() //目录类型
		platform, _ := strconv.ParseInt(args[1].String(), 10, 64)
		p := ""
		switch gtype {
		case "rom":
			p = Config.Platform[platform].RomPath
		case "thumb":
			p = Config.Platform[platform].ThumbPath
		case "video":
			p = Config.Platform[platform].VideoPath
		case "sim":
			exe := Config.Platform[platform].UseSim.Path
			p = filepath.Dir(exe)
		}
		if err := exec.Command(`explorer`, p).Start(); err != nil {
			if _, err := w.Call("errorBox", sciter.NewValue(err)); err != nil {
			}
		}
		return sciter.NullValue()
	})

	//切换模拟器
	w.DefineFunction("SwitchSim", func(args ...*sciter.Value) *sciter.Value {
		simId, _ := strconv.ParseInt(args[0].String(), 10, 64)
		platform, _ := strconv.ParseInt(args[1].String(), 10, 64)

		Config.Platform[platform].UseSim = Config.Platform[platform].SimList[simId]
		err := (&db.Simulator{}).UpdateDefault(platform, simId)

		if err != nil {
			if _, err := w.Call("errorBox", sciter.NewValue(err.Error())); err != nil {
			}
		}
		return sciter.NullValue()
	})

	//更新配置文件
	w.DefineFunction("UpdateConfig", func(args ...*sciter.Value) *sciter.Value {
		field := args[0].String()
		value := args[1].String()

		err := (&db.Config{}).UpdateField(field, value)

		if err != nil {
			if _, err := w.Call("errorBox", sciter.NewValue(err.Error())); err != nil {
			}
		}
		return sciter.NullValue()
	})

	//生成所有缓存
	w.DefineFunction("CreateRomCache", func(args ...*sciter.Value) *sciter.Value {
		//清理数据库
		db.DbClear()

		for k, _ := range Config.Platform {
			err := CreateRomCache(k)
			if err != nil {
				if _, err := w.Call("errorBox", sciter.NewValue(err.Error())); err != nil {
				}
			}
		}
		return sciter.NullValue()
	})

	//读取目录列表
	w.DefineFunction("GetMenuList", func(args ...*sciter.Value) *sciter.Value {
		platform, _ := strconv.ParseInt(args[0].String(), 10, 64)

		//读取数据库
		menu, err := (&db.Menu{}).GetByPlatform(platform) //从数据库中读取当前平台的分类目录

		if err != nil {
			if _, err := w.Call("errorBox", sciter.NewValue(err.Error())); err != nil {
			}
		}
		getjson, _ := json.Marshal(menu)
		return sciter.NewValue(string(getjson))
	})

	//读取游戏列表
	w.DefineFunction("GetGameList", func(args ...*sciter.Value) *sciter.Value {
		platform := strings.Trim(args[0].String(), " ")              //平台
		catname := strings.Trim(args[1].String(), " ")               //分类
		keyword := strings.Trim(args[2].String(), " ")               //关键字
		num := strings.Trim(args[3].String(), " ")                   //字母索引
		page, _ := strconv.Atoi(strings.Trim(args[4].String(), " ")) //分页数

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
		platform := strings.Trim(args[0].String(), " ")
		catname := strings.Trim(args[1].String(), " ")
		keyword := strings.Trim(args[2].String(), " ")
		count, _ := (&db.Rom{}).Count(platform, catname, keyword)
		return sciter.NewValue(strconv.Itoa(count))
	})

	//读取rom详情
	w.DefineFunction("GetGameDetail", func(args ...*sciter.Value) *sciter.Value {
		id := strings.Trim(args[0].String(), " ") //游戏id
		res := &RomDetail{}
		//游戏游戏详细数据
		info, err := (&db.Rom{}).GetById(id)
		if err != nil {
			if _, err := w.Call("errorBox", sciter.NewValue(err.Error())); err != nil {
			}
		}
		//子游戏列表
		sub, _ := (&db.Rom{}).GetSubRom(info.Platform, info.Name)
		res.Info = info
		res.Sublist = sub
		res.DocContent = getDocContent(info.DocPath)
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
			if _, err := w.Call("errorBox", sciter.NewValue(err)); err != nil {
			}
			return sciter.NewValue("0")
		}
		return sciter.NewValue(strconv.Itoa(int(id)))
	})

	//删除一个平台
	w.DefineFunction("DelPlatform", func(args ...*sciter.Value) *sciter.Value {
		idstr := args[0].String()
		id, err := strconv.Atoi(idstr)
		platform := &db.Platform{
			Id: int64(id),
		}

		//删除平台
		err = platform.Delete()
		if err != nil {
			if _, err := w.Call("errorBox", sciter.NewValue(err)); err != nil {
			}
			return sciter.NewValue("0")
		}
		return sciter.NewValue("1")
	})

	//更新平台信息
	w.DefineFunction("UpdatePlatform", func(args ...*sciter.Value) *sciter.Value {
		data := args[0].String()
		d := make(map[string]string)
		json.Unmarshal([]byte(data), &d)

		id, _ := strconv.ParseInt(d["id"], 10, 64)
		status, _ := strconv.ParseInt(d["status"], 10, 64)
		exts := strings.Split(d["exts"], ",")

		platform := &db.Platform{
			Id:        id,
			Name:      d["name"],
			RomExts:   exts,
			RomPath:   d["rom"],
			ThumbPath: d["thumb"],
			SnapPath:  d["snap"],
			VideoPath: d["video"],
			DocPath:   d["doc"],
			Romlist:   d["romlist"],
			Status:    status,
			Pinyin:    TextToPinyin(d["name"]),
		}

		err := platform.UpdateById()
		if err != nil {
			if _, err := w.Call("errorBox", sciter.NewValue(err)); err != nil {
			}
			return sciter.NewValue("0")
		}
		return sciter.NewValue("1")
	})

	//添加模拟器
	w.DefineFunction("AddSimulator", func(args ...*sciter.Value) *sciter.Value {
		data := args[0].String()
		d := make(map[string]string)
		json.Unmarshal([]byte(data), &d)
		pfId, _ := strconv.ParseInt(d["platform"], 10, 64)

		sim := &db.Simulator{
			Name:     d["name"],
			Platform: pfId,
			Path:     d["path"],
			Cmd:      d["cmd"],
			Pinyin:   TextToPinyin(d["name"]),
		}
		id, err := sim.Add()

		//更新默认模拟器
		def, _ := strconv.ParseInt(d["default"], 10, 64)

		if def == 1 {
			err = sim.UpdateDefault(pfId, id)
			if err != nil {
				if _, err := w.Call("errorBox", sciter.NewValue(err)); err != nil {
				}
				return sciter.NewValue("")
			}
		}
		sim.Id = id
		jsonData, _ := json.Marshal(&sim)
		return sciter.NewValue(string(jsonData))
	})

	//添加模拟器
	w.DefineFunction("UpdateSimulator", func(args ...*sciter.Value) *sciter.Value {
		data := args[0].String()
		d := make(map[string]string)
		json.Unmarshal([]byte(data), &d)
		id, _ := strconv.ParseInt(d["id"], 10, 64)
		pfId, _ := strconv.ParseInt(d["platform"], 10, 64)
		sim := &db.Simulator{
			Id:       id,
			Name:     d["name"],
			Platform: pfId,
			Path:     d["path"],
			Cmd:      d["cmd"],
			Pinyin:   TextToPinyin(d["name"]),
		}
		//更新模拟器
		err := sim.UpdateById()
		if err != nil {
			if _, err := w.Call("errorBox", sciter.NewValue(err)); err != nil {
			}
			return sciter.NewValue("")
		}

		jsonData, _ := json.Marshal(&sim)
		return sciter.NewValue(string(jsonData))
	})

	//读取平台详情
	w.DefineFunction("GetPlatformById", func(args ...*sciter.Value) *sciter.Value {
		id, _ := strconv.ParseInt(args[0].String(), 10, 64)

		//游戏游戏详细数据
		info, err := (&db.Platform{}).GetById(id)
		if err != nil {
			if _, err := w.Call("errorBox", sciter.NewValue(err.Error())); err != nil {
			}
		}
		jsonInfo, _ := json.Marshal(&info)
		return sciter.NewValue(string(jsonInfo))
	})

	//读取一个平台下的所有模拟器
	w.DefineFunction("GetPlatform", func(args ...*sciter.Value) *sciter.Value {
		//游戏游戏详细数据
		info, err := (&db.Platform{}).GetAll()
		if err != nil {
			if _, err := w.Call("errorBox", sciter.NewValue(err.Error())); err != nil {
			}
		}
		jsonInfo, _ := json.Marshal(&info)
		return sciter.NewValue(string(jsonInfo))
	})

	//读取模拟器详情
	w.DefineFunction("GetSimulatorById", func(args ...*sciter.Value) *sciter.Value {
		id, _ := strconv.ParseInt(args[0].String(), 10, 64)

		//游戏游戏详细数据
		info, err := (&db.Simulator{}).GetById(id)
		if err != nil {
			if _, err := w.Call("errorBox", sciter.NewValue(err.Error())); err != nil {
			}
		}
		jsonInfo, _ := json.Marshal(&info)
		return sciter.NewValue(string(jsonInfo))
	})

	//读取一个平台下的所有模拟器
	w.DefineFunction("GetSimulatorByPlatform", func(args ...*sciter.Value) *sciter.Value {
		id, _ := strconv.ParseInt(args[0].String(), 10, 64)
		//游戏游戏详细数据
		info, err := (&db.Simulator{}).GetByPlatform(id)
		if err != nil {
			if _, err := w.Call("errorBox", sciter.NewValue(err.Error())); err != nil {
			}
		}
		jsonInfo, _ := json.Marshal(&info)
		return sciter.NewValue(string(jsonInfo))
	})

	//设为我的最爱
	w.DefineFunction("SetFavorite", func(args ...*sciter.Value) *sciter.Value {
		platform, _ := strconv.ParseInt(args[0].String(), 10, 64)
		name := args[1].String()
		star, _ := strconv.ParseInt(args[2].String(), 10, 64)

		fav := &db.Favorite{
			Platform: platform,
			Name:     name,
			Star:     star,
		}

		err := errors.New("")
		if star == 0 {
			err = fav.Delete()
		} else {
			err = fav.UpSert()
		}

		//更新rom表
		rom := &db.Rom{
			Platform: platform,
			Name:     name,
			Star:     star,
		}

		err = rom.UpdateStar()
		//更新数据
		if err != nil {
			if _, err := w.Call("errorBox", sciter.NewValue(err.Error())); err != nil {
			}
			return sciter.NewValue("0")
		}

		//更新数据
		if err != nil {
			if _, err := w.Call("errorBox", sciter.NewValue(err.Error())); err != nil {
			}
			return sciter.NewValue("0")
		}
		return sciter.NewValue("1")
	})

	//更新rom图片
	w.DefineFunction("UpdateRomThumbs", func(args ...*sciter.Value) *sciter.Value {
		typeid, _ := strconv.ParseInt(args[0].String(), 10, 64)
		id, _ := strconv.ParseInt(args[1].String(), 10, 64)
		newpath := args[2].String()

		rom := &db.Rom{
			Id: id,
		}

		//设定新的文件名
		vo, _ := rom.GetById(args[1].String())

		//下载文件
		res, err := http.Get(newpath)
		if err != nil {
			if _, err := w.Call("errorBox", sciter.NewValue(err.Error())); err != nil {
			}
			return sciter.NewValue("")
		}

		//下载成功后，备份原文件
		oldFilePath := ""
		platformPath := ""
		//原图存在，则备份
		isset := false
		if typeid == 1 {
			isset = Exists(vo.ThumbPath)
			oldFilePath = vo.ThumbPath
			platformPath = Config.Platform[vo.Platform].ThumbPath
		} else {
			isset = Exists(vo.SnapPath)
			oldFilePath = vo.SnapPath
			platformPath = Config.Platform[vo.Platform].SnapPath
		}

		if platformPath == "" {
			if _, err := w.Call("errorBox", sciter.NewValue("当前平台没有设置图片目录，请先设置图片目录")); err != nil {
			}
			return sciter.NewValue("")
		}

		if isset == true {

			bakFolder := Config.RootPath + "bak " + separator
			//检测bak文件夹是否存在，不存在这创建bak目录
			folder := ExistsFolder(bakFolder)
			if folder == false {
				_ = os.Mkdir(bakFolder, os.ModePerm);
			}

			oldFileName := filepath.Base(oldFilePath)
			bakFileName := GetFileName(oldFileName) + "_" + strconv.Itoa(int(time.Now().Unix())) + path.Ext(oldFileName)
			err := os.Rename(oldFilePath, bakFolder + bakFileName)

			if err != nil {
				fmt.Println(err.Error())
			}
		}

		//生成新文件
		platformPathAbs, err := filepath.Abs(platformPath)                       //读取平台图片路径
		newFileName := platformPathAbs + separator + vo.Name + path.Ext(newpath) //生成新文件的完整绝路路径地址
		f, err := os.Create(newFileName)
		if err != nil {
			if _, err := w.Call("errorBox", sciter.NewValue(err.Error())); err != nil {
			}
			return sciter.NewValue("")
		}
		io.Copy(f, res.Body)

		if typeid == 1 {
			rom.ThumbPath = newFileName
		} else {
			rom.SnapPath = newFileName
		}

		//游戏游戏详细数据
		err = rom.UpdatePic()
		if err != nil {
			if _, err := w.Call("errorBox", sciter.NewValue(err.Error())); err != nil {
			}
			return sciter.NewValue("")
		}
		return sciter.NewValue(newFileName)
	})

}
