package kafka

import (
	"fmt"
	"testing"
)

func TestProducer(t *testing.T) {
	producer, err := NewSyncProducer(&ProducerConf{
		Topic:   "my_topic",
		Brokers: []string{"127.0.0.1:9092"},
	})
	if err != nil {
		panic(err)
	}
	err = producer.SendMessage("hello i am syncProducer")
	fmt.Println(err)
}
