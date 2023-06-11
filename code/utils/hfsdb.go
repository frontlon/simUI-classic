package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"
)

const HFSDB_GAME_URL = "https://db.hfsplay.fr/api/v1/games?limit=%v&offset=%v&search=%v&format=json"
const GET_TOKEN_URL = "https://www.simui.net/api/getHsfToken.php"

var hfsDbToken = ""

type HfsDbGame struct {
	Count    int             `json:"count,omitempty"`
	Next     string          `json:"next,omitempty"`
	Previous string          `json:"previous,omitempty"`
	Results  []HfsDbGameData `json:"results,omitempty"`
	Detail   string          `json:"detail"` //返回错误信息
}

type HfsDbGameData struct {
	Id               int64            `json:"id,omitempty"`
	Medias           []HfsDbGameMedia `json:"medias,omitempty"`
	NameEn           string           `json:"name_en,omitempty"`
	NameFr           string           `json:"name_fr,omitempty"`
	NameJp           string           `json:"name_jp,omitempty"`
	DescriptionEn    string           `json:"description_en,omitempty"`
	DescriptionFr    string           `json:"description_fr,omitempty"`
	DescriptionJp    string           `json:"description_jp,omitempty"`
	RegionEn         string           `json:"region_en,omitempty"`
	RegionFr         string           `json:"region_fr,omitempty"`
	RegionJp         string           `json:"region_jp,omitempty"`
	ReleasedAt_PAL   string           `json:"released_at_PAL,omitempty"`
	ReleasedAt_US    string           `json:"released_at_US,omitempty"`
	ReleasedAt_JPN   string           `json:"released_at_JPN,omitempty"`
	ReleasedAt_WORLD string           `json:"released_at_WORLD,omitempty"`
}

type HfsDbGameMedia struct {
	Id          int64  `json:"id,omitempty"`          //资源ID
	File        string `json:"file,omitempty"`        //文件/图片地址
	Type        string `json:"type,omitempty"`        //资源类型
	Description string `json:"description,omitempty"` //资源描述
	Region      string `json:"region,omitempty"`      //地区
	IsImage     bool   `json:"is_image,omitempty"`    //是否为图片
	Extension   string `json:"extension,omitempty"`   //扩展名
	ResX        int    `json:"res_x,omitempty"`       //宽度
	ResY        int    `json:"res_y,omitempty"`       //高度
}

// 搜索游戏资料
func GetHfsDbGameList(keyword string, limit, offset int) (*HfsDbGame, error) {
	keyword = url.QueryEscape(keyword)

	uri := fmt.Sprintf(HFSDB_GAME_URL, limit, offset, keyword)

	headers := map[string]string{
		"Authorization": "Token " + GetHfsDbUserToken(),
	}

	resp, err := HttpGet(uri, headers)

	if err != nil {
		return nil, err
	}

	var body HfsDbGame
	if err := json.Unmarshal(resp, &body); err != nil {
		return nil, err
	}

	//报错了
	if body.Detail != "" {
		//token失效：Invalid token.
		return nil, errors.New(body.Detail)
	}
	return &body, nil
}

// 获取token
func GetHfsDbUserToken() string {

	if hfsDbToken != "" {
		return hfsDbToken
	}

	from := "simui"
	now := ToString(time.Now().Unix())
	secret := RandStr(16, 16)
	salt := ToString(RandInt(10000000, 99999999))

	req := map[string]string{
		"secret": secret,
		"salt":   salt,
		"time":   now,
		"from":   from,
		"sign":   Md5(from + now + salt + Md5(secret+salt)),
	}

	resp, err := HttpPostForm(GET_TOKEN_URL, req, nil)

	if err != nil {
		return ""
	}

	var body map[string]string
	if err := json.Unmarshal(resp, &body); err != nil {
		return ""
	}

	if body["token"] != "" {
		hfsDbToken = body["token"]
	}
	return hfsDbToken
}
