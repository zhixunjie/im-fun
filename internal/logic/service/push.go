package service

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/zhixunjie/im-fun/internal/logic/model/request"
)

// PushUserKeys 发送消息（by kafka）
func (svc *Service) PushUserKeys(ctx context.Context, req *request.PushUserKeysReq) error {
	logHead := "PushUserKeys|"
	res, err := svc.dao.SessionGetByUserKeys(ctx, req.UserKeys)
	if err != nil {
		logrus.Errorf(logHead+"res=%v, err=%v", res, err)
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
		err = svc.dao.KafkaPushKeys(serverId, serverIdMap[serverId], req.SubId, []byte(req.Message))
		if err != nil {
			logrus.Errorf(logHead+"err=%v", err)
		}
	}

	return nil
}

// PushUserIds 发送消息（by kafka）
func (svc *Service) PushUserIds(ctx context.Context, req *request.PushUserIdsReq) error {
	logHead := "PushUserIds|"
	res, err := svc.dao.SessionGetByUserIds(ctx, req.UserIds)
	if err != nil {
		logrus.Errorf(logHead+"res=%v, err=%v", res, err)
		return err
	}

	// 重整数据：获取某个serverId下的userKey
	serverIdMap := make(map[string][]string)
	for userKey, serverId := range res {
		serverIdMap[serverId] = append(serverIdMap[serverId], userKey)
	}

	// 同一台机器的userKey一次性发送
	for serverId := range serverIdMap {
		err = svc.dao.KafkaPushKeys(serverId, serverIdMap[serverId], req.SubId, req.Message)
		if err != nil {
			logrus.Errorf(logHead+"err=%v", err)
		}
	}

	return nil
}

// PushUserRoom 发送消息（by kafka）
func (svc *Service) PushUserRoom(ctx context.Context, req *request.PushUserRoomReq) error {
	logHead := "BroadcastRoom|"
	err := svc.dao.KafkaPushRoom(req)
	if err != nil {
		logrus.Errorf(logHead+"err=%v", err)
		return err
	}
	return nil
}

// PushUserAll 发送消息
func (svc *Service) PushUserAll(ctx context.Context, req *request.PushUserAllReq) error {
	logHead := "PushUserAll|"
	err := svc.dao.KafkaPushAll(req)
	if err != nil {
		logrus.Errorf(logHead+"err=%v", err)
		return err
	}
	return nil
}
