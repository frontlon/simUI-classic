package components

import (
	"github.com/go-ini/ini"
	"simUI/code/config"
	"simUI/code/utils"
	"strings"
)

type Slnk struct {
	Type string `ini:"type"`
	Path string `ini:"path"`
	Cmd  string `ini:"cmd"`
}

//写slnk文件
func SaveSlnkFile(slnkPath, gameType, romPath, cmd string) error {
	cfg := ini.Empty()
	defaultSection := cfg.Section("")
	defaultSection.NewKey("type", gameType)
	defaultSection.NewKey("path", romPath)
	defaultSection.NewKey("cmd", cmd)
	err := cfg.SaveTo(slnkPath)
	if err != nil {
		return err
	}
	return nil
}

//读slnk文件
func GetSlnkFile(f string) (string, []string) {

	if !utils.FileExists(f) {
		return "", nil
	}

	//解析ini
	cfg, err := ini.Load(f)
	if err != nil {
		return "", nil
	}
	c := Slnk{}
	cfg.MapTo(&c)

	//解析真实路径
	p := c.Path
	if p != "" {
		p = strings.Trim(p, "\r\n")
		p = strings.Trim(p, "\r")
		p = strings.Trim(p, "\n")
		if utils.IsAbsPath(p) == false {
			p = config.Cfg.RootPath + p
		}
	}

	//解析启动参数
	args := []string{}
	if c.Cmd != "" {
		cmds := strings.Split(c.Cmd, " ")
		for _, v := range cmds {
			if strings.Trim(v, " ") == "" {
				continue
			}
			args = append(args, v)
		}
	}
	return p, args
}
