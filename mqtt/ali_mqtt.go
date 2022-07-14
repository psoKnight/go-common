package mqtt

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/onsmqtt"
)

/**
获取阿里巴巴 MQTT
*/
func NewAliMqtt(cfg *AliConfig) *AliMqtt {
	var err error
	am, err := onsmqtt.NewClientWithAccessKey(cfg.RegionID, cfg.AccessKey, cfg.SecretKey)
	if err != nil {
		cfg.Logger.Println(err)
		return nil
	}

	return &AliMqtt{Client: am, cfg: cfg}
}

/**
sha 1 和 Base 64 加密
*/
func sha1AndBase64Encrypt(secreKey, data []byte) string {
	mac := hmac.New(sha1.New, secreKey)
	mac.Write(data)
	macres := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(macres)
}

/**
生成Resources
*/
func (am *AliMqtt) generateResources(deviceID string, topics []string) string {
	upTopicPrefix, downTopicPrefix := am.cfg.GetMqttUpAndDownTopicPrefix()
	uptopic := fmt.Sprintf("%s/%s", upTopicPrefix, deviceID)
	downtopic := fmt.Sprintf("%s/%s", downTopicPrefix, deviceID)

	resArr := make([]string, 0, 2+len(topics))
	resArr = append(resArr, uptopic, downtopic)
	resArr = append(resArr, topics...)
	sort.Strings(resArr)

	resources := strings.Join(resArr, ",")
	return resources
}

/**
获取阿里巴巴MQTT token 账号/密码
*/
func (am *AliMqtt) GetAliMqttTokenUsernameAndPassword(deviceID string, topics []string) (username, password string, err error) {
	cfg := am.cfg

	request := onsmqtt.CreateApplyTokenRequest()
	request.Resources = am.generateResources(deviceID, nil)
	request.InstanceId = cfg.InstanceID
	request.ExpireTime = requests.Integer(strconv.Itoa(int(time.Now().Add(cfg.GetTokenExpireInterval()).UnixNano() / int64(time.Millisecond))))
	request.Actions = "R,W"

	response, err := am.ApplyToken(request)
	if err != nil {
		return "", "", err
	}

	username = fmt.Sprintf("Token|%s|%s", cfg.AccessKey, cfg.InstanceID)
	password = fmt.Sprintf("RW|%s", response.Token)
	return username, password, nil
}

/**
获取阿里巴巴MQTT signature 账号/密码
*/
func (am *AliMqtt) GetAliMQTTSignatureUsernameAndPassword(clientID string) (username, password string) {
	cfg := am.cfg
	username = fmt.Sprintf("Signature|%s|%s", cfg.AccessKey, cfg.InstanceID)
	password = sha1AndBase64Encrypt([]byte(cfg.SecretKey), []byte(clientID))
	return username, password
}

/**
根据server mode 获取阿里巴巴MQTT 用户名/密码
*/
func (am *AliMqtt) GetMQTTUsernameAndPassword(clientID string) (username, password string, err error) {
	return am.GetAliMqttTokenUsernameAndPassword(clientID, nil)
}

/**
根据server mode 生成阿里巴巴MQTT client id
*/
func (am *AliMqtt) GenMQTTClientID(deviceID string) string {
	return am.cfg.GetMqttGroupID() + deviceID
}

// 阿里巴巴MQTT 相关配置

/**
返回阿里巴巴 MQTT uptopic 和 downtopic prefix
*/
func (c *AliConfig) GetMqttUpAndDownTopicPrefix() (upTopicPrefix, downTopicPrefix string) {
	return c.UpTopicPrefix, c.DownTopicPrefix
}

/**
返回阿里巴巴MQTT group id
*/
func (c *AliConfig) GetMqttGroupID() string {
	return c.GroupID
}

/**
返回阿里巴巴MQTT broker url
*/
func (c *AliConfig) GetMqttBrokerURL() string {
	return c.BrokerURL
}

/**
返回阿里巴巴MQTT 配置
*/
func (c *AliConfig) GetMqttAliConfig() *AliConfig {
	return c
}

/**
返回阿里巴巴MQTT token 过期间隔
*/
func (c *AliConfig) GetTokenExpireInterval() time.Duration {
	return c.TokenExpireInterval
}
