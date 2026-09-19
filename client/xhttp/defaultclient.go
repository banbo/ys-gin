package xhttp

import (
	"context"
	"io"
	"net/http"
	"time"
)

// DefaultClient 带默认超时的http client
var DefaultClient = defaultHttpClient{
	client: http.Client{Timeout: 15 * time.Second},
}

// default http client
type defaultHttpClient HttpClient

// Do httpClient.Do()
func (c *defaultHttpClient) Do(method HttpMethod, url string, body io.Reader, options ...OptionFn) (*http.Response,
	[]byte, error) {
	return (*HttpClient)(c).Do(method, url, body, options...)
}

// DoWithContext httpClient.DoWithContext()
func (c *defaultHttpClient) DoWithContext(ctx context.Context, method HttpMethod, url string, body io.Reader,
	options ...OptionFn) (*http.Response, []byte, error) {
	return (*HttpClient)(c).DoWithContext(ctx, method, url, body, options...)
}

// DoStream httpClient.DoStream()
func (c *defaultHttpClient) DoStream(method HttpMethod, url string, body io.Reader, options ...OptionFn) (*http.Response,
	error) {
	return (*HttpClient)(c).DoStream(method, url, body, options...)
}

// DoStreamWithContext httpClient.DoStreamWithContext()
func (c *defaultHttpClient) DoStreamWithContext(ctx context.Context, method HttpMethod, url string, body io.Reader,
	options ...OptionFn) (*http.Response, error) {
	return (*HttpClient)(c).DoStreamWithContext(ctx, method, url, body, options...)
}
