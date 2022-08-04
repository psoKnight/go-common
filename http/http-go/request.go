package http_go

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

var (
	// DefaultDialer 默认拨号器
	DefaultDialer = &net.Dialer{Timeout: 2000 * time.Millisecond}
	// DefaultTransport 默认传输
	DefaultTransport http.RoundTripper = &http.Transport{Dial: DefaultDialer.Dial, Proxy: http.ProxyFromEnvironment}
	// DefaultClient 默认客户端
	DefaultClient = &http.Client{Transport: DefaultTransport}
	// proxyTransport
	proxyTransport http.RoundTripper
	// proxyClient
	proxyClient *http.Client
)

// Request ghttp 请求
type Request struct {
	headers           []headerElements
	Method            string
	Url               string
	Body              interface{}
	Query             interface{}
	cookies           []*http.Cookie
	Timeout           time.Duration
	ContentType       string
	XForwardedFor     string
	Accept            string
	UserAgent         string
	Host              string
	Insecure          bool
	TlsConfig         *tls.Config
	MaxRedirects      int
	RedirectHeaders   bool
	Proxy             string
	Compression       *compression
	BasicAuthUsername string
	BasicAuthPassword string
	CookieJar         http.CookieJar
	ShowDebug         bool
	OnBeforeRequest   func(goxhttp *Request, httpreq *http.Request)
}

// 请求取消器
type transportRequestCanceler interface {
	CancelRequest(*http.Request)
}

// NewRequest 新建请求
func (r Request) NewRequest() (*http.Request, error) {
	b, e := prepareRequestBody(r.Body)
	if e != nil {
		return nil, &Error{Err: e}
	}
	if r.Query != nil {
		param, e := paramParse(r.Query)
		if e != nil {
			return nil, &Error{Err: e}
		}
		// http://127.0.0.1?user={}
		r.Url = r.Url + "?" + param
	}

	var (
		bodyReader io.Reader
	)
	if b != nil && r.Compression != nil {
		buffer := bytes.NewBuffer([]byte{})
		readBuffer := bufio.NewReader(b)
		writer, err := r.Compression.writer(buffer)
		if err != nil {
			return nil, &Error{Err: err}
		}
		_, e = readBuffer.WriteTo(writer)
		writer.Close()
		if e != nil {
			return nil, &Error{Err: e}
		}
		bodyReader = buffer
	} else {
		bodyReader = b
	}

	req, err := http.NewRequest(r.Method, r.Url, bodyReader)
	if err != nil {
		return nil, err
	}
	// 向请求中添加header
	req.Host = r.Host

	r.addHeaders(req.Header)
	if r.Compression != nil {
		req.Header.Add("Content-Encoding", r.Compression.ContentEncoding)
		req.Header.Add("Accept-Encoding", r.Compression.ContentEncoding)
	}
	if r.headers != nil {
		for _, header := range r.headers {
			req.Header.Add(header.key, header.value)
		}
	}

	// 使用基本身份验证auth
	if r.BasicAuthUsername != "" {
		req.SetBasicAuth(r.BasicAuthUsername, r.BasicAuthPassword)
	}

	for _, c := range r.cookies {
		req.AddCookie(c)
	}
	return req, nil
}

// Do 发起请求
func (r Request) Do() (*Response, error) {
	var (
		client         = DefaultClient
		transport      = DefaultTransport
		resURL         string
		redirectFailed bool
	)

	r.Method = valueOrDefault(r.Method, "GET")

	// 使用旧的 cookiejar
	if r.CookieJar != nil {
		client = &http.Client{
			Transport: transport,
			Jar:       r.CookieJar,
		}
	}

	if len(r.Proxy) != 0 {
		proxyUrl, err := url.Parse(r.Proxy)
		if err != nil {
			return nil, &Error{Err: err}
		}

		// 如果指定，则需要重新构建
		if proxyTransport == nil || client.Jar != nil {
			proxyTransport = &http.Transport{Dial: DefaultDialer.Dial, Proxy: http.ProxyURL(proxyUrl)}
			proxyClient = &http.Client{Transport: proxyTransport, Jar: client.Jar}
		} else if proxyTransport, ok := proxyTransport.(*http.Transport); ok {
			proxyTransport.Proxy = http.ProxyURL(proxyUrl)
		}
		transport = proxyTransport
		client = proxyClient
	}

	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if len(via) > r.MaxRedirects {
			redirectFailed = true
			return errors.New("error redirecting, max redirects reached")
		}
		resURL = req.URL.String()
		// 默认不会重定向请求头，重新设置
		if r.RedirectHeaders {
			for key, val := range via[0].Header {
				req.Header[key] = val
			}
		}
		return nil
	}

	if transport, ok := transport.(*http.Transport); ok {
		if r.Insecure {
			if r.TlsConfig != nil {
				transport.TLSClientConfig = r.TlsConfig
			} else {
				transport.TLSClientConfig = &tls.Config{
					InsecureSkipVerify: true,
				}
			}
		} else if transport.TLSClientConfig != nil {
			// 默认
			transport.TLSClientConfig.InsecureSkipVerify = false
		}
	}

	req, err := r.NewRequest()
	if err != nil {
		return nil, &Error{Err: err}
	}

	timeout := false
	if r.Timeout > 0 {
		client.Timeout = r.Timeout
	}

	if r.ShowDebug {
		dump, err := httputil.DumpRequest(req, true)
		if err != nil {
			log.Println(err)
		}
		log.Println(string(dump))
	}

	if r.OnBeforeRequest != nil {
		r.OnBeforeRequest(&r, req)
	}
	res, err := client.Do(req)
	if err != nil {
		if t, ok := err.(reqtimeout); ok {
			timeout = t.Timeout()
		}
		if ue, ok := err.(*url.Error); ok {
			if t, ok := ue.Err.(reqtimeout); ok {
				timeout = t.Timeout()
			}
		}

		var (
			response *Response
		)
		// 重定向失败时的响应
		if redirectFailed {
			if res != nil {
				response = &Response{
					res,
					resURL,
					&Body{
						reader: res.Body,
					},
					req,
				}
			} else {
				response = &Response{
					res,
					resURL,
					nil,
					req,
				}
			}
		}

		// 重定向失败或者最大重定向数为0，返回 nil（无错误）
		if redirectFailed && r.MaxRedirects == 0 {
			return response, nil
		}
		return response, &Error{
			timeout: timeout,
			Err:     err,
		}
	}

	if r.Compression != nil && strings.Contains(res.Header.Get("Content-Encoding"), r.Compression.ContentEncoding) {
		compressedReader, err := r.Compression.reader(res.Body)
		if err != nil {
			return nil, &Error{Err: err}
		}
		return &Response{
				res, resURL,
				&Body{
					reader:           res.Body,
					compressedReader: compressedReader,
				},
				req,
			},
			nil
	}
	return &Response{
			res,
			resURL,
			&Body{
				reader: res.Body,
			},
			req,
		},
		nil
}

// CancelRequest 取消请求
func (r Response) CancelRequest() {
	cancelRequest(DefaultTransport, r.req)
}

// 取消请求
func cancelRequest(transport interface{}, r *http.Request) {
	if tp, ok := transport.(transportRequestCanceler); ok {
		tp.CancelRequest(r)
	}
}
