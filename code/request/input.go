package request

/**
 * 导出(分享)rom OutputRom
**/
type OutputRom struct {
	Save       string   `json:"save"`       //zip文件保存路径
	Platform   uint32   `json:"platform"`   //平台id
	Opt        string   `json:"opt"`        //导出类型
	Options    []string `json:"options"`    //导出选项
	Menus      []string `json:"menus"`      //目录列表
	Roms       []uint64 `json:"roms"`       //rom列表
	Simulators []uint32 `json:"simulators"` //模拟器列表
}
