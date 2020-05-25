package modules

import (
	"simUI/code/config"
	"simUI/code/utils"
)

type RomBase struct {
	RomName   string // rom文件名
	Name      string // 游戏名称
	Type      string // 类型
	Year      string // 年份
	Developer string // 开发商
	Publisher string // 发行商
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

	Baseinfo = map[string]*RomBase{}

	isUtf8 := false
	if len(records) > 0 {
		isUtf8 = utils.IsUTF8(records[0][0])
	} else {
		return Baseinfo, nil
	}

	for k, r := range records {

		if k == 0 {
			continue
		}

		if isUtf8 == false {
			r[0] = utils.ToUTF8(r[0])
			r[1] = utils.ToUTF8(r[1])
			r[2] = utils.ToUTF8(r[2])
			r[3] = utils.ToUTF8(r[3])
			r[4] = utils.ToUTF8(r[4])
			r[5] = utils.ToUTF8(r[5])
		}

		Baseinfo[r[0]] = &RomBase{
			RomName:   r[0],
			Name:      r[1],
			Type:      r[2],
			Year:      r[3],
			Developer: r[4],
			Publisher: r[5],
		}

	}

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
	head := []string{"rom名称", "名称", "游戏类型", "年份", "开发商", "发行商"}
	create = append(create, head)

	for _, v := range Baseinfo {
		d := []string{v.RomName, v.Name, v.Type, v.Year, v.Developer, v.Publisher}
		create = append(create, d)
	}

	if err := utils.WriteCsv(config.Cfg.Platform[platform].Rombase, create); err != nil {
		return err
	}

	return nil
}
