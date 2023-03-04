package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"testing"
	"time"
)

func TestConsumerGroup(t *testing.T) {
	fn := func(message *sarama.ConsumerMessage) {
		fmt.Printf("fn get message,%s\n", message.Value)
	}
	_, err := NewConsumerGroup(&ConsumerGroupConf{
		Topic:   "my_topic",
		Brokers: []string{"127.0.0.1:9092"},
		GroupId: "my_group",
	}, fn)
	if err != nil {
		panic(err)
	}
	time.Sleep(10 * time.Second)
}
