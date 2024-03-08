package data

import (
	"github.com/Shopify/sarama"
	pb "github.com/zhixunjie/im-fun/api/pb"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/request"
	"github.com/zhixunjie/im-fun/pkg/utils"
	"google.golang.org/protobuf/proto"
)

func (d *Data) KafkaSendToUsers(serverId string, tcpSessionIds []string, subId int32, msg []byte) (err error) {
	protoMsg := &pb.KafkaSendMsg{
		Type:          pb.KafkaSendMsg_ToUsers,
		SubId:         subId,
		ServerId:      serverId,
		TcpSessionIds: tcpSessionIds,
		Msg:           msg,
	}
	buf, err := proto.Marshal(protoMsg)
	if err != nil {
		return
	}
	err = d.KafkaProducer.SendProducerMessage(&sarama.ProducerMessage{
		Key:   sarama.StringEncoder(tcpSessionIds[0]),
		Topic: d.conf.Kafka[0].Topic,
		Value: sarama.ByteEncoder(buf),
	})

	return err
}

func (d *Data) KafkaSendToRoom(req *request.SendToRoomReq) (err error) {
	protoMsg := &pb.KafkaSendMsg{
		Type:   pb.KafkaSendMsg_ToRoom,
		SubId:  req.SubId,
		RoomId: utils.EncodeRoomKey(req.RoomType, req.RoomId),
		Msg:    []byte(req.Message),
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

func (d *Data) KafkaSendToAll(req *request.SendToAllReq) (err error) {
	protoMsg := &pb.KafkaSendMsg{
		Type:  pb.KafkaSendMsg_ToAll,
		SubId: req.SubId,
		Speed: req.Speed,
		Msg:   []byte(req.Message),
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
