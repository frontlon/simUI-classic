package modules

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-ini/ini"
	"simUI/code/components"
	"simUI/code/config"
	"simUI/code/db"
	"simUI/code/utils"
	"simUI/code/utils/go-sciter"
	"strings"
)

type MergeDb struct {
	Platform    *db.Platform
	RomCount    int64
	Simulators  []*db.Simulator
	FolderCheck map[string]MergeDbFolderCheck
}

type MergeDbFolderCheck struct {
	Status string //检测结果 suc正常 warn警告 err错误
	Desc   string
	Path   string //目标路径
}

/**
 * 导出rom配置
 */
func BackupRomConfig(p string, platform uint32) error {

	go func() {
		subRom := map[uint32]map[string]string{}
		romSetting := map[uint32]map[string]*db.RomSetting{}

		platformList := []uint32{}
		mode := ""
		if platform == 0 {
			mode = "multi" //多平台导出
			for platformId, _ := range config.Cfg.Platform {
				platformList = append(platformList, platformId)
			}
		} else {
			mode = "single" //单平台导出
			platformList = append(platformList, platform)
		}

		for _, platformId := range platformList {

			//读取rom数据
			romList, _ := (&db.Rom{}).GetByPlatform(platformId)

			romMap := map[string]string{}
			for _, v := range romList {
				romMap[v.FileMd5] = utils.GetFileName(v.RomPath)
			}

			//读取子游戏关系
			subList, _ := (&db.RomSubGame{}).GetByPlatform(platformId)
			subMap := map[string]string{}
			for _, v := range subList {

				if _, ok := romMap[v.FileMd5]; !ok {
					continue
				}
				if _, ok := romMap[v.Pname]; !ok {
					continue
				}

				subMap[utils.Base64Encode(romMap[v.FileMd5])] = utils.Base64Encode(romMap[v.Pname])
			}

			//读取rom配置
			settings, _ := (&db.RomSetting{}).GetByPlatform(platformId)
			settingsMap := map[string]*db.RomSetting{}
			for _, v := range settings {
				if _, ok := romMap[v.FileMd5]; !ok {
					continue
				}
				settingsMap[utils.Base64Encode(romMap[v.FileMd5])] = v
			}

			savePlatform := platformId
			if mode == "single" {
				savePlatform = 0
			}
			subRom[savePlatform] = subMap
			romSetting[savePlatform] = settingsMap
		}

		components.WriteRomConfigToIni(p, subRom, romSetting)

		if _, err := utils.Window.Call("CB_romConfigBackup"); err != nil {
		}

		fmt.Println("导出完成")
	}()
	return nil
}

/**
 * 导入rom配置
 */
func RestoreRomConfig(p string, platform uint32) error {

	f, err := ini.Load(p)
	if err != nil {
		return err
	}

	go func() {

		cfg := f.Sections()

		subRom := map[uint32]map[string]string{}
		romSetting := map[uint32][]map[string]string{}
		//解析数据
		for _, section := range cfg {
			sectionName := strings.Split(section.Name(), ".")
			if len(sectionName) < 2 {
				continue
			}
			savePlatform := platform
			if platform == 0 {
				//导入全部
				savePlatform = uint32(utils.ToInt(sectionName[1]))
			}

			romSetting[savePlatform] = []map[string]string{}
			if strings.Contains(section.Name(), "subGame") {
				//子游戏
				subMap := map[string]string{}
				for _, v := range section.Keys() {
					subMap[utils.Base64Decode(v.Name())] = utils.Base64Decode(v.Value())
				}
				subRom[savePlatform] = subMap
			} else if strings.Contains(section.Name(), "setting") {
				//rom配置
				data := map[string]string{}
				data["name"] = utils.Base64Decode(sectionName[2])
				for _, v := range section.Keys() {
					data[v.Name()] = v.Value()
				}
				romSetting[savePlatform] = append(romSetting[savePlatform], data)
			}
		}

		//开始写入rom_subgame数据
		components.SaveRomConfigSubRom(subRom)

		//开始写入rom_config数据
		components.SaveRomConfigSetting(romSetting)

		//清空rom表
		(&db.Rom{}).Truncate()

		if _, err := utils.Window.Call("CB_romConfigRestore"); err != nil {
		}

		fmt.Println("导入完成")
	}()
	return nil

}

// 合并数据库 - 读取数据
func GetMergeDbData(dbFile string) ([]MergeDb, error) {

	cuDbFile, _ := db.GetDbFileName()
	dbFile = strings.ReplaceAll(dbFile, "/", "\\")
	cuDbFile = strings.ReplaceAll(cuDbFile, "/", "\\")
	if dbFile == cuDbFile {
		return nil, errors.New(config.Cfg.Lang["SoftMergeMsgInputSelf"])
	}

	dbPath := utils.GetFilePath(dbFile)
	engine, err := db.CustomDBConn(dbFile)
	if err != nil {
		return nil, err
	}
	defer engine.Close()

	result := []MergeDb{}

	//读取平台
	platforms, err := engine.GetAllPlatform()
	if err != nil {
		return nil, err
	}

	if len(platforms) == 0 {
		return result, nil
	}

	//读取平台rom数
	romCountMap, err := engine.GetRomCount()
	//读取平台模拟器
	simulatorMap, err := engine.GetAllSimulator()

	for _, v := range platforms {
		var romCount int64 = 0
		simulators := []*db.Simulator{}

		if _, ok := romCountMap[v.Id]; ok {
			romCount = romCountMap[v.Id]
		}

		if _, ok := simulatorMap[v.Id]; ok {
			simulators = simulatorMap[v.Id]
		}

		resPath := config.GetPlatformResPath(v)
		folderCheck := map[string]MergeDbFolderCheck{}
		for k, pth := range resPath {
			srcRel := strings.Replace(pth, dbPath, "", 1) //相对路径
			src := dbPath + pth                           //绝对路径
			dst := config.Cfg.RootPath + srcRel           //绝对路径

			if srcRel == "" {
				//目录未设置
				folderCheck[k] = MergeDbFolderCheck{
					Status: "err",
					Desc:   config.Cfg.Lang["SoftMergeMsgSrcNotSet"],
					Path:   pth,
				}
			} else if k == "Rombase" && utils.FileExists(dst) {
				//资料文件检查
				folderCheck[k] = MergeDbFolderCheck{
					Status: "err",
					Desc:   config.Cfg.Lang["SoftMergeMsgDstExists"],
					Path:   pth,
				}
			} else if k == "Rombase" && utils.FileExists(src) {
				//资料文件检查
				folderCheck[k] = MergeDbFolderCheck{
					Status: "suc",
					Desc:   config.Cfg.Lang["SoftMergeMsgOK"],
					Path:   pth,
				}
			} else if utils.IsAbsPath(srcRel) {
				//绝对路径不导入
				folderCheck[k] = MergeDbFolderCheck{
					Status: "warn",
					Desc:   config.Cfg.Lang["SoftMergeMsgAbsPath"],
					Path:   pth,
				}
			} else if utils.IsDirEmpty(src) {
				//源目录为空，不导入
				folderCheck[k] = MergeDbFolderCheck{
					Status: "warn",
					Desc:   config.Cfg.Lang["SoftMergeMsgFIleNotExists"],
					Path:   pth,
				}
			} else if !utils.IsDirEmpty(dst) {
				//目标目录不为空，不导入
				folderCheck[k] = MergeDbFolderCheck{
					Status: "err",
					Desc:   config.Cfg.Lang["SoftMergeMsgDstNotEmpty"],
					Path:   pth,
				}
			} else {
				//正常导入
				folderCheck[k] = MergeDbFolderCheck{
					Status: "suc",
					Desc:   config.Cfg.Lang["SoftMergeMsgOK"],
					Path:   pth,
				}
			}
		}
		
		r := MergeDb{
			Platform:    v,
			RomCount:    romCount,
			Simulators:  simulators,
			FolderCheck: folderCheck,
		}
		result = append(result, r)
	}

	return result, nil
}

// 合并数据库
func MergeDB(dbFile string, platformIds []uint32, simulatorIds []string) error {

	//数据更新完成后，页面回调，更新页面DOM
	if _, err := utils.Window.Call("CB_romMergeStart", sciter.NewValue("")); err != nil {
		fmt.Println(err)
	}

	dbPath := utils.GetFilePath(dbFile)
	engine, err := db.CustomDBConn(dbFile)
	if err != nil {
		return err
	}
	defer func() {
		engine.Close()
		//数据合并完成后，页面回调，更新页面DOM
		if _, err := utils.Window.Call("CB_romMergeEnd", sciter.NewValue("")); err != nil {
			fmt.Println(err)
		}
	}()

	//读取平台
	platforms, err := engine.GetPlatformByIds(platformIds)
	if err != nil {
		return err
	}

	simulatorMap := map[uint32][]*db.Simulator{}
	simulators, _ := engine.GetAllSimulator()
	if simulators != nil {
		simulatorMap = simulators
	}
	oldSimNameMap := map[uint32]string{}
	for _, v := range simulatorMap {
		for _, b := range v {
			oldSimNameMap[b.Id] = b.Name
		}
	}

	for _, v := range platforms {
		//新建平台
		oldPlatformId := v.Id
		v.Id = 0
		v.Id, err = v.Add()
		if err != nil {
			return err
		}

		//拷贝资源
		resPath := config.GetPlatformResPath(v)
		for k, pth := range resPath {
			utils.Loading(config.Cfg.Lang["StartMergeRes"]+k, v.Name)
			utils.FolderCopy(components.GetSrcAndDstPath(dbPath, pth, k))
		}
		//拷贝模拟器
		for _, s := range simulatorMap[oldPlatformId] {
			utils.Loading(config.Cfg.Lang["StartMergeSimulator"]+s.Name, v.Name)
			utils.FolderCopy(components.GetSrcAndDstPath(dbPath, s.Path, "simulator"))
		}
		//添加模拟器
		if _, ok := simulatorMap[oldPlatformId]; ok {
			createSimulators := []*db.Simulator{}
			for k, _ := range simulatorMap[oldPlatformId] {
				if !utils.InSliceString(utils.ToString(simulatorMap[oldPlatformId][k].Id), simulatorIds) {
					continue
				}
				simulatorMap[oldPlatformId][k].Id = 0
				simulatorMap[oldPlatformId][k].Platform = v.Id
				createSimulators = append(createSimulators, simulatorMap[oldPlatformId][k])
			}
			(&db.Simulator{}).BatchAdd(createSimulators)
		}

		//读取新模拟器数据
		newSimulator, _ := (&db.Simulator{}).GetByPlatform(v.Id)
		newSimNameMap := map[string]uint32{}
		if newSimulator != nil {
			for _, b := range newSimulator {
				newSimNameMap[b.Name] = b.Id
			}
		}

		//添加rom setting
		utils.Loading(config.Cfg.Lang["StartMergeRomConfig"], v.Name)

		romSettingList, _ := engine.GetRomSettingByPlatform(oldPlatformId)
		if romSettingList != nil && len(romSettingList) > 0 {
			for k, _ := range romSettingList {
				romSettingList[k].Id = 0
				romSettingList[k].Platform = v.Id
				//更新模拟器独立配置
				if romSettingList[k].SimConf != "" && romSettingList[k].SimConf != "{}" {
					newSim := map[uint32]*db.SimConf{}
					oldSim := map[uint32]*db.SimConf{}
					_ = json.Unmarshal([]byte(romSettingList[k].SimConf), &oldSim)
					for sid, d := range oldSim {
						if _, ok := oldSimNameMap[sid]; ok {
							if _, ok = newSimNameMap[oldSimNameMap[sid]]; ok {
								newSim[newSimNameMap[oldSimNameMap[sid]]] = d
							}
						}
					}
					newSimByte, _ := json.Marshal(newSim)
					romSettingList[k].SimConf = string(newSimByte)
				}
			}
			(&db.RomSetting{}).BatchAdd(romSettingList)
		}

		//添加subgame
		subGameList, _ := engine.GetSubGameByPlatform(oldPlatformId)
		if romSettingList != nil && len(subGameList) > 0 {
			for k, _ := range subGameList {
				subGameList[k].Id = 0
				subGameList[k].Platform = v.Id
			}
			(&db.RomSubGame{}).BatchAdd(subGameList)
		}

	}

	return nil
}
