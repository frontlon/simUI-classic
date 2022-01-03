package modules

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"simUI/code/config"
	"simUI/code/db"
	"simUI/code/utils"
	"strings"
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
	RomFileSize   string              //rom文件大小
}

//运行游戏
func RunGame(romId uint64, simId uint32) error {

	//数据库中读取rom详情
	rom, err := (&db.Rom{}).GetById(romId)
	if err != nil {
		return err
	}

	romCmd, _ := (&db.Rom{}).GetSimConf(romId, simId)

	sim := &db.Simulator{}
	if simId == 0 {
		sim = config.Cfg.Platform[rom.Platform].UseSim
		if sim == nil {
			return errors.New(config.Cfg.Lang["SimulatorNotFound"])
		}
	} else {
		if config.Cfg.Platform[rom.Platform].SimList == nil {
			return errors.New(config.Cfg.Lang["SimulatorNotFound"])
		}
		sim = config.Cfg.Platform[rom.Platform].SimList[simId]
	}

	//如果是相对路径，转换成绝对路径
	if !strings.Contains(rom.RomPath, ":") {
		rom.RomPath = config.Cfg.Platform[rom.Platform].RomPath + config.Cfg.Separator + rom.RomPath
	}

	//解压zip包
	if (sim.Unzip == 1 && romCmd.Unzip == 2) || romCmd.Unzip == 1 {
		RomExts := strings.Split(config.Cfg.Platform[rom.Platform].RomExts, ",")
		rom.RomPath, err = UnzipRom(rom.RomPath, RomExts)
		if err != nil {
			return err
		}
		if rom.RomPath == "" {
			return errors.New(config.Cfg.Lang["UnzipExeNotFound"])
		}

		//如果指定了执行文件
		if romCmd.File != "" {
			rom.RomPath = utils.GetFileAbsPath(rom.RomPath) + "/" + romCmd.File
		}

	}

	//检测rom文件是否存在
	if utils.FileExists(rom.RomPath) == false {
		return errors.New(config.Cfg.Lang["RomNotFound"] + rom.RomPath)
	}

	//加载运行参数
	cmd := []string{}

	ext := utils.GetFileExt(rom.RomPath)

	//运行游戏前，先杀掉之前运行的程序
	if err = utils.KillGame(); err != nil {
		return err
	}

	simCmd := ""
	simLua := ""

	if romCmd.Cmd != "" {
		simCmd = romCmd.Cmd
	} else if sim.Cmd != "" {
		simCmd = sim.Cmd
	}

	if romCmd.Lua != "" {
		simLua = romCmd.Lua
	} else if sim.Lua != "" {
		simLua = sim.Lua
	}

	//如果是可执行程序，则不依赖模拟器直接运行
	if utils.InSliceString(ext, config.RUN_EXTS) {
		cmd = append(cmd, rom.RomPath)
	} else { //如果依赖模拟器
		//检测模拟器文件是否存在
		_, err = os.Stat(sim.Path)
		if err != nil {
			return errors.New(config.Cfg.Lang["SimulatorNotFound"])
		}

		if simCmd == "" {
			cmd = append(cmd, rom.RomPath)
		} else {
			//如果rom运行参数存在，则使用rom的参数
			cmd = strings.Split(simCmd, " ")
			filename := filepath.Base(rom.RomPath) //exe运行文件路径
			//替换变量
			for k, _ := range cmd {
				cmd[k] = strings.ReplaceAll(cmd[k], `{RomName}`, utils.GetFileName(filename))
				cmd[k] = strings.ReplaceAll(cmd[k], `{RomExt}`, utils.GetFileExt(filename))
				cmd[k] = strings.ReplaceAll(cmd[k], `{RomFullPath}`, rom.RomPath)
			}
		}
	}

	//运行lua脚本
	if simLua != "" {
		cmdStr := utils.SlicetoString(" ", cmd)
		callLua(sim.Lua, sim.Path, cmdStr)
	}

	//记录运行信息
	_ = (&db.Rom{}).UpdateRunNumAndTime(romId)

	//运行游戏
	err = utils.RunGame(sim.Path, cmd)

	return nil
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
	if err := utils.OpenFolderByWindow(fileName); err != nil {
		return err
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
	romName := utils.GetFileName(info.RomPath)
	sub, _ := (&db.Rom{}).GetSubRom(info.Platform, romName)

	res.Info = info
	res.Sublist = sub
	res.Simlist, _ = (&db.Simulator{}).GetByPlatform(info.Platform)

	//获取rom文件大小
	if res.Info.RomPath != "" {
		fi := config.Cfg.Platform[info.Platform].RomPath + config.Cfg.Separator + res.Info.RomPath
		f, err := os.Stat(fi)
		if err == nil {
			res.RomFileSize = utils.GetFileSizeString(f.Size())
		}
	}

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
	fmt.Println(content)

	return content
}

//更新模拟器独立参数
func UpdateRomCmd(id uint64, simId uint32, data map[string]string) error {

	if data["cmd"] == "" && data["unzip"] == "2" {
		//如果当前配置和模拟器默认配置一样，则删除该记录
		if err := (&db.Rom{}).DelSimConf(id, simId); err != nil {
			return err
		}
	} else {
		//开始更新
		if err := (&db.Rom{}).UpdateSimConf(id, simId, data["cmd"], uint8(utils.ToInt(data["unzip"])), data["file"], data["lua"]); err != nil {
			return err
		}
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

func UploadStrategyImages(id uint64, p string) (string, error) {

	//游戏游戏详细数据
	info, err := (&db.Rom{}).GetById(id)
	if err != nil {
		return "", err
	}

	strategyPath := config.Cfg.Platform[info.Platform].StrategyPath + config.Cfg.Separator + "images/"
	if strategyPath == "" {
		return "", nil
	}

	//先检查目录是否存在，不存在创建目录
	if !utils.FolderExists(strategyPath) {
		if err := utils.CreateDir(strategyPath); err != nil {
			return "", err
		}
	}

	//复制文件
	newFilename := utils.GetFileNameAndExt(p)
	newFile := strategyPath + newFilename
	if err := utils.FileCopy(p, newFile); err != nil {
		return "", err
	}
	return newFile, nil

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
func MoveRom(id uint64, newPlatform uint32, newFolder string) error {

	//读取rom详情
	rom, err := (&db.Rom{}).GetById(id)
	if err != nil {
		utils.WriteLog(err.Error())
	}

	//生成目录地址
	romName := utils.GetFileNameAndExt(rom.RomPath)
	oldFile := config.Cfg.Platform[rom.Platform].RomPath + config.Cfg.Separator + rom.RomPath
	newFile := ""

	if newFolder == "/" {
		newFile = config.Cfg.Platform[newPlatform].RomPath + config.Cfg.Separator + romName
	} else {
		newFile = config.Cfg.Platform[newPlatform].RomPath + config.Cfg.Separator + newFolder + config.Cfg.Separator + romName
	}

	//如果位置一样则不用移动
	if oldFile == newFile {
		return nil
	}

	//移动主rom文件
	if err := utils.FileMove(oldFile, newFile); err != nil {
	}

	//移动子rom
	subName := utils.GetFileName(romName)
	romFiles, _ := utils.ScanDirByKeyword(config.Cfg.Platform[rom.Platform].RomPath, subName+"__")

	for _, f := range romFiles {
		subRomName := utils.GetFileNameAndExt(f)
		newSubFile := ""
		if newFolder == "/" {
			newSubFile = config.Cfg.Platform[newPlatform].RomPath + config.Cfg.Separator + subRomName
		} else {
			newSubFile = config.Cfg.Platform[newPlatform].RomPath + config.Cfg.Separator + newFolder + config.Cfg.Separator + subRomName
		}

		if err := utils.FileMove(f, newSubFile); err != nil {
		}
	}

	//同平台下不用移动资源文件
	if rom.Platform == newPlatform {
		return nil
	}

	//开始移动资源文件
	romName = utils.GetFileName(filepath.Base(rom.RomPath))
	for resName, path := range config.GetResPath(newPlatform) {
		name := GetRomRes(resName, rom.Platform, romName)
		if name != "" {
			_ = utils.FileMove(name, path+config.Cfg.Separator+utils.GetFileNameAndExt(name))
		}
	}

	//移动攻略文件
	if config.Cfg.Platform[newPlatform].FilesPath == "" {
		return nil
	}
	files, _ := utils.ScanDirByKeyword(config.Cfg.Platform[rom.Platform].FilesPath, romName+"__")
	filesPath := config.Cfg.Platform[newPlatform].FilesPath
	for _, f := range files {
		_ = utils.FileMove(f, filesPath+config.Cfg.Separator+utils.GetFileNameAndExt(f))
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
	fileMd5 := GetRomMd5(utils.ToString(rom.Platform), rom.RomPath)
	infoMd5 := GetRomMd5(rom.Name, rom.RomPath, d["type"], d["year"], d["producer"], d["publisher"], d["country"], d["translate"], d["version"], d["name_en"], d["name_jp"], d["other_a"], d["other_b"], d["other_c"], d["other_d"])
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
		FileMd5:       fileMd5,
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

		romName := utils.GetFileName(rom.RomPath)

		base := &RomBase{
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

		fmt.Println("base", base)

		rombaseList[romName] = base

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
func UploadSubGameFile(id uint64, name string, p string) (string, error) {
	vo, _ := (&db.Rom{}).GetById(id)
	if config.Cfg.Platform[vo.Platform].RomPath == "" {
		return "", errors.New(config.Cfg.Lang["FilesMenuCanNotBeEmpty"])
	}
	ext := utils.GetFileExt(p)
	fileName := utils.GetFileName(vo.RomPath)
	newPath := config.Cfg.Platform[vo.Platform].RomPath + config.Cfg.Separator + fileName + "__" + name + ext

	rel := strings.Replace(newPath, config.Cfg.RootPath, "", 1)
	if rel == p {
		return p, nil
	}

	if(strings.Contains(newPath,filepath.Dir(p))){
		if err := utils.FileMove(p, newPath); err != nil {
			return "", err
		}
	}else{
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

	//删除旧缓存
	_ = (&db.Rom{}).DeleteSubRom(vo.Platform, vo.Name)

	//写入缓存
	create := []*db.Rom{}
	for _, r := range realFiles {
		title := utils.GetFileName(r)
		titles := strings.Split(title, "__")
		title = titles[1]
		r = strings.Replace(utils.AbsPath(r), config.Cfg.Platform[vo.Platform].RomPath+config.Cfg.Separator, "", 1)

		fileMd5 := GetRomMd5(utils.ToString(vo.Platform), r)
		infoMd5 := GetRomMd5(title, r, "", "", "", "", "", "")

		c := &db.Rom{
			RomPath:  r,
			Name:     title,
			Pinyin:   utils.TextToPinyin(title),
			FileMd5:  fileMd5,
			InfoMd5:  infoMd5,
			Pname:    utils.GetFileName(vo.RomPath),
			Platform: vo.Platform,
			Menu:     vo.Menu,
		}
		create = append(create, c)
	}

	if len(create) > 0 {
		(&db.Rom{}).BatchAdd(create, 0)
	}

	return nil
}
