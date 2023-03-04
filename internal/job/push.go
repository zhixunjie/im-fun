package job

import (
	"context"
	"fmt"
	"github.com/zhixunjie/im-fun/api/comet"
	pb "github.com/zhixunjie/im-fun/api/logic"
	"github.com/zhixunjie/im-fun/api/protocol"
	"github.com/zhixunjie/im-fun/pkg/buffer"
)

func (job *Job) push(ctx context.Context, pushMsg *pb.PushMsg) (err error) {
	switch pushMsg.Type {
	case pb.PushMsg_PUSH:
		err = job.pushKeys(pushMsg.SubId, pushMsg.ServerId, pushMsg.UserKeys, pushMsg.Msg)
	//case pb.PushMsg_ROOM:
	//	err = job.getRoom(pushMsg.Room).Push(pushMsg.Operation, pushMsg.Msg)
	//case pb.PushMsg_BROADCAST:
	//	err = job.broadcast(pushMsg.Operation, pushMsg.Msg, pushMsg.Speed)
	default:
		err = fmt.Errorf("no match push type: %s", pushMsg.Type)
	}
	return
}

func (job *Job) pushKeys(op int32, serverId string, userKeys []string, body []byte) (err error) {
	buf := buffer.NewWriterSize(len(body) + 64)
	p := &protocol.Proto{
		Ver:  1,
		Op:   op,
		Body: body,
	}
	p.WriteTo(buf)
	p.Body = buf.Buffer()
	p.Op = int32(protocol.OpRaw)
	var args = comet.PushMsgReq{
		UserKeys: userKeys,
		ProtoOp:  op,
		Proto:    p,
	}
	if c, ok := job.cometServers[serverId]; ok {
		if err = c.Push(&args); err != nil {
		}
	}
	return
}
