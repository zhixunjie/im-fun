package conf

func defaultConfig() *Config {
	return &Config{
		Debug: false,
		Discovery: &Discovery{
			Addr: "127.0.0.1:7171",
		},
		Kafka: DefaultKafka(),
	}
}
