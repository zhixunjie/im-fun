package dao

import (
	"fmt"
	"github.com/Shopify/sarama"
	"testing"
)

const (
	TEST_TOPIC  = "web_log"
	BROKER_ADDR = "127.0.0.1:9092"
)

func TestSyncProducer(t *testing.T) {
	// 1. new client
	client, err := sarama.NewSyncProducer([]string{BROKER_ADDR}, GetProducerConfig())
	if err != nil {
		fmt.Println("NewSyncProducer error: ", err)
		return
	}
	defer func() {
		if err = client.Close(); err != nil {
			fmt.Println("Close error：", err)
			return
		}
	}()

	// 2. build message
	// 构造消息
	kafkaMsg := &sarama.ProducerMessage{
		Topic: TEST_TOPIC,
		Value: sarama.StringEncoder("test kafka msg"),
	}

	// 3. send message
	// 报错：dial tcp 127.0.0.1:9092: connect: connection refused
	// 解决：重启kafka即可
	partition, offset, err := client.SendMessage(kafkaMsg)
	if err != nil {
		fmt.Println("SendMessage: ", err)
		return
	}

	fmt.Printf("Send kafka msg success,Partition:%d Offset:%d msg:%+v\n", partition, offset, kafkaMsg)
}
