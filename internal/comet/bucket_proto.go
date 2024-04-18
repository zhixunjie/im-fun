package comet

import (
	pb "github.com/zhixunjie/im-fun/api/pb"
	"github.com/zhixunjie/im-fun/api/protocol"
	"github.com/zhixunjie/im-fun/internal/comet/channel"
	"time"
)

// BroadcastToAllBucket
// 消息广播：广播到所有Bucket的所有Channel
func (s *TcpServer) BroadcastToAllBucket(proto *protocol.Proto, speed int) {
	// TODO 使用队列进行广播的处理
	go func() {
		for _, bucket := range s.Buckets() {
			bucket.broadcast(proto)
			// 如果100个连接（channel），speed等于5。每次广播一个bucket，广播后睡眠20秒。
			if speed > 0 {
				t := time.Duration(bucket.ChannelCount() / speed)
				time.Sleep(t * time.Second)
			}
		}
	}()
}

// broadcast
// 消息广播：发送给当前Bucket的所有Channel
func (b *Bucket) broadcast(p *protocol.Proto) {
	var ch *channel.Channel

	b.rwLock.RLock()
	for _, ch = range b.chs {
		_ = ch.Push(p)
	}
	b.rwLock.RUnlock()
}

// BroadcastRoom
// 房间广播：发送一个Proto到某个房间ID（即：SendToAllChan => Proto => ROOM）
func (b *Bucket) BroadcastRoom(req *pb.SendToRoomReq) {
	num := b.routineCounter.Add(1) % uint64(b.conf.RoutineAmount)
	b.routines[num] <- req
}

// ProcessProtoToRoom
// 处理发送都某个房间ID的消息 （Pop => ROOM => Proto => DEAL）
func (b *Bucket) ProcessProtoToRoom(index int) {
	for {
		req := <-b.routines[index]
		room := b.GetRoomById(req.RoomId)
		if room != nil {
			room.SendToAllChan(req.Proto)
		}
	}
}
