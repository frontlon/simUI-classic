package modules

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"github.com/go-ini/ini"
	"io"
	"os"
	"simUI/code/config"
	"simUI/code/db"
	"simUI/code/utils"
	"time"
)

var dir = ""

func InputPlatform(p string) error {

	//解压文件
	dir = "games" + config.Cfg.Separator + "input_" + utils.ToString(time.Now().Unix()) + config.Cfg.Separator
	if err := deCompress(p, dir); err != nil {
		return err
	}

	//处理配置文件
	inputIniConfig(dir + "/config.ini")

	return nil
}

//解压
func deCompress(zipFile, dest string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			return err
		}
		defer rc.Close()
		filename := dest + file.Name
		err = os.MkdirAll(utils.GetFileAbsPath(filename), 0755)
		if err != nil {
			return err
		}
		w, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer w.Close()
		_, err = io.Copy(w, rc)
		if err != nil {
			return err
		}
		w.Close()
		rc.Close()
	}
	return nil
}

//处理配置文件
func inputIniConfig(p string) error {
	//载入配置文件
	file, err := ini.LoadSources(ini.LoadOptions{IgnoreInlineComment: true}, p)
	if err != nil {
		return err
	}

	//导入平台数据
	platformId, err := inputPlatformData(file)
	if err != nil {
		return err
	}
	//导入模拟器数据
	simIds, err := inputSimulatormData(file, platformId)
	if err != nil {
		return err
	}
	//更新rom缓存
	config.InitConf()
	if err := CreateRomCache(platformId); err != nil {
		fmt.Println(err)
		return err
	}

	//写入rom数据
	starList := []string{}
	json.Unmarshal([]byte(utils.Base64Decode(file.Section("rom").Key("star").String())), &starList)
	hideList := []string{}
	json.Unmarshal([]byte(utils.Base64Decode(file.Section("rom").Key("hide").String())), &hideList)
	simList := map[string]string{}
	json.Unmarshal([]byte(utils.Base64Decode(file.Section("rom").Key("sim").String())), &simList)

	simIdList := map[string][]string{}
	json.Unmarshal([]byte(utils.Base64Decode(file.Section("rom").Key("simId").String())), &simIdList)
	//合并name，根据name读取ids
	nameList := []string{}
	for _, v := range starList {
		nameList = append(nameList, v)
	}
	for _, v := range hideList {
		nameList = append(nameList, v)
	}
	for k, _ := range simList {
		nameList = append(nameList, k)
	}
	for _, v := range simIdList {
		for _, b := range v {
			nameList = append(nameList, b)
		}
	}

	nameList = utils.SliceRemoveDuplicate(nameList)

	IdList, _ := (&db.Rom{}).GetIdsByNames(platformId, nameList)

	//更新star到数据库
	starIds := []uint64{}
	for _, v := range starList {
		starIds = append(starIds, IdList[v])
	}
	(&db.Rom{}).UpdateStarByIds(starIds, 1)

	//更新hide到数据库
	hideIds := []uint64{}
	for _, v := range hideList {
		hideIds = append(hideIds, IdList[v])
	}
	(&db.Rom{}).UpdateHideByIds(hideIds, 1)

	//更新simId到数据库
	for k, v := range simIdList {
		simListIds := []uint64{}
		for _, b := range v {
			simListIds = append(simListIds, IdList[b])
		}
		(&db.Rom{}).UpdateSimIdByIds(simListIds, simIds[k])
	}

	//更新sim_conf到数据库
	for romName, v := range simList {
		romData := map[string]map[string]string{}
		json.Unmarshal([]byte(v), &romData)

		romId := IdList[romName]

		//遍历rom
		simData := map[string]map[string]string{}
		for simName, b := range romData {
			simData[utils.ToString(simIds[simName])] = b
		}

		//写入数据库
		create, _ := json.Marshal(&simData)
		(&db.Rom{}).UpdateSimConfById(romId, string(create))
	}

	return nil
}

//导入平台数据
func inputPlatformData(file *ini.File) (uint32, error) {
	//载入平台数据
	ico := ""
	rom := ""
	thumb := ""
	snap := ""
	poster := ""
	packing := ""
	title := ""
	cassette := ""
	icon := ""
	gif := ""
	background := ""
	wallpaper := ""
	docs := ""
	strategies := ""
	video := ""
	files := ""
	rombase := ""
	platformSection := file.Section("platform")

	if platformSection.Key("ico").String() != "" {
		ico = dir + platformSection.Key("ico").String()
	}
	if platformSection.Key("rom").String() != "" {
		rom = dir + platformSection.Key("rom").String()
	}
	if platformSection.Key("thumb").String() != "" {
		thumb = dir + platformSection.Key("thumb").String()
	}
	if platformSection.Key("snap").String() != "" {
		snap = dir + platformSection.Key("snap").String()
	}
	if platformSection.Key("poster").String() != "" {
		poster = dir + platformSection.Key("poster").String()
	}
	if platformSection.Key("packing").String() != "" {
		packing = dir + platformSection.Key("packing").String()
	}
	if platformSection.Key("title").String() != "" {
		title = dir + platformSection.Key("title").String()
	}
	if platformSection.Key("cassette").String() != "" {
		cassette = dir + platformSection.Key("cassette").String()
	}
	if platformSection.Key("icon").String() != "" {
		icon = dir + platformSection.Key("icon").String()
	}
	if platformSection.Key("gif").String() != "" {
		gif = dir + platformSection.Key("gif").String()
	}
	if platformSection.Key("background").String() != "" {
		background = dir + platformSection.Key("background").String()
	}
	if platformSection.Key("wallpaper").String() != "" {
		wallpaper = dir + platformSection.Key("wallpaper").String()
	}
	if platformSection.Key("docs").String() != "" {
		docs = dir + platformSection.Key("docs").String()
	}
	if platformSection.Key("strategies").String() != "" {
		strategies = dir + platformSection.Key("strategies").String()
	}
	if platformSection.Key("video").String() != "" {
		video = dir + platformSection.Key("video").String()
	}
	if platformSection.Key("files").String() != "" {
		files = dir + platformSection.Key("files").String()
	}
	if platformSection.Key("rombase").String() != "" {
		rombase = dir + platformSection.Key("rombase").String()
	}

	platformDom := &db.Platform{
		Name:           platformSection.Key("name").String(),
		Icon:           ico,
		RomExts:        platformSection.Key("exts").String(),
		RomPath:        rom,
		ThumbPath:      thumb,
		SnapPath:       snap,
		PosterPath:     poster,
		PackingPath:    packing,
		TitlePath:      title,
		CassettePath:   cassette,
		IconPath:       icon,
		GifPath:        gif,
		BackgroundPath: background,
		WallpaperPath:  wallpaper,
		DocPath:        docs,
		StrategyPath:   strategies,
		VideoPath:      video,
		FilesPath:      files,
		Rombase:        rombase,
		Pinyin:         utils.TextToPinyin(platformSection.Key("name").String()),
		Sort:           0,
	}

	platformId, err := platformDom.Add()
	if err != nil {
		return 0, err
	}
	return platformId, nil
}

//导入模拟器数据
func inputSimulatormData(file *ini.File, platformId uint32) (map[string]uint32, error) {
	simSection := file.Section("simulator")
	simIds := map[string]uint32{}
	if simSection != nil {
		sims := simSection.ChildSections()
		for _, v := range sims {
			sim := &db.Simulator{
				Name:     v.Key("name").String(),
				Platform: uint32(utils.ToInt(platformId)),
				Path:     v.Key("path").String(),
				Cmd:      v.Key("cmd").String(),
				Unzip:    uint8(utils.ToInt(v.Key("unzip").String())),
				Default:  uint8(utils.ToInt(v.Key("default").String())),
				Pinyin:   utils.TextToPinyin(v.Key("name").String()),
			}
			simIds[v.Key("name").String()], _ = sim.Add()
		}
	}
	return simIds, nil
}
