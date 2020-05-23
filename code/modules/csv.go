package modules

import (
	"simUI/code/utils"
)

type CsvStruct struct {
	RomName   string //rom文件名
	EnName    string //英文名
	CnName    string //中文名
	Type      string //类型
	Platform  string //平台
	Year      string //年份
	Developer string //开发商
	Publisher string //发行商
}

//读取详情文件
func ReadDescFile(platform uint32) (map[string]*CsvStruct, error) {
	csvData := map[string]*CsvStruct{}

	filename := "1.csv"
	records, err := utils.ReadCsv(filename)
	if err != nil {
		return nil, err
	}
	for _, r := range records {
		csvData[r[0]] = &CsvStruct{
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
	delete(csvData, "rom名称") //删除第一列

	return csvData, nil
}

func WriteDescFile(platform uint32, newData *CsvStruct) error {
	filename := "1.csv"

	data, _ := ReadDescFile(platform) //读取老数据
	data[newData.RomName] = newData   //并入新数据

	//转换为切片
	create := [][]string{}

	//表头
	head := []string{"rom名称", "英文名", "中文名", "游戏类型", "平台", "年份", "开发商", "发行商"}
	create = append(create, head)

	for _, v := range data {
		d := []string{v.RomName, v.EnName, v.CnName, v.Type, v.Platform, v.Year, v.Developer, v.Publisher}
		create = append(create, d)
	}

	if err := utils.WriteCsv(filename, create); err != nil {
		return err
	}

	return nil
}
