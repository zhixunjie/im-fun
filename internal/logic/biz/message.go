package biz

import (
	"context"
	"encoding/json"
	"github.com/zhixunjie/im-fun/internal/logic/data"
	"github.com/zhixunjie/im-fun/internal/logic/model"
	"github.com/zhixunjie/im-fun/internal/logic/model/request"
	"github.com/zhixunjie/im-fun/internal/logic/model/response"
	"github.com/zhixunjie/im-fun/pkg/gen_id"
	"github.com/zhixunjie/im-fun/pkg/utils"
	"time"
)

type MessageUseCase struct {
	repo           *data.MessageRepo
	contactUseCase *ContactUseCase
}

func NewMessageUseCase(repo *data.MessageRepo, contactUseCase *ContactUseCase) *MessageUseCase {
	return &MessageUseCase{
		repo:           repo,
		contactUseCase: contactUseCase,
	}
}

// SendMessage 发送消息
func (bz *MessageUseCase) SendMessage(ctx context.Context, req *request.SendMsgReq) (resp response.SendMsgResp, err error) {
	currTimestamp := time.Now().Unix()
	contactUseCase := bz.contactUseCase

	// transform message
	msg, err := bz.transformMessage(ctx, req, currTimestamp)
	if err != nil {
		return
	}

	// transform contact（send）
	senderContact, err := contactUseCase.TransformSender(ctx, req, currTimestamp, msg.MsgId)
	if err != nil {
		return
	}

	// transform contact（receive）
	peerContact, err := contactUseCase.TransformPeer(ctx, req, currTimestamp, msg.MsgId)
	if err != nil {
		return
	}

	// DB操作
	_ = bz.repo.AddMsg(&msg)
	_ = contactUseCase.repo.AddOrUpdateContact(&senderContact)
	_ = contactUseCase.repo.AddOrUpdateContact(&peerContact)

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

func (bz *MessageUseCase) transformMessage(ctx context.Context, req *request.SendMsgReq, currTimestamp int64) (msg model.Message, err error) {
	mem := bz.repo.RedisClient

	// gen msg_id
	smallerId, largeId := utils.GetSortNum(req.SendId, req.PeerId)
	msgId, err := gen_id.MsgId(ctx, mem, largeId, currTimestamp)
	if err != nil {
		return
	}

	// gen version_id
	versionId, err := gen_id.MsgVersionId(ctx, mem, currTimestamp, smallerId, largeId)
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
		SeqId:         req.SeqId,
		MsgType:       int32(req.MsgBody.MsgType),
		Content:       string(bufContent),
		SessionId:     gen_id.SessionId(req.SendId, req.PeerId),
		SendId:        req.SendId,
		VersionId:     versionId,
		SortKey:       versionId, // sort_key的值等同于version_id
		Status:        model.MsgStatusNormal,
		HasRead:       model.MsgRead,
		InvisibleList: string(buf),
	}

	return
}

// FetchMessage 拉取消息
func (bz *MessageUseCase) FetchMessage(ctx context.Context, req *request.FetchMsgReq) (resp response.SendMsgResp, err error) {
	// https://redis.io/commands/zrevrangebyscore/
	// https://redis.io/commands/zcount/

	return
}
