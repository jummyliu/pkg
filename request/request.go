package request

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	// "go.elastic.co/apm"
	// "go.elastic.co/apm/module/apmhttp"
	"golang.org/x/net/http2"
)

// Options request options
type Options struct {
	method   string
	params   map[string]string
	data     []byte
	headers  map[string]string
	timeout  int
	ssl      bool
	http2    bool
	ctx      context.Context
	client   *http.Client
	proxy    func(*http.Request) (*url.URL, error)
	redirect func(*http.Request, []*http.Request) error
}

// Option request option
type Option func(*Options)

func initOptions(options ...Option) *Options {
	opts := &Options{
		method:   http.MethodGet,
		params:   map[string]string{},
		data:     []byte{},
		headers:  map[string]string{"Content-Type": "application/json; charset=UTF-8"},
		timeout:  15,
		ssl:      true,
		http2:    false,
		ctx:      context.Background(),
		client:   nil,
		proxy:    nil,
		redirect: nil,
	}
	for _, option := range options {
		option(opts)
	}
	return opts
}

// WithOptions accepts the whole options config.
func WithOptions(options Options) Option {
	return func(opts *Options) {
		*opts = options
	}
}

// WithMethod set request method.
func WithMethod(method string) Option {
	return func(opts *Options) {
		opts.method = method
	}
}

// WithParams set request url params.
func WithParams(params map[string]string) Option {
	return func(opts *Options) {
		opts.params = params
	}
}

// WithData set request data.
func WithData(data []byte) Option {
	return func(opts *Options) {
		opts.data = data
	}
}

// WithHeader set request header.
func WithHeader(header map[string]string) Option {
	return func(opts *Options) {
		opts.headers = header
	}
}

// WithTimeout set request timeout.
func WithTimeout(timeout int) Option {
	return func(opts *Options) {
		opts.timeout = timeout
	}
}

// WithSSL set request skip ssl verify.
func WithSSL(ssl bool) Option {
	return func(opts *Options) {
		opts.ssl = ssl
	}
}

// WithContext set context.
func WithContext(ctx context.Context) Option {
	return func(opts *Options) {
		opts.ctx = ctx
	}
}

// WithClient use custom http.Client
//
//	client will invalidate WithSSL, WithTimeout, WithProxy
func WithClient(client *http.Client) Option {
	return func(opts *Options) {
		opts.client = client
	}
}

// WithProxy specifies a function to return a proxy for a given
// Request. If the function returns a non-nil error, the
// request is aborted with the provided error.
//
// The proxy type is determined by the URL scheme. "http",
// "https", and "socks5" are supported. If the scheme is empty,
// "http" is assumed.
//
// If Proxy is nil or returns a nil *URL, no proxy is used.
func WithProxy(proxy string) Option {
	return func(opts *Options) {
		if proxy == "" {
			opts.proxy = nil
			return
		}
		proxy_url, _ := url.Parse(proxy)
		opts.proxy = http.ProxyURL(proxy_url)
	}
}

// WithProxyFn specifies a function to return a proxy for a given
// Request. If the function returns a non-nil error, the
// request is aborted with the provided error.
//
// The proxy type is determined by the URL scheme. "http",
// "https", and "socks5" are supported. If the scheme is empty,
// "http" is assumed.
//
// If Proxy is nil or returns a nil *URL, no proxy is used.
func WithProxyFn(fn func(*http.Request) (*url.URL, error)) Option {
	return func(opts *Options) {
		opts.proxy = fn
	}
}

// WithHTTP2 配置是否开启 HTTP2
func WithHTTP2(http2 bool) Option {
	return func(opts *Options) {
		opts.http2 = http2
	}
}

// WithRedirect 配置最大允许跳转的次数，并返回最近的一个response
func WithRedirect(max_redirect_times int) Option {
	if max_redirect_times < 0 {
		return func(opts *Options) {}
	}
	return func(opts *Options) {
		opts.redirect = func(req *http.Request, via []*http.Request) error {
			if len(via) > max_redirect_times {
				return http.ErrUseLastResponse
			}
			return nil
		}
	}
}

// DoRequest exec https? request and return []byte
func DoRequest(request_url string, options ...Option) (code int, respBuf []byte, respHeader map[string][]string, err error) {
	// exec the undercourse request
	resp, err := DoRequestUndercourse(request_url, options...)
	if err != nil {
		// 错误
		return -1, nil, nil, errors.New("response failure")
	}
	respBuf, err = io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, resp.Header, errors.New("read response failure")
	}
	defer resp.Body.Close()

	return resp.StatusCode, respBuf, resp.Header, nil
}

// DoRequestUndercourse exec https? request and return response
func DoRequestUndercourse(request_url string, options ...Option) (resp *http.Response, err error) {
	opts := initOptions(options...)
	// span, ctx := apm.StartSpan(opts.ctx, "dorequest", "custom")
	// defer span.End()

	var req *http.Request

	if len(opts.params) != 0 {
		// 当附带请求参数时,判断是否以?结尾
		if !strings.HasSuffix(request_url, "?") {
			// 若末尾不为?且已包含?则自动拼接&,类似url=/api/test?query=aaaa
			if strings.Contains(request_url, "?") {
				request_url = request_url + "&"
				// 若末尾不为?且不包含?则自动拼接?,类似url=/api/test
			} else {
				request_url = request_url + "?"
			}
		}
		url_values := url.Values{}
		for key, val := range opts.params {
			url_values.Add(key, val)
		}
		// url后拼接urlencode的参数
		request_url = request_url + url_values.Encode()
	}

	req, err = http.NewRequest(opts.method, request_url, bytes.NewBuffer(opts.data))
	// switch opts.method {
	// case http.MethodPost, http.MethodPut:
	// 	req, err = http.NewRequest(opts.method, request_url, bytes.NewBuffer(opts.data))
	// default:
	// 	req, err = http.NewRequest(opts.method, request_url, nil)
	// }
	if err != nil {
		return nil, errors.New("build request failure")
	}

	for key, val := range opts.headers {
		req.Header.Set(key, val)
	}
	client := opts.client
	if opts.client == nil {
		client = NewHttpClient(opts)
	}
	// client = apmhttp.WrapClient(client)
	// resp, err = client.Do(req.WithContext(ctx))
	resp, err = client.Do(req.WithContext(opts.ctx))
	return resp, err
}

// 生成HttpClient对象
func NewHttpClient(opts *Options) (client_http *http.Client) {
	tr := &http.Transport{
		Proxy:             opts.proxy,
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: opts.ssl},
		DisableKeepAlives: true,
	}

	if opts.http2 {
		http2.ConfigureTransport(tr)
	}

	client_http = &http.Client{
		Transport:     tr,
		Timeout:       time.Duration(opts.timeout) * time.Second,
		CheckRedirect: opts.redirect,
	}
	return
}
