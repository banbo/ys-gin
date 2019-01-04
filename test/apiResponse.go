package test

import (
	"encoding/json"
	"fmt"
)

type APIResponse struct {
	body []byte
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func NewAPIResponse(body []byte) (*APIResponse, error) {
	r := &APIResponse{body: body}
	err := json.Unmarshal(body, r)
	if err != nil {
		r = nil
		err = fmt.Errorf("响应结果错误：%s，响应原文：%s", err.Error(), string(body))
	}

	return r, err
}

//响应结果
func (r *APIResponse) Content() string {
	return string(r.body)
}
