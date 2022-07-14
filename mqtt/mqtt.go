package mqtt

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/errors"
	"reflect"
	"sync"
	"time"
)

/**
获取MQTT
*/
func NewMqtt(cfg *Config) (*Mqtt, error) {
	logger := cfg.Logger
	switch cfg.LogMode {
	case "debug":
		MQTT.DEBUG = logger
		fallthrough
	case "release":
		MQTT.WARN = logger
		MQTT.CRITICAL = logger
		fallthrough
	case "error":
		MQTT.ERROR = logger
	}
	mqttx := &Mqtt{
		config:           cfg,
		DefaultMsgCh:     make(chan Message, 1000),
		reconnectHandler: make(map[uintptr]func()),
		handlerMu:        sync.Mutex{},
	}
	if err := mqttx.connectMQTTBroker(); err != nil {
		cfg.Logger.Println(errors.Errorf("[mqtt]Connect broker err: %v.", err))
		return nil, err
	}

	// refresh 配置
	mqttx.refreshTimerCount = time.Now().Unix() + 600 // TODO 提供配置化
	mqttx.reconnectCh = make(chan struct{})

	go mqttx.keepAlive()

	return mqttx, nil
}

/**
获取client id
*/
func (mt *Mqtt) GetClientId() string {
	return mt.clientID
}

/**
注册重新连接处理程序
*/
func (mt *Mqtt) RegisterReconnectHandler(f func()) {
	mt.handlerMu.Lock()
	defer mt.handlerMu.Unlock()
	mt.reconnectHandler[reflect.ValueOf(f).Pointer()] = f
}

/***
取消注册重新连接处理程序
*/
func (mt *Mqtt) UnRegisterReconnectHandler(f func()) {
	mt.handlerMu.Lock()
	defer mt.handlerMu.Unlock()
	if _, ok := mt.reconnectHandler[reflect.ValueOf(f).Pointer()]; ok {
		delete(mt.reconnectHandler, reflect.ValueOf(f).Pointer())
	}
}

/**
返回一个布尔值，表示客户端是否与mqtt 代理有活动连接，即未处于断开连接或重新连接模式
*/
func (mt *Mqtt) IsConnectionOpen() bool {
	return mt.IsConnectionOpen()
}

/**
根据server mode 获取MQTT 用户名和密码
*/
func (mt *Mqtt) GetMqttUsernameAndPassword(clientId string) (username, password string) {
	return mt.config.AccessKey, mt.config.SecretKey
}

/**
订阅主题消息，收到的消息通过output 返回
*/
func (mt *Mqtt) Subscribe(topic string, qos byte, output chan<- Message) error {
	if token := mt.Client.Subscribe(topic, qos, func(c MQTT.Client, m MQTT.Message) {
		output <- m
	}); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

/**
订阅多个主题消息，收到的消息都经由output 返回
*/
func (mt *Mqtt) MultiSubscribe(filters map[string]byte, output chan<- Message) error {
	if token := mt.SubscribeMultiple(filters, func(c MQTT.Client, m MQTT.Message) {
		output <- m
	}); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

/**
默认订阅主题消息，收到的消息统一通过DefaultMsgCh 返回，可能包含其它主题的消息，需自己通过Message.Topic 分别处理
*/
func (mt *Mqtt) DefalutSubscribe(topic string, qos byte) error {
	return mt.Subscribe(topic, qos, mt.DefaultMsgCh)
}

/**
默认订阅多个主题消息，收到的消息统一通过DefaultMsgCh 返回，可能包含其它主题的消息，需自己通过Message.Topic 分别处理
*/
func (mt *Mqtt) DefalutMultiSubscribe(filters map[string]byte) error {
	return mt.MultiSubscribe(filters, mt.DefaultMsgCh)
}

/**
取消订阅
*/
func (mt *Mqtt) Unsubscribe(topics []string) error {
	if token := mt.Client.Unsubscribe(topics...); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

/**
会将具有指定QoS 和内容的消息发布到指定主题
*/
func (mt *Mqtt) Publish(topic string, qos byte, retained bool, payload interface{}) error {
	if token := mt.Client.Publish(topic, qos, retained, payload); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}
