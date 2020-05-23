package modules

import (
	"simUI/code/config"
	"simUI/code/utils"
)

type RomBase struct {
	RomName   string //rom文件名
	EnName    string //英文名
	CnName    string //中文名
	Type      string //类型
	Platform  string //平台
	Year      string //年份
	Developer string //开发商
	Publisher string //发行商
}

var Baseinfo map[string]*RomBase

//读取详情文件
func GetRomBase(platform uint32) (map[string]*RomBase, error) {
	if config.Cfg.Platform[platform].Rombase == "" {
		return map[string]*RomBase{}, nil
	}

	if Baseinfo != nil {
		return Baseinfo, nil
	}

	records, err := utils.ReadCsv(config.Cfg.Platform[platform].Rombase)
	if err != nil {
		return nil, err
	}
	for _, r := range records {
		Baseinfo[r[0]] = &RomBase{
			RomName:   r[0],
			EnName:    r[1],
			CnName:    r[2],
			Type:      r[3],
			Platform:  r[4],
			Year:      r[5],
			Developer: r[6],
			Publisher: r[7],
		}
	}
	delete(Baseinfo, "rom名称") //删除第一列

	return Baseinfo, nil
}

func WriteRomBaseFile(platform uint32, newData *RomBase) error {

	if config.Cfg.Platform[platform].Rombase == "" {
		return nil
	}
	
	Baseinfo, _ = GetRomBase(platform)  //读取老数据
	Baseinfo[newData.RomName] = newData //并入新数据

	//转换为切片
	create := [][]string{}

	//表头
	head := []string{"rom名称", "英文名", "中文名", "游戏类型", "平台", "年份", "开发商", "发行商"}
	create = append(create, head)

	for _, v := range Baseinfo {
		d := []string{v.RomName, v.EnName, v.CnName, v.Type, v.Platform, v.Year, v.Developer, v.Publisher}
		create = append(create, d)
	}

	if err := utils.WriteCsv(config.Cfg.Platform[platform].Rombase, create); err != nil {
		return err
	}

	return nil
}
