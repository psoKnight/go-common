package http_go

import (
	"net/http"
	"time"
)

// 标题元素
type headerElements struct {
	key   string
	value string
}

// SetConnectTimeout 设置连接超时
func SetConnectTimeout(duration time.Duration) {
	DefaultDialer.Timeout = duration
}

// AddHeader 添加header
func (r *Request) AddHeader(key string, value string) {
	if r.headers == nil {
		r.headers = []headerElements{}
	}
	r.headers = append(r.headers, headerElements{key: key, value: value})
}

// WithHeader 附加header
func (r Request) WithHeader(key string, value string) Request {
	r.AddHeader(key, value)
	return r
}

// AddCookie 添加cookie
func (r *Request) AddCookie(c *http.Cookie) {
	r.cookies = append(r.cookies, c)
}

// WithCookie 附加cookie
func (r Request) WithCookie(c *http.Cookie) Request {
	r.AddCookie(c)
	return r

}
