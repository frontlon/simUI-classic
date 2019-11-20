package main

import (
	"VirtualNesGUI/code/db"
	"VirtualNesGUI/code/utils"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

var constSeparator = "__"                                    //rom子分隔符
var DOC_EXTS = []string{".txt", ".md", ".html", ".htm"}      //doc文档支持的扩展名
var PIC_EXTS = []string{"png", "jpg", "gif", "jpeg", "bmp"}; //支持的图片类型

type RomDetail struct {
	Info            *db.Rom
	DocContent      string
	StrategyContent string
	Sublist         []*db.Rom
}

/**
 * 读取游戏介绍文本
 **/
func getDocContent(f string) string {
	if f == "" {
		return ""
	}
	text, err := ioutil.ReadFile(f)
	content := ""
	if err != nil {
		return content
	}
	//enc := mahonia.NewDecoder("gbk")
	//content = enc.ConvertString(string(text))
	content = string(text)
	return content
}

/**
 * 运行游戏
 **/
func runGame(exeFile string, cmd []string) error {

	//更改程序运行目录
	if err := os.Chdir(filepath.Dir(exeFile)); err != nil {
		return err
	}

	result := &exec.Cmd{}

	//这个写法牛不牛逼~但有更好的吗？有的话请告诉我。
	switch len(cmd) {
	case 0:
		result = exec.Command(exeFile)
	case 1:
		result = exec.Command(exeFile, cmd[0])
	case 2:
		result = exec.Command(exeFile, cmd[0], cmd[1])
	case 3:
		result = exec.Command(exeFile, cmd[0], cmd[1], cmd[2])
	case 4:
		result = exec.Command(exeFile, cmd[0], cmd[1], cmd[2], cmd[3])
	case 5:
		result = exec.Command(exeFile, cmd[0], cmd[1], cmd[2], cmd[3], cmd[4])
	case 6:
		result = exec.Command(exeFile, cmd[0], cmd[1], cmd[2], cmd[3], cmd[4], cmd[5])
	case 7:
		result = exec.Command(exeFile, cmd[0], cmd[1], cmd[2], cmd[3], cmd[4], cmd[5], cmd[6])
	case 8:
		result = exec.Command(exeFile, cmd[0], cmd[1], cmd[2], cmd[3], cmd[4], cmd[5], cmd[6], cmd[7])
	case 9:
		result = exec.Command(exeFile, cmd[0], cmd[1], cmd[2], cmd[3], cmd[4], cmd[5], cmd[6], cmd[7], cmd[8])
	case 10:
		result = exec.Command(exeFile, cmd[0], cmd[1], cmd[2], cmd[3], cmd[4], cmd[5], cmd[6], cmd[7], cmd[8], cmd[9])
	case 11:
		result = exec.Command(exeFile, cmd[0], cmd[1], cmd[2], cmd[3], cmd[4], cmd[5], cmd[6], cmd[7], cmd[8], cmd[9], cmd[10])
	case 12:
		result = exec.Command(exeFile, cmd[0], cmd[1], cmd[2], cmd[3], cmd[4], cmd[5], cmd[6], cmd[7], cmd[8], cmd[9], cmd[10], cmd[11])
	case 13:
		result = exec.Command(exeFile, cmd[0], cmd[1], cmd[2], cmd[3], cmd[4], cmd[5], cmd[6], cmd[7], cmd[8], cmd[9], cmd[10], cmd[11], cmd[12])
	case 14:
		result = exec.Command(exeFile, cmd[0], cmd[1], cmd[2], cmd[3], cmd[4], cmd[5], cmd[6], cmd[7], cmd[8], cmd[9], cmd[10], cmd[11], cmd[12], cmd[13])
	case 15:
		result = exec.Command(exeFile, cmd[0], cmd[1], cmd[2], cmd[3], cmd[4], cmd[5], cmd[6], cmd[7], cmd[8], cmd[9], cmd[10], cmd[11], cmd[12], cmd[13], cmd[14])
	case 16:
		result = exec.Command(exeFile, cmd[0], cmd[1], cmd[2], cmd[3], cmd[4], cmd[5], cmd[6], cmd[7], cmd[8], cmd[9], cmd[10], cmd[11], cmd[12], cmd[13], cmd[14], cmd[15])
	case 17:
		result = exec.Command(exeFile, cmd[0], cmd[1], cmd[2], cmd[3], cmd[4], cmd[5], cmd[6], cmd[7], cmd[8], cmd[9], cmd[10], cmd[11], cmd[12], cmd[13], cmd[14], cmd[15], cmd[16])
	case 18:
		result = exec.Command(exeFile, cmd[0], cmd[1], cmd[2], cmd[3], cmd[4], cmd[5], cmd[6], cmd[7], cmd[8], cmd[9], cmd[10], cmd[11], cmd[12], cmd[13], cmd[14], cmd[15], cmd[16], cmd[17])
	case 19:
		result = exec.Command(exeFile, cmd[0], cmd[1], cmd[2], cmd[3], cmd[4], cmd[5], cmd[6], cmd[7], cmd[8], cmd[9], cmd[10], cmd[11], cmd[12], cmd[13], cmd[14], cmd[15], cmd[16], cmd[17], cmd[18])
	case 20:
		result = exec.Command(exeFile, cmd[0], cmd[1], cmd[2], cmd[3], cmd[4], cmd[5], cmd[6], cmd[7], cmd[8], cmd[9], cmd[10], cmd[11], cmd[12], cmd[13], cmd[14], cmd[15], cmd[16], cmd[17], cmd[18], cmd[19])
	}

	if err := result.Start(); err != nil {
		return err
	}

	return nil
}

/**
 * 创建缓存
 **/
func CreateRomCache(platform uint32) ([]*db.Rom,[]string,error) {
	romlist := []*db.Rom{}
	uniqs := []string{}               //rom名称列表，用户清理rom表使用
	menuList := map[string]*db.Menu{}            //分类目录
	RomPath := Config.Platform[platform].RomPath //rom文件路径
	RomExt := Config.Platform[platform].RomExts  //rom扩展名
	RomAlias, _ := getRomAlias(platform)         //别名配置

	//进入循环，遍历文件
	if err := filepath.Walk(RomPath,
		func(p string, f os.FileInfo, err error) (error) {
			if f == nil {
				return err
			}

			//整理目录格式，并转换为数组
			newpath := strings.Replace(RomPath, "/", "\\", -1)
			newpath = strings.Replace(p, newpath, "", -1)
			if len(newpath) > 0 && newpath[0:1] == "\\" {
				newpath = strings.Replace(newpath, "\\", "", 1)
			}
			subpath := strings.Split(newpath, "\\")
			romExt := strings.ToLower(path.Ext(p)) //获取文件后缀

			//如果该文件是游戏rom
			if f.IsDir() == false && CheckPlatformExt(RomExt, romExt) {
				romName := utils.GetFileName(f.Name())
				title := romName
				//如果有别名配置，则读取别名
				if _, ok := RomAlias[title]; ok {
					title = RomAlias[title]
				}

				py := TextToPinyin(title)
				md5 := GetFileUniqId(f)
				//如果游戏名称存在分隔符，说明是子游戏
				menu := constMenuRootKey //无目录，读取默认参数
				//定义目录，如果有子目录，则记录子目录名称
				if len(subpath) > 1 {
					menu = subpath[0]
				}

				//如果游戏名称存在分隔符，说明是子游戏
				if strings.Contains(title, constSeparator) {

					//拆分文件名
					sub := strings.Split(title, constSeparator)

					//去掉扩展名，生成标题
					sinfo := &db.Rom{
						Name:     sub[1],
						Pname:    sub[0],
						RomPath:  p,
						Menu:     menu,
						Platform: platform,
						Star:     0,
						Pinyin:   py,
						Md5:      md5,
					}
					romlist = append(romlist, sinfo)
					uniqs = append(uniqs, md5) //游戏md5列表，用于删除不存在的rom
				} else { //不是子游戏
					//去掉扩展名，生成标题
					rinfo := &db.Rom{
						Menu:     menu,
						Name:     title,
						Platform: platform,
						RomPath:  p,
						Star:     0,
						Pinyin:   py,
						Md5:      md5,
					}

					//rom列表
					romlist = append(romlist, rinfo)
					//分类列表
					if menu != constMenuRootKey {
						menuList[menu] = &db.Menu{
							Platform: platform,
							Name:     menu,
							Pinyin:   TextToPinyin(menu),
						}
					}
					uniqs = append(uniqs, md5) //游戏md5列表，用于删除不存在的rom
				}
			}
			return nil
		}); err != nil {
	}





		/*
	menus := []string{}

	//删除当前平台下，不存在的rom
	if err := (&db.Rom{}).DeleteNotExists(platform, uniqs); err != nil {
	}

	//删除当前平台下不存在的菜单
	if err := (&db.Menu{}).DeleteNotExists(platform, menus); err != nil {
	}


	issetMd5, err :=  (&db.Rom{}).GetMd5ByMd5(platform,uniqs)
	if err != nil{
		return err
	}

	//取出需要写入数据库的rom数据。
	saveRomlist := []*db.Rom{}
	for _,v := range romlist{
		if utils.InSliceString(v.Md5,issetMd5) == false{
			saveRomlist = append(saveRomlist,v)
		}
	}

	//保存新数据到数据库rom表
	if len(saveRomlist) > 0 {
		for _, v := range saveRomlist {
			if err := v.Add(); err != nil {
				fmt.Println(err.Error())
			}
		}
	}

	//保存数据到数据库cate表
	if len(menuList) > 0 {
		for _, v := range menuList {
			if err := v.Add(); err != nil {
			}
		}

	}

	//这些变量较大，写入完成后清理变量
	romlist = []*db.Rom{}
	saveRomlist = []*db.Rom{}
	menuList = make(map[string]*db.Menu)
*/


	return romlist,uniqs,nil
}

//检查文件扩展名是否存在于该平台中
func CheckPlatformExt(exts []string, file string) bool {
	for _, v := range exts {
		if v == file {
			return true
		}
	}
	return false
}
