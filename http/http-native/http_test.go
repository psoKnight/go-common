package http_native

import (
	"encoding/json"
	"testing"
	"time"
)

type faceThresholdConfig struct {
	HorizontalAngleThreshold int     `json:"horizontal_angle_threshold"` // 人脸水平角度过滤阈值
	VerticalAngleThreshold   int     `json:"vertical_angle_threshold"`   // 人脸垂直角度过滤阈值
	RotationAngleThreshold   int     `json:"rotation_angle_threshold"`   // 人脸旋转角度过滤阈值
	FuzzinessThreshold       float64 `json:"fuzziness_threshold"`        // 模糊度过滤阈值
	ClusterFilterLevel       int     `json:"cluster_filter_level"`       // 聚类过滤等级
}

type faceClusterConfig struct {
	faceThresholdConfig
	CustomTag int `json:"custom_tag"` // 设置标志 0-默认/1-自定义
}

type res struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func TestGet(t *testing.T) {
	cfg := &HttpConfig{Timeout: time.Duration(30) * time.Second}
	client := NewHttp(cfg)
	get, err := client.Get("http://10.171.5.193/prisms/get/clusterConfig", nil, nil)
	if err != nil {
		t.Errorf("Get err: %v.", err)
		return
	}

	t.Log(string(get) == "")

	g := &faceThresholdConfig{}
	err = json.Unmarshal(get, g)
	if err != nil {
		t.Errorf("Json unmarshal err: %v.", err)
	}
	t.Logf("Get: %+v.", g)
}

func TestPostJson(t *testing.T) {
	cfg := &HttpConfig{Timeout: time.Duration(30) * time.Second}
	client := NewHttp(cfg)

	param := &faceClusterConfig{
		faceThresholdConfig: faceThresholdConfig{
			HorizontalAngleThreshold: 6,
			VerticalAngleThreshold:   66,
			RotationAngleThreshold:   6,
			FuzzinessThreshold:       66,
			ClusterFilterLevel:       1,
		},
		CustomTag: 0,
	}

	postJson, err := client.PostJson("http://10.171.5.193:20211/prisms/set/clusterConfig", param, nil)
	if err != nil {
		t.Errorf("Post json err: %v.", err)
		return
	}

	t.Log(string(postJson) == "")

	res := &res{}
	err = json.Unmarshal(postJson, res)
	if err != nil {
		t.Errorf("Json unmarshal err: %v.", err)
	}
	t.Logf("Post json res: %+v.", res)
}

func TestPostForm(t *testing.T) {}
func TestPostFile(t *testing.T) {}

func TestDelete(t *testing.T) {
	cfg := &HttpConfig{Timeout: time.Duration(30) * time.Second}
	client := NewHttp(cfg)
	get, err := client.Delete("http://10.117.48.122:9222/1,46f422b6b4", nil, nil, nil)
	if err != nil {
		t.Errorf("Get err: %v.", err)
		return
	}
	t.Log(string(get))
}
