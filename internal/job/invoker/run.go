package invoker

import (
	"context"
	"fmt"
	"github.com/zhixunjie/im-fun/api/pb"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"time"
)

func (c *CometInvoker) Run(i int) {
	logHead := fmt.Sprintf("Run[%v]|", i)

	// loop to send msg to allComet
	for {
		select {
		case <-c.ctx.Done():
			logging.Info(logHead + "routine exit,because ctx.Done()")
			return
		case msg := <-c.chUsers[i]:
			_, err := c.rpcClient.SendToUsers(context.Background(), msg)
			if err != nil {
				logging.Errorf(logHead+"rpcClient.SendToUsers(%s),serverId=%s,error=%v", msg, c.serverId, err)
			}
		case msg := <-c.chRoom[i]:
			_, err := c.rpcClient.SendToRoom(context.Background(), msg)
			if err != nil {
				logging.Errorf(logHead+"rpcClient.SendToRoom(%s),serverId=%s,error=%v", msg, c.serverId, err)
			}
		case msg := <-c.chAll:
			_, err := c.rpcClient.SendToAll(context.Background(), msg)
			if err != nil {
				logging.Errorf(logHead+"rpcClient.SendToAll(%s),serverId=%s,error=%v", msg, c.serverId, err)
			}
		}
	}
}

func (c *CometInvoker) SendToUsers(arg *pb.SendToUsersReq) (err error) {
	idx := c.counterToUser.Add(1) % c.RoutineNum
	c.chUsers[idx] <- arg
	return
}

func (c *CometInvoker) SendToRoom(arg *pb.SendToRoomReq) (err error) {
	idx := c.counterToRoom.Add(1) % c.RoutineNum
	c.chRoom[idx] <- arg
	return
}

func (c *CometInvoker) SendToAll(arg *pb.SendToAllReq) (err error) {
	c.chAll <- arg
	return
}

// Close 关闭（暂时没有任何的调用地方）
// - 确保所有channel的内容消费完毕后，再继续往下走
// - 如果超时时间内还没消费完，也继续往下走
func (c *CometInvoker) Close() (err error) {
	finish := make(chan bool)
	sendUserChan := c.chUsers
	sendRoomChan := c.chRoom
	sendAllChan := c.chAll

	go func() {
		for {
			n := len(sendAllChan)
			for _, ch := range sendUserChan {
				n += len(ch)
			}
			for _, ch := range sendRoomChan {
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
		logging.Info("close all CometInvoker finish")
	case <-time.After(5 * time.Second):
		format := "serverId=%v,sendUserChan=%d,sendRoomChan=%d,sendAllChan=%d"
		err = fmt.Errorf("close all CometInvoker timeout"+format, c.serverId, len(sendUserChan), len(sendRoomChan), len(sendAllChan))
		logging.Error(err)
	}
	c.Cancel()
	return
}
