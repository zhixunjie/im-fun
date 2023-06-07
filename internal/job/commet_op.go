package job

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/zhixunjie/im-fun/api/comet"
	"sync/atomic"
	"time"
)

func (c *Comet) SendToUserKeys(arg *comet.SendToUserKeysReq) (err error) {
	idx := atomic.AddUint64(&c.pushChanNum, 1) % c.routineNum
	c.chUserKeys[idx] <- arg
	return
}

func (c *Comet) SendToRoom(arg *comet.SendToRoomReq) (err error) {
	idx := atomic.AddUint64(&c.roomChanNum, 1) % c.routineNum
	c.chRoom[idx] <- arg
	return
}

func (c *Comet) SendToAll(arg *comet.SendToAllReq) (err error) {
	c.chAll <- arg
	return
}

func (c *Comet) Process(i int) {
	logHead := "Process|"

	// loop to send msg to allComet
	for {
		select {
		case <-c.ctx.Done():
			return
		case msg := <-c.chUserKeys[i]:
			_, err := c.rpcClient.SendToUserKeys(context.Background(), msg)
			if err != nil {
				logrus.Errorf(logHead+"conf.rpcClient.SendToUsers(%s),serverId=%s,error=%v",
					msg, c.serverId, err)
			}
		case msg := <-c.chRoom[i]:
			_, err := c.rpcClient.SendToRoom(context.Background(), msg)
			if err != nil {
				logrus.Errorf(logHead+"conf.rpcClient.BroadcastRoom(%s),serverId=%s,error=%v",
					msg, c.serverId, err)
			}
		case msg := <-c.chAll:
			_, err := c.rpcClient.SendToAll(context.Background(), msg)
			if err != nil {
				logrus.Errorf(logHead+"conf.rpcClient.Broadcast(%s),serverId=%s,error=%v",
					msg, c.serverId, err)
			}
		}
	}
}

func (c *Comet) Close() (err error) {
	finish := make(chan bool)
	closePushChan := c.chUserKeys
	closeRoomChan := c.chRoom
	closeBroadcastChan := c.chAll

	go func() {
		for {
			n := len(closeBroadcastChan)
			for _, ch := range closePushChan {
				n += len(ch)
			}
			for _, ch := range closeRoomChan {
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
		logrus.Info("close allComet finish")
	case <-time.After(5 * time.Second):
		err = fmt.Errorf("close allComet(server:%s push:%d room:%d broadcast:%d) timeout",
			c.serverId, len(closePushChan), len(closeRoomChan), len(closeBroadcastChan))
	}
	c.cancel()
	return
}
