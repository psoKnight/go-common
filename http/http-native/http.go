package http_native

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type HttpConfig struct {
	Transport     http.RoundTripper
	CheckRedirect func(req *http.Request, via []*http.Request) error
	Jar           http.CookieJar
	Timeout       time.Duration
}

type Http struct {
	cli *http.Client
	cfg *HttpConfig // 基础配置
}

type UploadFile struct {
	// 表单名称
	Name string
	// 文件全路径
	FilePath string
}

func NewHttp(cfg *HttpConfig) *Http {
	c := &Http{
		cfg: cfg,
	}
	cli := &http.Client{
		Transport:     cfg.Transport,
		CheckRedirect: cfg.CheckRedirect,
		Jar:           cfg.Jar,
		Timeout:       cfg.Timeout,
	}
	c.cli = cli

	return c
}

// GetClient 获取http client
func (hc *Http) GetClient() *http.Client {
	return hc.cli
}

// Get HTTP "GET" 请求
/**
reqUrl：待请求的url
reqParams：http req 参数
headers：请求头
*/
func (hc *Http) Get(reqUrl string, reqParams map[string]string, headers map[string]string) ([]byte, error) {

	urlParams := url.Values{}
	u, err := url.Parse(reqUrl)
	if err != nil {
		return []byte{}, err
	}

	for key, val := range reqParams {
		urlParams.Set(key, val)
	}

	// 如果参数中有中文参数，这个方法会进行URL encoded
	u.RawQuery = urlParams.Encode()

	// 得到完整的url，示例：http://xx?query
	urlPath := u.String()

	httpRequest, err := http.NewRequest("GET", urlPath, nil)
	if err != nil {
		return []byte{}, err
	}

	// 添加请求头
	if headers != nil {
		for k, v := range headers {
			httpRequest.Header.Add(k, v)
		}
	}

	// 发送请求
	resp, err := hc.cli.Do(httpRequest)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return res, nil

}

// PostJson HTTP "POST" 请求（Content-Type 是application/json）
/**
reqUrl：待请求的url
bodyParams：http body 参数
headers：请求头
*/
func (hc *Http) PostJson(reqUrl string, bodyParams interface{}, headers map[string]string) ([]byte, error) {
	return hc.post(reqUrl, bodyParams, nil, "application/json", nil, headers)
}

// PostForm HTTP "POST" 请求（Content-Type 是application/x-www-form-urlencoded）
/**
reqUrl：待请求的url
reqParams：http req 参数
headers：请求头
*/
func (hc *Http) PostForm(reqUrl string, reqParams map[string]string, headers map[string]string) ([]byte, error) {
	return hc.post(reqUrl, nil, reqParams, "application/x-www-form-urlencoded", nil, headers)
}

// PostFile HTTP "POST" 请求，上传文件（Content-Type 是multipart/form-data）
/**
reqUrl：待请求的url
reqParams：http req 参数
files：待上传的文件
headers：请求头
*/
func (hc *Http) PostFile(reqUrl string, reqParams map[string]string, files []UploadFile, headers map[string]string) ([]byte, error) {
	return hc.post(reqUrl, nil, reqParams, "multipart/form-data", files, headers)
}

func (hc *Http) post(reqUrl string, bodyParams interface{}, reqParams map[string]string, contentType string, files []UploadFile, headers map[string]string) ([]byte, error) {
	requestBody, realContentType, err := getReader(bodyParams, reqParams, contentType, files)
	if err != nil {
		return []byte{}, err
	}

	httpRequest, err := http.NewRequest("POST", reqUrl, requestBody)
	if err != nil {
		return []byte{}, err
	}

	// 添加请求头
	httpRequest.Header.Add("Content-Type", realContentType)
	if headers != nil {
		for k, v := range headers {
			httpRequest.Header.Add(k, v)
		}
	}

	// 发送请求
	resp, err := hc.cli.Do(httpRequest)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return response, nil
}

func getReader(reqJson interface{}, reqParams map[string]string, contentType string, files []UploadFile) (io.Reader, string, error) {
	if strings.Index(contentType, "json") > -1 {
		bytesData, err := json.Marshal(reqJson)
		if err != nil {
			return bytes.NewReader([]byte{}), "", err
		}

		return bytes.NewReader(bytesData), contentType, nil
	} else if files != nil {
		body := &bytes.Buffer{}

		// 文件写入 body
		writer := multipart.NewWriter(body)

		for _, uploadFile := range files {
			file, err := os.Open(uploadFile.FilePath)
			if err != nil {
				return bytes.NewReader([]byte{}), "", err
			}
			part, err := writer.CreateFormFile(uploadFile.Name, filepath.Base(uploadFile.FilePath))
			if err != nil {
				return bytes.NewReader([]byte{}), "", err
			}
			_, err = io.Copy(part, file)
			if err != nil {
				return bytes.NewReader([]byte{}), "", err
			}

			file.Close()
		}

		// 其他参数列表写入 body
		for k, v := range reqParams {
			if err := writer.WriteField(k, v); err != nil {
				if err != nil {
					return bytes.NewReader([]byte{}), "", err
				}
			}
		}

		if err := writer.Close(); err != nil {
			if err != nil {
				return bytes.NewReader([]byte{}), "", err
			}
		}

		// 上传文件需要自己专用的contentType
		return body, writer.FormDataContentType(), nil
	} else {
		urlValues := url.Values{}

		for key, val := range reqParams {
			urlValues.Set(key, val)
		}

		reqBody := urlValues.Encode()

		return strings.NewReader(reqBody), contentType, nil
	}
}

// Delete HTTP "Delete" 请求
/**
reqUrl：待请求的url
reqParams：http req 参数
bodyParams：http body 参数
headers：请求头
*/
func (hc *Http) Delete(reqUrl string, reqParams map[string]string, bodyParams map[string]string, headers map[string]string) ([]byte, error) {

	urlParams := url.Values{}
	u, err := url.Parse(reqUrl)
	if err != nil {
		return []byte{}, err
	}

	// 添加req params
	for key, val := range reqParams {
		urlParams.Set(key, val)
	}
	// 如果参数中有中文参数，这个方法会进行URL encoded
	u.RawQuery = urlParams.Encode()
	// 得到完整的url，示例：http://xx?query
	urlPath := u.String()

	// 添加body params
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	for bodyParamK, bodyParamV := range bodyParams {
		_ = writer.WriteField(bodyParamK, bodyParamV)
	}
	err = writer.Close()
	if err != nil {
		return []byte{}, err
	}

	httpRequest, err := http.NewRequest("DELETE", urlPath, payload)
	if err != nil {
		return []byte{}, err
	}

	// 添加请求头
	if headers != nil {
		for k, v := range headers {
			httpRequest.Header.Add(k, v)
		}
	}

	// 发送请求
	resp, err := hc.cli.Do(httpRequest)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return res, nil

}
