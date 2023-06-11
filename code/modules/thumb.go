package modules

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"simUI/code/components"
	"simUI/code/config"
	"simUI/code/db"
	"simUI/code/utils"
	"strings"
)

var ConstThumbSubSeparator = "__" //新资源子ID
var ConstThumbSubName = "_j5D"    //新资源子ID

type DownThumbs struct {
	Width  int
	Height int
	Ext    string
	ImgUrl string
	Type   string
}

// 编辑展示图片
func EditRomThumbs(typeName string, id uint64, sid string, picPath string, ext string) (string, string, error) {

	rom := &db.Rom{
		Id: id,
	}

	//设定新的文件名
	vo, err := rom.GetById(id)
	if err != nil {
		return "", "", err
	}

	//下载成功后，备份原文件
	res := config.GetResPath(vo.Platform)
	resPath := res[typeName] //当前资源目录

	if resPath == "" {
		return "", "", errors.New(config.Cfg.Lang["NoSetThumbDir"])
	}

	FileExt := ext
	if FileExt == "" {
		FileExt = utils.GetFileExt(picPath)
		if FileExt == "" {
			FileExt = ".jpg"
		}
	}

	//生成新文件
	platformPathAbs, err := filepath.Abs(resPath) //读取平台图片路径
	RomFileName := utils.GetFileName(vo.RomPath)

	newFileName := ""
	newFileNamePrefix := platformPathAbs + config.Cfg.Separator + RomFileName
	subName := ""
	if sid == "" {
		fmt.Println("aaaaaaaaaaaaaaaaaaa")
		//主rom
		subName = ""
		newFileName = newFileNamePrefix + FileExt //生成新文件的完整绝路路径地址
	} else if sid == ConstThumbSubName {
		fmt.Println("bbbbbbbb")
		//新图片
		for i := 1; i <= 999; i++ {
			subName = ConstThumbSubSeparator + utils.ToString(i)
			newFileName = newFileNamePrefix + subName + FileExt //生成新文件的完整绝路路径地址
			if !utils.FileExists(newFileName) {
				fmt.Println("===============")
				sid = utils.ToString(i)
				fmt.Println(sid)

				break
			}
		}
	} else {
		fmt.Println("cccccccccccccc")

		//子图片
		subName = ConstThumbSubSeparator + sid
		newFileName = newFileNamePrefix + subName + FileExt //生成新文件的完整绝路路径地址
	}

	//备份老图片
	if err := components.BackupOldPic(platformPathAbs, RomFileName+subName); err != nil {
		return "", "", err
	}

	//复制文件
	if ext != "" {
		//如果是网络下载
		if err := components.DownloadRomThumbs(picPath, newFileName); err != nil {
			return "", "", err
		}
	} else {
		//如果是本地图片
		if err := utils.FileCopy(picPath, newFileName); err != nil {
			return "", "", errors.New(config.Cfg.Lang["ResFileNotFound"])
		}
	}

	return newFileName, sid, nil
}

// 删除缩略图
func DeleteThumbs(typeName string, id uint64, sid string) error {

	vo, err := (&db.Rom{Id: id}).GetById(id)
	if err != nil {
		return err
	}

	res := config.GetResPath(vo.Platform)
	platformPath := res[typeName]
	if platformPath == "" {
		return errors.New(config.Cfg.Lang["NoSetThumbDir"])
	}

	fileName := utils.GetFileName(vo.RomPath)
	if sid != "" {
		fileName = fileName + ConstThumbSubSeparator + sid
	}

	//移动文件到备份文件夹
	if err := components.BackupOldPic(platformPath, fileName); err != nil {
		return err
	}

	//删除图片文件
	if err := components.DeleteResPic(platformPath, fileName); err != nil {
		return err
	}

	return nil
}

// 设置一张展示图为主图
func SetMasterThumbs(typeName string, id uint64, sid string) error {

	vo, err := (&db.Rom{Id: id}).GetById(id)
	if err != nil {
		return err
	}

	res := config.GetResPath(vo.Platform)
	platformPath := res[typeName]
	if platformPath == "" {
		return errors.New(config.Cfg.Lang["NoSetThumbDir"])
	}

	masterName := utils.GetFileName(vo.RomPath)
	tempName := utils.GetFileName(vo.RomPath) + ConstThumbSubSeparator + ConstThumbSubName
	slaveName := masterName + ConstThumbSubSeparator + sid

	//1.把master改为tmp，防止覆盖丢失
	for _, ext := range config.PIC_EXTS {
		masterPath := platformPath + config.Cfg.Separator + masterName + ext
		tempPath := platformPath + config.Cfg.Separator + tempName + ext
		if utils.FileExists(masterPath) {
			if err := utils.FileMove(masterPath, tempPath); err != nil {
				return err
			}
		}
	}

	//2.把slave改为master
	for _, ext := range config.PIC_EXTS {
		slavePath := platformPath + config.Cfg.Separator + slaveName + ext
		masterPath := platformPath + config.Cfg.Separator + masterName + ext
		if utils.FileExists(slavePath) {
			if err := utils.FileMove(slavePath, masterPath); err != nil {
				return err
			}
		}
	}

	//3.把temp改为slave
	for _, ext := range config.PIC_EXTS {
		tempPath := platformPath + config.Cfg.Separator + tempName + ext
		slavePath := platformPath + config.Cfg.Separator + slaveName + ext
		if utils.FileExists(tempPath) {
			if err := utils.FileMove(tempPath, slavePath); err != nil {
				return err
			}
		}
	}

	return nil
}

// 百度图片搜索
func SearchThumbsForBaidu(keyword string, page int) (map[string][]DownThumbs, error) {
	size := 30
	postUrl := config.Cfg.Default.SearchEngines
	num := page * size
	keyword = url.QueryEscape(keyword)
	postUrl = strings.ReplaceAll(postUrl, "{$keyword}", keyword)
	postUrl = strings.Replace(postUrl, "{$NumIndex}", utils.ToString(num), 1)
	postUrl = strings.Replace(postUrl, "{$pageNum}", utils.ToString(size), 1)

	//整理http请求体
	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, postUrl, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// 添加请求头
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	respMap := make(map[string]interface{})
	bstr := string(b)
	bstr = strings.ReplaceAll(bstr, `\'`, "'")
	err = json.Unmarshal([]byte(bstr), &respMap)

	//被识别为蜘蛛，被拦截
	if _, ok := respMap["antiFlag"]; ok {
		return nil, errors.New(respMap["message"].(string))
	}

	//请求成功，开始组装消息体
	respList := map[string][]DownThumbs{}
	list := []DownThumbs{}
	if _, ok := respMap["data"]; ok {
		for _, v := range respMap["data"].([]interface{}) {
			vo := v.(map[string]interface{})

			if _, ok = vo["thumbURL"]; !ok {
				continue
			}

			width := 0
			height := 0
			ext := ""
			if _, ok = vo["width"]; ok {
				width = utils.ToInt(vo["width"].(float64))
			}
			if _, ok = vo["height"]; ok {
				height = utils.ToInt(vo["height"].(float64))
			}
			if _, ok = vo["type"]; ok {
				ext = vo["type"].(string)
			}

			stu := DownThumbs{
				Width:  width,
				Height: height,
				Ext:    ext,
				ImgUrl: vo["thumbURL"].(string),
			}
			list = append(list, stu)
		}
	}

	if len(list) > 0 {
		respList["百度图片搜索结果"] = list
	}

	return respList, nil
}

// hfsDB图片搜索
func SearchThumbsForHfsDb(keyword string, page int) (map[string][]DownThumbs, error) {

	if !utils.IsEnglish(keyword) {
		keyword = utils.BaiduTranslate(keyword, "en")
	}
	limit := 10
	offset := page * limit
	body, err := components.GetHfsDbGameList(keyword, limit, offset)
	if err != nil {
		fmt.Println(err)
	}

	respList := map[string][]DownThumbs{}

	for _, v := range body.Results {
		medias := []DownThumbs{}

		if len(v.Medias) == 0 {
			continue
		}
		for _, m := range v.Medias {
			if m.IsImage == true {
				md := DownThumbs{
					Width:  m.ResX,
					Height: m.ResY,
					Ext:    m.Extension,
					ImgUrl: m.File,
					Type:   m.Type,
				}
				medias = append(medias, md)
			}
		}
		respList[v.NameEn] = medias
	}
	return respList, nil
}
