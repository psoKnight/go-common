//go:build !rocketmq_cgo
// +build !rocketmq_cgo

package rocketmq

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/apache/rocketmq-client-go/v2/rlog"
	"github.com/pkg/errors"
	"sync"
)

type Producer struct {
	rocketmq.Producer
	cfg   *RocketMQConfig
	mutex sync.Mutex
}

// NewProducer 新建消费者
func NewProducer(cfg *RocketMQConfig) (*Producer, error) {
	c := &Producer{
		cfg: cfg,
	}

	// 重定向rocketmq库的log输出和级别
	if cfg.Logger != nil {
		rlog.SetLogger(&loggerWrap{
			logger: cfg.Logger,
		})
		rlog.SetLogLevel(cfg.LogLevel)
	}

	newProducer, err := rocketmq.NewProducer(
		producer.WithNameServer(cfg.Endpoints),
		producer.WithCredentials(primitive.Credentials{
			AccessKey: cfg.AccessKey,
			SecretKey: cfg.SecretKey,
		}),
		producer.WithNamespace(cfg.InstanceID),
		producer.WithInstanceName("rocketmq_producer"),
		producer.WithRetry(c.cfg.RetryTimes),
	)
	if err != nil {
		return nil, err
	}

	if err := newProducer.Start(); err != nil {
		return nil, err
	}

	c.Producer = newProducer

	return c, nil
}

// SendMessageSync 发送同步消息，groupID 暂时无用
func (p *Producer) SendMessageSync(groupID string, msg *Message) error {
	transMsg := primitive.NewMessage(msg.Topic, msg.Body)
	if msg.Property != nil {
		transMsg.WithProperties(msg.Property)
	}
	if msg.Tags != "" {
		transMsg.WithTag(msg.Tags)
	}
	if msg.Keys != nil {
		transMsg.WithKeys(msg.Keys)
	}

	_, err := p.Producer.SendSync(context.Background(), transMsg)
	if err != nil {
		return err
	}
	return nil
}

// CreateTopic 创建topic
func (p *Producer) CreateTopic(groupID, topic, tag string) error {

	if err := p.SendMessageSync(groupID, &Message{
		Topic: topic,
		Tags:  tag,
		Body:  []byte(MsgBodyForCreateTopic),
	}); err != nil {
		return errors.Errorf("[rocketmq]create topic %s err: %v", topic, err)
	}
	return nil
}

// Close 关闭生产者
func (p *Producer) Close() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.Producer.Shutdown()
}
