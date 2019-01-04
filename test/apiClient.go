package test

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
)

type APIClient struct {
	ht *httptest.Server
	pm *APIParam
}

func NewAPIClient(handler http.Handler) *APIClient {
	return &APIClient{
		ht: httptest.NewServer(handler),
		pm: NewAPIParam(),
	}
}

func (c *APIClient) AddParam(key, value string) {
	c.pm.Add(key, value)
}

func (c *APIClient) Get(url string, cookies []*http.Cookie) (*APIResponse, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", c.ht.URL+"/"+url+"?"+c.pm.Encode(), nil)
	if err != nil {
		err = errors.New("请求错误：" + err.Error())
		return nil, err
	}

	// 加上cookies
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	httpRes, err := client.Do(req)
	if err != nil {
		err = errors.New("请求错误：" + err.Error())
		return nil, err
	}
	if httpRes.StatusCode != 200 {
		return nil, fmt.Errorf("请求错误：%d", httpRes.StatusCode)
	}

	return getAPIResponse(httpRes.Body)
}

func (c *APIClient) Post(url string, cookies []*http.Cookie) (*APIResponse, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", c.ht.URL+"/"+url, strings.NewReader(c.pm.Encode()))
	if err != nil {
		err = errors.New("请求错误：" + err.Error())
		return nil, err
	}

	// content-type
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 加上cookies
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	httpRes, err := client.Do(req)
	if err != nil {
		err = errors.New("请求错误：" + err.Error())
		return nil, err
	}
	if httpRes.StatusCode != 200 {
		return nil, fmt.Errorf("请求错误：%d", httpRes.StatusCode)
	}

	return getAPIResponse(httpRes.Body)
}

func getAPIResponse(r io.ReadCloser) (*APIResponse, error) {
	body, err := ioutil.ReadAll(r)
	r.Close()
	if err != nil {
		err = errors.New("读取结果错误：" + err.Error())
		return nil, err
	}

	return NewAPIResponse(body)
}
