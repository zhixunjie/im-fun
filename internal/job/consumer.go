package job

import (
	"context"
	"github.com/Shopify/sarama"
	pb "github.com/zhixunjie/im-fun/api/logic"
	"google.golang.org/protobuf/proto"
)

// Consume messages, watch signals
func (job *Job) Consume(msg *sarama.ConsumerMessage) {
	message := new(pb.PushMsg)
	if err := proto.Unmarshal(msg.Value, message); err != nil {
	}
	if err := job.push(context.Background(), message); err != nil {
	}
}

func (job *Job) Close() error {
	return nil
}
