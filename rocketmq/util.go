package rocketmq

import (
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"strings"
)

func (msg *Message) String() string {
	return fmt.Sprintf("[rocketrmq]Topic: %s, tags: %s, keys: %s, body: %s, property: %v.",
		msg.Topic, msg.Tags, msg.Keys, string(msg.Body), msg.Property)
}

func (msgExt *MessageExt) String() string {
	return fmt.Sprintf("[rocketrmq]Message=%s, MsgId=%s, OffsetMsgId=%s, StoreSize=%d, QueueOffset=%d, SysFlag=%d, "+
		"BornTimestamp=%d, BornHost='%s', StoreTimestamp=%d, StoreHost='%s', CommitLogOffset=%d, BodyCRC=%d, "+
		"ReconsumeTimes=%d, PreparedTransactionOffset=%d.", msgExt.Message.String(), msgExt.MsgId, msgExt.OffsetMsgId,
		msgExt.StoreSize, msgExt.QueueOffset, msgExt.SysFlag, msgExt.BornTimestamp, msgExt.BornHost,
		msgExt.StoreTimestamp, msgExt.StoreHost, msgExt.CommitLogOffset, msgExt.BodyCRC, msgExt.ReconsumeTimes,
		msgExt.PreparedTransactionOffset)
}

// 消息处理函数
// 如果消费处理成功返回nil；消费处理失败返回err，此时会触发消费重试。
type MessageExtHandler func(*MessageExt) error

// 消息转换
func convertToMessageExt(msg *primitive.MessageExt) *MessageExt {
	return &MessageExt{
		Message: Message{
			Topic:    msg.Topic,
			Tags:     msg.GetTags(),
			Keys:     strings.Split(msg.GetKeys(), primitive.PropertyKeySeparator),
			Body:     msg.Body,
			Property: msg.GetProperties(),
		},
		MsgId:                     msg.MsgId,
		OffsetMsgId:               msg.OffsetMsgId,
		StoreSize:                 msg.StoreSize,
		QueueOffset:               msg.QueueOffset,
		SysFlag:                   msg.SysFlag,
		BornTimestamp:             msg.BornTimestamp,
		BornHost:                  msg.BornHost,
		StoreTimestamp:            msg.StoreTimestamp,
		StoreHost:                 msg.StoreHost,
		CommitLogOffset:           msg.CommitLogOffset,
		BodyCRC:                   msg.BodyCRC,
		ReconsumeTimes:            msg.ReconsumeTimes,
		PreparedTransactionOffset: msg.PreparedTransactionOffset,
	}
}
