package job

import (
	"github.com/sirupsen/logrus"
	"github.com/zhixunjie/im-fun/api/comet"
	"github.com/zhixunjie/im-fun/api/protocol"
	"github.com/zhixunjie/im-fun/pkg/buffer"
)

func (job *Job) PushUserKeys(subId int32, serverId string, userKeys []string, message []byte) (err error) {
	logHead := "PushUserKeys|"

	// write to proto body（proto的body里面嵌套proto）
	// 这样写的好处见：Job.PushUserRoom
	proto := &protocol.Proto{
		Ver:  protocol.ProtoVersion,
		Op:   int32(protocol.OpBatchMsg),
		Body: message,
	}
	batchMessage := buffer.NewWriterSize(len(message) + 64)
	proto.WriteTo(batchMessage)
	proto.Body = batchMessage.Buffer()

	// push to comet
	if cm, ok := job.allComet[serverId]; ok {
		params := comet.PushUserKeysReq{
			UserKeys: userKeys,
			Proto:    proto,
			SubId:    subId,
		}
		if err = cm.PushUserKeys(&params); err != nil {
			logrus.Errorf(logHead+"PushToChannel err=%v,serverId=%v,params=%+v", err, serverId, params)
		}
	}
	return
}

func (job *Job) PushUserRoom(subId int32, roomId string, batchMessage []byte) (err error) {
	logHead := "PushUserRoom|"

	// write to proto batchMessage（proto的body里面嵌套proto）
	// - 通过嵌套的proto，使得Body里面能够存放多条的Proto（即：批量发送Proto）
	// - 客户端对于此类Op的数据包，会使用循环对Body进行解包
	proto := &protocol.Proto{
		Ver:  protocol.ProtoVersion,
		Op:   int32(protocol.OpBatchMsg),
		Body: batchMessage,
	}

	// push to every comet
	for serverId, cm := range job.allComet {
		params := comet.PushUserRoomReq{
			RoomId: roomId,
			Proto:  proto,
		}
		if err = cm.PushUserRoom(&params); err != nil {
			logrus.Errorf(logHead+"PushToChannel err=%v,serverId=%v,params=%+v", err, serverId, params)
		}
	}
	return
}

func (job *Job) PushUserAll(subId int32, speed int32, message []byte) (err error) {
	logHead := "PushUserAll|"

	// write to proto body（proto的body里面嵌套proto）
	// 通过嵌套的proto，使得Body里面能够存放多条的Proto（即：批量发送Proto）
	proto := &protocol.Proto{
		Ver:  protocol.ProtoVersion,
		Op:   int32(protocol.OpBatchMsg),
		Body: message,
	}

	// push to every comet
	speed = speed / int32(len(job.allComet))
	for serverId, cm := range job.allComet {
		params := comet.PushUserAllReq{
			Proto: proto,
			SubId: subId,
			Speed: speed,
		}
		if err = cm.PushUserAll(&params); err != nil {
			logrus.Errorf(logHead+"PushToChannel err=%v,serverId=%v,params=%+v", err, serverId, params)
		}
	}
	return
}

func (job *Job) Close() {
	if job.consumer != nil {
		job.consumer.Close()
	}
}
