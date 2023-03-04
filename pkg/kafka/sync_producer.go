package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

type SyncProducer struct {
	conf     *ProducerConf
	producer sarama.SyncProducer
}

func NewSyncProducer(conf *ProducerConf) (SyncProducer, error) {
	p := SyncProducer{
		conf:     conf,
		producer: nil,
	}
	producer, err := sarama.NewSyncProducer(conf.Brokers, getProducerConfig())
	if err != nil {
		logrus.Errorf("NewSyncProducer,err=%v,conf=%+v", err, p.conf)
		return p, err
	}
	p.producer = producer

	return p, nil
}

func (p *SyncProducer) SendMessage(string string) error {
	message := &sarama.ProducerMessage{
		Topic: p.conf.Topic,
		Value: sarama.StringEncoder(string),
	}
	partition, offset, err := p.producer.SendMessage(message)
	if err != nil {
		logrus.Errorf("send kafka msg failed(%v),partition:%d offset:%d msg:%+v\n",
			err, partition, offset, message)
		return err
	}
	logrus.Infof("send kafka msg success,partition:%d offset:%d msg:%+v\n",
		partition, offset, message)

	return nil
}

func (p *SyncProducer) Close() {
	if err := p.producer.Close(); err != nil {
		logrus.Errorf("Close error：err=%v", err)
		return
	}
}
