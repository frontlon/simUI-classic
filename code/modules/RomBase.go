package modules

import (
	"encoding/csv"
	"os"
	"simUI/code/config"
	"simUI/code/utils"
	"strings"
)

type RomBase struct {
	RomName   string // rom文件名
	Name      string // 游戏名称
	Type      string // 类型
	Year      string // 年份
	Publisher string // 出品公司
	Country   string // 国家
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

		if k == 0 || r[0] == "" {
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
			RomName:   strings.Trim(r[0], " "),
			Name:      strings.Trim(r[1], " "),
			Type:      strings.Trim(r[2], " "),
			Year:      strings.Trim(r[3], " "),
			Publisher: strings.Trim(r[4], " "),
			Country:   strings.Trim(r[5], " "),
		}
	}

	Baseinfo[platform] = des
	return des, nil
}

func WriteRomBaseFile(platform uint32, newData *RomBase) error {

	if config.Cfg.Platform[platform].Rombase == "" {
		return nil
	}

	info, _ := GetRomBase(platform) //读取老数据
	//如果全为空则删除当前记录
	if newData.Name == "" && newData.Publisher == "" && newData.Year == "" && newData.Type == "" && newData.Country == "" {
		delete(info, newData.RomName)
	} else {
		info[newData.RomName] = newData //并入新数据
	}
	//转换为切片
	f, err := os.Create(config.Cfg.Platform[platform].Rombase)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM，避免使用Microsoft Excel打开乱码

	writer := csv.NewWriter(f)

	//表头
	writer.Write([]string{config.Cfg.Lang["BaseRomName"], config.Cfg.Lang["BaseName"], config.Cfg.Lang["BaseType"], config.Cfg.Lang["BaseYear"], config.Cfg.Lang["BasePublisher"], config.Cfg.Lang["BaseCountry"]})

	for _, v := range info {

		writer.Write([]string{
			strings.Trim(v.RomName, " "),
			strings.Trim(v.Name, " "),
			strings.Trim(v.Type, " "),
			strings.Trim(v.Year, " "),
			strings.Trim(v.Publisher, " "),
			strings.Trim(v.Country, " "),
		})

	}
	writer.Flush() // 此时才会将缓冲区数据写入

	return nil
}
