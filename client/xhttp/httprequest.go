package xhttp

import (
	"context"
	"fmt"
	"io"
	"net/http"
	netUrl "net/url"
	"strings"
	"time"

	"github.com/banbo/ys-gin/errors"
)

// HttpMethod http method
type HttpMethod string

const (
	MethodGet     HttpMethod = http.MethodGet
	MethodHead    HttpMethod = http.MethodHead
	MethodPost    HttpMethod = http.MethodPost
	MethodPut     HttpMethod = http.MethodPut
	MethodPatch   HttpMethod = http.MethodPatch // RFC 5789
	MethodDelete  HttpMethod = http.MethodDelete
	MethodConnect HttpMethod = http.MethodConnect
	MethodOptions HttpMethod = http.MethodOptions
	MethodTrace   HttpMethod = http.MethodTrace
)

// MimeType mime type
type MimeType string

const (
	MimeTextHtml                      MimeType = "text/html"
	MimeTextPlain                     MimeType = "text/plain"
	MimeTextXml                       MimeType = "text/xml"
	MimeApplicationJson               MimeType = "application/json"
	MimeApplicationXWwwFormUrlencoded MimeType = "application/x-www-form-urlencoded"
	MimeMultipartFormData             MimeType = "multipart/form-data"
)

// HttpClient http 客户端（导出以便命名/缓存；Do/DoStream 等均协程安全）
type HttpClient struct {
	client http.Client
}

// NewHttpClient 创建新的http client
// 创建client应该复用而不是每次创建，创建的client是协程安全的（详见http.Client注释说明）
func NewHttpClient() *HttpClient {
	return new(HttpClient)
}

// WithTransport client transport
func (c *HttpClient) WithTransport(transport http.RoundTripper) *HttpClient {
	c.client.Transport = transport
	return c
}

// WithHTTP2 显式启用HTTP/2（同时保留HTTP/1.1）
// 自定义Transport若设置了Dial/DialTLS/TLSClientConfig，默认会关闭HTTP/2，
// 此时可用本方法重新开启；零值Transport本身已支持HTTP/2，无需调用。
func (c *HttpClient) WithHTTP2() *HttpClient {
	c.ensureTransport()
	if t, ok := c.client.Transport.(*http.Transport); ok {
		if t.Protocols == nil {
			t.Protocols = new(http.Protocols)
		}
		t.Protocols.SetHTTP1(true)
		t.Protocols.SetHTTP2(true)
	}
	return c
}

// ensureTransport 保证client拥有可修改的transport
func (c *HttpClient) ensureTransport() {
	if c.client.Transport != nil {
		return
	}
	if t, ok := http.DefaultTransport.(*http.Transport); ok {
		c.client.Transport = t.Clone()
	} else {
		c.client.Transport = &http.Transport{}
	}
}

// WithCheckRedirect client checkRedirect
func (c *HttpClient) WithCheckRedirect(checkRedirect func(req *http.Request, via []*http.Request) error) *HttpClient {
	c.client.CheckRedirect = checkRedirect
	return c
}

// WithJar client jar
func (c *HttpClient) WithJar(jar http.CookieJar) *HttpClient {
	c.client.Jar = jar
	return c
}

// WithTimeOut client timeout
func (c *HttpClient) WithTimeOut(timeout time.Duration) *HttpClient {
	c.client.Timeout = timeout
	return c
}

// Do 发http请求（读完整包，返回响应与响应体字节）
func (c *HttpClient) Do(method HttpMethod, url string, body io.Reader, options ...OptionFn) (*http.Response, []byte,
	error) {
	return c.DoWithContext(context.Background(), method, url, body, options...)
}

// DoWithContext 发http请求，支持context取消/超时（读完整包）
func (c *HttpClient) DoWithContext(ctx context.Context, method HttpMethod, url string, body io.Reader,
	options ...OptionFn) (*http.Response, []byte, error) {
	request, opt, err := c.buildRequest(ctx, method, url, body, options...)
	if err != nil {
		return nil, nil, err
	}
	response, err := c.client.Do(request)
	if err != nil {
		if response != nil {
			response.Body.Close()
		}
		return nil, nil, errors.NewSys(err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return response, nil, errors.NewSys(err)
	}
	// 可选：>=400 视为错误
	if opt.failOnHttpError && response.StatusCode >= http.StatusBadRequest {
		return response, responseBody, errors.NewSys(fmt.Sprintf("http status code: %d", response.StatusCode))
	}
	return response, responseBody, nil
}

// DoStream 发http请求并返回未读取的响应（Body 未关闭，由调用方负责读取/关闭）。
// 用于 SSE / 反向代理等流式场景；普通请求用 Do/DoWithContext。
func (c *HttpClient) DoStream(method HttpMethod, url string, body io.Reader, options ...OptionFn) (*http.Response,
	error) {
	return c.DoStreamWithContext(context.Background(), method, url, body, options...)
}

// DoStreamWithContext 发http请求并返回未读取的响应，支持context取消/超时。
// Body 由调用方负责读取/关闭。
func (c *HttpClient) DoStreamWithContext(ctx context.Context, method HttpMethod, url string, body io.Reader,
	options ...OptionFn) (*http.Response, error) {
	request, _, err := c.buildRequest(ctx, method, url, body, options...)
	if err != nil {
		return nil, err
	}
	response, err := c.client.Do(request)
	if err != nil {
		if response != nil {
			response.Body.Close()
		}
		return nil, errors.NewSys(err)
	}
	return response, nil
}

// buildRequest 组装 *http.Request（query/headers/透传头/自定义钩子）
func (c *HttpClient) buildRequest(ctx context.Context, method HttpMethod, url string, body io.Reader,
	options ...OptionFn) (*http.Request, *Option, error) {
	opt := &Option{
		headers:     make(map[string]string),
		queryParams: make(map[string]string),
	}
	for _, option := range options {
		option(opt)
	}

	// url加入query参数
	if len(opt.queryParams) > 0 {
		query := (&netUrl.URL{}).Query()
		for key, value := range opt.queryParams {
			query.Add(key, value)
		}
		if strings.Contains(url, "?") {
			url += "&" + query.Encode()
		} else {
			url += "?" + query.Encode()
		}
	}

	request, err := http.NewRequestWithContext(ctx, string(method), url, body)
	if err != nil {
		return nil, nil, errors.NewSys(err)
	}

	// 单值 header
	for key, value := range opt.headers {
		request.Header.Add(key, value)
	}
	// 多值 header 透传（反代场景）
	for key, values := range opt.requestHeader {
		for _, value := range values {
			request.Header.Add(key, value)
		}
	}
	// 自定义钩子（反代做最后加工：去 hop-by-hop、改 Host 等）
	if opt.requestHook != nil {
		opt.requestHook(request)
	}
	return request, opt, nil
}

// http请求参数
type Option struct {
	// headers
	headers map[string]string

	// query参数
	queryParams map[string]string

	// 透传的原始多值 header
	requestHeader http.Header

	// 自定义请求钩子
	requestHook func(*http.Request)

	// 是否将HTTP状态码>=400视为错误
	failOnHttpError bool
}

// OptionFn
type OptionFn func(*Option)

// WithHeaders headers
func WithHeaders(headers map[string]string) OptionFn {
	return func(o *Option) {
		for key, value := range headers {
			o.headers[key] = value
		}
	}
}

// WithHeader add headers one by one
func WithHeader(key string, value string) OptionFn {
	return func(o *Option) {
		o.headers[key] = value
	}
}

// WithQueryParams queryParams
func WithQueryParams(queryParams map[string]string) OptionFn {
	return func(o *Option) {
		for key, value := range queryParams {
			o.queryParams[key] = value
		}
	}
}

// WithQueryParam add queryParams one by one
func WithQueryParam(key string, value string) OptionFn {
	return func(o *Option) {
		o.queryParams[key] = value
	}
}

// WithRequestHeader 透传原始多值 http.Header（反向代理等）
func WithRequestHeader(h http.Header) OptionFn {
	return func(o *Option) {
		for key, values := range h {
			o.requestHeader = appendHeader(o.requestHeader, key, values)
		}
	}
}

// WithRequestHook 对组装好的 *http.Request 做最后加工（去 hop-by-hop 头、改 Host 等）
func WithRequestHook(fn func(*http.Request)) OptionFn {
	return func(o *Option) {
		o.requestHook = fn
	}
}

// WithFailOnHttpError 当HTTP状态码>=400时返回错误
func WithFailOnHttpError() OptionFn {
	return func(o *Option) {
		o.failOnHttpError = true
	}
}

// appendHeader 向 header 追加多值（惰性初始化）
func appendHeader(h http.Header, key string, values []string) http.Header {
	if h == nil {
		h = make(http.Header, 1)
	}
	for _, v := range values {
		h.Add(key, v)
	}
	return h
}
