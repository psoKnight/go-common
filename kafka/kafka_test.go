package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"testing"
	"time"
)

func TestSyncProducer(t *testing.T) {
	config := &Config{Endpoints: []string{"10.117.48.122:9092"}}

	kafka, err := NewKafka(config)
	if err != nil {
		t.Errorf("New kafka err: %v.", err)
		return
	}

	for i := 0; i <= 100000; i++ {
		s := &staff{
			Name: "张三",
			Age:  i,
		}

		marshal, err := json.Marshal(s)
		if err != nil {
			t.Errorf("Json marshal err: %v.", err)
			return
		}

		pid, offset, err := kafka.SyncSendMessage("test", "", string(marshal))
		if err != nil {
			t.Errorf("Sync send message err: %v.", err)
			continue
		}
		t.Log(i, pid, offset)

		time.Sleep(time.Second)
	}
}

func TestAsyncProducer(t *testing.T) {
	config := &Config{Endpoints: []string{"10.117.48.122:9092"}}

	kafka, err := NewKafka(config)
	if err != nil {
		t.Errorf("New kafka err: %v.", err)
		return
	}

	for i := 20; i <= 100000; i++ {
		s := &staff{
			Name: "张三",
			Age:  i,
		}

		marshal, err := json.Marshal(s)
		if err != nil {
			t.Errorf("Json marshal err: %v.", err)
			return
		}

		message, err := kafka.AsyncSendMessage("test", "", string(marshal), time.Duration(3)*time.Second)
		if err != nil {
			t.Errorf("Async send message err: %v.", err)
			continue
		}
		t.Log(message)

		time.Sleep(time.Second)
	}
}

func TestConsumer(t *testing.T) {
	config := &Config{Endpoints: []string{"10.117.48.122:9092"}}

	kafka, err := NewKafka(config)
	if err != nil {
		t.Errorf("New kafka err: %v.", err)
		return
	}

	ch := make(chan *sarama.ConsumerMessage)
	err = kafka.ConsumeMessage("test", "", -1, ch)
	if err != nil {
		t.Errorf("Kafka consume message err: %v.", err)
		return
	}

	for {
		select {
		case msg := <-ch:
			t.Logf("Partition: %d, offset: %d, key: %s, value: %s.", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
		case <-time.After(time.Duration(3) * time.Second):
			t.Log("After 1 second.")
		}
	}
}

func TestConsumerByGroup(t *testing.T) {
	config := &Config{Endpoints: []string{"10.117.48.122:9092"}}

	kafka, err := NewKafka(config)
	if err != nil {
		t.Errorf("New kafka err: %v.", err)
		return
	}

	handler := &ConsumerGroupHandler{} // 自定义handler
	err = kafka.ConsumeMessageByGroup([]string{"test"}, "group_b", -1, handler)
	if err != nil {
		t.Errorf("Kafka consume message err: %v.", err)
		return
	}
}

// ConsumerGroupHandler 实现github.com/Shopify/sarama/consumer_group.go/ConsumerGroupHandler 接口
type ConsumerGroupHandler struct {
}

func (ConsumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	//session.ResetOffset("test", 0, 0, "") // 重置偏移
	return nil
}

func (ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// 消费消息
func (cgh ConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// 获取消息
	for {
		select {
		case msg := <-claim.Messages():
			fmt.Println(fmt.Sprintf("Message topic:%q, partition:%d, offset:%d, key: %s, value: %s.",
				msg.Topic,
				msg.Partition,
				msg.Offset,
				string(msg.Key),
				string(msg.Value)))

			// 将消息标记为已使用
			//sess.MarkMessage(msg, "")

		case <-time.After(time.Duration(3) * time.Second):
			fmt.Println("After 1 second.")
		}
	}
	return nil
}

type staff struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}
