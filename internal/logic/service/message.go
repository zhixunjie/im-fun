package service

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/zhixunjie/im-fun/internal/logic/model"
	"github.com/zhixunjie/im-fun/internal/logic/model/request"
	"github.com/zhixunjie/im-fun/internal/logic/model/response"
	"github.com/zhixunjie/im-fun/pkg/gen_id"
	"github.com/zhixunjie/im-fun/pkg/utils"
	"time"
)

// SendMessage 发送消息
func (svc *Service) SendMessage(ctx context.Context, req *request.SendMsgReq) (resp response.SendMsgResp, err error) {
	currTimestamp := time.Now().Unix()

	// transform message
	msg, err := transformMessage(ctx, svc.dao.RedisClient, req, currTimestamp)
	if err != nil {
		return
	}

	// transform contact（send）
	senderContact, err := svc.transformSenderContact(ctx, req, currTimestamp, msg.MsgId)
	if err != nil {
		return
	}

	// transform contact（receive）
	peerContact, err := svc.transformPeerContact(ctx, req, currTimestamp, msg.MsgId)
	if err != nil {
		return
	}

	// DB操作
	_ = svc.dao.AddMsg(&msg)
	_ = svc.dao.AddOrUpdateContact(&senderContact)
	_ = svc.dao.AddOrUpdateContact(&peerContact)

	// build response
	resp = response.SendMsgResp{
		Data: response.SendMsgRespData{
			MsgId:        msg.MsgId,
			SeqId:        msg.SeqId,
			CreateTime:   msg.CreatedAt.Unix(),
			UpdateTime:   msg.UpdatedAt.Unix(),
			MsgVersionId: msg.VersionId,
			MsgSortKey:   msg.SortKey,
			UnreadCount:  0,
		},
	}

	return
}

func transformMessage(ctx context.Context, mem *redis.Client, req *request.SendMsgReq, currTimestamp int64) (msg model.Message, err error) {
	// get msg_id
	smallerId, largeId := utils.GetSortNum(req.SendId, req.PeerId)
	msgId, err := gen_id.GenerateMsgId(ctx, mem, largeId, currTimestamp)
	if err != nil {
		return
	}

	// get version_id
	versionId, err := gen_id.GetMsgVersionId(ctx, mem, currTimestamp, smallerId, largeId)
	if err != nil {
		return
	}

	// exchange：InvisibleList
	buf, err := json.Marshal(req.InvisibleList)
	if err != nil {
		return
	}

	// exchange：MsgContent
	bufContent, err := json.Marshal(req.MsgBody.MsgContent)
	if err != nil {
		return
	}

	// build message
	msg = model.Message{
		MsgId:         msgId,
		MsgType:       int32(req.MsgBody.MsgType),
		SessionId:     gen_id.GetSessionId(req.SendId, req.PeerId),
		SendId:        req.SendId,
		VersionId:     versionId,
		SortKey:       versionId, // sort_key的值等同于version_id
		Status:        model.MsgStatusNormal,
		Content:       string(bufContent),
		HasRead:       model.MsgRead,
		InvisibleList: string(buf),
		SeqId:         req.SeqId,
	}

	return
}

// FetchMessage 拉取消息
func (svc *Service) FetchMessage(ctx context.Context, req *request.FetchMsgReq) (resp response.SendMsgResp, err error) {

	return
}
