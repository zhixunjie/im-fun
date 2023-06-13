package service

import (
	"context"
	"github.com/zhixunjie/im-fun/internal/logic/model/request"
	"github.com/zhixunjie/im-fun/pkg/logging"
)

// SendToUserKeys 发送消息（by kafka）
func (svc *Service) SendToUserKeys(ctx context.Context, req *request.SendToUserKeysReq) error {
	logHead := "SendToUserKeys|"
	res, err := svc.dao.SessionGetByUserKeys(ctx, req.UserKeys)
	if err != nil {
		logging.Errorf(logHead+"res=%v, err=%v", res, err)
		return err
	}

	// 重整数据：获取某个serverId下的userKey
	serverIdMap := make(map[string][]string)
	for i, userKey := range req.UserKeys {
		serverId := res[i]
		serverIdMap[serverId] = append(serverIdMap[serverId], userKey)
	}

	// 同一台机器的userKey一次性发送
	for serverId := range serverIdMap {
		err = svc.dao.KafkaSendToUserKeys(serverId, serverIdMap[serverId], req.SubId, []byte(req.Message))
		if err != nil {
			logging.Errorf(logHead+"err=%v", err)
		}
	}

	return nil
}

// SendToUserIds 发送消息（by kafka）
func (svc *Service) SendToUserIds(ctx context.Context, req *request.SendToUserIdsReq) error {
	logHead := "SendToUserIds|"
	res, err := svc.dao.SessionGetByUserIds(ctx, req.UserIds)
	if err != nil {
		logging.Errorf(logHead+"res=%v, err=%v", res, err)
		return err
	}

	// 重整数据：获取某个serverId下的userKey
	serverIdMap := make(map[string][]string)
	for userKey, serverId := range res {
		serverIdMap[serverId] = append(serverIdMap[serverId], userKey)
	}

	// 同一台机器的userKey一次性发送
	for serverId := range serverIdMap {
		err = svc.dao.KafkaSendToUserKeys(serverId, serverIdMap[serverId], req.SubId, req.Message)
		if err != nil {
			logging.Errorf(logHead+"err=%v", err)
		}
	}

	return nil
}

// SendToRoom 发送消息（by kafka）
func (svc *Service) SendToRoom(ctx context.Context, req *request.SendToRoomReq) error {
	logHead := "SendToRoom|"
	err := svc.dao.KafkaSendToRoom(req)
	if err != nil {
		logging.Errorf(logHead+"err=%v", err)
		return err
	}
	return nil
}

// SendToAll 发送消息（by kafka）
func (svc *Service) SendToAll(ctx context.Context, req *request.SendToAllReq) error {
	logHead := "SendToAll|"
	err := svc.dao.KafkaSendToAll(req)
	if err != nil {
		logging.Errorf(logHead+"err=%v", err)
		return err
	}
	return nil
}
