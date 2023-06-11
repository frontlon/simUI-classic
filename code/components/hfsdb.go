package components

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"simUI/code/utils"
	"strconv"
	"time"
)

type HfsDb struct {
	Count    int         `json:"count,omitempty"`
	Next     string      `json:"next,omitempty"`
	Previous string      `json:"previous,omitempty"`
	Results  []HfsDbData `json:"results,omitempty"`
	Detail   string      `json:"detail"` //返回错误信息
}

type HfsDbData struct {
	Id               int64        `json:"id,omitempty"`
	Medias           []HfsDbMedia `json:"medias,omitempty"`
	NameEn           string       `json:"name_en,omitempty"`
	NameFr           string       `json:"name_fr,omitempty"`
	NameJp           string       `json:"name_jp,omitempty"`
	DescriptionEn    string       `json:"description_en,omitempty"`
	DescriptionFr    string       `json:"description_fr,omitempty"`
	DescriptionJp    string       `json:"description_jp,omitempty"`
	RegionEn         string       `json:"region_en,omitempty"`
	RegionFr         string       `json:"region_fr,omitempty"`
	RegionJp         string       `json:"region_jp,omitempty"`
	ReleasedAt_PAL   string       `json:"released_at_PAL,omitempty"`
	ReleasedAt_US    string       `json:"released_at_US,omitempty"`
	ReleasedAt_JPN   string       `json:"released_at_JPN,omitempty"`
	ReleasedAt_WORLD string       `json:"released_at_WORLD,omitempty"`
}

type HfsDbMedia struct {
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

func GetHfsDbGameList(keyword string, limit, offset int) (*HfsDb, error) {
	token, err := getHfsUserToken()
	if err != nil {
		return nil, err
	}

	gameUrl := `https://db.hfsplay.fr/api/v1/games?limit=%v&offset=%v&search=%v&format=json`
	keyword = url.QueryEscape(keyword)
	uri := fmt.Sprintf(gameUrl, limit, offset, keyword)
	headers := map[string]string{
		"Authorization": "Token " + token,
	}

	resp, err := utils.HttpGet(uri, headers)

	if err != nil {
		return nil, err
	}

	var body HfsDb
	if err := json.Unmarshal(resp, &body); err != nil {
		return nil, err
	}

	//报错了
	if body.Detail != "" {
		return nil, errors.New(body.Detail)
	}
	return &body, nil
}

// 读取token
func getHfsUserToken() (string, error) {
	uri := `https://www.simui.net/api/getHsfToken.php`
	tm := strconv.Itoa(int(time.Now().Unix()))
	salt := utils.ToString(utils.RandInt(8, 8))
	secret := utils.RandStr(16, 16)
	from := "simui"
	sign := utils.Md5(from + tm + salt + utils.Md5(secret+salt))
	req := map[string]string{
		"time":   tm,
		"salt":   salt,
		"secret": secret,
		"from":   from,
		"sign":   sign,
	}

	resp, err := utils.HttpPostForm(uri, req, nil)
	if err != nil {
		return "", err
	}

	respMap := map[string]string{}
	if err := json.Unmarshal(resp, &respMap); err != nil {
		return "", errors.New("token验证失败")
	}
	if _, ok := respMap["token"]; !ok {
		return "", errors.New("token验证失败")
	}
	if respMap["token"] == "" {
		return "", errors.New("token验证失败")
	}

	return respMap["token"], nil
}
