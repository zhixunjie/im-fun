package dao

import (
	"github.com/Shopify/sarama"
	pb "github.com/zhixunjie/im-fun/api/logic"
	"github.com/zhixunjie/im-fun/internal/logic/model/request"
	"github.com/zhixunjie/im-fun/pkg/utils"
	"google.golang.org/protobuf/proto"
)

func (d *Dao) KafkaSendToUserKeys(serverId string, userKeys []string, subId int32, msg []byte) (err error) {
	protoMsg := &pb.PushMsg{
		Type:     pb.PushMsg_UserKeys,
		SubId:    subId,
		ServerId: serverId,
		UserKeys: userKeys,
		Msg:      msg,
	}
	buf, err := proto.Marshal(protoMsg)
	if err != nil {
		return
	}
	err = d.KafkaProducer.SendProducerMessage(&sarama.ProducerMessage{
		Key:   sarama.StringEncoder(userKeys[0]),
		Topic: d.conf.Kafka[0].Topic,
		Value: sarama.ByteEncoder(buf),
	})

	return err
}

func (d *Dao) KafkaSendToRoom(req *request.SendToRoomReq) (err error) {
	protoMsg := &pb.PushMsg{
		Type:   pb.PushMsg_UserRoom,
		SubId:  req.SubId,
		RoomId: utils.EncodeRoomKey(req.RoomType, req.RoomId),
		Msg:    req.Message,
	}
	buf, err := proto.Marshal(protoMsg)
	if err != nil {
		return
	}
	err = d.KafkaProducer.SendProducerMessage(&sarama.ProducerMessage{
		Key:   sarama.StringEncoder(protoMsg.RoomId),
		Topic: d.conf.Kafka[0].Topic,
		Value: sarama.ByteEncoder(buf),
	})

	return err
}

func (d *Dao) KafkaSendToAll(req *request.SendToAllReq) (err error) {
	protoMsg := &pb.PushMsg{
		Type:  pb.PushMsg_UserAll,
		SubId: req.SubId,
		Speed: req.Speed,
		Msg:   req.Message,
	}
	buf, err := proto.Marshal(protoMsg)
	if err != nil {
		return
	}
	err = d.KafkaProducer.SendProducerMessage(&sarama.ProducerMessage{
		Key:   sarama.StringEncoder(protoMsg.RoomId),
		Topic: d.conf.Kafka[0].Topic,
		Value: sarama.ByteEncoder(buf),
	})

	return err
}
