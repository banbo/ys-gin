package test

import "net/url"

type APIParam struct {
	val url.Values
}

func NewAPIParam() *APIParam {
	return &APIParam{
		val: make(map[string][]string),
	}
}

func (p *APIParam) Add(key, value string) {
	p.val.Add(key, value)
}

func (p *APIParam) Encode() string {
	return p.val.Encode()
}

func (p *APIParam) Values() url.Values {
	return p.val
}

//实现 io.Reader 接口
//func (p *APIParam) Read(b []byte) (n int, err error) {
//	return 0, nil
//}
