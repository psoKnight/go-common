package http_go

import (
	"compress/gzip"
	"compress/zlib"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

// Response ghttp 响应
type Response struct {
	*http.Response
	Url  string
	Body *Body
	req  *http.Request
}

type Body struct {
	reader           io.ReadCloser
	compressedReader io.ReadCloser
}

// Read 读
func (b *Body) Read(p []byte) (int, error) {
	if b.compressedReader != nil {
		return b.compressedReader.Read(p)
	}
	return b.reader.Read(p)
}

// Close 关闭
func (b *Body) Close() error {
	err := b.reader.Close()
	if b.compressedReader != nil {
		return b.compressedReader.Close()
	}
	return err
}

// ToJson 解码为json 格式
func (b *Body) ToJson(o interface{}) error {
	return json.NewDecoder(b).Decode(o)
}

// ToString 解码为string 格式
func (b *Body) ToString() (string, error) {
	body, err := ioutil.ReadAll(b)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// Gzip gzip 压缩
func Gzip() *compression {
	reader := func(buffer io.Reader) (io.ReadCloser, error) {
		return gzip.NewReader(buffer)
	}
	writer := func(buffer io.Writer) (io.WriteCloser, error) {
		return gzip.NewWriter(buffer), nil
	}
	return &compression{
		writer:          writer,
		reader:          reader,
		ContentEncoding: "gzip",
	}
}

// Deflate deflate 压缩
func Deflate() *compression {
	reader := func(buffer io.Reader) (io.ReadCloser, error) {
		return zlib.NewReader(buffer)
	}
	writer := func(buffer io.Writer) (io.WriteCloser, error) {
		return zlib.NewWriter(buffer), nil
	}
	return &compression{
		writer:          writer,
		reader:          reader,
		ContentEncoding: "deflate",
	}
}

func Zlib() *compression {
	return Deflate()
}
