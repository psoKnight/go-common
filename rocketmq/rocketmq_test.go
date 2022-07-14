package rocketmq

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestRocketMQ(t *testing.T) {

	// 获取rocketmq
	rocketmqClient, err := NewRocketMQ(&RocketMQConfig{
		Endpoints:  []string{"127.0.0.1:9876"},
		BrokerAddr: "127.0.0.1:10911",
		InstanceID: "",
		AccessKey:  "",
		SecretKey:  "",
		RetryTimes: 0,
		LogLevel:   "error",
		Logger:     logrus.StandardLogger(),
	})
	if err != nil {
		t.Errorf("Rocketmq conect err: %v.", err)
		return
	}

	// 关闭rocketmq
	defer rocketmqClient.Close()

	groupId := "subscribe_group"
	topic := "subscribe_topic"
	tag := "subscribe_tag"

	// 创建topic
	if err := rocketmqClient.CreateTopicUseAdmin(topic); err != nil {
		t.Errorf("Rocketmq create topic use admin err: %v.", err)
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

		if err := rocketmqClient.SendMessageSync(groupId, msg); err != nil {
			t.Errorf("Rocketmq send sync message err: %v.", err)
		}
	}

	// 订阅消息
	handlerMessage := func(ext *MessageExt) error {
		t.Logf("Receive msg: %s.", string(ext.Body))
		return nil
	}
	if err := rocketmqClient.Subscribe(groupId, topic, tag, handlerMessage); err != nil {
		t.Errorf("Rocketmq subscribe err: %v.", err)
		return
	}

	// 取消订阅消息
	if err := rocketmqClient.UnSubscribe(groupId, topic); err != nil {
		t.Errorf("Rocketmq unsubcribe err: %v.", err)
		return
	}

	// 删除topic
	if err := rocketmqClient.DeleteTopic(topic); err != nil {
		t.Errorf("Rocketmq delete topic err: %v.", err)
		return
	}

	t.Log("Test rocketmq success!")
}
