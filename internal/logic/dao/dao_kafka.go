package dao

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/zhixunjie/im-fun/internal/logic/conf"
)

func newKafkaProducer(conf *conf.Kafka) (sarama.SyncProducer, error) {
	client, err := sarama.NewSyncProducer(conf.Brokers, getProducerConfig())
	if err != nil {
		fmt.Println("NewSyncProducer error: ", err)
		return client, err
	}

	return client, nil
}

func getProducerConfig() *sarama.Config {
	config := sarama.NewConfig()

	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner

	// 成功交付的消息将在success channel返回
	config.Producer.Return.Successes = true

	return config
}
