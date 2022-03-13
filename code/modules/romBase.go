package modules

import (
	"encoding/csv"
	"errors"
	"os"
	"simUI/code/config"
	"simUI/code/utils"
	"strings"
)

type RomBase struct {
	RomName   string // rom文件名
	Name      string // 中文名
	NameEN    string // 英文名
	NameJP    string // 日文名
	Type      string // 类型
	Year      string // 年份
	Producer  string // 制造商
	Publisher string // 出品公司
	Country   string // 国家
	Translate string // 汉化组
	Version   string // 版本
	OtherA    string // 其他内容a
	OtherB    string // 其他内容b
	OtherC    string // 其他内容c
	OtherD    string // 其他内容d

}

//读取游戏资料列表
func GetRomBaseList(platform uint32) (map[string]*RomBase, error) {

	if config.Cfg.Platform[platform].Rombase == "" {
		return map[string]*RomBase{}, nil
	}

	if !utils.FileExists(config.Cfg.Platform[platform].Rombase) {
		return map[string]*RomBase{}, nil
	}

	des := map[string]*RomBase{}
	records, err := utils.ReadCsv(config.Cfg.Platform[platform].Rombase)
	if err != nil {
		return map[string]*RomBase{}, nil //直接返回空，不返回错误
	}

	isUtf8 := false
	if len(records) > 0 {
		isUtf8 = utils.IsUTF8(records[0][0])
	} else {
		return nil, nil
	}

	for k, r := range records {

		if k == 0 || r[0] == "" {
			continue
		}

		create := []string{"", "", "", "", "", "", "", "", "", "", "", "", "", "", ""}
		createLen := len(create)
		i := 1
		for a, b := range r {
			if i > createLen {
				continue
			}
			i++
			create[a] = b
		}

		if isUtf8 == false {
			create[0] = utils.ToUTF8(create[0])
			create[1] = utils.ToUTF8(create[1])
			create[2] = utils.ToUTF8(create[2])
			create[3] = utils.ToUTF8(create[3])
			create[4] = utils.ToUTF8(create[4])
			create[5] = utils.ToUTF8(create[5])
			create[6] = utils.ToUTF8(create[6])
			create[7] = utils.ToUTF8(create[7])
			create[8] = utils.ToUTF8(create[8])
			create[9] = utils.ToUTF8(create[9])
			create[10] = utils.ToUTF8(create[10])
			create[11] = utils.ToUTF8(create[11])
			create[12] = utils.ToUTF8(create[12])
			create[13] = utils.ToUTF8(create[13])
			create[14] = utils.ToUTF8(create[14])
		}

		des[create[0]] = &RomBase{
			RomName:   create[0],
			Name:      strings.Trim(create[1], " "),
			Type:      strings.Trim(create[2], " "),
			Year:      strings.Trim(create[3], " "),
			Publisher: strings.Trim(create[4], " "),
			Country:   strings.Trim(create[5], " "),
			Translate: strings.Trim(create[6], " "),
			Version:   strings.Trim(create[7], " "),
			Producer:  strings.Trim(create[8], " "),
			NameEN:    strings.Trim(create[9], " "),
			NameJP:    strings.Trim(create[10], " "),
			OtherA:    strings.Trim(create[11], " "),
			OtherB:    strings.Trim(create[12], " "),
			OtherC:    strings.Trim(create[13], " "),
			OtherD:    strings.Trim(create[14], " "),
		}
	}

	return des, nil
}

//写csv文件
func WriteRomBaseFile(platform uint32, newData *RomBase) error {

	if config.Cfg.Platform[platform].Rombase == "" {
		return nil
	}

	info, _ := GetRomBaseList(platform) //读取老数据
	//如果全为空则删除当前记录
	if newData.Name == "" && newData.Publisher == "" && newData.Year == "" && newData.Type == "" && newData.Country == "" && newData.Translate == "" && newData.Version == "" {
		delete(info, newData.RomName)
	} else {
		info[newData.RomName] = newData //并入新数据
	}
	//转换为切片
	f, err := os.Create(config.Cfg.Platform[platform].Rombase)
	defer f.Close()
	if err != nil {
		return errors.New(config.Cfg.Lang["TipRombaseWriteError"])
	}

	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM，避免使用Microsoft Excel打开乱码

	writer := csv.NewWriter(f)

	//表头
	if err := writer.Write(getCsvTitle()); err != nil {
		return err
	}

	for _, v := range info {

		writer.Write([]string{
			v.RomName,
			strings.Trim(v.Name, " "),
			strings.Trim(v.Type, " "),
			strings.Trim(v.Year, " "),
			strings.Trim(v.Publisher, " "),
			strings.Trim(v.Country, " "),
			strings.Trim(v.Translate, " "),
			strings.Trim(v.Version, " "),
			strings.Trim(v.Producer, " "),
			strings.Trim(v.NameEN, " "),
			strings.Trim(v.NameJP, " "),
			strings.Trim(v.OtherA, " "),
			strings.Trim(v.OtherB, " "),
			strings.Trim(v.OtherC, " "),
			strings.Trim(v.OtherD, " "),
		})

	}
	writer.Flush() // 此时才会将缓冲区数据写入

	return nil
}

//覆盖csv文件
func CoverRomBaseFile(platform uint32, newData map[string]*RomBase) error {

	if config.Cfg.Platform[platform].Rombase == "" {
		return nil
	}

	//转换为切片
	f, err := os.Create(config.Cfg.Platform[platform].Rombase)
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM，避免使用Microsoft Excel打开乱码

	writer := csv.NewWriter(f)

	//表头
	writer.Write(getCsvTitle())

	for _, v := range newData {

		if v.Name == "" && v.Producer == "" && v.Publisher == "" && v.Year == "" && v.Type == "" && v.Country == "" && v.Translate == "" && v.Version == "" && v.NameEN == "" && v.NameJP == "" && v.OtherA == "" && v.OtherB == "" && v.OtherC == "" && v.OtherD == "" {
			continue
		}

		writer.Write([]string{
			v.RomName,
			strings.Trim(v.Name, " "),
			strings.Trim(v.Type, " "),
			strings.Trim(v.Year, " "),
			strings.Trim(v.Publisher, " "),
			strings.Trim(v.Country, " "),
			strings.Trim(v.Translate, " "),
			strings.Trim(v.Version, " "),
			strings.Trim(v.Producer, " "),
			strings.Trim(v.NameEN, " "),
			strings.Trim(v.NameJP, " "),
			strings.Trim(v.OtherA, " "),
			strings.Trim(v.OtherB, " "),
			strings.Trim(v.OtherC, " "),
			strings.Trim(v.OtherD, " "),
		})

	}
	writer.Flush() // 此时才会将缓冲区数据写入

	return nil
}

//创建一个新的csv文件
func CreateNewRomBaseFile(p string) error {

	//转换为切片
	f, err := os.Create(p)
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM，避免使用Microsoft Excel打开乱码

	writer := csv.NewWriter(f)

	//表头
	writer.Write(getCsvTitle())
	writer.Flush() // 此时才会将缓冲区数据写入
	return nil
}

//表头
func getCsvTitle() []string {
	return []string{
		config.Cfg.Lang["BaseRomName"],
		config.Cfg.Lang["BaseName"],
		config.Cfg.Lang["BaseType"],
		config.Cfg.Lang["BaseYear"],
		config.Cfg.Lang["BasePublisher"],
		config.Cfg.Lang["BaseCountry"],
		config.Cfg.Lang["BaseTranslate"],
		config.Cfg.Lang["BaseVersion"],
		config.Cfg.Lang["BaseProducer"],
		config.Cfg.Lang["BaseNameEN"],
		config.Cfg.Lang["BaseNameJP"],
		config.Cfg.Lang["BaseOther"] + "A",
		config.Cfg.Lang["BaseOther"] + "B",
		config.Cfg.Lang["BaseOther"] + "C",
		config.Cfg.Lang["BaseOther"] + "D",
	}
}
