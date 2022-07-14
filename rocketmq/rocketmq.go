package rocketmq

import (
	"context"
	"errors"
	"github.com/apache/rocketmq-client-go/v2/admin"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

/**
获取RocketMQ 客户端
*/
func NewRocketMQ(cfg *Config) (*RocketMq, error) {

	if cfg == nil {
		return nil, errors.New("[rocketmq]Config is nil.")
	}

	client := &RocketMq{
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

/**
RocketMQ 订阅消息
	需保持订阅关系一致，一个消费者groupId 下订阅的topic、tag 需保持一致，
所以针对不同的topic、tag 对用不同的groupId 区分，统一由rocketmq 模块启动start 和shutdown
*/
func (rm *RocketMq) Subscribe(groupId, topic, tag string, handler MessageExtHandler) error {
	// 发送创建topic 消息，保证topic 已经生成
	//if err := rm.producer.CreateTopic(groupId, topic, tag); err != nil {
	//	return err
	//}

	return rm.consumer.Subscribe(groupId, topic, tag, handler)
}

/**
RocketMQ 发送同步消息
*/
func (rm *RocketMq) SendMessageSync(groupId string, msg *Message) error {
	return rm.producer.SendMessageSync(groupId, msg)
}

/**
RocketMQ 取消订阅消息
*/
func (rm *RocketMq) UnSubscribe(groupId, topic string) error {
	return rm.consumer.UnSubscribe(groupId, topic)
}

/**
RocketMQ 关闭客户端
*/
func (rm *RocketMq) Close() {
	rm.producer.Close()
	rm.consumer.Close()
}

/**
RocketMQ 删除topic
*/
func (rm *RocketMq) Producer() *Producer {
	return rm.producer
}

/**
RocketMQ 获取消费者
*/
func (rm *RocketMq) Consumer() *Consumer {
	return rm.consumer
}

/**
RocketMQ 创建topic
*/
func (rm *RocketMq) CreateTopicUseAdmin(topic string) error {

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

/**
RocketMQ 删除topic
*/
func (rm *RocketMq) DeleteTopic(topic string) error {

	defaultAdmin, err := admin.NewAdmin(admin.WithResolver(primitive.NewPassthroughResolver(rm.cfg.Endpoints)))
	if err != nil {
		return err
	}

	defer defaultAdmin.Close()

	err = defaultAdmin.DeleteTopic(
		context.Background(),
		admin.WithTopicDelete(topic),
		//admin.WithBrokerAddrDelete(rm.config.BrokerAddr),
		//admin.WithNameSrvAddr(rm.config.Endpoints),
	)
	if err != nil {
		return err
	}

	return nil
}
