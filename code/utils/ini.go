package utils

import (
)

func SetConfig(platform string,pname string)  error {
	/*
	controller.Config.Platform[platform].Romlist

	langpath := Config.RootPath + "lang" + Config.Separator
	fpath := langpath + lang + ".ini"
	section := make(map[string]string)

	//如果默认语言不存在，则读取列表中的其他语言
	if !utils.FileExists(fpath) {
		if len(Config.LangList) > 0 {
			for langName, langFile := range Config.LangList {
				fpath = langpath + langFile
				//如果找到其他语言，则将第一项更新到数据库配置中
				if err := (&db.Config{}).UpdateField("lang", langName); err != nil {
					return section, err
				}
				break
			}
		}
	}

	file, err := ini.LoadSources(ini.LoadOptions{IgnoreInlineComment: true}, fpath)

	if err != nil {
		return section, err
	}

	section = file.Section("").KeysHash()
	return section, nil*/
}