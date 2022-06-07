package modules

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"simUI/code/compoments"
	"simUI/code/config"
	"simUI/code/db"
	"simUI/code/utils"
	"strings"
	"time"
)

var ConstSeparator = "__"     //rom子分隔符
var ConstMenuRootKey = "_7b9" //根子目录游戏的Menu参数

type RomDetail struct {
	Info          *db.Rom             //rom信息
	DocContent    string              //简介内容
	StrategyFiles []map[string]string //攻略文件列表
	AudioList     []map[string]string //音频文件列表
	Sublist       []*db.Rom           //子游戏
	Simlist       []*db.Simulator     //模拟器
}

//右键打开文件夹
func OpenFolder(id uint64, opt string, simId uint32) error {

	info, err := (&db.Rom{}).GetById(id)
	platform := config.Cfg.Platform[info.Platform] //读取当前平台信息
	if err != nil {
		return err
	}
	romName := utils.GetFileName(filepath.Base(info.RomPath)) //读取文件名
	fileName := ""
	switch opt {
	case "rom":
		fileName = platform.RomPath + config.Cfg.Separator + info.RomPath
	case "sim":
		if _, ok := platform.SimList[simId]; ok {
			fileName = platform.SimList[simId].Path
		}
	default:
		res := config.GetResPath(platform.Id)

		if res[opt] != "" {
			fileName = GetRomRes(opt, info.Platform, romName)
			if fileName == "" {
				fileName = res[opt]
			}
		}
	}
	if err := compoments.OpenFolderByWindow(fileName); err != nil {
		return err
	}
	return nil
}

//打开rom存储目录
func OpenRomPathFolder(platform uint32) error {

	fmt.Println(config.Cfg.Platform)
	fmt.Println(platform)

	folder := config.Cfg.Platform[platform].RomPath
	if folder == "" {
		return errors.New(config.Cfg.Lang["RomMenuCanNotBeExists"])
	}

	fmt.Println(folder)
	if err := compoments.OpenFolderByWindow(folder); err != nil {
		return err
	}
	return nil
}

//添加pc游戏
func AddPcGame(platform uint32, menu string, files []string) error {
	folder := config.Cfg.Platform[platform].RomPath + config.Cfg.Separator
	if menu != "" {
		folder = folder + menu + config.Cfg.Separator
	}

	contentTemp := `chcp 65001` + "\r\n"
	contentTemp += `cd /d %~dp0` + "\r\n"
	contentTemp += `cd "{PATH}"` + "\r\n"
	contentTemp += `start "" "{FILENAME}"` + "\r\n"

	for _, f := range files {
		p := folder + utils.GetFileName(f) + ".bat"
		f = strings.Replace(f, `file://`, "", 1)
		f = strings.Replace(f, `file:\\`, "", 1)
		//f = utils.Utf8ToGbk(f)
		filename := utils.GetFileNameAndExt(f)
		relPath := utils.GetRelPathByTowPath(p, f)
		content := strings.ReplaceAll(contentTemp, "{PATH}", relPath)
		content = strings.ReplaceAll(content, "{FILENAME}", filename)

		//如果文件已存在，则使用新名称
		if utils.FileExists(p) {
			path := utils.GetFilePath(p)
			name := utils.GetFileName(p) + utils.ToString(time.Now().Unix())
			ext := utils.GetFileExt(p)
			p = path + config.Cfg.Separator + name + ext
		}

		utils.OverlayWriteFile(p, content)
	}

	return nil
}

//读取rom详情
func GetGameDetail(id uint64) (*RomDetail, error) {

	res := &RomDetail{}
	//游戏游戏详细数据
	info, err := (&db.Rom{}).GetById(id)

	if err != nil {
		return res, err
	}

	//子游戏列表
	sub, _ := (&db.Rom{}).GetSubRom(info.Platform, info.FileMd5)

	res.Info = info
	res.Sublist = sub
	res.Simlist, _ = (&db.Simulator{}).GetByPlatform(info.Platform)

	for k, v := range res.Simlist {
		if res.Simlist[k].Path != "" {
			res.Simlist[k].Path, _ = filepath.Abs(v.Path)
		}
	}

	//读取文档内容
	romBaseName := utils.GetFileName(filepath.Base(info.RomPath)) //生成新文件的完整绝路路径地址
	if config.Cfg.Platform[info.Platform].DocPath != "" {
		docFileName := ""
		for _, v := range config.DOC_EXTS {
			docFileName = config.Cfg.Platform[info.Platform].DocPath + config.Cfg.Separator + romBaseName + v
			res.DocContent = GetDocContent(docFileName)
			if res.DocContent != "" {
				break
			}
		}
	}

	//攻略文件列表
	res.StrategyFiles, _ = GetStrategyFile(id)

	//音频文件列表
	res.AudioList, _ = GetAudioList(id)

	return res, nil
}

//读取游戏攻略内容
func GetGameDoc(t string, id uint64) (string, error) {

	//游戏游戏详细数据
	info, err := (&db.Rom{}).GetById(id)

	if err != nil {
		return "", err
	}

	//如果没有执行运行的文件，则读取文档内容
	romName := utils.GetFileName(filepath.Base(info.RomPath)) //生成新文件的完整绝路路径地址
	strategy := ""
	for _, v := range config.DOC_EXTS {
		strategyFileName := ""
		if t == "strategy" {
			strategyFileName = config.Cfg.Platform[info.Platform].StrategyPath + config.Cfg.Separator + romName + v
		} else if t == "doc" {
			strategyFileName = config.Cfg.Platform[info.Platform].DocPath + config.Cfg.Separator + romName + v
		}
		strategy = GetDocContent(strategyFileName)
		if strategy != "" {
			break
		}
	}

	strategy = strings.ReplaceAll(strategy, `<img src="`, `<img src="`+config.Cfg.RootPath)
	strategy = strings.Trim(strategy, "\t")
	strategy = strings.Trim(strategy, "\n\r")
	strategy = strings.Trim(strategy, "\r")
	strategy = strings.Trim(strategy, "\n")

	return strategy, nil
}

//更新游戏攻略内容
func SetGameDoc(t string, id uint64, content string) error {

	//游戏游戏详细数据
	info, err := (&db.Rom{}).GetById(id)

	if err != nil {
		return err
	}

	//如果没有执行运行的文件，则读取文档内容
	romName := utils.GetFileName(filepath.Base(info.RomPath)) //生成新文件的完整绝路路径地址
	newExt := ""
	Filename := ""
	for _, v := range config.DOC_EXTS {
		strategyFileName := ""
		if t == "strategy" {
			strategyFileName = config.Cfg.Platform[info.Platform].StrategyPath + config.Cfg.Separator + romName + v
			newExt = config.Cfg.Platform[info.Platform].StrategyPath + config.Cfg.Separator + romName + ".txt"
		} else if t == "doc" {
			strategyFileName = config.Cfg.Platform[info.Platform].DocPath + config.Cfg.Separator + romName + v
			newExt = config.Cfg.Platform[info.Platform].DocPath + config.Cfg.Separator + romName + ".txt"
		}

		if utils.FileExists(strategyFileName) {
			Filename = strategyFileName
			break
		}
	}

	if Filename == "" {
		Filename = newExt
	}

	if !utils.FileExists(Filename) {
		if err := utils.CreateFile(Filename); err != nil {
			return err
		}
	}

	if !utils.IsUTF8(content) {
		content = utils.ToUTF8(content)
	}

	//替换图片路径为相对路径
	content = strings.ReplaceAll(content, config.Cfg.RootPath, "")
	if err := utils.OverlayWriteFile(Filename, content); err != nil {
		return err
	}

	return nil
}

//删除游戏攻略内容
func DelGameDoc(t string, id uint64) error {

	//游戏游戏详细数据
	info, err := (&db.Rom{}).GetById(id)

	if err != nil {
		return err
	}

	//如果没有执行运行的文件，则读取文档内容
	romName := utils.GetFileName(filepath.Base(info.RomPath)) //生成新文件的完整绝路路径地址
	res := config.GetResPath(info.Platform)
	Filename := ""
	for _, v := range config.DOC_EXTS {
		strategyFileName := res[t] + config.Cfg.Separator + romName + v

		if utils.FileExists(strategyFileName) {
			Filename = strategyFileName
			break
		}
	}

	if Filename != "" && utils.FileExists(Filename) {
		if err := utils.FileDelete(Filename); err != nil {
			return err
		}
	}

	return nil
}

/**
 * 读取游戏介绍文本
 **/
func GetDocContent(f string) string {
	if f == "" {
		return ""
	}
	text, err := ioutil.ReadFile(f)
	content := ""
	if err != nil {
		return content
	}
	content = string(text)

	if !utils.IsUTF8(content) {
		content = utils.ToUTF8(content)
	}

	content = strings.Trim(content, "\t")
	content = strings.Trim(content, "\n\r")
	content = strings.Trim(content, "\r")
	content = strings.Trim(content, "\n")
	content = strings.ReplaceAll(content, "\r\n", "<br>")

	return content
}

//更新模拟器独立参数
func UpdateRomCmd(id uint64, simId uint32, data map[string]string) error {

	//开始更新
	if err := (&db.Rom{}).UpdateSimConf(id, simId, data["cmd"], uint8(utils.ToInt(data["unzip"])), data["file"], data["lua"]); err != nil {
		return err
	}

	return nil
}

//读取rom以及相关资源
func DeleteRomAndRes(id uint64, deleteRes int) error {

	//游戏游戏详细数据
	info, err := (&db.Rom{}).GetById(id)
	if err != nil {
		return err
	}

	fname := utils.GetFileName(info.RomPath)
	platform := config.Cfg.Platform[info.Platform]

	//删除rom文件
	go func() {

		//删除主游戏
		_ = utils.FileDelete(platform.RomPath + config.Cfg.Separator + info.RomPath)

		//删除子游戏
		romFiles, _ := utils.ScanDirByKeyword(platform.RomPath, fname+"__")
		for _, f := range romFiles {
			_ = utils.FileDelete(f)
		}
	}()

	//不删除资源文件
	if deleteRes == 0 {
		return nil
	}

	//删除资源文件
	exts := config.GetResExts()
	go func() {
		for t, path := range config.GetResPath(platform.Id) {

			for _, v := range exts[t] {
				_ = utils.FileDelete(path + config.Cfg.Separator + fname + v)
			}

			resFiles, _ := utils.ScanDirByKeyword(path, fname+"__")
			for _, f := range resFiles {
				utils.FileDelete(f)
			}
		}
	}()

	return nil
}

//读取rom资源
func GetRomRes(typ string, pf uint32, romName string) string {

	fileName := ""
	resName := ""
	res := config.GetResPath(pf)
	types := config.GetResExts()

	if res[typ] != "" {
		for _, v := range types[typ] {
			fileName = res[typ] + config.Cfg.Separator + romName + v
			if utils.FileExists(fileName) {
				resName = fileName
				break
			}
		}
	}

	return resName
}

//移动rom及资源文件
func MoveRom(id uint64, newFolder string) error {

	//读取rom详情
	rom, err := (&db.Rom{}).GetById(id)
	if err != nil {
		utils.WriteLog(err.Error())
	}

	//生成目录地址
	romPath := config.Cfg.Platform[rom.Platform].RomPath + config.Cfg.Separator
	romName := utils.GetFileNameAndExt(rom.RomPath)
	oldFile := romPath + rom.RomPath
	newFile := ""

	if newFolder == ConstMenuRootKey {
		newFile = romPath + romName
	} else {
		newFile = romPath + newFolder + config.Cfg.Separator + romName
	}

	//如果位置一样则不用移动
	if oldFile == newFile {
		return nil
	}

	//移动主rom文件
	if err := utils.FileMove(oldFile, newFile); err != nil {
	}

	//更新数据库
	newFile = strings.Replace(newFile, romPath, "", -1)
	(&db.Rom{Id: id, RomPath: newFile, Menu: newFolder}).UpdateRomPath()
	return nil
}

//编辑rom基础信息
func SetRomBaseName(id uint64, name string) error {

	rom, _ := (&db.Rom{}).GetById(id)

	if config.Cfg.Platform[rom.Platform].Rombase == "" {
		return errors.New(config.Cfg.Lang["RomBaseFileNotFound"])
	}

	rombase := GetRomBaseById(rom.Platform, utils.GetFileName(rom.RomPath))
	romName := utils.GetFileName(rom.RomPath)

	create := &RomBase{}

	if rombase != nil {
		create = rombase
		create.Name = name
	} else {
		create = &RomBase{
			RomName: romName,
			Name:    name,
		}
	}

	//写入配置文件
	if err := WriteRomBaseFile(rom.Platform, create); err != nil {
		return err
	}

	if name == "" || name == romName {
		name = romName
	}

	//更新到数据库
	infoMd5 := GetRomMd5(rom.Name, rom.RomPath, create.Type, create.Year, create.Producer, create.Publisher, create.Country, create.Translate, create.Version, create.NameEN, create.NameJP, create.OtherA, create.OtherB, create.OtherC, create.OtherD, rom.Size)
	dbRom := &db.Rom{
		Name:          name,
		BaseType:      create.Type,
		BaseYear:      create.Year,
		BaseProducer:  create.Producer,
		BasePublisher: create.Publisher,
		BaseCountry:   create.Country,
		BaseTranslate: create.Translate,
		BaseVersion:   create.Version,
		BaseNameEn:    create.NameEN,
		BaseNameJp:    create.NameJP,
		BaseOtherA:    create.OtherA,
		BaseOtherB:    create.OtherB,
		BaseOtherC:    create.OtherC,
		BaseOtherD:    create.OtherD,
		InfoMd5:       infoMd5,
	}
	if err := dbRom.UpdateRomBase(id); err != nil {
		return err
	}

	return nil
}

//编辑rom基础信息
func SetRomBase(d map[string]string) (*db.Rom, error) {

	rom, _ := (&db.Rom{}).GetById(uint64(utils.ToInt(d["id"])))

	if config.Cfg.Platform[rom.Platform].Rombase == "" {
		return nil, errors.New(config.Cfg.Lang["RomBaseFileNotFound"])
	}

	romName := utils.GetFileName(rom.RomPath)
	romBase := &RomBase{
		RomName:   romName,
		Name:      d["name"],
		Type:      d["type"],
		Year:      d["year"],
		Producer:  d["producer"],
		Publisher: d["publisher"],
		Country:   d["country"],
		Translate: d["translate"],
		Version:   d["version"],
		NameEN:    d["name_en"],
		NameJP:    d["name_jp"],
		OtherA:    d["other_a"],
		OtherB:    d["other_b"],
		OtherC:    d["other_c"],
		OtherD:    d["other_d"],
	}

	//写入配置文件
	if err := WriteRomBaseFile(rom.Platform, romBase); err != nil {
		return nil, err
	}

	name := d["name"]
	if name == "" || d["name"] == romName {
		name = romName
	}

	//更新到数据库
	infoMd5 := GetRomMd5(rom.Name, rom.RomPath, d["type"], d["year"], d["producer"], d["publisher"], d["country"], d["translate"], d["version"], d["name_en"], d["name_jp"], d["other_a"], d["other_b"], d["other_c"], d["other_d"], rom.Size)
	dbRom := &db.Rom{
		Name:          name,
		BaseType:      d["type"],
		BaseYear:      d["year"],
		BaseProducer:  d["producer"],
		BasePublisher: d["publisher"],
		BaseCountry:   d["country"],
		BaseTranslate: d["translate"],
		BaseVersion:   d["version"],
		BaseNameEn:    d["name_en"],
		BaseNameJp:    d["name_jp"],
		BaseOtherA:    d["other_a"],
		BaseOtherB:    d["other_b"],
		BaseOtherC:    d["other_c"],
		BaseOtherD:    d["other_d"],
		InfoMd5:       infoMd5,
	}
	if err := dbRom.UpdateRomBase(uint64(utils.ToInt(d["id"]))); err != nil {
		return nil, err
	}

	return dbRom, nil
}

//批量编辑rom基础信息
func BatchSetRomBase(data []map[string]string) error {
	ids := []uint64{}
	for _, v := range data {
		ids = append(ids, uint64(utils.ToInt(v["id"])))
	}

	voList, _ := (&db.Rom{}).GetByIds(ids)
	romList := map[uint64]*db.Rom{}
	platform := uint32(0)
	for _, v := range voList {
		if platform == 0 {
			platform = v.Platform
		}
		ids = append(ids, uint64(utils.ToInt(v.Id)))
		romList[v.Id] = v
	}

	//读取出所有资料
	rombaseList, _ := GetRomBaseList(platform)
	for _, d := range data {

		rom := romList[uint64(utils.ToInt(d["id"]))]

		base := &RomBase{
			RomName:   utils.GetFileName(rom.RomPath),
			Name:      d["name"],
			Type:      d["type"],
			Year:      d["year"],
			Producer:  d["producer"],
			Publisher: d["publisher"],
			Country:   d["country"],
			Translate: d["translate"],
			Version:   d["version"],
			NameEN:    d["name_en"],
			NameJP:    d["name_jp"],
			OtherA:    d["other_a"],
			OtherB:    d["other_b"],
			OtherC:    d["other_c"],
			OtherD:    d["other_d"],
		}

		rombaseList[utils.GetFileName(rom.RomPath)] = base
	}

	//写入配置文件
	if err := CoverRomBaseFile(platform, rombaseList); err != nil {
		return err
	}

	return nil
}

/**
 * 上传子rom文件
 **/
func UploadSubGameFile(pid uint64, name string, p string) (string, error) {
	parent, _ := (&db.Rom{}).GetById(pid)
	if config.Cfg.Platform[parent.Platform].RomPath == "" {
		return "", errors.New(config.Cfg.Lang["FilesMenuCanNotBeEmpty"])
	}
	fileNameAndExt := utils.GetFileNameAndExt(p)
	path := ""

	if utils.IsAbsPath(parent.RomPath) {
		path = utils.GetFilePath(parent.RomPath)
	} else {
		path = config.Cfg.Platform[parent.Platform].RomPath + utils.GetFilePath(parent.RomPath)
	}

	newPath := path + config.Cfg.Separator + fileNameAndExt

	if utils.FileExists(newPath) {
		fmt.Println("重复游戏，跳过")
		return p, nil
	}

	if strings.Contains(newPath, filepath.Dir(p)) {
		if err := utils.FileMove(p, newPath); err != nil {
			return "", err
		}
	} else {
		if err := utils.FileCopy(p, newPath); err != nil {
			return "", err
		}
	}

	relPath := strings.Replace(newPath, config.Cfg.RootPath, "", -1)
	return relPath, nil
}

/**
 * 更新子rom数据
 **/
func UpdateSubGameFiles(id uint64, data string) error {
	vo, _ := (&db.Rom{}).GetById(id)
	if config.Cfg.Platform[vo.Platform].RomPath == "" {
		return errors.New(config.Cfg.Lang["FilesMenuCanNotBeEmpty"])
	}

	//整理需要删除的文件
	d := []map[string]string{}
	json.Unmarshal([]byte(data), &d)
	newData := []string{}
	for _, v := range d {
		newData = append(newData, v["path"])
	}

	//处理文件
	romName := utils.GetFileName(vo.RomPath)
	exists, _ := utils.ScanDirByKeyword(config.Cfg.Platform[vo.Platform].RomPath, romName+"__")
	//真正存在的文件
	realFiles := []string{}
	for _, v := range exists {
		rel := strings.Replace(v, config.Cfg.RootPath, "", 1)
		if !utils.InSliceString(rel, newData) {
			utils.FileDelete(v)
		} else {
			realFiles = append(realFiles, rel)
		}
	}

	//写入缓存
	create := []*db.Rom{}
	for _, p := range realFiles {
		title := utils.GetFileName(p)
		titles := strings.Split(title, "__")
		title = titles[1]
		relPath := strings.Replace(utils.AbsPath(p), config.Cfg.Platform[vo.Platform].RomPath+config.Cfg.Separator, "", 1)

		f, err := os.Open(p)
		defer f.Close()
		if err != nil {
			fmt.Println("err:打开文件失败", p, err)
			continue
		}
		stat, _ := f.Stat()
		fileMd5 := utils.ToString(stat.ModTime().UnixNano())

		if fileMd5 == "0" || len(fileMd5) < 19 || fileMd5[len(fileMd5)-4:] == "0000" {
			t := time.Now()
			os.Chtimes(p, t, t)
			fileMd5 = utils.ToString(t.UnixNano())
		}

		fileMd5 = utils.CreateRomUniqId(fileMd5, stat.Size())
		infoMd5 := GetRomMd5(title, relPath)

		c := &db.Rom{
			RomPath:  relPath,
			Name:     title,
			Pinyin:   utils.TextToPinyin(title),
			FileMd5:  fileMd5,
			InfoMd5:  infoMd5,
			Pname:    vo.FileMd5,
			Platform: vo.Platform,
			Menu:     vo.Menu,
		}
		create = append(create, c)
	}

	//删除旧缓存
	_ = (&db.Rom{}).DeleteSubRom(vo.Platform, vo.FileMd5)

	if len(create) > 0 {
		(&db.Rom{}).BatchAdd(create, 0)
	}

	return nil
}

/**
 * 读取过滤器列表
 **/
func GetFilter(platform uint32) (map[string][]string, error) {

	volist, _ := (&db.Filter{}).GetByPlatform(platform)

	//填充数据
	filterList := map[string][]string{}
	for _, v := range volist {
		if _, ok := filterList[v.Type]; ok {
			filterList[v.Type] = append(filterList[v.Type], v.Name)
		} else {
			filterList[v.Type] = []string{v.Name}
		}
	}

	//补全不存在的数据
	types := []string{"base_type", "base_year", "base_producer", "base_publisher", "base_country", "base_translate", "base_version"}
	for _, t := range types {
		if _, ok := filterList[t]; !ok {
			filterList[t] = []string{}
		}
	}
	return filterList, nil
}
