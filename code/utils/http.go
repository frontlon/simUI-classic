package utils

import (
	"fmt"
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
		fmt.Println(err)
		return ""
	}

	resp, err := client.Do(req)
	if resp != nil {
		fmt.Println(err)
		defer resp.Body.Close()
	}
	if err != nil {
		fmt.Println(err)
		return ""
	}

	if resp.StatusCode != 200 {
		fmt.Println("StatusCode", resp.StatusCode)
		return ""
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return string(body)
}
