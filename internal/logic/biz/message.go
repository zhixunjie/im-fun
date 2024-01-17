package biz

import (
	"context"
	"encoding/json"
	"github.com/zhixunjie/im-fun/internal/logic/data"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/model"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/query"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/request"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/response"
	"github.com/zhixunjie/im-fun/pkg/gen_id"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"github.com/zhixunjie/im-fun/pkg/utils"
	"go.uber.org/zap"
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
	logHead := "SendMessage|"
	currTimestamp := time.Now().Unix()
	contactUseCase := bz.contactUseCase

	// transform message
	msg, err := bz.transformMessage(ctx, req, currTimestamp)
	if err != nil {
		return
	}

	// transform contact（send）
	senderContact, err := contactUseCase.TransformSender(ctx, req, currTimestamp, msg.MsgID)
	if err != nil {
		return
	}

	// transform contact（receive）
	peerContact, err := contactUseCase.TransformPeer(ctx, req, currTimestamp, msg.MsgID)
	if err != nil {
		return
	}

	// DB操作
	err = bz.repo.Db.Transaction(func(tx *query.Query) error {
		var errTx error

		// 1. add message
		errTx = bz.repo.AddMsg(tx, msg)
		if errTx != nil {
			logging.Error(logHead+"AddMsg error=%v", err)
			return errTx
		}
		// 2. add contact(sender)
		errTx = contactUseCase.repo.AddOrUpdateContact(tx, senderContact)
		if errTx != nil {
			logging.Error(logHead+"AddOrUpdateContact error=%v", err)
			return errTx
		}
		// 3. add contact(peer)
		errTx = contactUseCase.repo.AddOrUpdateContact(tx, peerContact)
		if errTx != nil {
			logging.Error(logHead+"AddOrUpdateContact error=%v", err)
			return errTx
		}

		return nil
	})
	if err != nil {
		logging.Error(logHead+"mysql tx error", zap.Error(err))
		return
	}
	logging.Info(logHead + "mysql tx success")

	// build response
	resp = response.SendMsgResp{
		Data: response.SendMsgRespData{
			MsgId:        msg.MsgID,
			SeqId:        msg.SeqID,
			CreateTime:   msg.CreatedAt.Unix(),
			UpdateTime:   msg.UpdatedAt.Unix(),
			MsgVersionId: msg.VersionID,
			MsgSortKey:   msg.SortKey,
			UnreadCount:  0,
		},
	}

	return
}

func (bz *MessageUseCase) transformMessage(ctx context.Context, req *request.SendMsgReq, currTimestamp int64) (msg *model.Message, err error) {
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
	msg = &model.Message{
		MsgID:         msgId,
		SeqID:         req.SeqId,
		MsgType:       uint32(req.MsgBody.MsgType),
		Content:       string(bufContent),
		SessionID:     gen_id.SessionId(req.SendId, req.PeerId),
		SenderID:      req.SendId,
		VersionID:     versionId,
		SortKey:       versionId, // sort_key的值等同于version_id
		Status:        model.MsgStatusNormal,
		HasRead:       model.MsgRead,
		InvisibleList: string(buf),
	}

	return
}

// FetchMessage 拉取消息
func (bz *MessageUseCase) FetchMessage(ctx context.Context, req *request.FetchMsgReq) (resp response.FetchMsgResp, err error) {
	// https://redis.io/commands/zrevrangebyscore/
	// https://redis.io/commands/zcount/

	return
}
