package request

/**
 * 读取游戏列表 GetGameList
 * 读取游戏数量 GetGameCount
**/
type GetGameList struct {
	ShowHide      uint8  `json:"showHide"`      //是否隐藏
	ShowSubGame   uint8  `json:"showSubGame"`   //是否显示子游戏
	Platform      uint32 `json:"platform"`      //平台
	Catname       string `json:"catname"`       //分类
	Keyword       string `json:"keyword"`       //关键字
	Num           string `json:"num"`           //字母索引
	Page          int    `json:"page"`          //分页数
	BaseType      string `json:"baseType"`      //资料 - 游戏类型
	BasePublisher string `json:"basePublisher"` //资料 - 发布者
	BaseYear      string `json:"baseYear"`      //资料 - 发型年份
	BaseCountry   string `json:"baseCountry"`   //资料 - 国家
	BaseTranslate string `json:"baseTranslate"` //资料 - 语言
	BaseVersion   string `json:"baseVersion"`   //资料 - 版本
	BaseProducer  string `json:"baseProducer"`  //资料 - 制作商
}
