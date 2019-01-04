package test

import (
	"testing"

	"net/http"
	"github.com/banbo/ys-gin/constant"
)

func Test_Test_Add(t *testing.T) {
	c := getAPIClient()

	// 请求参数
	c.AddParam("name", "test1")
	c.AddParam("age", "18")
	c.AddParam("sign", "1c2cec0b1f799d8dffc649ee60d3e7634e153fc5")

	// cookie
	cookies := []*http.Cookie{
		LOGIN_COOKIE,
	}

	// 发送
	resp, err := c.Post("test/add", cookies)
	if err != nil {
		t.Fatal(err)
	}

	// 判断错误码
	if resp.Code != constant.RESPONSE_CODE_OK {
		t.Fatalf(resp.Msg)
	}

	// 记下新增的生成的id，供下面的测试用例使用
	data, _ := resp.Data.(map[string]interface{})
	lastUID = data["uid"].(string)
}

func Test_Test_List(t *testing.T) {
	c := getAPIClient()

	// 请求参数
	c.AddParam("is_page", "true")
	c.AddParam("page_index", "1")
	c.AddParam("page_size", "10")

	// cookie
	cookies := []*http.Cookie{
		LOGIN_COOKIE,
	}

	// 发送
	resp, err := c.Get("test/list", cookies)
	if err != nil {
		t.Fatal(err)
	}

	// 判断错误码
	if resp.Code != constant.RESPONSE_CODE_OK {
		t.Fatalf(resp.Msg)
	}

	t.Log(resp.Data)
}

func Test_Test_Get(t *testing.T) {
	c := getAPIClient()

	// 请求参数
	c.AddParam("uid", lastUID)

	// cookie
	cookies := []*http.Cookie{
		LOGIN_COOKIE,
	}

	// 发送
	resp, err := c.Get("test/get", cookies)
	if err != nil {
		t.Fatal(err)
	}

	// 判断错误码
	if resp.Code != constant.RESPONSE_CODE_OK {
		t.Fatalf(resp.Msg)
	}

	t.Log(resp.Data)
}
