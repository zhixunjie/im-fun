package biz

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/samber/lo"
	"github.com/zhixunjie/im-fun/internal/logic/api"
	"github.com/zhixunjie/im-fun/internal/logic/data"
	"github.com/zhixunjie/im-fun/internal/logic/data/cache"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/format"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/model"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/request"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/response"
	"github.com/zhixunjie/im-fun/pkg/gen_id"
	"github.com/zhixunjie/im-fun/pkg/goredis/distrib_lock"
	k "github.com/zhixunjie/im-fun/pkg/goredis/key"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"github.com/zhixunjie/im-fun/pkg/routine"
	"github.com/zhixunjie/im-fun/pkg/utils"
	"gorm.io/gorm"
	"math"
	"sort"
	"time"
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

// SendSimpleCustomMessage 简化接口：发送自定义消息
func (b *MessageUseCase) SendSimpleCustomMessage(ctx context.Context, sender, receiver *gen_id.ComponentId, d string) (rsp response.MessageSendRsp, err error) {
	return b.Send(ctx, &request.MessageSendReq{
		SeqId:    uint64(gen_id.SeqId()),
		Sender:   sender,
		Receiver: receiver,
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
	sender := req.Sender
	receiver := req.Receiver
	logHead := fmt.Sprintf("Send|sender=%v,receiver=%v|", sender, receiver)

	// check limit
	err = b.checkMessageSend(ctx, req)
	if err != nil {
		return
	}

	// 1. create sender's contact if not exists
	var senderContact, peerContact *model.Contact
	if b.needCreateContact(logHead, sender) {
		if !lo.Contains(req.InvisibleList, req.Sender.Id()) {
			senderContact, err = b.repoContact.CreateNotExists(logHead, &model.BuildContactParams{
				Owner: sender,
				Peer:  receiver,
			})
			if err != nil {
				return
			}
		}
	}
	// 2. create receiver's contact if not exists
	if b.needCreateContact(logHead, receiver) {
		if !lo.Contains(req.InvisibleList, req.Receiver.Id()) {
			peerContact, err = b.repoContact.CreateNotExists(logHead, &model.BuildContactParams{
				Owner: receiver,
				Peer:  sender,
			})
			if err != nil {
				return
			}
		}
	}
	// 3. build && create message（无扩散）
	msg, err := b.build(ctx, logHead, req, sender, receiver)
	if err != nil {
		return
	}
	currMsgId := msg.MsgID

	routine.Go(ctx, func() {
		// 增加未读数（先save db，再incr cache，保证尽快执行）
		if !lo.Contains(req.InvisibleList, req.Receiver.Id()) {
			_ = b.repoMessage.IncrUnreadAfterSend(ctx, logHead, receiver, sender, 1)
		}
		// update contact's info（写扩散）
		if senderContact != nil {
			err = b.repoContact.UpdateLastMsgId(ctx, logHead, senderContact.ID, sender, currMsgId, model.PeerNotAck)
			if err != nil {
				return
			}
		}
		// update contact's info（写扩散）
		if peerContact != nil {
			err = b.repoContact.UpdateLastMsgId(ctx, logHead, peerContact.ID, receiver, currMsgId, model.PeerAcked)
			if err != nil {
				return
			}
			if err != nil {
				return
			}
		}
	})

	return response.MessageSendRsp{
		Data: response.SendMsgRespData{
			MsgId:       msg.MsgID,
			SeqId:       msg.SeqID,
			VersionId:   msg.VersionID,
			SortKey:     msg.SortKey,
			SessionId:   msg.SessionID,
			UnreadCount: 0,
		},
	}, nil
}

// Fetch 拉取消息列表
func (b *MessageUseCase) Fetch(ctx context.Context, req *request.MessageFetchReq) (rsp response.MessageFetchRsp, err error) {
	logHead := fmt.Sprintf("Fetch|req=%v", req)
	pivotVersionId := req.VersionId
	limit := 50

	// get: contact info
	owner := req.Owner
	peer := req.Peer
	contactInfo, err := b.repoContact.Info(logHead, owner, peer)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = api.ErrContactNotExists
		}
		logging.Error(logHead+"repoContact Info,err=%v", err)
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

	// get: message list
	sessionId := gen_id.SessionId(owner, peer)
	list, err := b.repoMessage.RangeList(&model.FetchMsgRangeParams{
		FetchType:           req.FetchType,
		SessionId:           sessionId,
		LastDelMsgVersionId: lastDelMsgVersionId,
		PivotVersionId:      pivotVersionId,
		Limit:               limit,
		Owner:               owner,
		Peer:                peer,
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

			if lo.Contains(tmpList, req.Owner.Id()) {
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

	// 异步协程：减少未读数（先read db，再decr cache）
	routine.Go(ctx, func() {
		_ = b.repoMessage.DecrUnreadAfterFetch(ctx, logHead, owner, peer, int64(len(retList)))
	})

	// sort: 返回之前进行重新排序
	sort.Sort(response.MessageSortByVersion(retList))

	rsp.Data = response.FetchMsgData{
		MsgList:       retList,
		NextVersionId: nextVersionId,
		HasMore:       len(list) == limit,
	}

	return
}

// Recall 某条消息撤回（两边的聊天记录都需要撤回）
// 核心：更新status和version_id（需要通知对方，所以需要更新version_id）
func (b *MessageUseCase) Recall(ctx context.Context, req *request.MessageRecallReq) (rsp response.MessageWithdrawRsp, err error) {
	logHead := "Recall|"
	senderId := req.Sender

	// invoke common method
	err = b.updateMsgIdStatus(ctx, logHead, req.MsgId, model.MsgStatusWithdraw, senderId)
	if err != nil {
		return
	}
	return
}

// DelBothSide 某条消息删除（两边的聊天记录都需要删除）
// 核心：更新status和version_id（需要通知对方，所以需要更新version_id）
func (b *MessageUseCase) DelBothSide(ctx context.Context, req *request.MessageDelBothSideReq) (rsp response.DelBothSideRsp, err error) {
	logHead := "DelBothSide|"
	senderId := req.Sender

	// invoke common method
	err = b.updateMsgIdStatus(ctx, logHead, req.MsgId, model.MsgStatusDeleted, senderId)
	if err != nil {
		return
	}
	return
}

// DelOneSide 某条消息删除（只有一边的聊天记录是不可见，另外一边可见）
// 核心：更新 invisible_list
func (b *MessageUseCase) DelOneSide(ctx context.Context, req *request.MessageDelOneSideReq) (rsp response.DelOneSideRsp, err error) {
	return
}

// ClearHistory 清空聊天记录（批量清空）
// 核心：更新Contact的 last_del_msg_id 为 last_msg_id
func (b *MessageUseCase) ClearHistory(ctx context.Context, req *request.ClearHistoryReq) (rsp response.ClearHistoryRsp, err error) {
	logHead := fmt.Sprintf("ClearHistory|")
	lastDelMsgId := req.MsgId
	owner := req.Owner
	peer := req.Peer
	mem := b.repoMessage.RedisClient

	// 聊天记录的清空：
	// 1. 指定msgId进行清空
	// 2. 获取contact记录的最后一条消息进行清空
	if lastDelMsgId == 0 {
		// get: contact info
		var contactInfo *model.Contact
		contactInfo, err = b.repoContact.Info(logHead, owner, peer)
		if err != nil {
			logging.Error(logHead+"repoContact Info,err=%v", err)
			return
		}
		lastDelMsgId = contactInfo.LastMsgID
	}

	// contact: gen version_id
	versionId, err := gen_id.ContactVersionId(ctx, &gen_id.ContactVerParams{
		Mem:   mem,
		Owner: owner,
	})
	if err != nil {
		logging.Errorf(logHead+"gen VersionId error=%v", err)
		return
	}

	// update to db
	affectedRow, err := b.repoContact.UpdateLastDelMsg(logHead, lastDelMsgId, versionId, owner, peer)
	if err != nil {
		logging.Errorf(logHead+"UpdateLastDelMsg error=%v", err)
		return
	}
	if affectedRow == 0 {
		err = errors.New("affectedRow not allow")
		logging.Errorf(logHead+"UpdateLastDelMsg error=%v", err)
		return
	}

	return
}

// build 构建消息体
func (b *MessageUseCase) build(ctx context.Context, logHead string, req *request.MessageSendReq, sender, receiver *gen_id.ComponentId) (msg *model.Message, err error) {
	logHead += "build|"
	mem := b.repoMessage.RedisClient

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

	// message: gen session id
	sessionId := gen_id.SessionId(sender, receiver)

	// note: 同一个消息timeline的版本变动，需要加锁
	// 保证数据库记录中的 msg_id 与 session_id 是递增的
	lockKey := cache.TimelineMessageLock.Format(k.M{"session_id": sessionId})
	redisSpinLock := distrib_lock.NewSpinLock(mem, lockKey, 5*time.Second, &distrib_lock.SpinOption{Interval: 20 * time.Millisecond, Times: 20})
	if err = redisSpinLock.AcquireWithTimes(); err != nil {
		logging.Errorf(logHead+"acquire fail,lockKey=%v,err=%v", lockKey, err)
		return
	}
	defer func() { _ = redisSpinLock.Release() }()
	logging.Infof(logHead+"acquire success,lockKey=%v", lockKey)

	// message: gen msg_id
	msgId, err := gen_id.MsgId(ctx, mem, sender, receiver)
	if err != nil {
		logging.Errorf(logHead+"gen MsgId error=%v", err)
		return
	}

	// message: gen version_id
	versionId, err := gen_id.MsgVersionId(ctx, &gen_id.MsgVerParams{
		Mem: mem,
		Id1: sender,
		Id2: receiver,
	})
	if err != nil {
		logging.Errorf(logHead+"gen VersionId error=%v", err)
		return
	}

	// build message
	msg = &model.Message{
		MsgID:         msgId,                         // 唯一id（服务端）
		SeqID:         req.SeqId,                     // 唯一id（客户端）
		MsgType:       uint32(req.MsgBody.MsgType),   // 消息类型
		Content:       string(bContent),              // 消息内容
		SessionID:     sessionId,                     // 会话ID
		SenderID:      req.Sender.Id(),               // 发送者ID
		SenderType:    uint32(req.Sender.Type()),     // 发送者的用户类型
		VersionID:     versionId,                     // 版本ID
		SortKey:       versionId,                     // sort_key的值等同于version_id
		Status:        uint32(model.MsgStatusNormal), // 状态正常
		HasRead:       uint32(model.MsgRead),         // 已读（功能还没做好）
		InvisibleList: string(bInvisibleList),        // 不可见的列表
	}
	err = b.repoMessage.Create(logHead, msg)
	if err != nil {
		return
	}

	return
}

// 双方通信时，判断是否需要创建对方的Contact
func (b *MessageUseCase) needCreateContact(logHead string, contactId *gen_id.ComponentId) bool {
	logHead += fmt.Sprintf("doNotNeedCreateContact,contactId=%v|", contactId)

	typeArr := []gen_id.ContactIdType{
		gen_id.TypeRobot,
		gen_id.TypeSystem,
	}

	// 如果用户是指定类型，那么不需要创建他的contact信息（比如：机器人）
	if lo.Contains(typeArr, contactId.Type()) || contactId.IsGroup() {
		logging.Info(logHead + "do not need create contact")
		return false
	}

	return true
}

// 限制：发送者和接受者的类型
func (b *MessageUseCase) checkMessageSend(ctx context.Context, req *request.MessageSendReq) error {
	allowSenderType := []gen_id.ContactIdType{
		gen_id.TypeUser,
		gen_id.TypeRobot,
		gen_id.TypeSystem,
	}

	allowReceiverType := []gen_id.ContactIdType{
		gen_id.TypeUser,
		gen_id.TypeRobot,
		gen_id.TypeGroup,
	}

	if !lo.Contains(allowSenderType, req.Sender.Type()) {
		return api.ErrSenderTypeNotAllow
	}
	if !lo.Contains(allowReceiverType, req.Receiver.Type()) {
		return api.ErrReceiverTypeNotAllow
	}

	return nil
}

// updateMsgIdStatus 通用的方法，用于更新消息的状态和版本ID
func (b *MessageUseCase) updateMsgIdStatus(ctx context.Context, logHead string, msgId model.BigIntType, status model.MsgStatus, senderId *gen_id.ComponentId) (err error) {
	logHead += fmt.Sprintf("updateMsgIdStatus,msgId=%v,status=%v|", msgId, status)

	// get: message
	msgInfo, err := b.repoMessage.Info(msgId)
	if err != nil {
		logging.Errorf(logHead+"repoMessage Info error=%v", err)
		return
	}
	sessionId := msgInfo.SessionID

	// save to db
	fn := func(versionId uint64) (err error) {
		// update to db
		err = b.repoMessage.UpdateMsgVerAndStatus(logHead, msgId, versionId, status)
		if err != nil {
			return
		}
		return
	}
	err = b.updateMsgVersion(ctx, logHead, sessionId, senderId, fn)
	if err != nil {
		return
	}

	return
}

type Fn func(versionId uint64) (err error)

// updateMsgVersion 加锁 => 生成 version_id => 执行回调函数
func (b *MessageUseCase) updateMsgVersion(ctx context.Context, logHead string, sessionId string, sender *gen_id.ComponentId, fn Fn) (err error) {
	mem := b.repoMessage.RedisClient

	// get: 接收者的信息
	var receiverId *gen_id.ComponentId
	parseResult := gen_id.ParseSessionId(sessionId)
	switch parseResult.Prefix {
	case gen_id.PrefixPair: // 1对1的timeline
		switch {
		case parseResult.IdArr[0].Equal(sender):
			receiverId = parseResult.IdArr[1]
		case parseResult.IdArr[1].Equal(sender):
			receiverId = parseResult.IdArr[0]
		default:
			err = errors.New("can not find peer")
			return
		}
	case gen_id.PrefixGroup: // 群组的timeline
		receiverId = parseResult.IdArr[0]
	default:
		err = errors.New("can not find peer")
		return
	}

	// note: 同一个消息timeline的版本变动，需要加锁
	lockKey := cache.TimelineMessageLock.Format(k.M{"session_id": sessionId})
	redisSpinLock := distrib_lock.NewSpinLock(mem, lockKey, 5*time.Second, &distrib_lock.SpinOption{Interval: 20 * time.Millisecond, Times: 20})
	if err = redisSpinLock.AcquireWithTimes(); err != nil {
		logging.Errorf(logHead+"acquire fail,lockKey=%v,err=%v", lockKey, err)
		return
	}
	defer func() { _ = redisSpinLock.Release() }()
	logging.Infof(logHead+"acquire success,lockKey=%v", lockKey)

	// message: gen version_id
	versionId, err := gen_id.MsgVersionId(ctx, &gen_id.MsgVerParams{
		Mem: mem,
		Id1: sender,
		Id2: receiverId,
	})
	if err != nil {
		logging.Errorf(logHead+"gen VersionId error=%v", err)
		return
	}

	return fn(versionId)
}
