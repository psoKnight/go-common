package rocketmq

import (
	"context"
	"errors"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/rlog"
)

/**
消费者
*/
func NewConsumer(config *Config) *Consumer {

	// 重定向log 输出和级别
	if config.Logger != nil {
		rlog.SetLogger(&loggerWrap{
			logger: config.Logger,
		})
		rlog.SetLogLevel(config.LogLevel)
	}

	return &Consumer{
		cfg:          config,
		consumerMaps: make(map[string]rocketmq.PushConsumer),
	}
}

/**
消费者订阅消息
*/
func (c *Consumer) Subscribe(groupID, topic, tag string, handler MessageExtHandler) error {
	pushConsumer, err := c.getPushConsumer(groupID)
	if err != nil {
		return err
	}

	selector := consumer.MessageSelector{}
	if tag != "" {
		selector.Type = consumer.TAG
		selector.Expression = tag
	}

	if err = pushConsumer.Subscribe(topic, selector,
		func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			for _, msg := range msgs {
				if handler == nil {
					c.cfg.Logger.Warnf("[rocketmq]Receive msg with no handler: %s.", msg.String())
				}

				if string(msg.Body) == MsgBodyForCreateTopic {
					continue
				}

				if err := handler(convertToMessageExt(msg)); err != nil {
					c.cfg.Logger.Errorf("[rocketmq]Handle msg %s err: %v.", msg.String(), err)
				}
			}
			return consumer.ConsumeSuccess, nil
		}); err != nil {
		return err
	}

	if err = pushConsumer.Start(); err != nil {
		return err
	}

	return nil
}

/**
消费者取消订阅消息
*/
func (c *Consumer) UnSubscribe(groupID, topic string) error {
	pushConsumer, err := c.getPushConsumer(groupID)
	if err != nil {
		return err
	}
	err = pushConsumer.Shutdown()
	if err != nil {
		return err
	}

	{
		c.mutex.Lock()
		defer c.mutex.Unlock()
		delete(c.consumerMaps, groupID)
	}
	return nil
}

/**
关闭消费者
*/
func (c *Consumer) Close() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, consumer := range c.consumerMaps {
		consumer.Shutdown()
	}
}

/**
获取消费者
*/
func (c *Consumer) getPushConsumer(groupID string) (rocketmq.PushConsumer, error) {
	if groupID == "" {
		return nil, errors.New("[rocketmq]Group id is empty.")
	}

	pushConsumer, ok := c.consumerMaps[groupID]
	if !ok {
		c.mutex.Lock()
		defer c.mutex.Unlock()
		pushConsumer, err := rocketmq.NewPushConsumer(
			consumer.WithNameServer(c.cfg.Endpoints),
			consumer.WithCredentials(primitive.Credentials{
				AccessKey: c.cfg.AccessKey,
				SecretKey: c.cfg.SecretKey,
			}),
			consumer.WithNamespace(c.cfg.InstanceID),
			consumer.WithGroupName(groupID),
			consumer.WithConsumerModel(consumer.Clustering),
			consumer.WithInstance(c.cfg.InstanceID+"rocketmq_consumer"),
			consumer.WithRetry(c.cfg.RetryTimes),
		)
		if err != nil {
			return nil, err
		}

		c.consumerMaps[groupID] = pushConsumer

		return pushConsumer, nil
	}

	return pushConsumer, nil
}
