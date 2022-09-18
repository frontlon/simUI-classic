package modules

import (
	"errors"
	"fmt"
	"io/ioutil"
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
func OpenFolder(id uint64, opt string, otherId string) error {

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
		if _, ok := platform.SimList[uint32(utils.ToInt(otherId))]; ok {
			fileName = platform.SimList[uint32(utils.ToInt(otherId))].Path
		}
	default:
		res := config.GetResPath(platform.Id)

		if res[opt] != "" {
			if otherId != "" {
				romName = romName + ConstSeparator + otherId
			}
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

//添加独立游戏（PC/PS3）
func AddIndieGame(platform uint32, menu string, files []string) error {

	ext := ".path"

	//补丁逻辑：先检查下平台扩展名是否存在.path扩展名
	exts := config.Cfg.Platform[platform].RomExts
	if !strings.Contains(exts, ext) {
		exts = exts + "," + ext
		(&db.Platform{Id: platform}).UpdateFieldById("rom_exts", exts)
		config.Cfg.Platform[platform].RomExts = exts
	}

	root := config.Cfg.Platform[platform].RomPath + config.Cfg.Separator
	folder := ""
	if menu != "" {
		folder = root + menu + config.Cfg.Separator
	}

	//如果目录不存在，则放到根目录中
	if !utils.FolderExists(folder) {
		folder = root
	}

	for _, f := range files {
		filename := utils.GetFileName(f)
		p := folder + filename + ext
		f = strings.Replace(f, `file://`, "", 1)
		f = strings.Replace(f, `file:\\`, "", 1)

		//转换为相对路径
		rootPath := strings.ReplaceAll(config.Cfg.RootPath, "\\", "/")
		relPath := strings.Replace(f, rootPath, "", 1)

		//如果文件已存在，则使用新名称
		if utils.FileExists(p) {
			path := utils.GetFilePath(p)
			name := utils.GetFileName(p) + utils.ToString(time.Now().Unix())
			ext := utils.GetFileExt(p)
			p = path + config.Cfg.Separator + name + ext
		}

		if err := utils.OverlayWriteFile(p, relPath); err != nil {
			return err
		}
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
				res.DocContent = strings.Trim(res.DocContent, "\t")
				res.DocContent = strings.Trim(res.DocContent, "\n\r")
				res.DocContent = strings.Trim(res.DocContent, "\r")
				res.DocContent = strings.Trim(res.DocContent, "\n")
				res.DocContent = strings.ReplaceAll(res.DocContent, "\r\n", "<br>")
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
func GetGameDoc(t string, id uint64, toHtml uint8) (string, error) {

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
	if toHtml == 1 {
		strategy = strings.Trim(strategy, "\t")
		strategy = strings.Trim(strategy, "\n\r")
		strategy = strings.Trim(strategy, "\r")
		strategy = strings.Trim(strategy, "\n")
		strategy = strings.ReplaceAll(strategy, "\r\n", "<br>")
	}
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

//删除rom以及相关资源
func DeleteRomAndRes(id uint64, deleteRes int) error {

	//游戏游戏详细数据
	info, err := (&db.Rom{}).GetById(id)
	if err != nil {
		return err
	}

	subInfo, err := (&db.Rom{}).GetSubRom(info.Platform, info.FileMd5)
	fileMd5s := []string{info.FileMd5}
	if len(subInfo) > 0 {
		for _, v := range subInfo {
			fileMd5s = append(fileMd5s, v.FileMd5)
		}
	}
	fname := utils.GetFileName(info.RomPath)
	platform := config.Cfg.Platform[info.Platform]

	//删除数据库
	(&db.Rom{}).DeleteById(id)
	(&db.Rom{}).DeleteSubRom(info.Platform, info.FileMd5)
	(&db.RomSubGame{}).DeleteByFileMd5s(info.Platform, fileMd5s)
	(&db.RomSetting{}).DeleteByFileMd5s(info.Platform, fileMd5s)

	//删除rom文件
	go func() {
		romFiles, _ := utils.ScanMasterSlaveFiles(platform.RomPath, fname)
		for _, f := range romFiles {
			_ = utils.FileDelete(f)
		}
	}()

	//不删除资源文件
	if deleteRes == 0 {
		return nil
	}

	//删除资源文件
	go func() {
		resPaths := config.GetResPath(platform.Id)
		for _, path := range resPaths {
			romFiles, _ := utils.ScanMasterSlaveFiles(path, fname)
			for _, f := range romFiles {
				_ = utils.FileDelete(f)
			}
		}

		//删除音乐文件夹
		audioPath := resPaths["audio"] + config.Cfg.Separator + fname
		_ = utils.DeleteDir(audioPath)
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
		return err
	}

	//读取menu信息
	virtual := int8(0)
	if newFolder == ConstMenuRootKey {
		virtual = 1
	} else {
		menu := (&db.Menu{}).GetByName(rom.Platform, newFolder)
		if menu != nil {
			virtual = menu.Virtual
		}
	}

	//如果是虚拟目录，则只更新数据库
	if virtual == 1 {
		//更新数据库
		(&db.Rom{Id: id, RomPath: rom.RomPath, Menu: newFolder}).UpdateRomPath()
		//更新menu信息
		(&db.RomSetting{
			Platform: rom.Platform,
			FileMd5:  rom.FileMd5,
			Menu:     newFolder,
		}).UpdateMenu()

		return nil
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

	//更新数据库
	newRelFile := strings.Replace(newFile, romPath, "", -1)
	(&db.Rom{Id: id, RomPath: newRelFile, Menu: newFolder}).UpdateRomPath()

	//清除menu信息
	(&db.RomSetting{
		Platform: rom.Platform,
		FileMd5:  rom.FileMd5,
	}).ClearMenu()

	//如果位置一样则不用移动
	if oldFile == newFile {
		return nil
	}

	//移动主rom文件
	if err := utils.FileMove(oldFile, newFile); err != nil {
	}

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
	infoMd5 := utils.GetRomMd5(rom.Name, rom.RomPath, create.Type, create.Year, create.Producer, create.Publisher, create.Country, create.Translate, create.Version, create.NameEN, create.NameJP, create.OtherA, create.OtherB, create.OtherC, create.OtherD, rom.Score, rom.Size)
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
		Score:     d["score"],
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
	infoMd5 := utils.GetRomMd5(rom.Name, rom.RomPath, d["type"], d["year"], d["producer"], d["publisher"], d["country"], d["translate"], d["version"], d["name_en"], d["name_jp"], d["other_a"], d["other_b"], d["other_c"], d["other_d"], d["score"], rom.Size)
	dbRom := &db.Rom{
		Name:          name,
		BaseType:      d["type"],
		BaseYear:      d["year"],
		BaseProducer:  d["producer"],
		BasePublisher: d["publisher"],
		BaseCountry:   d["country"],
		BaseTranslate: d["translate"],
		BaseVersion:   d["version"],
		Score:         d["score"],
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
 * 读取过滤器列表
 **/
func GetFilter(platform uint32) (map[string][]string, error) {

	volist := []*db.Filter{}
	if platform == 0 {
		volist, _ = (&db.Filter{}).GetAll()
	} else {
		volist, _ = (&db.Filter{}).GetByPlatform(platform)
	}

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
	types := []string{"base_type", "base_year", "base_producer", "base_publisher", "base_country", "base_translate", "base_version", "score", "complete"}
	for _, t := range types {
		if _, ok := filterList[t]; !ok {
			filterList[t] = []string{}
		}
	}
	return filterList, nil
}

/**
 * 设置游戏评分
 **/
func SetScore(id uint64, score string) error {

	//更新csv文件
	rom, _ := (&db.Rom{}).GetById(id)

	if config.Cfg.Platform[rom.Platform].Rombase == "" {
		return errors.New(config.Cfg.Lang["RomBaseFileNotFound"])
	}

	rombase := GetRomBaseById(rom.Platform, utils.GetFileName(rom.RomPath))
	romName := utils.GetFileName(rom.RomPath)

	create := &RomBase{}

	if rombase != nil {
		create = rombase
		create.Score = score
	} else {
		create = &RomBase{
			RomName: romName,
			Score:   score,
		}
	}
	if err := WriteRomBaseFile(rom.Platform, create); err != nil {
		return err
	}

	//更新数据库
	if err := (&db.Rom{}).UpdateScore(id, score); err != nil {
		return err
	}

	return nil
}

/**
 * 设置通关状态
 **/
func SetComplete(id uint64, status uint8) error {

	//更新csv文件
	rom, _ := (&db.Rom{}).GetById(id)

	(&db.Rom{}).UpdateComplete(id, status)

	(&db.RomSetting{
		Platform: rom.Platform,
		FileMd5:  rom.FileMd5,
		Complete: status,
	}).UpdateComplete()

	return nil
}
