package mqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"log"
	"reflect"
	"sync"
	"time"
)

type MQTTConfig struct {
	BrokerUrl string `json:"broker_url"`
	GroupId   string `json:"group_id"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`

	LogMode string      `json:"log_mode"`
	Logger  *log.Logger `json:"logger"`
}

type MQTT struct {
	cli      mqtt.Client
	cfg      *MQTTConfig
	clientId string

	// DefaultMsgCh 默认的接收消息的channel
	DefaultMsgCh     chan Message
	reconnectHandler map[uintptr]func()
	handlerMu        sync.Mutex

	// refresh 刷新机制相关
	refreshMu         sync.Mutex
	reconnectCh       chan struct{}
	refreshTimerCount int64
}

type (
	Message mqtt.Message
)

// NewMQTT 新建mqtt
func NewMQTT(cfg *MQTTConfig) (*MQTT, error) {
	logger := cfg.Logger
	switch cfg.LogMode {
	case "debug":
		mqtt.DEBUG = logger
		fallthrough
	case "release":
		mqtt.WARN = logger
		mqtt.CRITICAL = logger
		fallthrough
	case "error":
		mqtt.ERROR = logger
	}

	mqttx := &MQTT{
		cfg:              cfg,
		DefaultMsgCh:     make(chan Message, 1000),
		reconnectHandler: make(map[uintptr]func()),
		handlerMu:        sync.Mutex{},
	}
	if err := mqttx.connectMQTTBroker(); err != nil {
		return nil, err
	}

	// refresh 配置
	mqttx.refreshTimerCount = time.Now().Unix() + 600 // TODO 提供配置化
	mqttx.reconnectCh = make(chan struct{})

	go mqttx.keepAlive()

	return mqttx, nil
}

// GetClient 获取client
func (mt *MQTT) GetClient() mqtt.Client {
	return mt.cli
}

// GetClientId 获取client ID
func (mt *MQTT) GetClientId() string {
	return mt.clientId
}

// RegisterReconnectHandler 注册重新连接处理程序
func (mt *MQTT) RegisterReconnectHandler(f func()) {
	mt.handlerMu.Lock()
	defer mt.handlerMu.Unlock()
	mt.reconnectHandler[reflect.ValueOf(f).Pointer()] = f
}

// UnRegisterReconnectHandler 取消注册重新连接处理程序
func (mt *MQTT) UnRegisterReconnectHandler(f func()) {
	mt.handlerMu.Lock()
	defer mt.handlerMu.Unlock()
	if _, ok := mt.reconnectHandler[reflect.ValueOf(f).Pointer()]; ok {
		delete(mt.reconnectHandler, reflect.ValueOf(f).Pointer())
	}
}

// GetMQTTUsernameAndPassword 根据server mode 获取MQTT 用户名和密码
func (mt *MQTT) GetMQTTUsernameAndPassword(clientId string) (username, password string) {
	return mt.cfg.AccessKey, mt.cfg.SecretKey
}

// Subscribe 订阅主题消息，收到的消息通过output 返回
func (mt *MQTT) Subscribe(topic string, qos byte, output chan<- Message) error {
	if token := mt.cli.Subscribe(topic, qos, func(c mqtt.Client, m mqtt.Message) {
		output <- m
	}); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// MultiSubscribe 订阅多个主题消息，收到的消息都经由output 返回
func (mt *MQTT) MultiSubscribe(filters map[string]byte, output chan<- Message) error {
	if token := mt.cli.SubscribeMultiple(filters, func(c mqtt.Client, m mqtt.Message) {
		output <- m
	}); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// DefalutSubscribe 默认订阅主题消息，收到的消息统一通过DefaultMsgCh 返回，可能包含其它主题的消息，需自己通过Message.Topic 分别处理
func (mt *MQTT) DefalutSubscribe(topic string, qos byte) error {
	return mt.Subscribe(topic, qos, mt.DefaultMsgCh)
}

// DefalutMultiSubscribe 默认订阅多个主题消息，收到的消息统一通过DefaultMsgCh 返回，可能包含其它主题的消息，需自己通过Message.Topic 分别处理
func (mt *MQTT) DefalutMultiSubscribe(filters map[string]byte) error {
	return mt.MultiSubscribe(filters, mt.DefaultMsgCh)
}

// Unsubscribe 取消订阅
func (mt *MQTT) Unsubscribe(topics []string) error {
	if token := mt.cli.Unsubscribe(topics...); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// Publish 会将具有指定QoS 和内容的消息发布到指定主题
func (mt *MQTT) Publish(topic string, qos byte, retained bool, payload interface{}) error {
	if token := mt.cli.Publish(topic, qos, retained, payload); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// Mqtt 连接broker
func (mt *MQTT) connectMQTTBroker() error {
	cfg := mt.cfg
	mt.clientId = cfg.GroupId + uuid.NewV4().String()

	opts := mqtt.NewClientOptions()
	opts.AddBroker(cfg.BrokerUrl).SetClientID(mt.clientId).SetUsername(cfg.AccessKey).SetPassword(cfg.SecretKey)
	opts.SetMaxReconnectInterval(10 * time.Second).SetCleanSession(false).SetResumeSubs(true).SetKeepAlive(10 * time.Second)
	opts.SetConnectionLostHandler(mt.connectionLostHandler).SetOnConnectHandler(mt.onConnectHandler)
	var client mqtt.Client = mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return errors.Errorf("[mqtt]connect broker %s using client id %s err: %v.", cfg.BrokerUrl, mt.clientId, token.Error())
	}
	mt.cli = client
	return nil
}

// 保持live
func (mt *MQTT) keepAlive() {
	timer := time.NewTimer(1 * time.Minute)
	defer timer.Stop()
	mu := sync.Mutex{}
	resetFunc := func() {
		mu.Lock()
		defer mu.Unlock()
		mt.cfg.Logger.Println("[mqtt]reconnected start.")
		// 首先刷新订阅
		mt.onConnectHandler(mt.cli)

		// 其次断开连接
		mt.cli.Disconnect(250)
		time.Sleep(500 * time.Millisecond)
		// 最后重新连接
		if err := mt.connectMQTTBroker(); err != nil {
			mt.cfg.Logger.Println(err)
		}
		mt.cfg.Logger.Println("[mqtt]reconnect end.")
	}
	for {
		select {
		case <-timer.C:
			mt.cfg.Logger.Printf("[mqtt]timer trigger, time now: %d, refresh timer count: %d.", time.Now().Unix(), mt.refreshTimerCount)
			mt.refreshMu.Lock()
			if time.Now().Unix() > mt.refreshTimerCount {
				resetFunc()
				mt.refreshTimerCount = time.Now().Unix() + 600
			}
			mt.refreshMu.Unlock()
			timer.Reset(1 * time.Minute)
		case <-mt.reconnectCh:
			mt.cfg.Logger.Println("[mqtt]<-reconnectCh.")
			resetFunc()
		}
	}
}

// 当客户端连接时调用，在初始连接时和自动重新连接时
func (mt *MQTT) onConnectHandler(cli mqtt.Client) {
	mt.handlerMu.Lock()
	defer mt.handlerMu.Unlock()
	for _, f := range mt.reconnectHandler {
		f()
	}
}

// 会在客户端意外失去与mqtt 代理的连接的情况下执行
func (mt *MQTT) connectionLostHandler(cli mqtt.Client, err error) {
	mt.cfg.Logger.Printf("[mqtt]%s connect failed, err: %v.", mt.clientId, err)
}
