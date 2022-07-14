package mqtt

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/onsmqtt"
)

type AliMQTTConfig struct {
	RegionId   string `json:"region_id"`
	InstanceId string `json:"instance_id"`
	BrokerUrl  string `json:"broker_url"`
	GroupId    string `json:"group_id"`
	AccessKey  string `json:"access_key"`
	SecretKey  string `json:"secret_key"`

	UpTopicPrefix       string        `json:"up_topic_prefix"`
	DownTopicPrefix     string        `json:"down_topic_prefix"`
	TokenExpireInterval time.Duration `json:"token_expire_interval"` // token 过期间隔

	Logger *log.Logger `json:"logger"`
}

type AliMQTT struct {
	cli *onsmqtt.Client
	cfg *AliMQTTConfig
}

// NewAliMQTT 新建阿里巴巴mqtt
func NewAliMQTT(cfg *AliMQTTConfig) (*AliMQTT, error) {
	am, err := onsmqtt.NewClientWithAccessKey(cfg.RegionId, cfg.AccessKey, cfg.SecretKey)
	if err != nil {
		return nil, err
	}

	return &AliMQTT{cli: am, cfg: cfg}, nil
}

// GetAliClient 获取阿里巴巴client
func (amt *AliMQTT) GetAliClient() *onsmqtt.Client {
	return amt.cli
}

// sha1 和Base64 加密
func sha1AndBase64Encrypt(secreKey, data []byte) string {
	mac := hmac.New(sha1.New, secreKey)
	mac.Write(data)
	macres := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(macres)
}

// 生成Resources
func (amt *AliMQTT) generateResources(deviceId string, topics []string) string {
	upTopicPrefix, downTopicPrefix := amt.cfg.UpTopicPrefix, amt.cfg.DownTopicPrefix
	upTopic := fmt.Sprintf("%s/%s", upTopicPrefix, deviceId)
	downTopic := fmt.Sprintf("%s/%s", downTopicPrefix, deviceId)

	resArr := make([]string, 0, 2+len(topics))
	resArr = append(resArr, upTopic, downTopic)
	resArr = append(resArr, topics...)
	sort.Strings(resArr)

	resources := strings.Join(resArr, ",")
	return resources
}

// GetAliMQTTTokenUsernameAndPassword 获取阿里巴巴mqtt token 账号/密码
func (amt *AliMQTT) GetAliMQTTTokenUsernameAndPassword(deviceId string, topics []string) (username, password string, err error) {
	cfg := amt.cfg

	request := onsmqtt.CreateApplyTokenRequest()
	request.Resources = amt.generateResources(deviceId, nil)
	request.InstanceId = cfg.InstanceId
	request.ExpireTime = requests.Integer(strconv.Itoa(int(time.Now().Add(cfg.TokenExpireInterval).UnixNano() / int64(time.Millisecond))))
	request.Actions = "R,W"

	response, err := amt.cli.ApplyToken(request)
	if err != nil {
		return "", "", err
	}

	username = fmt.Sprintf("TOKEN|%s|%s", cfg.AccessKey, cfg.InstanceId)
	password = fmt.Sprintf("RW|%s", response.Token)
	return username, password, nil
}

// GetAliMQTTSignatureUsernameAndPassword 获取阿里巴巴mqtt signature 账号/密码
func (amt *AliMQTT) GetAliMQTTSignatureUsernameAndPassword(clientId string) (username, password string) {
	cfg := amt.cfg
	username = fmt.Sprintf("SIGNATURE|%s|%s", cfg.AccessKey, cfg.InstanceId)
	password = sha1AndBase64Encrypt([]byte(cfg.SecretKey), []byte(clientId))
	return username, password
}

// GetMQTTUsernameAndPassword 根据server mode 获取阿里巴巴MQTT 用户名/密码
func (amt *AliMQTT) GetMQTTUsernameAndPassword(clientId string) (username, password string, err error) {
	return amt.GetAliMQTTTokenUsernameAndPassword(clientId, nil)
}
