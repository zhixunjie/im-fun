package connect

import (
	pb "github.com/zhixunjie/im-fun/api/comet"
	"github.com/zhixunjie/im-fun/api/protocol"
	"github.com/zhixunjie/im-fun/internal/comet/connect/channel"
	"sync/atomic"
	"time"
)

// BroadcastToAllBucket 消息广播：广播到所有Bucket的所有Channel
func BroadcastToAllBucket(srv *Server, proto *protocol.Proto, speed int) {
	// TODO 使用队列进行广播的处理
	go func() {
		for _, bucket := range srv.Buckets() {
			bucket.Broadcast(proto)
			if speed > 0 {
				// 如果100个连接（channel），速度等于5。
				// 那么就会，每次广播一个bucket，广播完毕后，睡眠20秒。
				t := time.Duration(bucket.ChannelCount() / speed)
				time.Sleep(t * time.Second)
			}
		}
	}()
}

// Broadcast 消息广播：发送给当前Bucket的所有Channel
func (b *Bucket) Broadcast(p *protocol.Proto) {
	var ch *channel.Channel

	b.rwLock.RLock()
	for _, ch = range b.chs {
		_ = ch.Push(p)
	}
	b.rwLock.RUnlock()
}

// SendToTheRoom 发送一个Proto到某个房间ID（PushToAllChan => Proto => ROOM）
func (b *Bucket) SendToTheRoom(req *pb.BroadcastRoomReq) {
	num := atomic.AddUint64(&b.routineCounter, 1) % b.conf.RoutineAmount
	b.routines[num] <- req
}

// ProcessProtoToRoom 处理发送都某个房间ID的消息 （Pop => ROOM => Proto => DEAL）
func (b *Bucket) ProcessProtoToRoom(c chan *pb.BroadcastRoomReq) {
	for {
		req := <-c
		room := b.GetRoomById(req.RoomId)
		if room != nil {
			room.PushToAllChan(req.Proto)
		}
	}
}
