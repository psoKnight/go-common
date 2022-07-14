package mqtt

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/onsmqtt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
	"sync"
	"time"
)

// MQTT 通用配置
type Config struct {
	BrokerURL string
	GroupID   string
	AccessKey string
	SecretKey string

	LogMode string

	Logger *log.Logger
}

// MQTT
type Mqtt struct {
	MQTT.Client
	config   *Config
	clientID string

	// DefaultMsgCh 默认的接收消息的channel
	DefaultMsgCh     chan Message
	reconnectHandler map[uintptr]func()
	handlerMu        sync.Mutex

	// refresh 刷新机制相关
	refreshMu         sync.Mutex
	reconnectCh       chan struct{}
	refreshTimerCount int64
}

// 收到的消息处理类型
type (
	Message MQTT.Message
)

// Alibaba 通用配置
type AliConfig struct {
	RegionID   string
	InstanceID string
	BrokerURL  string
	GroupID    string
	AccessKey  string
	SecretKey  string

	UpTopicPrefix       string
	DownTopicPrefix     string
	TokenExpireInterval time.Duration // token 过期间隔

	Logger *log.Logger
}

// Alibaba MQTT
type AliMqtt struct {
	*onsmqtt.Client
	cfg *AliConfig
}
