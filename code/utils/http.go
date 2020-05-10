package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

//get请求
func GetHttp(uri string) (string, error) {

	if uri == "" {
		return "", nil
	}

	client := &http.Client{}
	client.Timeout = 3 * time.Second
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	resp, err := client.Do(req)
	if resp != nil {
		fmt.Println(err)
		defer resp.Body.Close()
	}
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	if resp.StatusCode != 200 {
		fmt.Println("StatusCode", resp.StatusCode)
		return "", errors.New(ToString(resp.StatusCode))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
