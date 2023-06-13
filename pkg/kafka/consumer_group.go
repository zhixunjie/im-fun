package kafka

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/zhixunjie/im-fun/pkg/logging"
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
		logging.Errorf("NewConsumerGroup,err=%v,conf=%+v", err, conf)
		return cg, err
	}
	cg.consumerGroup = consumerGroup

	// get errors
	go func() {
		var newErr error
		for newErr = range consumerGroup.Errors() {
			logging.Errorf("consumerGroup.Errors,err=%v,conf=%+v", newErr, conf)
		}
	}()

	handler := consumerGroupHandler{fn: fn}
	ctx := context.Background()
	topics := []string{conf.Topic}
	go func() {
		logHead := "KafkaConsumer|"
		var newErr error
		for {
			//fmt.Println("waiting for message......")
			logging.Infof(logHead+"conf=%+v,waiting for message......", conf)
			newErr = consumerGroup.Consume(ctx, topics, handler)
			if newErr != nil {
				logging.Errorf(logHead+"err=%v,conf=%+v", err, conf)
				return
			}
		}
	}()
	return cg, nil
}

func (c *ConsumerGroup) Close() {
	if err := c.consumerGroup.Close(); err != nil {
		logging.Errorf("Close error：err=%v,conf=%+v", err, c.conf)
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
		logging.Infof("Message topic:%q partition:%d offset:%d", msg.Topic, msg.Partition, msg.Offset)
		func() {
			defer func() {
				if err := recover(); err != nil {
					logging.Errorf("ConsumeClaim recover err=%v", err)
				}
			}()
			h.fn(msg)
		}()
		sess.MarkMessage(msg, "")
	}
	return nil
}
