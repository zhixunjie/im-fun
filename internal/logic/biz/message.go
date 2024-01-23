package biz

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

// SendCustomMessage 发送自定义消息
func (b *MessageUseCase) SendCustomMessage(ctx context.Context, sender, receiver *gen_id.ComponentId, d string) (rsp response.MessageSendRsp, err error) {
	return b.Send(ctx, &request.MessageSendReq{
		SeqId:        uint64(gen_id.SeqId()),
		SenderId:     sender.Id(),
		SenderType:   model.ContactIdType(sender.Type()),
		ReceiverId:   receiver.Id(),
		ReceiverType: model.ContactIdType(receiver.Type()),
		MsgBody: format.MsgBody{
			MsgType: format.MsgTypeCustom,
			MsgContent: &format.MsgContent{
				CustomContent: &format.CustomContent{
					Data: d,
				},
			},
		},
	})
}

// Send 发送消息
func (b *MessageUseCase) Send(ctx context.Context, req *request.MessageSendReq) (rsp response.MessageSendRsp, err error) {
	logHead := fmt.Sprintf("Send,SenderId=%v,ReceiverId=%v|", req.SenderId, req.ReceiverId)
	senderId := gen_id.NewComponentId(req.SenderId, uint32(req.SenderType))
	receiverId := gen_id.NewComponentId(req.ReceiverId, uint32(req.ReceiverType))

	// check limit
	err = b.checkMessageSend(ctx, req)
	if err != nil {
		return
	}

	// 1. build message
	var msg *model.Message
	msg, err = b.build(ctx, logHead, req, senderId, receiverId)
	if err != nil {
		return
	}

	// 2. build contact（sender's contact）
	var senderContact, peerContact *model.Contact
	if !lo.Contains[uint64](req.InvisibleList, req.SenderId) && b.canCreateContact(logHead, senderId) {
		senderContact, err = b.repoContact.Build(ctx, logHead, &model.BuildContactParams{
			OwnerId:   senderId,
			PeerId:    receiverId,
			LastMsgId: msg.MsgID,
			PeerAck:   model.PeerNotAck,
		})
		if err != nil {
			return
		}
	}

	// 3. build contact（receiver's contact）
	if !lo.Contains[uint64](req.InvisibleList, req.ReceiverId) && b.canCreateContact(logHead, receiverId) {
		peerContact, err = b.repoContact.Build(ctx, logHead, &model.BuildContactParams{
			OwnerId:   receiverId,
			PeerId:    senderId,
			LastMsgId: msg.MsgID,
			PeerAck:   model.PeerAcked,
		})
		if err != nil {
			return
		}
	}

	// 5. save to db
	err = b.repoMessage.Db.Transaction(func(tx *query.Query) error {
		var errTx error

		// 5.1 add message
		errTx = b.repoMessage.Create(logHead, tx, msg)
		if errTx != nil {
			return errTx
		}
		// 5.2 add contact(sender)
		if senderContact != nil {
			errTx = b.repoContact.Edit(logHead, tx, senderContact)
			if errTx != nil {
				return errTx
			}
		}

		// 5.3 add contact(peer)
		if peerContact != nil {
			errTx = b.repoContact.Edit(logHead, tx, peerContact)
			if errTx != nil {
				return errTx
			}
		}

		return nil
	})
	if err != nil {
		logging.Error(logHead + "mysql tx error,err=%v")
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
			SessionId:   msg.SessionID,
			UnreadCount: 0,
		},
	}

	return
}

// build 构建消息体
func (b *MessageUseCase) build(ctx context.Context, logHead string, req *request.MessageSendReq, senderId, receiverId *gen_id.ComponentId) (msg *model.Message, err error) {
	logHead += "Build|"
	mem := b.repoMessage.RedisClient

	// gen msg_id
	smallerId, largeId := gen_id.Sort(senderId, receiverId)
	msgId, err := b.repoMessage.GenMsgId(ctx, smallerId, largeId)
	if err != nil {
		logging.Errorf(logHead+"gen MsgId error=%v", err)
		return
	}

	// message: gen version_id
	versionId, err := gen_id.VersionId(ctx, &gen_id.GenVersionParams{
		Mem:            mem,
		GenVersionType: gen_id.GenVersionTypeMsg,
		SmallerId:      smallerId,
		LargerId:       largeId,
	})
	if err != nil {
		logging.Errorf(logHead+"gen VersionId error=%v", err)
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
	sessionId := b.repoMessage.GenSessionId(smallerId, largeId)
	msg = &model.Message{
		MsgID:         msgId,
		SeqID:         req.SeqId,
		MsgType:       uint32(req.MsgBody.MsgType),   // 消息类型
		Content:       string(bContent),              // 消息内容
		SessionID:     sessionId,                     // 会话ID
		SenderID:      req.SenderId,                  // 发送者ID
		VersionID:     versionId,                     // 版本ID
		SortKey:       versionId,                     // sort_key的值等同于version_id
		Status:        uint32(model.MsgStatusNormal), // 状态正常
		HasRead:       uint32(model.MsgRead),         // 已读（功能还没做好）
		InvisibleList: string(bInvisibleList),        // 不可见的列表
	}

	return
}

// Fetch 拉取消息
func (b *MessageUseCase) Fetch(ctx context.Context, req *request.MessageFetchReq) (rsp response.MessageFetchRsp, err error) {
	logHead := fmt.Sprintf("Fetch,req=%v", req)
	pivotVersionId := req.VersionId
	limit := 50

	// get: contact info
	ownerId := gen_id.NewComponentId(req.OwnerId, uint32(req.OwnerType))
	peerId := gen_id.NewComponentId(req.PeerId, uint32(req.PeerType))
	contactInfo, err := b.repoContact.Info(logHead, ownerId, peerId)
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
	smallId, largerId := gen_id.Sort(ownerId, peerId)
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
	case model.FetchTypeBackward: // 1. 拉取历史消息
		nextVersionId = minVersionId
	case model.FetchTypeForward: // 2. 拉取最新消息
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

// 双方通信时，判断是否需要创建对方的Contact
func (b *MessageUseCase) canCreateContact(logHead string, contactId *gen_id.ComponentId) bool {
	logHead += fmt.Sprintf("doNotNeedCreateContact,contactId=%v|", contactId)

	typeArr := []uint32{
		uint32(model.ContactIdTypeRobot),
		uint32(model.ContactIdTypeSystem),
		uint32(model.ContactIdTypeGroup),
	}

	// 如果用户是指定类型，那么不需要创建他的contact信息（比如：机器人）
	if lo.Contains(typeArr, contactId.Type()) {
		logging.Info(logHead + "do not need create contact")
		return false
	}

	return true
}

// 限制：发送者和接受者的类型
func (b *MessageUseCase) checkMessageSend(ctx context.Context, req *request.MessageSendReq) error {
	allowSenderType := []model.ContactIdType{
		model.ContactIdTypeUser,
		model.ContactIdTypeRobot,
		model.ContactIdTypeSystem,
	}

	allowReceiverType := []model.ContactIdType{
		model.ContactIdTypeUser,
		model.ContactIdTypeRobot,
		model.ContactIdTypeGroup,
	}

	if !lo.Contains(allowSenderType, req.SenderType) {
		return errors.New("checkMessageSend not allow")
	}
	if !lo.Contains(allowReceiverType, req.ReceiverType) {
		return errors.New("checkMessageSend not allow")
	}

	return nil
}
