package biz

import (
	"context"
	"encoding/json"
	"github.com/samber/lo"
	"github.com/zhixunjie/im-fun/internal/logic/data"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/format"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/model"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/query"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/request"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/response"
	"github.com/zhixunjie/im-fun/pkg/gen_id"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"github.com/zhixunjie/im-fun/pkg/utils"
	"go.uber.org/zap"
	"math"
	"sort"
)

type MessageUseCase struct {
	repoMessage *data.MessageRepo
	repoContact *data.ContactRepo
}

func NewMessageUseCase(repoMessage *data.MessageRepo, repoContact *data.ContactRepo) *MessageUseCase {
	return &MessageUseCase{
		repoMessage: repoMessage,
		repoContact: repoContact,
	}
}

// Send 发送消息
func (b *MessageUseCase) Send(ctx context.Context, req *request.MessageSendReq) (rsp response.MessageSendRsp, err error) {
	logHead := "SendMessage|"

	// 1. build message
	msg, err := b.Build(ctx, req)
	if err != nil {
		return
	}

	// 2. build contact（sender）
	var senderContact, peerContact *model.Contact
	if !lo.Contains[uint64](req.InvisibleList, req.SenderId) {
		senderContact, err = b.repoContact.Build(ctx, &model.BuildContactParams{
			LastMsgId:    msg.MsgID,
			OwnerId:      req.SenderId,
			PeerId:       req.ReceiverId,
			InitPeerType: req.ReceiverContactPeerType,
			InitPeerAck:  uint32(model.PeerNotAck),
		})
		if err != nil {
			return
		}
	}

	// 3. build contact（receive）
	if !lo.Contains[uint64](req.InvisibleList, req.ReceiverId) {
		peerContact, err = b.repoContact.Build(ctx, &model.BuildContactParams{
			LastMsgId:    msg.MsgID,
			OwnerId:      req.ReceiverId,
			PeerId:       req.SenderId,
			InitPeerType: req.SenderContactPeerType,
			InitPeerAck:  uint32(model.PeerAcked),
		})
		if err != nil {
			return
		}
	}

	// 5. save to db
	err = b.repoMessage.Db.Transaction(func(tx *query.Query) error {
		var errTx error

		// 5.1 add message
		errTx = b.repoMessage.Create(tx, msg)
		if errTx != nil {
			logging.Error(logHead+"AddMsg error=%v", err)
			return errTx
		}
		// 5.2 add contact(sender)
		if senderContact != nil {
			errTx = b.repoContact.Edit(tx, senderContact)
			if errTx != nil {
				logging.Error(logHead+"EditContact error=%v", err)
				return errTx
			}
		}

		// 5.3 add contact(peer)
		if peerContact != nil {
			errTx = b.repoContact.Edit(tx, peerContact)
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

	// 6.build response
	rsp = response.MessageSendRsp{
		Data: response.SendMsgRespData{
			MsgId:       msg.MsgID,
			SeqId:       msg.SeqID,
			VersionId:   msg.VersionID,
			SortKey:     msg.SortKey,
			UnreadCount: 0,
		},
	}

	return
}

// Build 构建消息体
func (b *MessageUseCase) Build(ctx context.Context, req *request.MessageSendReq) (msg *model.Message, err error) {
	logHead := "Build|"
	mem := b.repoMessage.RedisClient

	// gen msg_id
	smallerId, largeId := utils.SortNum(req.SenderId, req.ReceiverId)
	msgId, err := gen_id.MsgId(ctx, mem, largeId)
	if err != nil {
		return
	}

	// gen version_id
	versionId, err := gen_id.VersionId(ctx, &gen_id.GenVersionParams{
		Mem:            mem,
		GenVersionType: gen_id.GenVersionTypeMsg,
		SmallerId:      smallerId,
		LargerId:       largeId,
	})
	if err != nil {
		logging.Errorf(logHead+"MsgVersionId error=%v", err)
		return
	}

	// exchange：InvisibleList
	var bInvisibleList []byte
	if len(req.InvisibleList) > 0 {
		bInvisibleList, err = json.Marshal(req.InvisibleList)
		if err != nil {
			logging.Errorf(logHead+"Marshal error=%v", err)
			return
		}
	}

	// exchange：MsgContent
	bContent, err := json.Marshal(req.MsgBody.MsgContent)
	if err != nil {
		logging.Errorf(logHead+"Marshal error=%v", err)
		return
	}

	// build message
	msg = &model.Message{
		MsgID:         msgId,
		SeqID:         req.SeqId,
		MsgType:       uint32(req.MsgBody.MsgType),
		Content:       string(bContent),
		SessionID:     gen_id.SessionId(req.SenderId, req.ReceiverId), // 会话ID
		SenderID:      req.SenderId,                                   // 发送者ID
		VersionID:     versionId,                                      // 版本ID
		SortKey:       versionId,                                      // sort_key的值等同于version_id
		Status:        uint32(model.MsgStatusNormal),                  // 状态正常
		HasRead:       uint32(model.MsgRead),                          // 已读（功能还没做好）
		InvisibleList: string(bInvisibleList),
	}

	return
}

// Fetch 拉取消息
func (b *MessageUseCase) Fetch(ctx context.Context, req *request.MessageFetchReq) (rsp response.MessageFetchRsp, err error) {
	//logHead := "Fetch|"
	pivotVersionId := req.VersionId
	limit := 50

	// get: contact info
	contactInfo, err := b.repoContact.Info(req.OwnerId, req.PeerId)
	if err != nil {
		return
	}

	// get: last msg info
	var lastDelMsgVersionId uint64
	var messageInfo *model.Message
	if contactInfo.LastDelMsgID > 0 {
		messageInfo, err = b.repoMessage.Info(contactInfo.LastDelMsgID)
		if err != nil {
			return
		}
		lastDelMsgVersionId = messageInfo.VersionID
	}

	// set pivotVersionId
	switch req.FetchType {
	case model.FetchTypeBackward: // 拉取历史消息
	case model.FetchTypeForward: // 拉取最新消息
		// 避免：拉取最新消息时拉到已删除消息
		if pivotVersionId < lastDelMsgVersionId {
			pivotVersionId = lastDelMsgVersionId
		}
	default:
		return
	}

	// get: message list
	smallId, largerId := utils.SortNum(req.OwnerId, req.PeerId)
	list, err := b.repoMessage.RangeList(&model.FetchMsgRangeParams{
		FetchType:           req.FetchType,
		SmallerId:           smallId,
		LargerId:            largerId,
		PivotVersionId:      pivotVersionId,
		LastDelMsgVersionId: lastDelMsgVersionId,
		Limit:               limit,
	})
	if err != nil {
		return
	}

	// rebuild list
	var retList []*response.MsgEntity
	minVersionId := uint64(math.MaxUint64)
	maxVersionId := uint64(0)
	var tmpList []uint64
	for _, item := range list {
		minVersionId = utils.Min(minVersionId, item.VersionID)
		maxVersionId = utils.Max(maxVersionId, item.VersionID)

		// 过滤：invisible list
		if len(item.InvisibleList) > 0 {
			err = json.Unmarshal([]byte(item.InvisibleList), &tmpList)
			if err != nil {
				continue
			}

			if lo.Contains[uint64](tmpList, req.OwnerId) {
				continue
			}
		}

		// build message list
		msgContent := new(format.MsgContent)
		_ = json.Unmarshal([]byte(item.Content), msgContent)
		retList = append(retList, &response.MsgEntity{
			MsgID: item.MsgID,
			SeqID: item.SeqID,
			MsgBody: format.MsgBody{
				MsgType:    format.MsgType(item.MsgType),
				MsgContent: msgContent,
			},
			SessionID: item.SessionID,
			SenderID:  item.SenderID,
			VersionID: item.VersionID,
			SortKey:   item.SortKey,
			Status:    model.MsgStatus(item.Status),
			HasRead:   model.MsgReadStatus(item.HasRead),
		})
	}

	// get: nextVersionId
	var nextVersionId uint64
	switch req.FetchType {
	case model.FetchTypeBackward: // 拉取历史消息
		nextVersionId = minVersionId
	case model.FetchTypeForward: // 拉取最新消息
		nextVersionId = maxVersionId
	default:
		return
	}

	// sort: 按照sort_key排序（从小到大）
	sort.Slice(retList, func(i, j int) bool {
		return retList[i].SortKey < retList[j].SortKey
	})

	rsp.Data = response.FetchMsgData{
		MsgList:       retList,
		NextVersionId: nextVersionId,
		HasMore:       len(list) == limit,
	}

	return
}
