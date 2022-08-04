package http_go

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestSimpleRequest(t *testing.T) {
	req := Request{
		Url:           "https://www.baidu.com",
		Method:        "GET",
		XForwardedFor: "127.0.0.1",
		ContentType:   "application/json",
	}
	resp, err := req.Do()
	if err != nil {
		t.Error(err)
		return
	}
	defer resp.Body.Close()
	result, err := resp.Body.ToString()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("Body to string: %s.", result)
}

func TestQueryStruct(t *testing.T) {
	// request-> Get http://127.0.0.1:8080?name=xc&password=xc
	type User struct {
		Name     string
		Password string
	}
	user := User{
		Name:     "xc",
		Password: "xc",
	}
	res, err := Request{
		Url:   "http://127.0.0.1:8080",
		Query: user,
	}.Do()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res)
}

func TestQueryJson(t *testing.T) {
	// Get http://127.0.0.1:8080?name=xc
	type User struct {
		Name     string `json:"name"`
		Password string `json:"-"`
		Sex      string `json:"sex,omitempty"`
	}
	user := User{
		Name:     "xc",
		Password: "xc",
		Sex:      "",
	}
	res, err := Request{
		Url:   "http://127.0.0.1:8080",
		Query: user,
	}.Do()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res)
}

func TestQuerySquash(t *testing.T) {
	// Get http://127.0.0.1:8080?id=1&name=xc&password=xc&sex=1
	type User struct {
		Name     string `json:"name"`
		Password string `json:"password"`
		Sex      string `json:"sex"`
	}
	type GoPher struct {
		Id   int32 `json:"id"`
		User `json:",squash"`
	}
	gopher := GoPher{
		User: User{
			Name:     "xc",
			Password: "xc",
			Sex:      "1",
		},
		Id: 1,
	}
	fmt.Println(gopher)
	res, err := Request{
		Url:   "http://127.0.0.1:8080",
		Query: gopher,
	}.Do()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res)
}

func TestHeader(t *testing.T) {
	// Post http://127.0.0.1:8080?name=xc&password=xc
	type User struct {
		Name     string
		Password string
	}
	user := User{
		Name:     "xc",
		Password: "xc",
	}
	req := &Request{
		Method:      "POST",
		Url:         "http://127.0.0.1:8080",
		Query:       user,
		ContentType: "application/json",
	}
	req.AddHeader("X-Custom", "haha")
	res, err := req.Do()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res)
}

func TestCookie(t *testing.T) {
	res, err := Request{
		Url: "http://www.baidu.com",
	}.WithCookie(&http.Cookie{Name: "c1", Value: "v1"}).Do()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res)
}

func TestSetTimeOut(t *testing.T) {
	res, err := Request{
		Url:     "http://www.baidu.com",
		Timeout: 100 * time.Millisecond,
	}.Do()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res)
}

func TestResToJson(t *testing.T) {
	// 解析json，转换为相应结构体
	type User struct {
		Name     string `json:"name"`
		Password string `json:"password"`
		Sex      string `json:"sex"`
	}
	var user User
	res, err := Request{
		Url:     "http://127.0.0.1:8080",
		Timeout: 100 * time.Millisecond,
	}.Do()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res.Body.ToJson(&user))
}

func TestProxy(t *testing.T) {
	res, err := Request{
		Url:     "http://127.0.0.1:8080",
		Timeout: 100 * time.Millisecond,
		Proxy:   "http://127.0.0.1:8088",
	}.Do()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res)
}

func TestDebug(t *testing.T) {
	res, err := Request{
		Url:       "http://127.0.0.1:8080",
		Timeout:   100 * time.Millisecond,
		ShowDebug: true,
	}.Do()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res)
}

func TestGzip(t *testing.T) {
	res, err := Request{
		Method:      "POST",
		Url:         "http://www.baidu.com",
		Compression: Gzip(),
	}.Do()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res)
}

func TestDeflate(t *testing.T) {
	res, err := Request{
		Method:      "POST",
		Url:         "http://www.baidu.com",
		Compression: Deflate(),
	}.Do()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res)
}
