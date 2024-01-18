package biz

import (
	"context"
	"encoding/json"
	"github.com/samber/lo"
	"github.com/zhixunjie/im-fun/internal/logic/data"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/model"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/query"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/request"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/response"
	"github.com/zhixunjie/im-fun/pkg/gen_id"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"github.com/zhixunjie/im-fun/pkg/utils"
	"go.uber.org/zap"
)

type MessageUseCase struct {
	repo        *data.MessageRepo
	repoContact *data.ContactRepo
}

func NewMessageUseCase(repo *data.MessageRepo, repoContact *data.ContactRepo) *MessageUseCase {
	return &MessageUseCase{
		repo:        repo,
		repoContact: repoContact,
	}
}

// SendMessage 发送消息
func (b *MessageUseCase) SendMessage(ctx context.Context, req *request.SendMsgReq) (resp response.SendMsgResp, err error) {
	logHead := "SendMessage|"

	// build message
	msg, err := b.BuildMessage(ctx, req)
	if err != nil {
		return
	}

	// build contact（sender）
	var senderContact, peerContact *model.Contact
	if !lo.Contains[uint64](req.InvisibleList, req.SendId) {
		senderContact, err = b.repoContact.BuildContact(ctx, &model.BuildContactParams{
			MsgId:    msg.MsgID,
			OwnerId:  req.SendId,
			PeerId:   req.PeerId,
			PeerType: req.PeerType,
			PeerAck:  model.PeerNotAck,
		})
		if err != nil {
			return
		}
	}

	// build contact（receive）
	if !lo.Contains[uint64](req.InvisibleList, req.PeerId) {
		peerContact, err = b.repoContact.BuildContact(ctx, &model.BuildContactParams{
			MsgId:    msg.MsgID,
			OwnerId:  req.PeerId,
			PeerId:   req.SendId,
			PeerType: req.SenderType,
			PeerAck:  model.PeerAck,
		})
		if err != nil {
			return
		}
	}

	// DB操作
	err = b.repo.Db.Transaction(func(tx *query.Query) error {
		var errTx error

		// 1. add message
		errTx = b.repo.AddMsg(tx, msg)
		if errTx != nil {
			logging.Error(logHead+"AddMsg error=%v", err)
			return errTx
		}
		// 2. add contact(sender)
		if senderContact != nil {
			errTx = b.repoContact.EditContact(tx, senderContact)
			if errTx != nil {
				logging.Error(logHead+"EditContact error=%v", err)
				return errTx
			}
		}

		// 3. add contact(peer)
		if peerContact != nil {
			errTx = b.repoContact.EditContact(tx, peerContact)
			if errTx != nil {
				logging.Error(logHead+"EditContact error=%v", err)
				return errTx
			}
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

func (b *MessageUseCase) BuildMessage(ctx context.Context, req *request.SendMsgReq) (msg *model.Message, err error) {
	mem := b.repo.RedisClient

	// gen msg_id
	smallerId, largeId := utils.SortNum(req.SendId, req.PeerId)
	msgId, err := gen_id.MsgId(ctx, mem, largeId)
	if err != nil {
		return
	}

	// gen version_id
	versionId, err := gen_id.MsgVersionId(ctx, mem, smallerId, largeId)
	if err != nil {
		return
	}

	// exchange：InvisibleList
	var buf []byte
	if len(req.InvisibleList) > 0 {
		buf, err = json.Marshal(req.InvisibleList)
		if err != nil {
			return
		}
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
		SessionID:     gen_id.SessionId(req.SendId, req.PeerId), // 会话ID
		SenderID:      req.SendId,                               // 发送者ID
		VersionID:     versionId,                                // 版本ID
		SortKey:       versionId,                                // sort_key的值等同于version_id
		Status:        model.MsgStatusNormal,                    // 状态正常
		HasRead:       model.MsgRead,                            // 已读（功能还没做好）
		InvisibleList: string(buf),
	}

	return
}

// FetchMessage 拉取消息
func (b *MessageUseCase) FetchMessage(ctx context.Context, req *request.FetchMsgReq) (resp response.FetchMsgResp, err error) {
	// https://redis.io/commands/zrevrangebyscore/
	// https://redis.io/commands/zcount/

	return
}
