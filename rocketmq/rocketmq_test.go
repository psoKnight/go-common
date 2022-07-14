package rocketmq

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestRocketMQ(t *testing.T) {

	// 获取rocketmq
	rocketmqClient, errNRM := NewRocketMQ(&Config{
		Endpoints:  []string{"127.0.0.1:9876"},
		BrokerAddr: "127.0.0.1:10911",
		InstanceID: "",
		AccessKey:  "",
		SecretKey:  "",
		RetryTimes: 0,
		LogLevel:   "error",
		Logger:     logrus.StandardLogger(),
	})
	if errNRM != nil {
		t.Errorf("Rocketmq conect err: %v.", errNRM)
		return
	}

	// 关闭rocketmq
	defer rocketmqClient.Close()

	groupId := "subscribe_group"
	topic := "subscribe_topic"
	tag := "subscribe_tag"

	// 创建topic
	errCTUA := rocketmqClient.CreateTopicUseAdmin(topic)
	if errCTUA != nil {
		t.Errorf("Rocketmq create topic use admin err: %v.", errCTUA)
		return
	}

	// 生产消息

	for i := 0; i <= 1000; i++ {

		body := []byte(fmt.Sprintf("This is the %d mq message content.", i))

		msg := &Message{
			Topic:    topic,
			Tags:     tag,
			Body:     body,
			Property: map[string]string{"TraceId": uuid.NewV4().String()},
		}

		if errSMS := rocketmqClient.SendMessageSync(groupId, msg); errSMS != nil {
			t.Errorf("Rocketmq send sync message err: %v.", errSMS)
		}
	}

	// 订阅消息
	handlerMessage := func(ext *MessageExt) error {
		t.Logf("Receive msg: %s.", string(ext.Body))
		return nil
	}
	if errCS := rocketmqClient.Subscribe(groupId, topic, tag, handlerMessage); errCS != nil {
		t.Errorf("Rocketmq subscribe err: %v.", errCS)
		return
	}

	// 取消订阅消息
	if errCUS := rocketmqClient.UnSubscribe(groupId, topic); errCUS != nil {
		t.Errorf("Rocketmq unsubcribe err: %v.", errCUS)
		return
	}

	// 删除topic
	errDT := rocketmqClient.DeleteTopic(topic)
	if errDT != nil {
		t.Errorf("Rocketmq delete topic err: %v.", errDT)
		return
	}

	t.Log("Test rocketmq success!")

}
