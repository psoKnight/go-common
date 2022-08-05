package kafka

import (
	"context"
	"errors"
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

// Consumer 消费者
type Consumer struct {
	cfg *Config
	sarama.Consumer
}

// NewConsumer 新建消费者
func NewConsumer(cfg *Config) (*Consumer, error) {

	c := &Consumer{
		cfg: cfg,
	}

	// 消费者配置
	consumConfig := sarama.NewConfig()
	consumConfig.Consumer.Return.Errors = true
	consumConfig.Consumer.Offsets.AutoCommit.Enable = false // 取消自动指定提交消息
	// 设置使用的kafka版本,如果低于V0_10_0_0版本,消息中的timestrap 没有作用。需要消费和生产同时配置
	// 注意，版本设置不对的话，kafka 会返回很奇怪的错误，并且无法成功发送消息
	//consumConfig.Version = sarama.V0_10_2_1

	consumer, err := sarama.NewConsumer(cfg.Endpoints, consumConfig)
	if err != nil {
		return nil, err
	}
	c.Consumer = consumer
	return c, nil
}

// ConsumeMessage 消费者消费数据
/**
offsetType
	-1：OffsetNewest 代表日志头偏移量，即将分配给将要生成到分区的下一条消息的偏移量
	-2：OffsetOldest 代表代理上可用于分区的最旧偏移量
*/
func (c *Consumer) ConsumeMessage(topic, key string, offsetType int64, ch chan *sarama.ConsumerMessage) error {
	if topic == "" {
		return errors.New("[kafka]topic is '', please check")
	}
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic

	if key != "" {
		msg.Key = sarama.StringEncoder(key)
	}

	// 获取偏移方式
	if offsetType == -1 {
		offsetType = sarama.OffsetNewest
	} else if offsetType == -2 {
		offsetType = sarama.OffsetOldest
	} else {
		offsetType = -1 // 默认
	}

	// 设置分区
	partitions, err := c.Consumer.Partitions(topic)
	if err != nil {
		return err
	}

	logrus.Infof("[kafka]topic '%s' all partitions: %v.", topic, partitions)

	// 循环分区
	for partition := range partitions {
		pc, err := c.Consumer.ConsumePartition(topic, int32(partition), offsetType)
		if err != nil {
			logrus.Infof("[kafka]topic '%s' consume partition '%d' err: %v, continue.", topic, partition, err)
			continue
		}

		go func(pc sarama.PartitionConsumer, c chan *sarama.ConsumerMessage) {
			for msg := range pc.Messages() {
				c <- msg
			}
			defer pc.AsyncClose()
		}(pc, ch)

	}

	return nil
}

// ConsumeMessageByGroup 消费者消费数据(通过消费组)
/**
offsetType
	-1：OffsetNewest 代表日志头偏移量，即将分配给将要生成到分区的下一条消息的偏移量
	-2：OffsetOldest 代表代理上可用于分区的最旧偏移量
*/
func (c *Consumer) ConsumeMessageByGroup(topics []string, group string, offsetType int64, handler sarama.ConsumerGroupHandler) error {
	if len(topics) == 0 {
		return errors.New("[kafka]len topic is 0, please check")
	}

	// 消费者配置
	consumConfig := sarama.NewConfig()
	consumConfig.Consumer.Return.Errors = true
	consumConfig.Consumer.Offsets.AutoCommit.Enable = false // 取消自动指定提交消息
	//consuConfig.Version = sarama.V0_10_2_1

	// 获取偏移方式
	if offsetType == -1 {
		offsetType = sarama.OffsetNewest
	} else if offsetType == -2 {
		offsetType = sarama.OffsetOldest
	} else {
		offsetType = sarama.OffsetNewest // 默认
	}
	consumConfig.Consumer.Offsets.Initial = offsetType

	gr, err := sarama.NewConsumerGroup(c.cfg.Endpoints, group, consumConfig)
	if err != nil {
		return nil
	}
	defer gr.Close()

	ctx := context.Background()

	// 启动kafka 消费组模式
	err = gr.Consume(ctx, topics, handler)
	if err != nil {
		return err
	}

	return nil
}

// Close 关闭消费者
func (c *Consumer) Close() error {
	err := c.Consumer.Close()
	if err != nil {
		return err
	}

	return nil
}
