package kafka

import "github.com/Shopify/sarama"

type ProducerConf struct {
	Topic   string
	Brokers []string
}
type ConsumerGroupConf struct {
	Topic   string
	Brokers []string
	GroupId string
}

func getProducerConfig() *sarama.Config {
	config := sarama.NewConfig()

	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner

	// 成功交付的消息将在success channel返回
	config.Producer.Return.Successes = true

	return config
}

func GetConsumerGroupConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	return config
}
