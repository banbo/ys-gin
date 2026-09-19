package xhttp

import (
	"context"
	"io"
	"net/http"
	"time"
)

// DefaultClient 带默认超时的http client
var DefaultClient = defaultHttpClient{
	http.Client{
		Timeout: 15 * time.Second,
	},
}

// default http client
type defaultHttpClient httpClient

// Do httpClient.Do()
func (c *defaultHttpClient) Do(method HttpMethod, url string, body io.Reader, options ...OptionFn) (*http.Response,
	[]byte, error) {
	return (*httpClient)(c).Do(method, url, body, options...)
}

// DoWithContext httpClient.DoWithContext()
func (c *defaultHttpClient) DoWithContext(ctx context.Context, method HttpMethod, url string, body io.Reader,
	options ...OptionFn) (*http.Response, []byte, error) {
	return (*httpClient)(c).DoWithContext(ctx, method, url, body, options...)
}
