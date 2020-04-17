package controller

import (
	"VirtualNesGUI/code/db"
	"VirtualNesGUI/code/utils"
	"encoding/json"
	"fmt"
	"github.com/go-ini/ini"
)

//备份配置
func backupConfig(p string) error {
	if p == "" {
		return nil
	}
	iniCfg := ini.Empty()

	config, _ := (&db.Config{}).Get()
	confJson, _ := json.Marshal(config)
	confEnc := utils.Base64Encode(string(confJson))
	if _, err := iniCfg.Section("").NewKey("config", confEnc); err != nil {
		return err
	}

	platform, _ := (&db.Platform{}).GetAll()
	platformJson, _ := json.Marshal(platform)
	platformEnc := utils.Base64Encode(string(platformJson))
	if _, err := iniCfg.Section("").NewKey("platform", platformEnc); err != nil {
		return err
	}

	shortcut, _ := (&db.Shortcut{}).GetAll()
	shortcutJson, _ := json.Marshal(shortcut)
	shortcutEnc := utils.Base64Encode(string(shortcutJson))
	if _, err := iniCfg.Section("").NewKey("shortcut", shortcutEnc); err != nil {
		return err
	}

	simulator, _ := (&db.Simulator{}).GetAll()
	simulatorJson, _ := json.Marshal(simulator)
	simulatorEnc := utils.Base64Encode(string(simulatorJson))
	if _, err := iniCfg.Section("").NewKey("simulator", simulatorEnc); err != nil {
		return err
	}

	if err := iniCfg.SaveTo(p); err != nil {
		return err
	}

	return nil

}

//还原配置

func restoreConfig(p string) error {
	if p == "" {
		return nil
	}

	//创建数据
	file, err := ini.LoadSources(ini.LoadOptions{IgnoreInlineComment: true}, p)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	section := file.Section("Alias").KeysHash()

	//清空表

	if err := (&db.Rom{}).Truncate(); err != nil {
		return err
	}

	if err := (&db.Menu{}).Truncate(); err != nil {
		return err
	}

	if err := (&db.Config{}).Truncate(); err != nil {
		return err
	}
	if err := (&db.Platform{}).Truncate(); err != nil {
		return err
	}
	if err := (&db.Shortcut{}).Truncate(); err != nil {
		return err
	}
	if err := (&db.Simulator{}).Truncate(); err != nil {
		return err
	}

	if section["config"] != "" {
		configDb := &db.Config{}
		confDec := utils.Base64Decode(section["config"])
		if err := json.Unmarshal([]byte(confDec), &configDb); err != nil {
			return err
		}
		//复写数据
		configDb.Add()
	}

	if section["platform"] != "" {
		platform := []*db.Platform{}
		platformDec := utils.Base64Decode(section["platform"])
		if err := json.Unmarshal([]byte(platformDec), &platform); err != nil {
			return err
		}
		(&db.Platform{}).BatchAdd(platform)

	}

	if section["shortcut"] == "" {
		shortcut := []*db.Shortcut{}
		shortcutDec := utils.Base64Decode(section["shortcut"])
		if err := json.Unmarshal([]byte(shortcutDec), &shortcut); err != nil {
			return err
		}
		(&db.Shortcut{}).BatchAdd(shortcut)

	}

	if section["simulator"] == "" {
		simulator := []*db.Simulator{}
		simulatorDec := utils.Base64Decode(section["simulator"])
		if err := json.Unmarshal([]byte(simulatorDec), &simulator); err != nil {
			return err
		}
		(&db.Simulator{}).BatchAdd(simulator)
	}

	return nil

}
