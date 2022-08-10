package seaweed_fs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
	"time"
)

func TestSeaweed(t *testing.T) {

	// 获取fs
	fs := NewSeaweedFs("http://10.117.48.122:9333", time.Second*10)

	// 获取fs assign
	assign, err := fs.GetAssign("http://10.117.48.122:9333/dir/assign")
	//assign, err := fs.GetAssign("http://10.117.48.122:9333/dir/assign?ttl=3m&collenction=3mc")
	if err != nil {
		t.Errorf("Get assign err: %v.", err)
		return
	} else {
		t.Logf("Assign: %s.", string(assign))
	}
	assi := &Assign{}
	err = json.Unmarshal(assign, assi)
	if err != nil {
		t.Errorf("Json unmarshal err: %v.", err)
		return
	}

	// 存储对象
	putObject, err := fs.PutObject(fmt.Sprintf("http://%s/%s", assi.Url, assi.Fid), "./test.svg")
	if err != nil {
		t.Errorf("Put object err: %v.", err)
		return
	} else {
		t.Logf("Put object res: %s.", string(putObject))
	}

	p := &PutObjectRes{}
	err = json.Unmarshal(putObject, p)
	if err != nil {
		t.Errorf("Json unmarshal 'putObject' err: %v.", err)
		return
	}

	/**
	查看文件
	访问http://10.117.48.122:9222/3,0428f566d1
	*/

	// 获取对象
	object, err := fs.GetObject(fmt.Sprintf("http://%s/%s", assi.Url, assi.Fid))
	if err != nil {
		t.Errorf("Get object err: %v.", err)
		return
	}
	err = ioutil.WriteFile("ladder.svg", object, 000666)
	if err != nil {
		t.Errorf("Write File err: %v.", err)
		return
	}

	// 移除对象
	err = fs.RemoveObject(fmt.Sprintf("http://%s/%s", assi.Url, assi.Fid))
	if err != nil {
		t.Errorf("Remove object err: %v.", err)
		return
	}
}

type Assign struct {
	Fid       string `json:"fid"`
	Url       string `json:"url"`
	PublicUrl string `json:"publicUrl"`
	Count     int    `json:"count"`
}

type PutObjectRes struct {
	Name string `json:"name"`
	Size int    `json:"size"`
	ETag string `json:"eTag"`
}
