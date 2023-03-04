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

func (p *SyncProducer) SendStringMessage(topic, value string) error {
	message := &sarama.ProducerMessage{
		Key:   sarama.StringEncoder(Uuid()),
		Topic: topic,
		Value: sarama.StringEncoder(value),
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

func (p *SyncProducer) SendByteMessage(topic string, value []byte) error {
	message := &sarama.ProducerMessage{
		Key:   sarama.StringEncoder(Uuid()),
		Topic: topic,
		Value: sarama.ByteEncoder(value),
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

func (p *SyncProducer) SendProducerMessage(msg *sarama.ProducerMessage) error {
	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		logrus.Errorf("send kafka msg failed(%v),partition:%d offset:%d msg:%+v\n",
			err, partition, offset, msg)
		return err
	}
	logrus.Infof("send kafka msg success,partition:%d offset:%d msg:%+v\n",
		partition, offset, msg)

	return nil
}

func (p *SyncProducer) Close() {
	if err := p.producer.Close(); err != nil {
		logrus.Errorf("Close error：err=%v", err)
		return
	}
}
