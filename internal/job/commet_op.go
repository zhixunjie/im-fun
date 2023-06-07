package job

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/zhixunjie/im-fun/api/comet"
	"sync/atomic"
	"time"
)

func (c *Comet) PushUserKeys(arg *comet.PushUserKeysReq) (err error) {
	idx := atomic.AddUint64(&c.pushChanNum, 1) % c.routineNum
	c.userKeysChan[idx] <- arg
	return
}

func (c *Comet) PushUserRoom(arg *comet.PushUserRoomReq) (err error) {
	idx := atomic.AddUint64(&c.roomChanNum, 1) % c.routineNum
	c.userRoomChan[idx] <- arg
	return
}

func (c *Comet) PushUserAll(arg *comet.PushUserAllReq) (err error) {
	c.userAllChan <- arg
	return
}

func (c *Comet) Process(i int) {
	logHead := "Process|"

	// loop to send msg to allComet
	for {
		select {
		case <-c.ctx.Done():
			return
		case msg := <-c.userKeysChan[i]:
			_, err := c.rpcClient.PushUserKeys(context.Background(), msg)
			if err != nil {
				logrus.Errorf(logHead+"conf.rpcClient.SendToUsers(%s),serverId=%s,error=%v",
					msg, c.serverId, err)
			}
		case msg := <-c.userRoomChan[i]:
			_, err := c.rpcClient.PushUserRoom(context.Background(), msg)
			if err != nil {
				logrus.Errorf(logHead+"conf.rpcClient.BroadcastRoom(%s),serverId=%s,error=%v",
					msg, c.serverId, err)
			}
		case msg := <-c.userAllChan:
			_, err := c.rpcClient.PushUserAll(context.Background(), msg)
			if err != nil {
				logrus.Errorf(logHead+"conf.rpcClient.Broadcast(%s),serverId=%s,error=%v",
					msg, c.serverId, err)
			}
		}
	}
}

func (c *Comet) Close() (err error) {
	finish := make(chan bool)
	closePushChan := c.userKeysChan
	closeRoomChan := c.userRoomChan
	closeBroadcastChan := c.userAllChan

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
