package kafka

import "github.com/Shopify/sarama"

type ProducerConf struct {
	Topic   string   `yaml:"topic"`
	Brokers []string `yaml:"brokers"`
}
type ConsumerGroupConf struct {
	Topic   string   `yaml:"topic"`
	Brokers []string `yaml:"brokers"`
	GroupId string   `yaml:"groupId"`
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
