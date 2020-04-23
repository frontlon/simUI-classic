package utils

import (
	"io/ioutil"
	"net/http"
	"time"
)

//get请求
func GetHttp(uri string) string {

	if uri == "" {
		return ""
	}

	client := &http.Client{}
	client.Timeout = 3 * time.Second
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return ""
	}

	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil || resp.StatusCode != 200 {
		return ""
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return string(body)
}
