package rocketmq

// Message 通用消息配置（生产）
type Message struct {
	Topic    string
	Tags     string
	Keys     []string
	Body     []byte
	Property map[string]string
}

// MessageExt 通用消息配置（消费）
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
