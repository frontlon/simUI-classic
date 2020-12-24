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
	Translate string // 汉化组
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
		return nil, nil //直接返回空，不返回错误
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

		create := []string{"","","","","","",""}

		for a,b := range r{
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
		}

		des[create[0]] = &RomBase{
			RomName:   strings.Trim(create[0], " "),
			Name:      strings.Trim(create[1], " "),
			Type:      strings.Trim(create[2], " "),
			Year:      strings.Trim(create[3], " "),
			Publisher: strings.Trim(create[4], " "),
			Country:   strings.Trim(create[5], " "),
			Translate: strings.Trim(create[6], " "),
		}
	}

	Baseinfo[platform] = des
	return des, nil
}

//写csv文件
func WriteRomBaseFile(platform uint32, newData *RomBase) error {

	if config.Cfg.Platform[platform].Rombase == "" {
		return nil
	}

	info, _ := GetRomBase(platform) //读取老数据
	//如果全为空则删除当前记录
	if newData.Name == "" && newData.Publisher == "" && newData.Year == "" && newData.Type == "" && newData.Country == "" && newData.Translate == "" {
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
	writer.Write(getCsvTitle())

	for _, v := range info {

		writer.Write([]string{
			strings.Trim(v.RomName, " "),
			strings.Trim(v.Name, " "),
			strings.Trim(v.Type, " "),
			strings.Trim(v.Year, " "),
			strings.Trim(v.Publisher, " "),
			strings.Trim(v.Country, " "),
			strings.Trim(v.Translate, " "),
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
		panic(err)
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
	}
}
