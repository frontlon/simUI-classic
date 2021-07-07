package modules

import (
	"archive/zip"
	"encoding/json"
	"github.com/go-ini/ini"
	"io"
	"os"
	"simUI/code/config"
	"simUI/code/db"
	"simUI/code/utils"
	"strings"
)

var outputCfg = ini.Empty()
var outputConfigFile = "./cache/config.ini"

func OutputPlatform(platformId uint32, p string) error {
	setIniPlatform((platformId))
	setIniSimulator(platformId)
	setIniRomSet(platformId)
	err := outputCfg.SaveTo(outputConfigFile);
	defer utils.FileDelete(outputConfigFile)
	if err != nil {
		return err
	}
	packFiles(platformId, p)
	return nil
}

//设置ini - 平台信息
func setIniPlatform(platformId uint32) {
	platform, _ := (&db.Platform{}).GetById(platformId)
	section := "platform"

	paths := config.GetResPath(platformId)

	if (platform.RomPath != "") {
		platform.RomPath = strings.ReplaceAll(platform.RomPath, "\\", "/")
	}

	outputCfg.Section(section).Key("name").SetValue(platform.Name)
	outputCfg.Section(section).Key("icon").SetValue(utils.GetFileNameAndExt(platform.IconPath))
	outputCfg.Section(section).Key("exts").SetValue(platform.RomExts)
	outputCfg.Section(section).Key("rombase").SetValue(utils.GetFileNameAndExt(platform.Rombase))
	outputCfg.Section(section).Key("rom").SetValue(platform.RomPath)
	outputCfg.Section(section).Key("ico").SetValue(utils.GetFileNameAndExt(platform.Icon))

	for k, v := range paths {
		if v != "" {
			v = strings.ReplaceAll(v, "\\", "/")
			varr := strings.Split(v, "/")
			v = varr[len(varr)-1]
		}
		outputCfg.Section(section).Key(k).SetValue(v)
	}
}

//设置ini - 模拟器信息
func setIniSimulator(platformId uint32) {

	sims, _ := (&db.Simulator{}).GetByPlatform(platformId)
	for k, v := range sims {

		if v.Path != "" {
			p := strings.ReplaceAll(v.Path, "\\", "/")
			varr := strings.Split(p, "/")
			if (len(varr) > 1) {
				v.Path = varr[len(varr)-2] + config.Cfg.Separator + varr[len(varr)-1]
			} else {
				v.Path = varr[len(varr)-1]
			}
		}

		outputCfg.Section("simulator." + utils.ToString(k)).Key("name").SetValue(v.Name)
		outputCfg.Section("simulator." + utils.ToString(k)).Key("path").SetValue(v.Path)
		outputCfg.Section("simulator." + utils.ToString(k)).Key("cmd").SetValue(v.Cmd)
		outputCfg.Section("simulator." + utils.ToString(k)).Key("default").SetValue(utils.ToString(v.Default))
		outputCfg.Section("simulator." + utils.ToString(k)).Key("unzip").SetValue(utils.ToString(v.Unzip))
	}
}

//设置ini - rom信息
func setIniRomSet(platformId uint32) {
	roms, _ := (&db.Rom{}).GetByPlatform(platformId)
	stars := []string{}
	hides := []string{}
	sims := map[string]string{}

	for _, v := range roms {
		if v.Star == 1 {
			stars = append(stars, v.Name)
		}
		if v.Hide == 1 {
			hides = append(hides, v.Name)
		}
		if v.SimConf != "{}" && v.SimConf != "" {
			sims[v.Name] = v.SimConf
		}
	}

	starJson, _ := json.Marshal(stars)
	hideJson, _ := json.Marshal(hides)
	simJson, _ := json.Marshal(sims)
	starStr := utils.Base64Encode(string(starJson))
	hideStr := utils.Base64Encode(string(hideJson))
	simStr := utils.Base64Encode(string(simJson))

	outputCfg.Section("rom").Key("star").SetValue(starStr)
	outputCfg.Section("rom").Key("hide").SetValue(hideStr)
	outputCfg.Section("rom").Key("sim").SetValue(simStr)
}

//打包文件
func packFiles(platformId uint32, p string) {
	platform, _ := (&db.Platform{}).GetById(platformId)
	files := map[string]*os.File{}

	//压缩资源
	files["rombase"], _ = os.Open(platform.Rombase)
	files["ico"], _ = os.Open(platform.Icon)
	res := config.GetResPath(platformId)
	for k, v := range res {
		files[k], _ = os.Open(v)
	}

	compreFile, _ := os.Create("./1/" + platform.Name + ".zip")
	zw := zip.NewWriter(compreFile)
	defer zw.Close()

	for k, file := range files {
		if k == "rombase" || k == "ico" {
			k = ""
		}
		err := compress_zip(file, k, zw)
		if err != nil {
			//return err
		}
		defer file.Close()
	}

	//模拟器文件
	sims, _ := (&db.Simulator{}).GetByPlatform(platformId)
	simList := []*os.File{}
	for _, v := range sims {

		a, _ := os.Open(utils.GetFileAbsPath(v.Path))
		simList = append(simList, a)
	}
	for _, file := range simList {

		err := compress_zip(file, "simulator/", zw)
		if err != nil {
			//return err
		}
		defer file.Close()
	}

	//ini配置文件
	c, _ := os.Open(outputConfigFile)
	_ = compress_zip(c, "", zw)

	defer c.Close()
}

func compress_zip(file *os.File, prefix string, zw *zip.Writer) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		if strings.Contains(prefix, "/") {
			prefix = prefix + "/" + info.Name()
		} else {
			prefix = info.Name()
		}
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			f, err := os.Open(file.Name() + "/" + fi.Name())
			if err != nil {
				return err
			}
			err = compress_zip(f, prefix, zw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := zip.FileInfoHeader(info)
		header.Name = prefix + "/" + header.Name
		if err != nil {
			return err
		}
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
