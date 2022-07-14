package mqtt

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

/**
Mqtt 连接broker
*/
func (c *Mqtt) connectMQTTBroker() error {
	cfg := c.config
	c.clientID = cfg.GroupID + uuid.NewV4().String()

	opts := MQTT.NewClientOptions()
	opts.AddBroker(cfg.BrokerURL).SetClientID(c.clientID).SetUsername(cfg.AccessKey).SetPassword(cfg.SecretKey)
	opts.SetMaxReconnectInterval(10 * time.Second).SetCleanSession(false).SetResumeSubs(true).SetKeepAlive(10 * time.Second)
	opts.SetConnectionLostHandler(c.connectionLostHandler).SetOnConnectHandler(c.onConnectHandler)
	var client MQTT.Client = MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return errors.Errorf("[mqtt]Connect broker %s using client id %s err: %v.", cfg.BrokerURL, c.clientID, token.Error())
	}
	c.Client = client
	return nil
}

/**
MQTT 保持live
*/
func (c *Mqtt) keepAlive() {
	timer := time.NewTimer(1 * time.Minute)
	defer timer.Stop()
	mu := sync.Mutex{}
	resetFunc := func() {
		mu.Lock()
		defer mu.Unlock()
		logrus.Infof("mqtt reconnected start!")
		// 首先刷新订阅
		c.onConnectHandler(c.Client)

		// 其次断开连接
		c.Client.Disconnect(250)
		time.Sleep(500 * time.Millisecond)
		// 最后重新连接
		if err := c.connectMQTTBroker(); err != nil {
			logrus.Error(err)
		}
		logrus.Infof("[mqtt]Mqtt reconnected end.")
	}
	for {
		select {
		case <-timer.C:
			logrus.Infof("[mqtt]Timer trigger, time.Now().Unix(): %+v, refreshTimerCount: %v.", time.Now().Unix(), c.refreshTimerCount)
			c.refreshMu.Lock()
			if time.Now().Unix() > c.refreshTimerCount {
				resetFunc()
				c.refreshTimerCount = time.Now().Unix() + 600
			}
			c.refreshMu.Unlock()
			timer.Reset(1 * time.Minute)
		case <-c.reconnectCh:
			logrus.Infof("[mqtt]<-reconnectCh.")
			resetFunc()
		}
	}
}

/**
当客户端连接时调用
在初始连接时和自动重新连接时
*/
func (c *Mqtt) onConnectHandler(client MQTT.Client) {
	c.handlerMu.Lock()
	defer c.handlerMu.Unlock()
	for _, f := range c.reconnectHandler {
		f()
	}
}

/**
会在客户端意外失去与MQTT 代理的连接的情况下执行
*/
func (c *Mqtt) connectionLostHandler(client MQTT.Client, err error) {
	c.config.Logger.Printf("mqtt %s connect lost %s", c.clientID, err)
}
