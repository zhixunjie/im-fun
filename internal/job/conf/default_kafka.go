package conf

func DefaultKafka() []Kafka {
	return []Kafka{
		{
			Topic:   "im_push",
			Brokers: []string{"127.0.0.1:9092"},
		},
	}
}

// Kafka .
type Kafka struct {
	Topic   string
	Brokers []string
}
