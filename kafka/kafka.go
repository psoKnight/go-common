package kafka

import (
	"errors"
	"github.com/Shopify/sarama"
	"time"
)

// Kafka 结构
type Kafka struct {
	producer *Producer
	consumer *Consumer
	cfg      *Config
}

// Config 通用对象
type Config struct {
	Endpoints []string
}

func NewKafka(cfg *Config) (*Kafka, error) {

	if cfg == nil {
		return nil, errors.New("[kafka]config is nil")
	}
	client := &Kafka{
		cfg: cfg,
	}

	producer, err := NewProducer(cfg)
	if err != nil {
		return nil, err
	}
	client.producer = producer

	consumer, err := NewConsumer(cfg)
	if err != nil {
		return nil, err
	}
	client.consumer = consumer

	return client, nil
}

// GetSyncProducer 获取Kafka SyncProducer
func (k *Kafka) GetSyncProducer() sarama.SyncProducer {
	return k.producer.SyncProducer
}

// GetAsyncProducer 获取Kafka AsyncProducer
func (k *Kafka) GetAsyncProducer() sarama.AsyncProducer {
	return k.producer.AsyncProducer
}

// GetConsumer 获取Kafka Consumer
func (k *Kafka) GetConsumer() sarama.Consumer {
	return k.consumer.Consumer
}

// SyncSendMessage 同步生产数据
func (k *Kafka) SyncSendMessage(topic, key, value string) (int32, int64, error) {
	return k.producer.SyncSendMessage(topic, key, value)
}

// AsyncSendMessage 异步生产数据
func (k *Kafka) AsyncSendMessage(topic, key, value string, timeOut time.Duration) (*sarama.ProducerMessage, error) {
	return k.producer.AsyncSendMessage(topic, key, value, timeOut)
}

// ConsumeMessage 消费数据
func (k *Kafka) ConsumeMessage(topic, key string, offsetType int64, ch chan *sarama.ConsumerMessage) error {
	return k.consumer.ConsumeMessage(topic, key, offsetType, ch)
}

// ConsumeMessageByGroup 消费数据(通过消费组)
func (k *Kafka) ConsumeMessageByGroup(topics []string, group string, offsetType int64, handler sarama.ConsumerGroupHandler) error {
	return k.consumer.ConsumeMessageByGroup(topics, group, offsetType, handler)
}

// Close 关闭kafka
func (k *Kafka) Close() error {
	_ = k.producer.Close()
	if err := k.consumer.Close(); err != nil {
		return err
	}
	return nil
}
