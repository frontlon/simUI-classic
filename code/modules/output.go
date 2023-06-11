package modules

import (
	"archive/zip"
	"io/ioutil"
	"os"
	"path/filepath"
	"simUI/code/components"
	"simUI/code/config"
	"simUI/code/db"
	"simUI/code/request"
	"simUI/code/utils"
	"strings"
)

/**
 * 导出(分享)rom
**/
func OutputRom(request *request.OutputRom) error {

	if request.Opt == "platform" {
		//导出全部平台
		outputByPlatform(request)
	} else {
		//导出指定rom
		roms := []*db.Rom{}
		hasSubgame := utils.InSliceString("subgame", request.Options)
		if request.Opt == "menu" {
			//导出指定目录
			if hasSubgame == true {
				roms, _ = (&db.Rom{}).GetMasterAndSubByMenus(request.Platform, request.Menus)
			} else {
				roms, _ = (&db.Rom{}).GetMasterByMenus(request.Platform, request.Menus)

			}
		} else {
			//导出指定rom
			if hasSubgame == true {
				roms, _ = (&db.Rom{}).GetMasterAndSubByMasterIds(request.Roms)
			} else {
				roms, _ = (&db.Rom{}).GetByIds(request.Roms)
			}
		}
		outputById(request, roms)
	}
	return nil
}

/**
 * 导出整个平台
 **/
func outputByPlatform(request *request.OutputRom) {
	res := config.GetResPath(request.Platform)
	res["roms"] = config.Cfg.Platform[request.Platform].RomPath

	compreFile, _ := os.Create(request.Save)
	zw := zip.NewWriter(compreFile)
	defer zw.Close()

	//压缩rom和资源
	for _, opt := range request.Options {
		resPath := res[opt]
		if resPath == "" {
			continue
		}

		files, err := ioutil.ReadDir(resPath)
		if err != nil {
			return
		}

		for _, p := range files {
			openF, err := os.Open(resPath + config.Cfg.Separator + p.Name())
			defer openF.Close()
			if err == nil {
				_ = components.CompressZip(openF, opt, zw)
			}
		}
	}

	//压缩模拟器
	if len(request.Simulators) > 0 {
		for _, simId := range request.Simulators {
			if _, ok := config.Cfg.Platform[request.Platform].SimList[simId]; !ok {
				continue
			}

			//把资源路径变为map
			pathMap := map[string]bool{}
			for _, v := range res {
				pathMap[v+config.Cfg.Separator] = true
			}

			sim := config.Cfg.Platform[request.Platform].SimList[simId]

			if err := filepath.Walk(utils.GetFilePath(sim.Path),
				func(p string, f os.FileInfo, err error) error {

					//资源目录跳过
					pPath := utils.GetFilePath(p)
					if _, ok := pathMap[pPath]; ok {
						return nil
					}

					//开始压缩文件
					fo, err := os.Open(p)
					defer func() {
						fo.Close()
						utils.FileDelete(p)
					}()
					if err != nil {
						return err
					}
					romOpt := "simulator"
					pPathArr := strings.Split(pPath, config.Cfg.Separator)
					if len(pPathArr) > 0 {
						romOpt = romOpt + config.Cfg.Separator + pPathArr[len(pPathArr)-1]
					}

					_ = components.CompressZip(fo, romOpt, zw)

					return nil
				}); err != nil {

			}
		}
	}

	//压缩资料文件
	if config.Cfg.Platform[request.Platform].Rombase != "" {
		tmp := "./cache/rombase.csv"
		utils.FileCopy(config.Cfg.Platform[request.Platform].Rombase, tmp)
		f, err := os.Open(tmp)
		defer func() {
			f.Close()
			utils.FileDelete(tmp)
		}()
		if err == nil {
			_ = components.CompressZip(f, "", zw)
		}
	}

	//生成ini文件
	if utils.InSliceString("subgame", request.Options) {
		subgames, _ := (&db.RomSubGame{}).GetByPlatformToMap(request.Platform)
		if err := components.CreateOutputIni(subgames); err != nil {
			return
		}
		f, err := os.Open(components.TmpOutputIniPath)
		defer func() {
			f.Close()
			utils.FileDelete(components.TmpOutputIniPath)
		}()
		if err == nil {
			_ = components.CompressZip(f, "", zw)
		}
	}

}

/**
 * 根据id导出rom
 **/
func outputById(request *request.OutputRom, roms []*db.Rom) {

	//zipPath string, platform uint32, options []string

	res := config.GetResPath(request.Platform)
	res["roms"] = config.Cfg.Platform[request.Platform].RomPath

	//生成rom_name map
	romMap := map[string]bool{}
	for _, v := range roms {
		if v.Pname != "" {
			continue
		}
		name := utils.GetFileName(v.RomPath)
		romMap[name] = true
	}

	//生成压缩包文件
	compreFile, _ := os.Create(request.Save)
	zw := zip.NewWriter(compreFile)
	defer zw.Close()

	//遍历资源类型，开始压缩
	for _, opt := range request.Options {

		resPath := res[opt]
		if resPath == "" {
			continue
		}

		//读取所有资源
		resMap, _ := components.GetAllRomThumbsToMap(resPath, romMap)

		for _, v := range roms {
			//压缩资源文件
			fileName := utils.GetFileName(v.RomPath)
			if _, ok := resMap[fileName]; ok {
				for _, sp := range resMap[fileName] {
					f, err := os.Open(sp)
					defer f.Close()
					if err == nil {
						romOpt := opt
						relPath := strings.Replace(sp, res[opt]+config.Cfg.Separator, "", 1)
						relPath = utils.GetFilePath(relPath)
						if relPath != "" {
							romOpt = romOpt + config.Cfg.Separator + relPath
						}
						_ = components.CompressZip(f, romOpt, zw)
					}
				}
			}
		}
	}

	//压缩csv资料文件
	oldRomBase, err := GetRomBaseList(request.Platform)
	newRomBase := map[string]*RomBase{}
	if err != nil {
		return
	}

	for romName, _ := range romMap {
		if _, ok := oldRomBase[romName]; ok {
			newRomBase[romName] = oldRomBase[romName]
		}
	}
	//写入csv文件
	csvFile := "./cache/rombase.csv"
	WriteDataToFile(csvFile, newRomBase)
	csvf, err := os.Open(csvFile)
	defer func() {
		csvf.Close()
		utils.FileDelete(csvFile)
	}()
	if err == nil {
		_ = components.CompressZip(csvf, "", zw)
	}

	//生成ini文件
	if utils.InSliceString("subgame", request.Options) {
		subgames := map[string]string{}
		for _, v := range roms {
			if v.Pname != "" {
				subgames[v.FileMd5] = v.Pname
			}
		}
		if err := components.CreateOutputIni(subgames); err != nil {
			return
		}
		f, err := os.Open(components.TmpOutputIniPath)
		defer func() {
			f.Close()
			utils.FileDelete(components.TmpOutputIniPath)
		}()
		if err == nil {
			_ = components.CompressZip(f, "", zw)
		}
	}

}

// 压缩资料文件
func compressRombase() {}
