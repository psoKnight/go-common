package rocketmq

import (
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/sirupsen/logrus"
	"sync"
)

// 生产者
type Consumer struct {
	cfg          *Config
	consumerMaps map[string]rocketmq.PushConsumer
	mutex        sync.Mutex
}

// 消费者
type Producer struct {
	rocketmq.Producer
	cfg   *Config
	mutex sync.Mutex
}

// 通用对象
type Config struct {
	Endpoints  []string
	BrokerAddr string
	InstanceID string
	AccessKey  string
	SecretKey  string
	RetryTimes int
	LogLevel   string
	Logger     *logrus.Logger
}

// RocketMQ
type RocketMq struct {
	producer *Producer
	consumer *Consumer
	cfg      *Config
}

// 通用消息配置（生产）
type Message struct {
	Topic    string
	Tags     string
	Keys     []string
	Body     []byte
	Property map[string]string
}

// 通用消息配置（消费）
type MessageExt struct {
	Message
	MsgId                     string
	OffsetMsgId               string
	StoreSize                 int32
	QueueOffset               int64
	SysFlag                   int32
	BornTimestamp             int64
	BornHost                  string
	StoreTimestamp            int64
	StoreHost                 string
	CommitLogOffset           int64
	BodyCRC                   int32
	ReconsumeTimes            int32
	PreparedTransactionOffset int64
}
