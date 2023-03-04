package job

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/zhixunjie/im-fun/api/comet"
	"sync/atomic"
	"time"
)

func (c *Comet) Push(arg *comet.PushMsgReq) (err error) {
	idx := atomic.AddUint64(&c.pushChanNum, 1) % c.routineNum
	c.pushChan[idx] <- arg
	return
}

func (c *Comet) BroadcastRoom(arg *comet.BroadcastRoomReq) (err error) {
	idx := atomic.AddUint64(&c.roomChanNum, 1) % c.routineNum
	c.roomChan[idx] <- arg
	return
}

func (c *Comet) Broadcast(arg *comet.BroadcastReq) (err error) {
	c.broadcastChan <- arg
	return
}

func (c *Comet) Process(i int) {
	logHead := "Process|"
	pushChan := c.pushChan[i]
	roomChan := c.roomChan[i]
	broadcastChan := c.broadcastChan

	// loop to send msg to comet
	for {
		select {
		case <-c.ctx.Done():
			return
		case broadcast := <-broadcastChan:
			_, err := c.rpcClient.Broadcast(context.Background(), &comet.BroadcastReq{
				Proto:   broadcast.Proto,
				ProtoOp: broadcast.ProtoOp,
				Speed:   broadcast.Speed,
			})
			if err != nil {
				logrus.Errorf(logHead+"c.rpcClient.Broadcast(%s, reply) serverId:%s error(%v)", broadcast, c.serverId, err)
			}
		case room := <-roomChan:
			_, err := c.rpcClient.BroadcastRoom(context.Background(), &comet.BroadcastRoomReq{
				RoomId: room.RoomId,
				Proto:  room.Proto,
			})
			if err != nil {
				logrus.Errorf(logHead+"c.rpcClient.BroadcastRoom(%s, reply) serverId:%s error(%v)", room, c.serverId, err)
			}
		case push := <-pushChan:
			_, err := c.rpcClient.PushMsg(context.Background(), &comet.PushMsgReq{
				UserKeys: push.UserKeys,
				Proto:    push.Proto,
				ProtoOp:  push.ProtoOp,
			})
			if err != nil {
				logrus.Errorf(logHead+"c.rpcClient.PushMsg(%s, reply) serverId:%s error(%v)", push, c.serverId, err)
			}
		}
	}
}

func (c *Comet) Close() (err error) {
	finish := make(chan bool)
	go func() {
		for {
			n := len(c.broadcastChan)
			for _, ch := range c.pushChan {
				n += len(ch)
			}
			for _, ch := range c.roomChan {
				n += len(ch)
			}
			if n == 0 {
				finish <- true
				return
			}
			time.Sleep(time.Second)
		}
	}()
	select {
	case <-finish:
		logrus.Info("close comet finish")
	case <-time.After(5 * time.Second):
		err = fmt.Errorf("close comet(server:%s push:%d room:%d broadcast:%d) timeout",
			c.serverId, len(c.pushChan), len(c.roomChan), len(c.broadcastChan))
	}
	c.cancel()
	return
}
