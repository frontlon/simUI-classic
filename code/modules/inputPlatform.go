package modules

import (
	"archive/zip"
	"fmt"
	"github.com/go-ini/ini"
	"io"
	"os"
	"simUI/code/db"
	"simUI/code/utils"
	"time"
)

var dir = ""

func InputPlatform(p string) error {

	//解压文件
	dir = "games/input_" + utils.ToString(time.Now().Unix()) + "/"
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
	fmt.Println(platformDom)

	platformId, err := platformDom.Add()
	if err != nil {
		return err
	}

	fmt.Println("============", platformId)

	/*simSection := file.Section("simulator")
	if simSection != nil{
		sims := simSection.ChildSections()
		for _,v := range sims{
			sim := &db.Simulator{}
			fmt.Println(v.Key("name"))
		}
	}*/

	return nil
}
