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
	Platform  string // 平台
	Publisher string // 出品公司
}

var Baseinfo map[uint32]map[string]*RomBase

//读取详情文件
func GetRomBase(platform uint32) (map[string]*RomBase, error) {

	if config.Cfg.Platform[platform].Rombase == "" {
		return map[string]*RomBase{}, nil
	}

	if Baseinfo[platform] != nil {
		return Baseinfo[platform], nil
	}

	Baseinfo = map[uint32]map[string]*RomBase{}

	des := map[string]*RomBase{}
	records, err := utils.ReadCsv(config.Cfg.Platform[platform].Rombase)
	if err != nil {
		return nil, err
	}

	isUtf8 := false
	if len(records) > 0 {
		isUtf8 = utils.IsUTF8(records[0][0])
	} else {
		return Baseinfo[platform], nil
	}

	for k, r := range records {

		if k == 0 || r[0] == ""{
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

		des[r[0]] = &RomBase{
			RomName:   r[0],
			Name:      r[1],
			Type:      r[2],
			Year:      r[3],
			Platform:  r[4],
			Publisher: r[5],
		}
	}
	Baseinfo[platform] = des
	return des, nil
}

func WriteRomBaseFile(platform uint32, newData *RomBase) error {

	if config.Cfg.Platform[platform].Rombase == "" {
		return nil
	}

	info, _ := GetRomBase(platform)  //读取老数据

	//如果全为空则删除当前记录
	if newData.Name == "" && newData.Platform == "" && newData.Publisher == "" && newData.Year == "" && newData.Type == ""{
		delete(info,newData.RomName)
	}else{
		info[newData.RomName] = newData //并入新数据
	}
	//转换为切片
	create := [][]string{}

	//表头
	head := []string{"rom名称", "游戏名称", "游戏类型", "游戏平台", "发行年份", "出品公司"}
	create = append(create, head)

	for _, v := range info {
		d := []string{v.RomName, v.Name, v.Type, v.Platform, v.Year, v.Publisher}
		create = append(create, d)
	}

	if err := utils.WriteCsv(config.Cfg.Platform[platform].Rombase, create); err != nil {
		return err
	}

	return nil
}
