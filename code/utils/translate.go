package utils

import (
	"encoding/json"
)

type BaiduTrans struct {
	From        string `json:"from"`
	To          string `json:"to"`
	TransResult []struct {
		Src string `json:"src"`
		Dst string `json:"dst"`
	} `json:"trans_result"`
}

// 百度翻译 to: zh/en/cht
func BaiduTranslate(keyword string, to string) string {
	postUrl := "https://fanyi-api.baidu.com/api/trans/vip/translate?"
	appId := "20230421001650804"
	scriet := "FzLUusAtVas4v1FJJphE"
	salt := ToString(RandInt(8, 8))
	sign := appId + keyword + salt + scriet

	data := map[string]string{
		"q":     keyword,
		"from":  "auto",
		"to":    to,
		"appid": appId,
		"salt":  salt,
		"sign":  Md5(sign),
	}
	resp, err := HttpPostForm(postUrl, data, nil)
	if err != nil {
		WriteLog("BaiduTranslate Error:" + err.Error())
		return ""
	}

	var body BaiduTrans
	json.Unmarshal(resp, &body)
	if len(body.TransResult) == 0 {
		return ""
	}

	return body.TransResult[0].Dst
}
