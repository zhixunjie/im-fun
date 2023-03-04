package service

import (
	"context"
	"encoding/json"
	dao2 "github.com/zhixunjie/im-fun/internal/logic/dao"
	"github.com/zhixunjie/im-fun/internal/logic/model"
	"github.com/zhixunjie/im-fun/internal/logic/model/request"
	"github.com/zhixunjie/im-fun/internal/logic/model/response"
	"github.com/zhixunjie/im-fun/pkg/gen_id"
	"github.com/zhixunjie/im-fun/pkg/utils"
	"time"
)

// SendMessage 发送消息
func (svc *Service) SendMessage(ctx context.Context, req *request.SendMsgReq) (response.SendMsgResp, error) {
	var resp response.SendMsgResp
	currTimestamp := time.Now().Unix()

	// transform message
	msg, err := transformMessage(ctx, req, currTimestamp)
	if err != nil {
		return resp, err
	}

	// transform contact（sender）
	senderContact, err := transformSenderContact(ctx, req, currTimestamp, msg.MsgId)
	if err != nil {
		return resp, err
	}

	// transform contact（receiver）
	receiverContact, err := transformReceiverContact(ctx, req, currTimestamp, msg.MsgId)
	if err != nil {
		return resp, err
	}

	// DB操作
	_ = dao2.AddMsg(&msg)
	_ = dao2.AddOrUpdateContact(&senderContact)
	_ = dao2.AddOrUpdateContact(&receiverContact)

	return response.SendMsgResp{
		Data: response.SendMsgRespData{
			MsgId:        msg.MsgId,
			SeqId:        msg.SeqId,
			CreateTime:   msg.CreatedAt.Unix(),
			UpdateTime:   msg.UpdatedAt.Unix(),
			MsgVersionId: msg.VersionId,
			MsgSortKey:   msg.SortKey,
			UnreadCount:  0,
		},
	}, nil
}

func transformMessage(ctx context.Context, req *request.SendMsgReq, currTimestamp int64) (model.Message, error) {
	var defaultRet model.Message

	// get msg_id
	smallerId, largeId := utils.GetSortNum(req.SendId, req.ReceiveId)
	msgId, err := gen_id.GenerateMsgId(ctx, largeId, currTimestamp)
	if err != nil {
		return defaultRet, err
	}

	// get version_id
	versionId, err := gen_id.GetMsgVersionId(ctx, currTimestamp, smallerId, largeId)
	if err != nil {
		return defaultRet, err
	}

	// exchange
	buf, err := json.Marshal(req.InvisibleList)
	if err != nil {
		return defaultRet, err
	}

	// build message
	msg := model.Message{
		MsgId:         msgId,
		MsgType:       req.MsgType,
		SessionId:     gen_id.GetSessionId(req.SendId, req.ReceiveId),
		SendId:        req.SendId,
		VersionId:     versionId,
		SortKey:       versionId, // sort_key的值等同于version_id
		Status:        model.MsgStatusNormal,
		Content:       req.Content,
		HasRead:       model.MsgRead,
		InvisibleList: string(buf),
		SeqId:         req.SeqId,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	return msg, err
}

// FetchMessage 拉取消息
func (svc *Service) FetchMessage(ctx context.Context, req *request.FetchMsgReq) (response.SendMsgResp, error) {
	var resp response.SendMsgResp

	return resp, nil
}
