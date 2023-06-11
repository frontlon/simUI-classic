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
	Score     string // 评分
}

var RomBaseList map[uint32]map[string]*RomBase

//读取游戏资料列表
func GetRomBaseList(platform uint32) (map[string]*RomBase, error) {

	//如果已经读取，则直接返回
	if len(RomBaseList[platform]) > 0 {
		return RomBaseList[platform], nil
	}

	//开始整理数据
	RomBaseList = map[uint32]map[string]*RomBase{}
	RomBaseList[platform] = map[string]*RomBase{}

	if config.Cfg.Platform[platform].Rombase == "" {
		return RomBaseList[platform], nil
	}

	if !utils.FileExists(config.Cfg.Platform[platform].Rombase) {
		return RomBaseList[platform], nil
	}

	records, err := utils.ReadCsv(config.Cfg.Platform[platform].Rombase)
	if err != nil {
		return RomBaseList[platform], nil //直接返回空，不返回错误
	}

	isUtf8 := false
	if len(records) > 0 {
		isUtf8 = utils.IsUTF8(records[0][0])
	} else {
		return RomBaseList[platform], nil
	}

	creates := [][]string{}
	for k, r := range records {

		if k == 0 {
			continue
		}

		create := []string{"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""}
		createLen := len(create)
		i := 1
		for a, b := range r {
			if i > createLen {
				continue
			}
			i++
			create[a] = b
		}

		//转换成utf-8编码
		if isUtf8 == false {
			for ck, cv := range create {
				create[ck] = utils.ToUTF8(cv)
			}
		}
		creates = append(creates, create)
	}

	//老版本csv没有Id，需要转换
	/*if len(creates) > 1 && len(records[0]) < 16 && utils.ToInt(creates[0][0]) == 0 {
		creates, _ = upgradeRomBase(platform, creates)
	}*/

	for _, create := range creates {
		if create[0] == "" {
			continue
		}
		lists := &RomBase{}
		if len(create) >= 1 {
			lists.RomName = strings.Trim(create[0], " ")
		}
		if len(create) >= 2 {
			lists.Name = strings.Trim(create[1], " ")
		}
		if len(create) >= 3 {
			lists.Type = strings.Trim(create[2], " ")
		}
		if len(create) >= 4 {
			lists.Year = strings.Trim(create[3], " ")
		}
		if len(create) >= 5 {
			lists.Publisher = strings.Trim(create[4], " ")
		}
		if len(create) >= 6 {
			lists.Country = strings.Trim(create[5], " ")
		}
		if len(create) >= 7 {
			lists.Translate = strings.Trim(create[6], " ")
		}
		if len(create) >= 8 {
			lists.Version = strings.Trim(create[7], " ")
		}
		if len(create) >= 9 {
			lists.Producer = strings.Trim(create[8], " ")
		}
		if len(create) >= 10 {
			lists.NameEN = strings.Trim(create[9], " ")
		}
		if len(create) >= 11 {
			lists.NameJP = strings.Trim(create[10], " ")
		}
		if len(create) >= 12 {
			lists.OtherA = strings.Trim(create[11], " ")
		}
		if len(create) >= 13 {
			lists.OtherB = strings.Trim(create[12], " ")
		}
		if len(create) >= 14 {
			lists.OtherC = strings.Trim(create[13], " ")
		}
		if len(create) >= 15 {
			lists.OtherD = strings.Trim(create[14], " ")
		}
		if len(create) >= 16 {
			lists.Score = strings.Trim(create[15], " ")
		}
		RomBaseList[platform][create[0]] = lists
	}
	return RomBaseList[platform], nil
}

//读取一个rom的资料信息
func GetRomBaseById(platform uint32, id string) *RomBase {
	romlist, _ := GetRomBaseList(platform)
	if _, ok := romlist[id]; ok {
		return RomBaseList[platform][id]
	}
	return nil
}

//写csv文件
func WriteRomBaseFile(platform uint32, newData *RomBase) error {

	if config.Cfg.Platform[platform].Rombase == "" {
		return nil
	}

	info, _ := GetRomBaseList(platform) //读取老数据
	//如果全为空则删除当前记录
	info[newData.RomName] = newData //并入新数据

	//写入csv文件
	if err := WriteDataToFile(config.Cfg.Platform[platform].Rombase, info); err != nil {
		return err
	}

	//更新全局变量
	RomBaseList[platform][newData.RomName] = newData

	return nil
}

//覆盖csv文件
func CoverRomBaseFile(platform uint32, newData map[string]*RomBase) error {

	if config.Cfg.Platform[platform].Rombase == "" {
		return nil
	}

	//写入csv文件
	if err := WriteDataToFile(config.Cfg.Platform[platform].Rombase, newData); err != nil {
		return err
	}

	//更新全局变量
	RomBaseList[platform] = map[string]*RomBase{}
	RomBaseList[platform] = newData

	return nil
}

//写入csv文件
func WriteDataToFile(filePath string, data map[string]*RomBase) error {

	if filePath == "" {
		return nil
	}

	//转换为切片
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM，避免使用Microsoft Excel打开乱码

	writer := csv.NewWriter(f)

	//表头
	writer.Write(getCsvTitle())

	for _, v := range data {

		str := strings.Join([]string{v.Name, v.Producer, v.Publisher, v.Year, v.Type, v.Country, v.Translate, v.Version, v.NameEN, v.NameJP, v.OtherA, v.OtherB, v.OtherC, v.OtherD, v.Score}, "")

		if str == "" {
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
			strings.Trim(v.Score, " "),
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

//csv 升级
/*func upgradeRomBase(platform uint32, records [][]string) ([][]string, error) {
	fmt.Println("进入版本转换，platformId：", platform)

	if len(records) < 2 {
		return [][]string{}, nil
	}

	//遍历文件，读取文件名和file_md5
	RomExt := strings.Split(config.Cfg.Platform[platform].RomExts, ",") //rom扩展名
	RomExtMap := map[string]bool{}
	for _, v := range RomExt {
		RomExtMap[v] = true
	}
	romsMap := map[string]string{}
	if err := filepath.Walk(config.Cfg.Platform[platform].RomPath,
		func(p string, f os.FileInfo, err error) error {
			romExists := false                     //rom是否存在
			romExt := strings.ToLower(path.Ext(p)) //获取文件后缀
			if _, ok := RomExtMap[romExt]; ok {
				romExists = true //rom存在
			}
			if f.IsDir() == false && romExists == true {
				romsMap[utils.GetFileName(p)] = utils.CreateRomUniqId(f.ModTime().UnixNano())
			}
			return nil
		}); err != nil {
		return records, nil
	}

	//加入id字段
	for k, record := range records {
		id := ""
		if _, ok := romsMap[record[0]]; ok {
			id = romsMap[record[0]]
		}
		records[k] = append([]string{id}, record...)
	}

	//开始写入文件
	f, err := os.Create(config.Cfg.Platform[platform].Rombase)
	if err != nil {
		return records, err
	}
	defer f.Close()

	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM，避免使用Microsoft Excel打开乱码

	writer := csv.NewWriter(f)

	writer.Write(getCsvTitle())

	for _, r := range records {
		writer.Write(r)
	}
	writer.Flush() // 此时才会将缓冲区数据写入

	return records, nil
}*/

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
		config.Cfg.Lang["Score"],
	}
}
