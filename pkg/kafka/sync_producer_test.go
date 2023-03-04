package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"testing"
)

func TestSendStringMessage(t *testing.T) {
	producer, err := NewSyncProducer(&ProducerConf{
		Topic:   "my_topic",
		Brokers: []string{"127.0.0.1:9092"},
	})
	if err != nil {
		panic(err)
	}
	err = producer.SendStringMessage(producer.conf.Topic, "hello i am syncProducer")
	fmt.Println(err)
}

func TestSendByteMessage(t *testing.T) {
	producer, err := NewSyncProducer(&ProducerConf{
		Topic:   "my_topic",
		Brokers: []string{"127.0.0.1:9092"},
	})
	if err != nil {
		panic(err)
	}
	err = producer.SendByteMessage(producer.conf.Topic, []byte("hello i am syncProducer"))
	fmt.Println(err)
}

func TestSendProducerMessage(t *testing.T) {
	producer, err := NewSyncProducer(&ProducerConf{
		Topic:   "my_topic",
		Brokers: []string{"127.0.0.1:9092"},
	})
	if err != nil {
		panic(err)
	}
	err = producer.SendProducerMessage(&sarama.ProducerMessage{
		Key:   sarama.StringEncoder(Uuid()),
		Topic: producer.conf.Topic,
		Value: sarama.ByteEncoder("hello i am syncProducer"),
	})
	fmt.Println(err)
}
