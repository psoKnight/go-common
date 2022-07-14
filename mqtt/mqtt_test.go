package mqtt

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

func TestMqtt(t *testing.T) {
	// 获取MQTT
	cfg := &Config{
		BrokerURL: "10.171.8.115:1883",
		GroupID:   "group",
		AccessKey: "root",
		SecretKey: "50ICqCiAFio",
		LogMode:   "error",
		Logger:    log.New(os.Stderr, "", log.LstdFlags),
	}

	mqtt, err := NewMqtt(cfg)
	if err != nil {
		t.Errorf("Mqtt conenct err: %v.", err)
		return
	}

	// 订阅消息
	topic := fmt.Sprintf("meglink/test/%s", mqtt.GetClientId())
	if err := mqtt.Subscribe(topic, 1, mqtt.DefaultMsgCh); err != nil {
		t.Errorf("Subscribe err: %v.", err)
		return
	}

	// 发布消息
	go func() {
		if err := mqtt.Publish(topic, 1, false, "{I send a msg to mqtt.}"); err != nil {
			t.Errorf("Publish err: %v.", err)
			return
		}
	}()

	expire := time.NewTimer(5 * time.Second)
	select {
	case <-expire.C: // 设置超时时间
		t.Errorf("Get msg time out.")
	case msg := <-mqtt.DefaultMsgCh:
		t.Logf("Receive a msg: %v.", string(msg.Payload()))
	}

	// 取消订阅
	err = mqtt.Unsubscribe([]string{topic})
	if err != nil {
		t.Errorf("Unsubscribe err: %v.", err)
		return
	}

}
