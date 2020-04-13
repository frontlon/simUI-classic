package db

import (
	"fmt"
	"time"
)

//读取所有数据
func CreateRomTemp() {
	//result := getDb().Exec("CREATE TEMPORARY TABLE rom_temp(md5 CHAR(50));")
	result := getDb().Exec(`CREATE TEMPORARY TABLE temp_rom(md5 CHAR(32));`)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
}

func AddRomTemp(platform uint32,md5s []*Rom) {




	sta := time.Now().Unix()
	//禁用同步模式
	getDb().Exec("PRAGMA synchronous = OFF;")

	fmt.Println("准备执行")


	st := []string{}
	//sql := "BEGIN;"
	for _,v := range md5s{
		st = append(st,v.Md5)
		//sql += "INSERT INTO temp_rom (md5) VALUES ('"+v.Md5+"');"
	}

	//sql  += "COMMIT;"
	sta2 := time.Now().Unix()
	end := sta2 - sta
	fmt.Println("数据组装完成，执行时间：",end)



	volist :=[]*Rom{}
	result := getDb().Table("rom").Select("md5").Where("md5 not in (?)",st).Find(&volist)

	fmt.Println("错误信息",result.Error)
	sta3 := time.Now().Unix()
	end = sta3 - sta2

	fmt.Println("写库完成，执行时间",end)


}
