package httprequest

import (
	"net/http"
	"testing"
	"time"
)

// TestDo_Baidu
func TestDo_Baidu(t *testing.T) {
	_, responseData, err := DefaultClient.Do(MethodGet,
		"https://www.baidu.com",
		nil)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(responseData))
}

// TestDo_Concurrent 测试并发请求
func TestDo_Concurrent(t *testing.T) {
	num := 100
	chanData := make(chan error, num)

	client := NewHttpClient().
		WithTimeOut(10 * time.Second).
		WithTransport(&http.Transport{
			MaxIdleConnsPerHost: 10,
		})

	for i := 0; i < num; i++ {
		go func() {
			_, _, err := client.Do(MethodGet,
				"https://www.baidu.com",
				nil,
				WithHeader("Accept", string(MimeApplicationJson)),
				WithQueryParam("page", "1"),
				WithQueryParam("size", "1"))
			if err != nil {
				chanData <- err
				return
			}

			chanData <- nil
		}()
	}

	for i := 0; i < num; i++ {
		if err := <-chanData; err != nil {
			t.Fatal(err)
		}
	}
}
