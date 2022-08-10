package seaweed_fs

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Seaweed interface {
	GetAssign(assignUrl string) ([]byte, error)
	PutObject(putUrl, filePath string) ([]byte, error)
	GetObject(getUrl string) ([]byte, error)
	RemoveObject(removeUrl string) error
}

type seaweedFs struct {
	serverURL   string
	httpTimeout time.Duration // HTTP 超时时间，默认30s 超时
}

// NewSeaweedFs New seaweed-fs 对象
func NewSeaweedFs(serverURL string, httpTimeout time.Duration) Seaweed {
	if httpTimeout == time.Duration(0) {
		httpTimeout = time.Second * time.Duration(30)
	}
	return &seaweedFs{serverURL: serverURL, httpTimeout: httpTimeout}
}

// GetAssign 获取分配
/**assignUrl：获取assign 的url
 */
func (s *seaweedFs) GetAssign(assignUrl string) ([]byte, error) {

	body, err := http.Get(assignUrl)
	if err != nil {
		return nil, err
	}
	defer body.Body.Close()
	return ioutil.ReadAll(body.Body)
}

// PutObject 存储对象
/**
putUrl：存储对象的url
filePath：对象目录
*/
func (s *seaweedFs) PutObject(putUrl, filePath string) ([]byte, error) {

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	part1, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part1, file)
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: s.httpTimeout}
	req, err := http.NewRequest("POST", putUrl, payload)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

// GetObject 获取对象
/**
getUrl：获取对象的url
*/
func (s *seaweedFs) GetObject(getUrl string) ([]byte, error) {
	body, err := http.Get(getUrl)
	if err != nil {
		return nil, err
	}
	defer body.Body.Close()
	return ioutil.ReadAll(body.Body)
}

// RemoveObject 删除对象
/**
removeUrl：移除对象的url
*/
func (s *seaweedFs) RemoveObject(removeUrl string) error {
	client := &http.Client{Timeout: s.httpTimeout}
	req, err := http.NewRequest("DELETE", removeUrl, nil)
	if err != nil {
		return err
	}

	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
