package xhttp

import (
	"context"
	"net/http"
	"net/http/httptest"
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

// TestDoWithContext 测试context取消/超时
func TestDoWithContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	_, _, err := DefaultClient.DoWithContext(ctx, MethodGet, server.URL, nil)
	if err == nil {
		t.Fatal("expected context deadline error")
	}
}

// TestDo_FailOnHttpError 测试>=400是否视为错误
func TestDo_FailOnHttpError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	// 默认不视为错误
	resp, _, err := DefaultClient.Do(MethodGet, server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}

	// 开启后视为错误
	_, _, err = DefaultClient.Do(MethodGet, server.URL, nil, WithFailOnHttpError())
	if err == nil {
		t.Fatal("expected http status error")
	}
}
