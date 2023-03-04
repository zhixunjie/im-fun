package kafka

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

type ConsumerGroup struct {
	conf          *ConsumerGroupConf
	consumerGroup sarama.ConsumerGroup
}

func NewConsumerGroup(conf *ConsumerGroupConf, fn Callback) (*ConsumerGroup, error) {
	cg := &ConsumerGroup{
		conf:          conf,
		consumerGroup: nil,
	}

	// 设置消费组的信息
	consumerGroup, err := sarama.NewConsumerGroup(conf.Brokers, conf.GroupId, GetConsumerGroupConfig())
	if err != nil {
		logrus.Errorf("NewConsumerGroup,err=%v,conf=%+v", err, conf)
		return cg, err
	}
	cg.consumerGroup = consumerGroup

	// get errors
	go func() {
		var err error
		for err = range consumerGroup.Errors() {
			logrus.Errorf("consumerGroup.Errors,err=%v,conf=%+v", err, conf)
		}
	}()

	handler := consumerGroupHandler{fn: fn}
	ctx := context.Background()
	topics := []string{conf.Topic}
	go func() {
		var err error
		for {
			//fmt.Println("waiting for message......")
			logrus.Infof("waiting for message......,conf=%+v", conf)
			err = consumerGroup.Consume(ctx, topics, handler)
			if err != nil {
				logrus.Errorf("err=%v,conf=%+v", err, conf)
				return
			}
		}
	}()
	return cg, nil
}

func (c *ConsumerGroup) Close() {
	if err := c.consumerGroup.Close(); err != nil {
		logrus.Errorf("Close error：err=%v,conf=%+v", err, c.conf)
		return
	}
}

type Callback func(message *sarama.ConsumerMessage)

type consumerGroupHandler struct {
	fn Callback
}

func (consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil

}
func (consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}
func (h consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		logrus.Infof("Message topic:%q partition:%d offset:%d\n", msg.Topic, msg.Partition, msg.Offset)
		func() {
			defer func() {
				if err := recover(); err != nil {
					logrus.Errorf("ConsumeClaim recover err=%v", err)
				}
			}()
			h.fn(msg)
		}()
		sess.MarkMessage(msg, "")
	}
	return nil
}
