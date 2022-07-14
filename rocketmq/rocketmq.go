package rocketmq

import (
	"context"
	"errors"
	"github.com/apache/rocketmq-client-go/v2/admin"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/sirupsen/logrus"
)

type RocketMQConfig struct {
	Endpoints  []string       `json:"endpoints"`
	BrokerAddr string         `json:"broker_addr"`
	InstanceID string         `json:"instance_id"`
	AccessKey  string         `json:"access_key"`
	SecretKey  string         `json:"secret_key"`
	RetryTimes int            `json:"retry_times"`
	LogLevel   string         `json:"log_level"`
	Logger     *logrus.Logger `json:"logger"`
}

type RocketMQ struct {
	producer *Producer
	consumer *Consumer
	cfg      *RocketMQConfig
}

// NewRocketMQ 新建rocketmq
func NewRocketMQ(cfg *RocketMQConfig) (*RocketMQ, error) {

	if cfg == nil {
		return nil, errors.New("[rocketmq]config is nil")
	}

	client := &RocketMQ{
		cfg: cfg,
	}

	producer, err := NewProducer(cfg)
	if err != nil {
		return nil, err
	}
	client.producer = producer

	client.consumer = NewConsumer(cfg)

	return client, nil
}

// Close 关闭客户端
func (rm *RocketMQ) Close() {
	rm.producer.Close()
	rm.consumer.Close()
}

// GetProducer 获取producer
func (rm *RocketMQ) GetProducer() *Producer {
	return rm.producer
}

// GetConsumer 获取consumer
func (rm *RocketMQ) GetConsumer() *Consumer {
	return rm.consumer
}

// Subscribe 订阅消息
/**
	需保持订阅关系一致，一个消费者groupId 下订阅的topic、tag 需保持一致，
所以针对不同的topic、tag 对用不同的groupId 区分，统一由rocketmq 模块启动start 和shutdown
*/
func (rm *RocketMQ) Subscribe(groupId, topic, tag string, handler MessageExtHandler) error {
	// 发送创建topic 消息，保证topic 已经生成
	//if err := rm.producer.CreateTopic(groupId, topic, tag); err != nil {
	//	return err
	//}

	return rm.consumer.Subscribe(groupId, topic, tag, handler)
}

// SendMessageSync 发送同步消息
func (rm *RocketMQ) SendMessageSync(groupId string, msg *Message) error {
	return rm.producer.SendMessageSync(groupId, msg)
}

// UnSubscribe 取消订阅消息
func (rm *RocketMQ) UnSubscribe(groupId, topic string) error {
	return rm.consumer.UnSubscribe(groupId, topic)
}

// CreateTopicUseAdmin 创建topic
func (rm *RocketMQ) CreateTopicUseAdmin(topic string) error {

	defaultAdmin, err := admin.NewAdmin(admin.WithResolver(primitive.NewPassthroughResolver(rm.cfg.Endpoints)))
	if err != nil {
		return err
	}

	defer defaultAdmin.Close()

	err = defaultAdmin.CreateTopic(
		context.Background(),
		admin.WithTopicCreate(topic),
		admin.WithBrokerAddrCreate(rm.cfg.BrokerAddr),
		//admin.WithBrokerAddrCreate(rm.config.BrokerAddr),
	)
	if err != nil {
		return err
	}

	return nil
}

// DeleteTopic 删除topic
func (rm *RocketMQ) DeleteTopic(topic string) error {

	defaultAdmin, err := admin.NewAdmin(admin.WithResolver(primitive.NewPassthroughResolver(rm.cfg.Endpoints)))
	if err != nil {
		return err
	}

	defer defaultAdmin.Close()

	if err = defaultAdmin.DeleteTopic(
		context.Background(),
		admin.WithTopicDelete(topic),
		//admin.WithBrokerAddrDelete(rm.config.BrokerAddr),
		//admin.WithNameSrvAddr(rm.config.Endpoints),
	); err != nil {
		return err
	}

	return nil
}
