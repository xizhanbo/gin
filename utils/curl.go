package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

// @Description 支持任意方式的http请求
// @Param url 请求地址
// @Param method 请求方式
// @Param headers 请求header头设置
// @Param params 请求地址栏后需要拼接参数操作
// @Param data 请求报文
// @Return 返回类型 "错误信息"
func HttpDo(url, method string, params, headers map[string]string, data []byte) (interface{}, error) {

	//自定义cient
	client := &http.Client{
		Timeout: 5 * time.Second, // 超时时间：5秒
	}
	//http.post等方法只是在NewRequest上又封装来了一层而已
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, errors.New("new request is fail: %v \n")
	}
	req.Header.Set("Content-type", "application/json")

	//add params
	q := req.URL.Query()
	if params != nil {
		for key, val := range params {
			q.Add(key, val)
		}
		req.URL.RawQuery = q.Encode()
	}

	//add headers
	if headers != nil {
		for key, val := range headers {
			req.Header.Add(key, val)
		}
	}

	resp, err := client.Do(req) //  默认的resp ,err :=  http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	result := make(map[string]interface{}, 2)
	result["code"] = resp.StatusCode
	bodyData := make(map[string]interface{}, 0)
	err = json.Unmarshal(body, &bodyData)
	result["data"] = bodyData

	return result, nil
}
