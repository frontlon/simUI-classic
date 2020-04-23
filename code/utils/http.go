package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

//get请求
func GetHttp(uri string) string {
	resp, err := http.Get(uri)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return ""
	}
	return string(body)
}
