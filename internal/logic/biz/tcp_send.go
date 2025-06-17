package biz

import (
	"context"
	"errors"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/request"
	"github.com/zhixunjie/im-fun/pkg/logging"
)

// SendToUsers 发送消息（by kafka）
func (bz *Biz) SendToUsers(ctx context.Context, req *request.SendToUsersReq) error {
	logHead := "SendToUsers|"
	serverIds, err := bz.data.GetServerIds(ctx, req.TcpSessionIds)
	if err != nil {
		logging.Errorf(logHead+"res=%v, err=%v", serverIds, err)
		return err
	}

	// 重整数据：获取某个serverId下的tcpSessionId
	serverIdMap := make(map[string][]string)
	for i, tcpSessionId := range req.TcpSessionIds {
		serverId := serverIds[i]
		serverIdMap[serverId] = append(serverIdMap[serverId], tcpSessionId)
	}

	// 把同一台机器的请求聚合到一起
	for serverId := range serverIdMap {
		err = bz.data.KafkaSendToUsers(serverId, serverIdMap[serverId], req.SubId, []byte(req.Message))
		if err != nil {
			logging.Errorf(logHead+"err=%v", err)
		}
	}

	return nil
}

// SendToUsersByIds 发送消息（by kafka）
func (bz *Biz) SendToUsersByIds(ctx context.Context, req *request.SendToUsersByIdsReq) error {
	logHead := "SendToUsersByIds|"

	// get: data
	mSession, err := bz.data.GetSessionByUniIds(ctx, req.UniIds)
	if err != nil {
		logging.Errorf(logHead+"res=%v, err=%v", mSession, err)
		return err
	}
	if len(mSession) == 0 {
		err = errors.New("can not find any session by uniIds")
		logging.Errorf(logHead+"err=%v", err)
		return err
	}

	// 重整数据：获取某个serverId下的tcpSessionId
	m := make(map[string][]string)
	for tcpSessionId, serverId := range mSession {
		m[serverId] = append(m[serverId], tcpSessionId)
	}

	// 把同一台机器的请求聚合在一起
	for serverId := range m {
		kafkaErr := bz.data.KafkaSendToUsers(serverId, m[serverId], req.SubId, []byte(req.Message))
		if kafkaErr != nil {
			logging.Errorf(logHead+"kafkaErr=%v", kafkaErr)
			continue
		}
	}

	return nil
}

// SendToRoom 发送消息（by kafka）
func (bz *Biz) SendToRoom(ctx context.Context, req *request.SendToRoomReq) error {
	logHead := "SendToRoom|"
	err := bz.data.KafkaSendToRoom(req)
	if err != nil {
		logging.Errorf(logHead+"err=%v", err)
		return err
	}
	return nil
}

// SendToAll 发送消息（by kafka）
func (bz *Biz) SendToAll(ctx context.Context, req *request.SendToAllReq) error {
	logHead := "SendToAll|"
	err := bz.data.KafkaSendToAll(req)
	if err != nil {
		logging.Errorf(logHead+"err=%v", err)
		return err
	}
	return nil
}
