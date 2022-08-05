package kafka

import (
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"time"
)

// Producer 生产者
type Producer struct {
	sarama.SyncProducer  // 同步生产者
	sarama.AsyncProducer // 异步生产者
	cfg                  *Config
}

// NewProducer 新建生产者
func NewProducer(cfg *Config) (*Producer, error) {
	p := &Producer{
		cfg: cfg,
	}

	// 新建同步消费者
	syncConfig := sarama.NewConfig()
	syncConfig.Producer.RequiredAcks = sarama.WaitForAll
	syncConfig.Producer.Partitioner = sarama.NewRandomPartitioner
	syncConfig.Producer.Return.Successes = true
	syncConfig.Producer.Return.Errors = true
	// 设置使用的kafka版本,如果低于V0_10_0_0版本,消息中的timestrap 没有作用。需要消费和生产同时配置
	// 注意，版本设置不对的话，kafka 会返回很奇怪的错误，并且无法成功发送消息
	//syncConfig.Version = sarama.V0_10_2_1
	syncProducer, err := sarama.NewSyncProducer(cfg.Endpoints, syncConfig)
	if err != nil {
		return nil, err
	}
	p.SyncProducer = syncProducer

	// 新建异步消费者
	asynConfig := sarama.NewConfig()
	asynConfig.Producer.RequiredAcks = sarama.WaitForAll
	asynConfig.Producer.Partitioner = sarama.NewRandomPartitioner
	asynConfig.Producer.Return.Successes = true
	asynConfig.Producer.Return.Errors = true
	// 设置使用的kafka版本,如果低于V0_10_0_0版本,消息中的timestrap 没有作用。需要消费和生产同时配置
	// 注意，版本设置不对的话，kafka 会返回很奇怪的错误，并且无法成功发送消息
	//asynConfig.Version = sarama.V0_10_2_1
	asyncProducer, err := sarama.NewAsyncProducer(cfg.Endpoints, asynConfig)
	if err != nil {
		return nil, err
	}
	p.AsyncProducer = asyncProducer

	return p, nil
}

// SyncSendMessage 同步生产者发送消息
func (p *Producer) SyncSendMessage(topic, key, value string) (int32, int64, error) {
	if topic == "" {
		return 0, 0, errors.New("[kafka]topic is '', please check")
	}

	msg := &sarama.ProducerMessage{}
	msg.Topic = topic

	if key != "" {
		msg.Key = sarama.StringEncoder(key)
	}
	if value != "" {
		msg.Value = sarama.StringEncoder(value)
	}

	pid, offset, err := p.SyncProducer.SendMessage(msg)
	if err != nil {
		return 0, 0, err
	}

	return pid, offset, nil
}

// AsyncSendMessage 异步生产者发送消息
func (p *Producer) AsyncSendMessage(topic, key, value string, timeOut time.Duration) (*sarama.ProducerMessage, error) {
	if topic == "" {
		return nil, errors.New("[kafka]topic is '', please check")
	}

	msg := &sarama.ProducerMessage{}
	msg.Topic = topic

	if key != "" {
		msg.Key = sarama.StringEncoder(key)
	}
	if value != "" {
		msg.Value = sarama.StringEncoder(value)
	}

	// 确认超时时间
	if timeOut == time.Duration(0) {
		timeOut = time.Duration(60) * time.Second // 默认60s
	}

	// 使用通道发送
	p.AsyncProducer.Input() <- msg

	for {
		select {
		case success := <-p.Successes():
			return success, nil
		case fail := <-p.Errors():
			return nil, fail.Err
		case <-time.After(timeOut):
			return nil, errors.New(fmt.Sprintf("[kafka]async send message time out, send msg topic: %s, key: %s, value: %s", topic, key, value))
		}
	}

}

// Close 关闭生产者
func (p *Producer) Close() error {
	_ = p.SyncProducer.Close()

	if err := p.AsyncProducer.Close(); err != nil {
		return err
	}

	return nil
}
